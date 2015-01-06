// This file was generated by counterfeiter
package fake_bbs

import (
	"sync"
	"time"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
)

type FakeConvergerBBS struct {
	NewConvergeLockStub        func(convergerID string, interval time.Duration) ifrit.Runner
	newConvergeLockMutex       sync.RWMutex
	newConvergeLockArgsForCall []struct {
		convergerID string
		interval    time.Duration
	}
	newConvergeLockReturns struct {
		result1 ifrit.Runner
	}
	ConvergeLRPsStub        func(lager.Logger, time.Duration)
	convergeLRPsMutex       sync.RWMutex
	convergeLRPsArgsForCall []struct {
		arg1 lager.Logger
		arg2 time.Duration
	}
	ConvergeTasksStub        func(logger lager.Logger, timeToClaim, convergenceInterval, timeToResolve time.Duration)
	convergeTasksMutex       sync.RWMutex
	convergeTasksArgsForCall []struct {
		logger              lager.Logger
		timeToClaim         time.Duration
		convergenceInterval time.Duration
		timeToResolve       time.Duration
	}
}

func (fake *FakeConvergerBBS) NewConvergeLock(convergerID string, interval time.Duration) ifrit.Runner {
	fake.newConvergeLockMutex.Lock()
	fake.newConvergeLockArgsForCall = append(fake.newConvergeLockArgsForCall, struct {
		convergerID string
		interval    time.Duration
	}{convergerID, interval})
	fake.newConvergeLockMutex.Unlock()
	if fake.NewConvergeLockStub != nil {
		return fake.NewConvergeLockStub(convergerID, interval)
	} else {
		return fake.newConvergeLockReturns.result1
	}
}

func (fake *FakeConvergerBBS) NewConvergeLockCallCount() int {
	fake.newConvergeLockMutex.RLock()
	defer fake.newConvergeLockMutex.RUnlock()
	return len(fake.newConvergeLockArgsForCall)
}

func (fake *FakeConvergerBBS) NewConvergeLockArgsForCall(i int) (string, time.Duration) {
	fake.newConvergeLockMutex.RLock()
	defer fake.newConvergeLockMutex.RUnlock()
	return fake.newConvergeLockArgsForCall[i].convergerID, fake.newConvergeLockArgsForCall[i].interval
}

func (fake *FakeConvergerBBS) NewConvergeLockReturns(result1 ifrit.Runner) {
	fake.NewConvergeLockStub = nil
	fake.newConvergeLockReturns = struct {
		result1 ifrit.Runner
	}{result1}
}

func (fake *FakeConvergerBBS) ConvergeLRPs(arg1 lager.Logger, arg2 time.Duration) {
	fake.convergeLRPsMutex.Lock()
	fake.convergeLRPsArgsForCall = append(fake.convergeLRPsArgsForCall, struct {
		arg1 lager.Logger
		arg2 time.Duration
	}{arg1, arg2})
	fake.convergeLRPsMutex.Unlock()
	if fake.ConvergeLRPsStub != nil {
		fake.ConvergeLRPsStub(arg1, arg2)
	}
}

func (fake *FakeConvergerBBS) ConvergeLRPsCallCount() int {
	fake.convergeLRPsMutex.RLock()
	defer fake.convergeLRPsMutex.RUnlock()
	return len(fake.convergeLRPsArgsForCall)
}

func (fake *FakeConvergerBBS) ConvergeLRPsArgsForCall(i int) (lager.Logger, time.Duration) {
	fake.convergeLRPsMutex.RLock()
	defer fake.convergeLRPsMutex.RUnlock()
	return fake.convergeLRPsArgsForCall[i].arg1, fake.convergeLRPsArgsForCall[i].arg2
}

func (fake *FakeConvergerBBS) ConvergeTasks(logger lager.Logger, timeToClaim time.Duration, convergenceInterval time.Duration, timeToResolve time.Duration) {
	fake.convergeTasksMutex.Lock()
	fake.convergeTasksArgsForCall = append(fake.convergeTasksArgsForCall, struct {
		logger              lager.Logger
		timeToClaim         time.Duration
		convergenceInterval time.Duration
		timeToResolve       time.Duration
	}{logger, timeToClaim, convergenceInterval, timeToResolve})
	fake.convergeTasksMutex.Unlock()
	if fake.ConvergeTasksStub != nil {
		fake.ConvergeTasksStub(logger, timeToClaim, convergenceInterval, timeToResolve)
	}
}

func (fake *FakeConvergerBBS) ConvergeTasksCallCount() int {
	fake.convergeTasksMutex.RLock()
	defer fake.convergeTasksMutex.RUnlock()
	return len(fake.convergeTasksArgsForCall)
}

func (fake *FakeConvergerBBS) ConvergeTasksArgsForCall(i int) (lager.Logger, time.Duration, time.Duration, time.Duration) {
	fake.convergeTasksMutex.RLock()
	defer fake.convergeTasksMutex.RUnlock()
	return fake.convergeTasksArgsForCall[i].logger, fake.convergeTasksArgsForCall[i].timeToClaim, fake.convergeTasksArgsForCall[i].convergenceInterval, fake.convergeTasksArgsForCall[i].timeToResolve
}

var _ bbs.ConvergerBBS = new(FakeConvergerBBS)
