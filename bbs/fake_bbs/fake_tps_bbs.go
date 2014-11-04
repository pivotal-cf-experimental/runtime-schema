// This file was generated by counterfeiter
package fake_bbs

import (
	"sync"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
)

type FakeTPSBBS struct {
	GetActualLRPsByProcessGuidStub        func(string) ([]models.ActualLRP, error)
	getActualLRPsByProcessGuidMutex       sync.RWMutex
	getActualLRPsByProcessGuidArgsForCall []struct {
		arg1 string
	}
	getActualLRPsByProcessGuidReturns struct {
		result1 []models.ActualLRP
		result2 error
	}
}

func (fake *FakeTPSBBS) GetActualLRPsByProcessGuid(arg1 string) ([]models.ActualLRP, error) {
	fake.getActualLRPsByProcessGuidMutex.Lock()
	fake.getActualLRPsByProcessGuidArgsForCall = append(fake.getActualLRPsByProcessGuidArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.getActualLRPsByProcessGuidMutex.Unlock()
	if fake.GetActualLRPsByProcessGuidStub != nil {
		return fake.GetActualLRPsByProcessGuidStub(arg1)
	} else {
		return fake.getActualLRPsByProcessGuidReturns.result1, fake.getActualLRPsByProcessGuidReturns.result2
	}
}

func (fake *FakeTPSBBS) GetActualLRPsByProcessGuidCallCount() int {
	fake.getActualLRPsByProcessGuidMutex.RLock()
	defer fake.getActualLRPsByProcessGuidMutex.RUnlock()
	return len(fake.getActualLRPsByProcessGuidArgsForCall)
}

func (fake *FakeTPSBBS) GetActualLRPsByProcessGuidArgsForCall(i int) string {
	fake.getActualLRPsByProcessGuidMutex.RLock()
	defer fake.getActualLRPsByProcessGuidMutex.RUnlock()
	return fake.getActualLRPsByProcessGuidArgsForCall[i].arg1
}

func (fake *FakeTPSBBS) GetActualLRPsByProcessGuidReturns(result1 []models.ActualLRP, result2 error) {
	fake.GetActualLRPsByProcessGuidStub = nil
	fake.getActualLRPsByProcessGuidReturns = struct {
		result1 []models.ActualLRP
		result2 error
	}{result1, result2}
}

var _ bbs.TPSBBS = new(FakeTPSBBS)