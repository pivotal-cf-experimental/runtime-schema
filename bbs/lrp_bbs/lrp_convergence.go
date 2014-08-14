package lrp_bbs

import (
	"sync"

	"github.com/cloudfoundry-incubator/delta_force/delta_force"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/shared"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/cloudfoundry/storeadapter"
	"github.com/pivotal-golang/lager"
)

type compareAndSwappableDesiredLRP struct {
	OldIndex      uint64
	NewDesiredLRP models.DesiredLRP
}

func (bbs *LRPBBS) ConvergeLRPs() {
	actualsByProcessGuid, err := bbs.pruneActualsWithMissingExecutors()
	if err != nil {
		bbs.logger.Error("failed-to-fetch-and-prune-actual-lrps", err)
		return
	}

	desiredLRPRoot, err := bbs.store.ListRecursively(shared.DesiredLRPSchemaRoot)
	if err != nil && err != storeadapter.ErrorKeyNotFound {
		bbs.logger.Error("failed-to-fetch-desired-lrps", err)
		return
	}

	var desiredLRPsToCAS []compareAndSwappableDesiredLRP
	var keysToDelete []string
	knownDesiredProcessGuids := map[string]bool{}

	for _, node := range desiredLRPRoot.ChildNodes {
		desiredLRP, err := models.NewDesiredLRPFromJSON(node.Value)

		if err != nil {
			bbs.logger.Info("pruning-invalid-desired-lrp-json", lager.Data{
				"error":   err.Error(),
				"payload": node.Value,
			})
			keysToDelete = append(keysToDelete, node.Key)
			continue
		}

		knownDesiredProcessGuids[desiredLRP.ProcessGuid] = true
		actualLRPsForDesired := actualsByProcessGuid[desiredLRP.ProcessGuid]

		if bbs.needsReconciliation(desiredLRP, actualLRPsForDesired) {
			desiredLRPsToCAS = append(desiredLRPsToCAS, compareAndSwappableDesiredLRP{
				OldIndex:      node.Index,
				NewDesiredLRP: desiredLRP,
			})
		}
	}

	stopLRPInstances := bbs.instancesToStop(knownDesiredProcessGuids, actualsByProcessGuid)

	bbs.store.Delete(keysToDelete...)
	bbs.batchCompareAndSwapDesiredLRPs(desiredLRPsToCAS)
	err = bbs.RequestStopLRPInstances(stopLRPInstances)
	if err != nil {
		bbs.logger.Error("failed-to-request-stops", err)
	}
}

func (bbs *LRPBBS) instancesToStop(knownDesiredProcessGuids map[string]bool, actualsByProcessGuid map[string][]models.ActualLRP) []models.StopLRPInstance {
	var stopLRPInstances []models.StopLRPInstance

	for processGuid, actuals := range actualsByProcessGuid {
		if !knownDesiredProcessGuids[processGuid] {
			for _, actual := range actuals {
				bbs.logger.Info("detected-undesired-process", lager.Data{
					"process-guid":  processGuid,
					"instance-guid": actual.InstanceGuid,
					"index":         actual.Index,
				})

				stopLRPInstances = append(stopLRPInstances, models.StopLRPInstance{
					ProcessGuid:  processGuid,
					InstanceGuid: actual.InstanceGuid,
					Index:        actual.Index,
				})
			}
		}
	}

	return stopLRPInstances
}

func (bbs *LRPBBS) needsReconciliation(desiredLRP models.DesiredLRP, actualLRPsForDesired []models.ActualLRP) bool {
	var actuals delta_force.ActualInstances
	for _, actualLRP := range actualLRPsForDesired {
		actuals = append(actuals, delta_force.ActualInstance{
			Index: actualLRP.Index,
			Guid:  actualLRP.InstanceGuid,
		})
	}
	result := delta_force.Reconcile(desiredLRP.Instances, actuals)

	if len(result.IndicesToStart) > 0 {
		bbs.logger.Info("detected-missing-instance", lager.Data{
			"process-guid":      desiredLRP.ProcessGuid,
			"desired-instances": desiredLRP.Instances,
			"missing-indices":   result.IndicesToStart,
		})
	}

	if len(result.GuidsToStop) > 0 {
		bbs.logger.Info("detected-extra-instance", lager.Data{
			"process-guid":      desiredLRP.ProcessGuid,
			"desired-instances": desiredLRP.Instances,
			"extra-guids":       result.GuidsToStop,
		})
	}

	if len(result.IndicesToStopAllButOne) > 0 {
		bbs.logger.Info("detected-duplicate-instance", lager.Data{
			"process-guid":       desiredLRP.ProcessGuid,
			"desired-instances":  desiredLRP.Instances,
			"duplicated-indices": result.IndicesToStopAllButOne,
		})
	}

	return !result.Empty()
}

func (bbs *LRPBBS) pruneActualsWithMissingExecutors() (map[string][]models.ActualLRP, error) {
	keysToDelete := []string{}
	actualsByProcessGuid := map[string][]models.ActualLRP{}

	dirsToKeep := map[string]struct{}{}
	indexDirsToPossiblyDelete := map[string]struct{}{}
	processDirsToPossiblyDelete := map[string]struct{}{}

	executorState, err := bbs.store.ListRecursively(shared.ExecutorSchemaRoot)
	if err == storeadapter.ErrorKeyNotFound {
		executorState = storeadapter.StoreNode{}
	} else if err != nil {
		bbs.logger.Error("failed-to-get-executors", err)
		return nil, err
	}

	actualRoot, err := bbs.store.ListRecursively(shared.ActualLRPSchemaRoot)
	if err == storeadapter.ErrorKeyNotFound {
		return actualsByProcessGuid, nil
	} else if err != nil {
		bbs.logger.Error("failed-to-get-actual-lrps", err)
		return nil, err
	}

	for _, processDir := range actualRoot.ChildNodes {
		if len(processDir.ChildNodes) == 0 {
			keysToDelete = append(keysToDelete, processDir.Key)
			continue
		}

		for _, indexDir := range processDir.ChildNodes {
			if len(processDir.ChildNodes) == 0 {
				keysToDelete = append(keysToDelete, indexDir.Key)
				processDirsToPossiblyDelete[processDir.Key] = struct{}{}
				continue
			}

			for _, instanceNode := range indexDir.ChildNodes {
				actual, err := models.NewActualLRPFromJSON(instanceNode.Value)
				if err != nil {
					keysToDelete = append(keysToDelete, instanceNode.Key)
					processDirsToPossiblyDelete[processDir.Key] = struct{}{}
					indexDirsToPossiblyDelete[indexDir.Key] = struct{}{}
					continue
				}

				if _, ok := executorState.Lookup(actual.ExecutorID); ok {
					actualsByProcessGuid[actual.ProcessGuid] = append(actualsByProcessGuid[actual.ProcessGuid], actual)

					dirsToKeep[processDir.Key] = struct{}{}
					dirsToKeep[indexDir.Key] = struct{}{}
				} else {
					bbs.logger.Info("detected-actual-with-missing-executor", lager.Data{
						"actual":      actual,
						"executor-id": actual.ExecutorID,
					})
					keysToDelete = append(keysToDelete, shared.ActualLRPSchemaPath(actual.ProcessGuid, actual.Index, actual.InstanceGuid))

					processDirsToPossiblyDelete[processDir.Key] = struct{}{}
					indexDirsToPossiblyDelete[indexDir.Key] = struct{}{}
				}
			}
		}
	}

	indexDirsToDelete := []string{}
	for key := range indexDirsToPossiblyDelete {
		if _, ok := dirsToKeep[key]; !ok {
			indexDirsToDelete = append(indexDirsToDelete, key)
		}
	}

	processDirsToDelete := []string{}
	for key := range processDirsToPossiblyDelete {
		if _, ok := dirsToKeep[key]; !ok {
			processDirsToDelete = append(processDirsToDelete, key)
		}
	}

	bbs.store.DeleteLeaves(keysToDelete...)
	bbs.store.DeleteLeaves(indexDirsToDelete...)
	bbs.store.DeleteLeaves(processDirsToDelete...)

	return actualsByProcessGuid, nil
}

func (bbs *LRPBBS) batchCompareAndSwapDesiredLRPs(desiredLRPsToCAS []compareAndSwappableDesiredLRP) {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(desiredLRPsToCAS))
	for _, desiredLRPToCAS := range desiredLRPsToCAS {
		desiredLRP := desiredLRPToCAS.NewDesiredLRP
		newStoreNode := storeadapter.StoreNode{
			Key:   shared.DesiredLRPSchemaPath(desiredLRP),
			Value: desiredLRP.ToJSON(),
		}

		go func(desiredLRPToCAS compareAndSwappableDesiredLRP, newStoreNode storeadapter.StoreNode) {
			err := bbs.store.CompareAndSwapByIndex(desiredLRPToCAS.OldIndex, newStoreNode)
			if err != nil {
				bbs.logger.Error("failed-to-compare-and-swap", err)
			}

			waitGroup.Done()
		}(desiredLRPToCAS, newStoreNode)
	}

	waitGroup.Wait()
}
