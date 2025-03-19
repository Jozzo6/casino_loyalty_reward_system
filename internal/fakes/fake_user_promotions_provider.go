// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"context"
	"sync"

	userpromotion "github.com/Jozzo6/casino_loyalty_reward_system/internal/component/user_promotion"
	"github.com/Jozzo6/casino_loyalty_reward_system/internal/types"
	"github.com/google/uuid"
)

type FakeUserPromotionProvider struct {
	AddPromotionStub        func(context.Context, types.UserPromotion) (types.UserPromotion, error)
	addPromotionMutex       sync.RWMutex
	addPromotionArgsForCall []struct {
		arg1 context.Context
		arg2 types.UserPromotion
	}
	addPromotionReturns struct {
		result1 types.UserPromotion
		result2 error
	}
	addPromotionReturnsOnCall map[int]struct {
		result1 types.UserPromotion
		result2 error
	}
	AddWelcomePromotionStub        func(context.Context, uuid.UUID) (types.UserPromotion, error)
	addWelcomePromotionMutex       sync.RWMutex
	addWelcomePromotionArgsForCall []struct {
		arg1 context.Context
		arg2 uuid.UUID
	}
	addWelcomePromotionReturns struct {
		result1 types.UserPromotion
		result2 error
	}
	addWelcomePromotionReturnsOnCall map[int]struct {
		result1 types.UserPromotion
		result2 error
	}
	ClaimPromotionStub        func(context.Context, uuid.UUID) error
	claimPromotionMutex       sync.RWMutex
	claimPromotionArgsForCall []struct {
		arg1 context.Context
		arg2 uuid.UUID
	}
	claimPromotionReturns struct {
		result1 error
	}
	claimPromotionReturnsOnCall map[int]struct {
		result1 error
	}
	DeleteUserPromotionStub        func(context.Context, uuid.UUID) error
	deleteUserPromotionMutex       sync.RWMutex
	deleteUserPromotionArgsForCall []struct {
		arg1 context.Context
		arg2 uuid.UUID
	}
	deleteUserPromotionReturns struct {
		result1 error
	}
	deleteUserPromotionReturnsOnCall map[int]struct {
		result1 error
	}
	GetUserPromotionByIDStub        func(context.Context, uuid.UUID) (types.UserPromotion, error)
	getUserPromotionByIDMutex       sync.RWMutex
	getUserPromotionByIDArgsForCall []struct {
		arg1 context.Context
		arg2 uuid.UUID
	}
	getUserPromotionByIDReturns struct {
		result1 types.UserPromotion
		result2 error
	}
	getUserPromotionByIDReturnsOnCall map[int]struct {
		result1 types.UserPromotion
		result2 error
	}
	GetUserPromotionsStub        func(context.Context, uuid.UUID) ([]types.UserPromotion, error)
	getUserPromotionsMutex       sync.RWMutex
	getUserPromotionsArgsForCall []struct {
		arg1 context.Context
		arg2 uuid.UUID
	}
	getUserPromotionsReturns struct {
		result1 []types.UserPromotion
		result2 error
	}
	getUserPromotionsReturnsOnCall map[int]struct {
		result1 []types.UserPromotion
		result2 error
	}
	ListenToRegisterEventStub        func(context.Context) error
	listenToRegisterEventMutex       sync.RWMutex
	listenToRegisterEventArgsForCall []struct {
		arg1 context.Context
	}
	listenToRegisterEventReturns struct {
		result1 error
	}
	listenToRegisterEventReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeUserPromotionProvider) AddPromotion(arg1 context.Context, arg2 types.UserPromotion) (types.UserPromotion, error) {
	fake.addPromotionMutex.Lock()
	ret, specificReturn := fake.addPromotionReturnsOnCall[len(fake.addPromotionArgsForCall)]
	fake.addPromotionArgsForCall = append(fake.addPromotionArgsForCall, struct {
		arg1 context.Context
		arg2 types.UserPromotion
	}{arg1, arg2})
	stub := fake.AddPromotionStub
	fakeReturns := fake.addPromotionReturns
	fake.recordInvocation("AddPromotion", []interface{}{arg1, arg2})
	fake.addPromotionMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeUserPromotionProvider) AddPromotionCallCount() int {
	fake.addPromotionMutex.RLock()
	defer fake.addPromotionMutex.RUnlock()
	return len(fake.addPromotionArgsForCall)
}

func (fake *FakeUserPromotionProvider) AddPromotionCalls(stub func(context.Context, types.UserPromotion) (types.UserPromotion, error)) {
	fake.addPromotionMutex.Lock()
	defer fake.addPromotionMutex.Unlock()
	fake.AddPromotionStub = stub
}

func (fake *FakeUserPromotionProvider) AddPromotionArgsForCall(i int) (context.Context, types.UserPromotion) {
	fake.addPromotionMutex.RLock()
	defer fake.addPromotionMutex.RUnlock()
	argsForCall := fake.addPromotionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeUserPromotionProvider) AddPromotionReturns(result1 types.UserPromotion, result2 error) {
	fake.addPromotionMutex.Lock()
	defer fake.addPromotionMutex.Unlock()
	fake.AddPromotionStub = nil
	fake.addPromotionReturns = struct {
		result1 types.UserPromotion
		result2 error
	}{result1, result2}
}

func (fake *FakeUserPromotionProvider) AddPromotionReturnsOnCall(i int, result1 types.UserPromotion, result2 error) {
	fake.addPromotionMutex.Lock()
	defer fake.addPromotionMutex.Unlock()
	fake.AddPromotionStub = nil
	if fake.addPromotionReturnsOnCall == nil {
		fake.addPromotionReturnsOnCall = make(map[int]struct {
			result1 types.UserPromotion
			result2 error
		})
	}
	fake.addPromotionReturnsOnCall[i] = struct {
		result1 types.UserPromotion
		result2 error
	}{result1, result2}
}

func (fake *FakeUserPromotionProvider) AddWelcomePromotion(arg1 context.Context, arg2 uuid.UUID) (types.UserPromotion, error) {
	fake.addWelcomePromotionMutex.Lock()
	ret, specificReturn := fake.addWelcomePromotionReturnsOnCall[len(fake.addWelcomePromotionArgsForCall)]
	fake.addWelcomePromotionArgsForCall = append(fake.addWelcomePromotionArgsForCall, struct {
		arg1 context.Context
		arg2 uuid.UUID
	}{arg1, arg2})
	stub := fake.AddWelcomePromotionStub
	fakeReturns := fake.addWelcomePromotionReturns
	fake.recordInvocation("AddWelcomePromotion", []interface{}{arg1, arg2})
	fake.addWelcomePromotionMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeUserPromotionProvider) AddWelcomePromotionCallCount() int {
	fake.addWelcomePromotionMutex.RLock()
	defer fake.addWelcomePromotionMutex.RUnlock()
	return len(fake.addWelcomePromotionArgsForCall)
}

func (fake *FakeUserPromotionProvider) AddWelcomePromotionCalls(stub func(context.Context, uuid.UUID) (types.UserPromotion, error)) {
	fake.addWelcomePromotionMutex.Lock()
	defer fake.addWelcomePromotionMutex.Unlock()
	fake.AddWelcomePromotionStub = stub
}

func (fake *FakeUserPromotionProvider) AddWelcomePromotionArgsForCall(i int) (context.Context, uuid.UUID) {
	fake.addWelcomePromotionMutex.RLock()
	defer fake.addWelcomePromotionMutex.RUnlock()
	argsForCall := fake.addWelcomePromotionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeUserPromotionProvider) AddWelcomePromotionReturns(result1 types.UserPromotion, result2 error) {
	fake.addWelcomePromotionMutex.Lock()
	defer fake.addWelcomePromotionMutex.Unlock()
	fake.AddWelcomePromotionStub = nil
	fake.addWelcomePromotionReturns = struct {
		result1 types.UserPromotion
		result2 error
	}{result1, result2}
}

func (fake *FakeUserPromotionProvider) AddWelcomePromotionReturnsOnCall(i int, result1 types.UserPromotion, result2 error) {
	fake.addWelcomePromotionMutex.Lock()
	defer fake.addWelcomePromotionMutex.Unlock()
	fake.AddWelcomePromotionStub = nil
	if fake.addWelcomePromotionReturnsOnCall == nil {
		fake.addWelcomePromotionReturnsOnCall = make(map[int]struct {
			result1 types.UserPromotion
			result2 error
		})
	}
	fake.addWelcomePromotionReturnsOnCall[i] = struct {
		result1 types.UserPromotion
		result2 error
	}{result1, result2}
}

func (fake *FakeUserPromotionProvider) ClaimPromotion(arg1 context.Context, arg2 uuid.UUID) error {
	fake.claimPromotionMutex.Lock()
	ret, specificReturn := fake.claimPromotionReturnsOnCall[len(fake.claimPromotionArgsForCall)]
	fake.claimPromotionArgsForCall = append(fake.claimPromotionArgsForCall, struct {
		arg1 context.Context
		arg2 uuid.UUID
	}{arg1, arg2})
	stub := fake.ClaimPromotionStub
	fakeReturns := fake.claimPromotionReturns
	fake.recordInvocation("ClaimPromotion", []interface{}{arg1, arg2})
	fake.claimPromotionMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeUserPromotionProvider) ClaimPromotionCallCount() int {
	fake.claimPromotionMutex.RLock()
	defer fake.claimPromotionMutex.RUnlock()
	return len(fake.claimPromotionArgsForCall)
}

func (fake *FakeUserPromotionProvider) ClaimPromotionCalls(stub func(context.Context, uuid.UUID) error) {
	fake.claimPromotionMutex.Lock()
	defer fake.claimPromotionMutex.Unlock()
	fake.ClaimPromotionStub = stub
}

func (fake *FakeUserPromotionProvider) ClaimPromotionArgsForCall(i int) (context.Context, uuid.UUID) {
	fake.claimPromotionMutex.RLock()
	defer fake.claimPromotionMutex.RUnlock()
	argsForCall := fake.claimPromotionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeUserPromotionProvider) ClaimPromotionReturns(result1 error) {
	fake.claimPromotionMutex.Lock()
	defer fake.claimPromotionMutex.Unlock()
	fake.ClaimPromotionStub = nil
	fake.claimPromotionReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeUserPromotionProvider) ClaimPromotionReturnsOnCall(i int, result1 error) {
	fake.claimPromotionMutex.Lock()
	defer fake.claimPromotionMutex.Unlock()
	fake.ClaimPromotionStub = nil
	if fake.claimPromotionReturnsOnCall == nil {
		fake.claimPromotionReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.claimPromotionReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeUserPromotionProvider) DeleteUserPromotion(arg1 context.Context, arg2 uuid.UUID) error {
	fake.deleteUserPromotionMutex.Lock()
	ret, specificReturn := fake.deleteUserPromotionReturnsOnCall[len(fake.deleteUserPromotionArgsForCall)]
	fake.deleteUserPromotionArgsForCall = append(fake.deleteUserPromotionArgsForCall, struct {
		arg1 context.Context
		arg2 uuid.UUID
	}{arg1, arg2})
	stub := fake.DeleteUserPromotionStub
	fakeReturns := fake.deleteUserPromotionReturns
	fake.recordInvocation("DeleteUserPromotion", []interface{}{arg1, arg2})
	fake.deleteUserPromotionMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeUserPromotionProvider) DeleteUserPromotionCallCount() int {
	fake.deleteUserPromotionMutex.RLock()
	defer fake.deleteUserPromotionMutex.RUnlock()
	return len(fake.deleteUserPromotionArgsForCall)
}

func (fake *FakeUserPromotionProvider) DeleteUserPromotionCalls(stub func(context.Context, uuid.UUID) error) {
	fake.deleteUserPromotionMutex.Lock()
	defer fake.deleteUserPromotionMutex.Unlock()
	fake.DeleteUserPromotionStub = stub
}

func (fake *FakeUserPromotionProvider) DeleteUserPromotionArgsForCall(i int) (context.Context, uuid.UUID) {
	fake.deleteUserPromotionMutex.RLock()
	defer fake.deleteUserPromotionMutex.RUnlock()
	argsForCall := fake.deleteUserPromotionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeUserPromotionProvider) DeleteUserPromotionReturns(result1 error) {
	fake.deleteUserPromotionMutex.Lock()
	defer fake.deleteUserPromotionMutex.Unlock()
	fake.DeleteUserPromotionStub = nil
	fake.deleteUserPromotionReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeUserPromotionProvider) DeleteUserPromotionReturnsOnCall(i int, result1 error) {
	fake.deleteUserPromotionMutex.Lock()
	defer fake.deleteUserPromotionMutex.Unlock()
	fake.DeleteUserPromotionStub = nil
	if fake.deleteUserPromotionReturnsOnCall == nil {
		fake.deleteUserPromotionReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteUserPromotionReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeUserPromotionProvider) GetUserPromotionByID(arg1 context.Context, arg2 uuid.UUID) (types.UserPromotion, error) {
	fake.getUserPromotionByIDMutex.Lock()
	ret, specificReturn := fake.getUserPromotionByIDReturnsOnCall[len(fake.getUserPromotionByIDArgsForCall)]
	fake.getUserPromotionByIDArgsForCall = append(fake.getUserPromotionByIDArgsForCall, struct {
		arg1 context.Context
		arg2 uuid.UUID
	}{arg1, arg2})
	stub := fake.GetUserPromotionByIDStub
	fakeReturns := fake.getUserPromotionByIDReturns
	fake.recordInvocation("GetUserPromotionByID", []interface{}{arg1, arg2})
	fake.getUserPromotionByIDMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeUserPromotionProvider) GetUserPromotionByIDCallCount() int {
	fake.getUserPromotionByIDMutex.RLock()
	defer fake.getUserPromotionByIDMutex.RUnlock()
	return len(fake.getUserPromotionByIDArgsForCall)
}

func (fake *FakeUserPromotionProvider) GetUserPromotionByIDCalls(stub func(context.Context, uuid.UUID) (types.UserPromotion, error)) {
	fake.getUserPromotionByIDMutex.Lock()
	defer fake.getUserPromotionByIDMutex.Unlock()
	fake.GetUserPromotionByIDStub = stub
}

func (fake *FakeUserPromotionProvider) GetUserPromotionByIDArgsForCall(i int) (context.Context, uuid.UUID) {
	fake.getUserPromotionByIDMutex.RLock()
	defer fake.getUserPromotionByIDMutex.RUnlock()
	argsForCall := fake.getUserPromotionByIDArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeUserPromotionProvider) GetUserPromotionByIDReturns(result1 types.UserPromotion, result2 error) {
	fake.getUserPromotionByIDMutex.Lock()
	defer fake.getUserPromotionByIDMutex.Unlock()
	fake.GetUserPromotionByIDStub = nil
	fake.getUserPromotionByIDReturns = struct {
		result1 types.UserPromotion
		result2 error
	}{result1, result2}
}

func (fake *FakeUserPromotionProvider) GetUserPromotionByIDReturnsOnCall(i int, result1 types.UserPromotion, result2 error) {
	fake.getUserPromotionByIDMutex.Lock()
	defer fake.getUserPromotionByIDMutex.Unlock()
	fake.GetUserPromotionByIDStub = nil
	if fake.getUserPromotionByIDReturnsOnCall == nil {
		fake.getUserPromotionByIDReturnsOnCall = make(map[int]struct {
			result1 types.UserPromotion
			result2 error
		})
	}
	fake.getUserPromotionByIDReturnsOnCall[i] = struct {
		result1 types.UserPromotion
		result2 error
	}{result1, result2}
}

func (fake *FakeUserPromotionProvider) GetUserPromotions(arg1 context.Context, arg2 uuid.UUID) ([]types.UserPromotion, error) {
	fake.getUserPromotionsMutex.Lock()
	ret, specificReturn := fake.getUserPromotionsReturnsOnCall[len(fake.getUserPromotionsArgsForCall)]
	fake.getUserPromotionsArgsForCall = append(fake.getUserPromotionsArgsForCall, struct {
		arg1 context.Context
		arg2 uuid.UUID
	}{arg1, arg2})
	stub := fake.GetUserPromotionsStub
	fakeReturns := fake.getUserPromotionsReturns
	fake.recordInvocation("GetUserPromotions", []interface{}{arg1, arg2})
	fake.getUserPromotionsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeUserPromotionProvider) GetUserPromotionsCallCount() int {
	fake.getUserPromotionsMutex.RLock()
	defer fake.getUserPromotionsMutex.RUnlock()
	return len(fake.getUserPromotionsArgsForCall)
}

func (fake *FakeUserPromotionProvider) GetUserPromotionsCalls(stub func(context.Context, uuid.UUID) ([]types.UserPromotion, error)) {
	fake.getUserPromotionsMutex.Lock()
	defer fake.getUserPromotionsMutex.Unlock()
	fake.GetUserPromotionsStub = stub
}

func (fake *FakeUserPromotionProvider) GetUserPromotionsArgsForCall(i int) (context.Context, uuid.UUID) {
	fake.getUserPromotionsMutex.RLock()
	defer fake.getUserPromotionsMutex.RUnlock()
	argsForCall := fake.getUserPromotionsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeUserPromotionProvider) GetUserPromotionsReturns(result1 []types.UserPromotion, result2 error) {
	fake.getUserPromotionsMutex.Lock()
	defer fake.getUserPromotionsMutex.Unlock()
	fake.GetUserPromotionsStub = nil
	fake.getUserPromotionsReturns = struct {
		result1 []types.UserPromotion
		result2 error
	}{result1, result2}
}

func (fake *FakeUserPromotionProvider) GetUserPromotionsReturnsOnCall(i int, result1 []types.UserPromotion, result2 error) {
	fake.getUserPromotionsMutex.Lock()
	defer fake.getUserPromotionsMutex.Unlock()
	fake.GetUserPromotionsStub = nil
	if fake.getUserPromotionsReturnsOnCall == nil {
		fake.getUserPromotionsReturnsOnCall = make(map[int]struct {
			result1 []types.UserPromotion
			result2 error
		})
	}
	fake.getUserPromotionsReturnsOnCall[i] = struct {
		result1 []types.UserPromotion
		result2 error
	}{result1, result2}
}

func (fake *FakeUserPromotionProvider) ListenToRegisterEvent(arg1 context.Context) error {
	fake.listenToRegisterEventMutex.Lock()
	ret, specificReturn := fake.listenToRegisterEventReturnsOnCall[len(fake.listenToRegisterEventArgsForCall)]
	fake.listenToRegisterEventArgsForCall = append(fake.listenToRegisterEventArgsForCall, struct {
		arg1 context.Context
	}{arg1})
	stub := fake.ListenToRegisterEventStub
	fakeReturns := fake.listenToRegisterEventReturns
	fake.recordInvocation("ListenToRegisterEvent", []interface{}{arg1})
	fake.listenToRegisterEventMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeUserPromotionProvider) ListenToRegisterEventCallCount() int {
	fake.listenToRegisterEventMutex.RLock()
	defer fake.listenToRegisterEventMutex.RUnlock()
	return len(fake.listenToRegisterEventArgsForCall)
}

func (fake *FakeUserPromotionProvider) ListenToRegisterEventCalls(stub func(context.Context) error) {
	fake.listenToRegisterEventMutex.Lock()
	defer fake.listenToRegisterEventMutex.Unlock()
	fake.ListenToRegisterEventStub = stub
}

func (fake *FakeUserPromotionProvider) ListenToRegisterEventArgsForCall(i int) context.Context {
	fake.listenToRegisterEventMutex.RLock()
	defer fake.listenToRegisterEventMutex.RUnlock()
	argsForCall := fake.listenToRegisterEventArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeUserPromotionProvider) ListenToRegisterEventReturns(result1 error) {
	fake.listenToRegisterEventMutex.Lock()
	defer fake.listenToRegisterEventMutex.Unlock()
	fake.ListenToRegisterEventStub = nil
	fake.listenToRegisterEventReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeUserPromotionProvider) ListenToRegisterEventReturnsOnCall(i int, result1 error) {
	fake.listenToRegisterEventMutex.Lock()
	defer fake.listenToRegisterEventMutex.Unlock()
	fake.ListenToRegisterEventStub = nil
	if fake.listenToRegisterEventReturnsOnCall == nil {
		fake.listenToRegisterEventReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.listenToRegisterEventReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeUserPromotionProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addPromotionMutex.RLock()
	defer fake.addPromotionMutex.RUnlock()
	fake.addWelcomePromotionMutex.RLock()
	defer fake.addWelcomePromotionMutex.RUnlock()
	fake.claimPromotionMutex.RLock()
	defer fake.claimPromotionMutex.RUnlock()
	fake.deleteUserPromotionMutex.RLock()
	defer fake.deleteUserPromotionMutex.RUnlock()
	fake.getUserPromotionByIDMutex.RLock()
	defer fake.getUserPromotionByIDMutex.RUnlock()
	fake.getUserPromotionsMutex.RLock()
	defer fake.getUserPromotionsMutex.RUnlock()
	fake.listenToRegisterEventMutex.RLock()
	defer fake.listenToRegisterEventMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeUserPromotionProvider) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ userpromotion.UserPromotionProvider = new(FakeUserPromotionProvider)
