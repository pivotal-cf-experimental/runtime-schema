package bbserrors

import (
	"errors"
	"fmt"

	"github.com/cloudfoundry-incubator/runtime-schema/models"
)

var (
	ErrActualLRPCannotBeUnclaimed = errors.New("cannot unclaim actual LRP")
	ErrActualLRPCannotBeClaimed   = errors.New("cannot claim actual LRP")
	ErrActualLRPCannotBeStarted   = errors.New("cannot start actual LRP")
	ErrActualLRPCannotBeCrashed   = errors.New("cannot crash actual LRP")
	ErrActualLRPCannotBeRemoved   = errors.New("cannot remove actual LRP")
	ErrActualLRPCannotBeEvacuated = errors.New("cannot evacuate actual LRP")

	ErrStoreResourceNotFound             = errors.New("the requested resource could not be found in the store")
	ErrStoreExpectedNonCollectionRequest = errors.New("unable to access single (non-collection) store resource, request body refers to a collection")
	ErrStoreExpectedCollectionRequest    = errors.New("unable to access collection of store resources, request body refers to a single object")
	ErrStoreTimeout                      = errors.New("store request timed out")
	ErrStoreInvalidFormat                = errors.New("store request has invalid format")
	ErrStoreInvalidTTL                   = errors.New("store request has invalid TTL")
	ErrStoreResourceExists               = errors.New("the requested store resource already exists")
	ErrStoreComparisonFailed             = errors.New("store resource comparison failed")

	ErrNoDomain      = errors.New("no domain given")
	ErrNoProcessGuid = errors.New("no process guid given")
	ErrNoCellID      = errors.New("no cell id given")

	ErrServiceUnavailable = errors.New("service unavailable")

	ErrTaskRunningOnDifferentCell = errors.New("task running on different cell")

	ErrActualLRPCannotBeFailed = errors.New("cannot set placement error on actual LRP")
)

func NewTaskStateTransitionError(from, to models.TaskState) TaskStateTransitionError {
	return TaskStateTransitionError{from, to}
}

type TaskStateTransitionError struct {
	from models.TaskState
	to   models.TaskState
}

func (e TaskStateTransitionError) Error() string {
	return fmt.Sprint("Cannot transition from ", e.from, " to ", e.to)
}

func NewTaskCannotBeResolvedError(from models.TaskState) taskCannotBeResolvedError {
	return taskCannotBeResolvedError{from}
}

type taskCannotBeResolvedError struct {
	from models.TaskState
}

func (e taskCannotBeResolvedError) Error() string {
	return fmt.Sprint("Cannot resolve task from ", e.from, " state")
}
