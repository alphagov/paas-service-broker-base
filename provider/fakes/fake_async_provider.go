// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"context"
	"sync"

	"github.com/alphagov/paas-service-broker-base/provider"
	"github.com/pivotal-cf/brokerapi/domain"
)

type FakeAsyncProvider struct {
	BindStub        func(context.Context, provider.BindData) (*domain.Binding, error)
	bindMutex       sync.RWMutex
	bindArgsForCall []struct {
		arg1 context.Context
		arg2 provider.BindData
	}
	bindReturns struct {
		result1 *domain.Binding
		result2 error
	}
	bindReturnsOnCall map[int]struct {
		result1 *domain.Binding
		result2 error
	}
	DeprovisionStub        func(context.Context, provider.DeprovisionData) (*domain.DeprovisionServiceSpec, error)
	deprovisionMutex       sync.RWMutex
	deprovisionArgsForCall []struct {
		arg1 context.Context
		arg2 provider.DeprovisionData
	}
	deprovisionReturns struct {
		result1 *domain.DeprovisionServiceSpec
		result2 error
	}
	deprovisionReturnsOnCall map[int]struct {
		result1 *domain.DeprovisionServiceSpec
		result2 error
	}
	GetBindingStub        func(context.Context, provider.GetBindData) (*domain.GetBindingSpec, error)
	getBindingMutex       sync.RWMutex
	getBindingArgsForCall []struct {
		arg1 context.Context
		arg2 provider.GetBindData
	}
	getBindingReturns struct {
		result1 *domain.GetBindingSpec
		result2 error
	}
	getBindingReturnsOnCall map[int]struct {
		result1 *domain.GetBindingSpec
		result2 error
	}
	LastBindingOperationStub        func(context.Context, provider.LastOperationData) (*domain.LastOperation, error)
	lastBindingOperationMutex       sync.RWMutex
	lastBindingOperationArgsForCall []struct {
		arg1 context.Context
		arg2 provider.LastOperationData
	}
	lastBindingOperationReturns struct {
		result1 *domain.LastOperation
		result2 error
	}
	lastBindingOperationReturnsOnCall map[int]struct {
		result1 *domain.LastOperation
		result2 error
	}
	LastOperationStub        func(context.Context, provider.LastOperationData) (*domain.LastOperation, error)
	lastOperationMutex       sync.RWMutex
	lastOperationArgsForCall []struct {
		arg1 context.Context
		arg2 provider.LastOperationData
	}
	lastOperationReturns struct {
		result1 *domain.LastOperation
		result2 error
	}
	lastOperationReturnsOnCall map[int]struct {
		result1 *domain.LastOperation
		result2 error
	}
	ProvisionStub        func(context.Context, provider.ProvisionData) (*domain.ProvisionedServiceSpec, error)
	provisionMutex       sync.RWMutex
	provisionArgsForCall []struct {
		arg1 context.Context
		arg2 provider.ProvisionData
	}
	provisionReturns struct {
		result1 *domain.ProvisionedServiceSpec
		result2 error
	}
	provisionReturnsOnCall map[int]struct {
		result1 *domain.ProvisionedServiceSpec
		result2 error
	}
	UnbindStub        func(context.Context, provider.UnbindData) (*domain.UnbindSpec, error)
	unbindMutex       sync.RWMutex
	unbindArgsForCall []struct {
		arg1 context.Context
		arg2 provider.UnbindData
	}
	unbindReturns struct {
		result1 *domain.UnbindSpec
		result2 error
	}
	unbindReturnsOnCall map[int]struct {
		result1 *domain.UnbindSpec
		result2 error
	}
	UpdateStub        func(context.Context, provider.UpdateData) (*domain.UpdateServiceSpec, error)
	updateMutex       sync.RWMutex
	updateArgsForCall []struct {
		arg1 context.Context
		arg2 provider.UpdateData
	}
	updateReturns struct {
		result1 *domain.UpdateServiceSpec
		result2 error
	}
	updateReturnsOnCall map[int]struct {
		result1 *domain.UpdateServiceSpec
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeAsyncProvider) Bind(arg1 context.Context, arg2 provider.BindData) (*domain.Binding, error) {
	fake.bindMutex.Lock()
	ret, specificReturn := fake.bindReturnsOnCall[len(fake.bindArgsForCall)]
	fake.bindArgsForCall = append(fake.bindArgsForCall, struct {
		arg1 context.Context
		arg2 provider.BindData
	}{arg1, arg2})
	fake.recordInvocation("Bind", []interface{}{arg1, arg2})
	fake.bindMutex.Unlock()
	if fake.BindStub != nil {
		return fake.BindStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.bindReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAsyncProvider) BindCallCount() int {
	fake.bindMutex.RLock()
	defer fake.bindMutex.RUnlock()
	return len(fake.bindArgsForCall)
}

func (fake *FakeAsyncProvider) BindCalls(stub func(context.Context, provider.BindData) (*domain.Binding, error)) {
	fake.bindMutex.Lock()
	defer fake.bindMutex.Unlock()
	fake.BindStub = stub
}

func (fake *FakeAsyncProvider) BindArgsForCall(i int) (context.Context, provider.BindData) {
	fake.bindMutex.RLock()
	defer fake.bindMutex.RUnlock()
	argsForCall := fake.bindArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeAsyncProvider) BindReturns(result1 *domain.Binding, result2 error) {
	fake.bindMutex.Lock()
	defer fake.bindMutex.Unlock()
	fake.BindStub = nil
	fake.bindReturns = struct {
		result1 *domain.Binding
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) BindReturnsOnCall(i int, result1 *domain.Binding, result2 error) {
	fake.bindMutex.Lock()
	defer fake.bindMutex.Unlock()
	fake.BindStub = nil
	if fake.bindReturnsOnCall == nil {
		fake.bindReturnsOnCall = make(map[int]struct {
			result1 *domain.Binding
			result2 error
		})
	}
	fake.bindReturnsOnCall[i] = struct {
		result1 *domain.Binding
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) Deprovision(arg1 context.Context, arg2 provider.DeprovisionData) (*domain.DeprovisionServiceSpec, error) {
	fake.deprovisionMutex.Lock()
	ret, specificReturn := fake.deprovisionReturnsOnCall[len(fake.deprovisionArgsForCall)]
	fake.deprovisionArgsForCall = append(fake.deprovisionArgsForCall, struct {
		arg1 context.Context
		arg2 provider.DeprovisionData
	}{arg1, arg2})
	fake.recordInvocation("Deprovision", []interface{}{arg1, arg2})
	fake.deprovisionMutex.Unlock()
	if fake.DeprovisionStub != nil {
		return fake.DeprovisionStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.deprovisionReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAsyncProvider) DeprovisionCallCount() int {
	fake.deprovisionMutex.RLock()
	defer fake.deprovisionMutex.RUnlock()
	return len(fake.deprovisionArgsForCall)
}

func (fake *FakeAsyncProvider) DeprovisionCalls(stub func(context.Context, provider.DeprovisionData) (*domain.DeprovisionServiceSpec, error)) {
	fake.deprovisionMutex.Lock()
	defer fake.deprovisionMutex.Unlock()
	fake.DeprovisionStub = stub
}

func (fake *FakeAsyncProvider) DeprovisionArgsForCall(i int) (context.Context, provider.DeprovisionData) {
	fake.deprovisionMutex.RLock()
	defer fake.deprovisionMutex.RUnlock()
	argsForCall := fake.deprovisionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeAsyncProvider) DeprovisionReturns(result1 *domain.DeprovisionServiceSpec, result2 error) {
	fake.deprovisionMutex.Lock()
	defer fake.deprovisionMutex.Unlock()
	fake.DeprovisionStub = nil
	fake.deprovisionReturns = struct {
		result1 *domain.DeprovisionServiceSpec
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) DeprovisionReturnsOnCall(i int, result1 *domain.DeprovisionServiceSpec, result2 error) {
	fake.deprovisionMutex.Lock()
	defer fake.deprovisionMutex.Unlock()
	fake.DeprovisionStub = nil
	if fake.deprovisionReturnsOnCall == nil {
		fake.deprovisionReturnsOnCall = make(map[int]struct {
			result1 *domain.DeprovisionServiceSpec
			result2 error
		})
	}
	fake.deprovisionReturnsOnCall[i] = struct {
		result1 *domain.DeprovisionServiceSpec
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) GetBinding(arg1 context.Context, arg2 provider.GetBindData) (*domain.GetBindingSpec, error) {
	fake.getBindingMutex.Lock()
	ret, specificReturn := fake.getBindingReturnsOnCall[len(fake.getBindingArgsForCall)]
	fake.getBindingArgsForCall = append(fake.getBindingArgsForCall, struct {
		arg1 context.Context
		arg2 provider.GetBindData
	}{arg1, arg2})
	fake.recordInvocation("GetBinding", []interface{}{arg1, arg2})
	fake.getBindingMutex.Unlock()
	if fake.GetBindingStub != nil {
		return fake.GetBindingStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getBindingReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAsyncProvider) GetBindingCallCount() int {
	fake.getBindingMutex.RLock()
	defer fake.getBindingMutex.RUnlock()
	return len(fake.getBindingArgsForCall)
}

func (fake *FakeAsyncProvider) GetBindingCalls(stub func(context.Context, provider.GetBindData) (*domain.GetBindingSpec, error)) {
	fake.getBindingMutex.Lock()
	defer fake.getBindingMutex.Unlock()
	fake.GetBindingStub = stub
}

func (fake *FakeAsyncProvider) GetBindingArgsForCall(i int) (context.Context, provider.GetBindData) {
	fake.getBindingMutex.RLock()
	defer fake.getBindingMutex.RUnlock()
	argsForCall := fake.getBindingArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeAsyncProvider) GetBindingReturns(result1 *domain.GetBindingSpec, result2 error) {
	fake.getBindingMutex.Lock()
	defer fake.getBindingMutex.Unlock()
	fake.GetBindingStub = nil
	fake.getBindingReturns = struct {
		result1 *domain.GetBindingSpec
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) GetBindingReturnsOnCall(i int, result1 *domain.GetBindingSpec, result2 error) {
	fake.getBindingMutex.Lock()
	defer fake.getBindingMutex.Unlock()
	fake.GetBindingStub = nil
	if fake.getBindingReturnsOnCall == nil {
		fake.getBindingReturnsOnCall = make(map[int]struct {
			result1 *domain.GetBindingSpec
			result2 error
		})
	}
	fake.getBindingReturnsOnCall[i] = struct {
		result1 *domain.GetBindingSpec
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) LastBindingOperation(arg1 context.Context, arg2 provider.LastOperationData) (*domain.LastOperation, error) {
	fake.lastBindingOperationMutex.Lock()
	ret, specificReturn := fake.lastBindingOperationReturnsOnCall[len(fake.lastBindingOperationArgsForCall)]
	fake.lastBindingOperationArgsForCall = append(fake.lastBindingOperationArgsForCall, struct {
		arg1 context.Context
		arg2 provider.LastOperationData
	}{arg1, arg2})
	fake.recordInvocation("LastBindingOperation", []interface{}{arg1, arg2})
	fake.lastBindingOperationMutex.Unlock()
	if fake.LastBindingOperationStub != nil {
		return fake.LastBindingOperationStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.lastBindingOperationReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAsyncProvider) LastBindingOperationCallCount() int {
	fake.lastBindingOperationMutex.RLock()
	defer fake.lastBindingOperationMutex.RUnlock()
	return len(fake.lastBindingOperationArgsForCall)
}

func (fake *FakeAsyncProvider) LastBindingOperationCalls(stub func(context.Context, provider.LastOperationData) (*domain.LastOperation, error)) {
	fake.lastBindingOperationMutex.Lock()
	defer fake.lastBindingOperationMutex.Unlock()
	fake.LastBindingOperationStub = stub
}

func (fake *FakeAsyncProvider) LastBindingOperationArgsForCall(i int) (context.Context, provider.LastOperationData) {
	fake.lastBindingOperationMutex.RLock()
	defer fake.lastBindingOperationMutex.RUnlock()
	argsForCall := fake.lastBindingOperationArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeAsyncProvider) LastBindingOperationReturns(result1 *domain.LastOperation, result2 error) {
	fake.lastBindingOperationMutex.Lock()
	defer fake.lastBindingOperationMutex.Unlock()
	fake.LastBindingOperationStub = nil
	fake.lastBindingOperationReturns = struct {
		result1 *domain.LastOperation
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) LastBindingOperationReturnsOnCall(i int, result1 *domain.LastOperation, result2 error) {
	fake.lastBindingOperationMutex.Lock()
	defer fake.lastBindingOperationMutex.Unlock()
	fake.LastBindingOperationStub = nil
	if fake.lastBindingOperationReturnsOnCall == nil {
		fake.lastBindingOperationReturnsOnCall = make(map[int]struct {
			result1 *domain.LastOperation
			result2 error
		})
	}
	fake.lastBindingOperationReturnsOnCall[i] = struct {
		result1 *domain.LastOperation
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) LastOperation(arg1 context.Context, arg2 provider.LastOperationData) (*domain.LastOperation, error) {
	fake.lastOperationMutex.Lock()
	ret, specificReturn := fake.lastOperationReturnsOnCall[len(fake.lastOperationArgsForCall)]
	fake.lastOperationArgsForCall = append(fake.lastOperationArgsForCall, struct {
		arg1 context.Context
		arg2 provider.LastOperationData
	}{arg1, arg2})
	fake.recordInvocation("LastOperation", []interface{}{arg1, arg2})
	fake.lastOperationMutex.Unlock()
	if fake.LastOperationStub != nil {
		return fake.LastOperationStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.lastOperationReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAsyncProvider) LastOperationCallCount() int {
	fake.lastOperationMutex.RLock()
	defer fake.lastOperationMutex.RUnlock()
	return len(fake.lastOperationArgsForCall)
}

func (fake *FakeAsyncProvider) LastOperationCalls(stub func(context.Context, provider.LastOperationData) (*domain.LastOperation, error)) {
	fake.lastOperationMutex.Lock()
	defer fake.lastOperationMutex.Unlock()
	fake.LastOperationStub = stub
}

func (fake *FakeAsyncProvider) LastOperationArgsForCall(i int) (context.Context, provider.LastOperationData) {
	fake.lastOperationMutex.RLock()
	defer fake.lastOperationMutex.RUnlock()
	argsForCall := fake.lastOperationArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeAsyncProvider) LastOperationReturns(result1 *domain.LastOperation, result2 error) {
	fake.lastOperationMutex.Lock()
	defer fake.lastOperationMutex.Unlock()
	fake.LastOperationStub = nil
	fake.lastOperationReturns = struct {
		result1 *domain.LastOperation
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) LastOperationReturnsOnCall(i int, result1 *domain.LastOperation, result2 error) {
	fake.lastOperationMutex.Lock()
	defer fake.lastOperationMutex.Unlock()
	fake.LastOperationStub = nil
	if fake.lastOperationReturnsOnCall == nil {
		fake.lastOperationReturnsOnCall = make(map[int]struct {
			result1 *domain.LastOperation
			result2 error
		})
	}
	fake.lastOperationReturnsOnCall[i] = struct {
		result1 *domain.LastOperation
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) Provision(arg1 context.Context, arg2 provider.ProvisionData) (*domain.ProvisionedServiceSpec, error) {
	fake.provisionMutex.Lock()
	ret, specificReturn := fake.provisionReturnsOnCall[len(fake.provisionArgsForCall)]
	fake.provisionArgsForCall = append(fake.provisionArgsForCall, struct {
		arg1 context.Context
		arg2 provider.ProvisionData
	}{arg1, arg2})
	fake.recordInvocation("Provision", []interface{}{arg1, arg2})
	fake.provisionMutex.Unlock()
	if fake.ProvisionStub != nil {
		return fake.ProvisionStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.provisionReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAsyncProvider) ProvisionCallCount() int {
	fake.provisionMutex.RLock()
	defer fake.provisionMutex.RUnlock()
	return len(fake.provisionArgsForCall)
}

func (fake *FakeAsyncProvider) ProvisionCalls(stub func(context.Context, provider.ProvisionData) (*domain.ProvisionedServiceSpec, error)) {
	fake.provisionMutex.Lock()
	defer fake.provisionMutex.Unlock()
	fake.ProvisionStub = stub
}

func (fake *FakeAsyncProvider) ProvisionArgsForCall(i int) (context.Context, provider.ProvisionData) {
	fake.provisionMutex.RLock()
	defer fake.provisionMutex.RUnlock()
	argsForCall := fake.provisionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeAsyncProvider) ProvisionReturns(result1 *domain.ProvisionedServiceSpec, result2 error) {
	fake.provisionMutex.Lock()
	defer fake.provisionMutex.Unlock()
	fake.ProvisionStub = nil
	fake.provisionReturns = struct {
		result1 *domain.ProvisionedServiceSpec
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) ProvisionReturnsOnCall(i int, result1 *domain.ProvisionedServiceSpec, result2 error) {
	fake.provisionMutex.Lock()
	defer fake.provisionMutex.Unlock()
	fake.ProvisionStub = nil
	if fake.provisionReturnsOnCall == nil {
		fake.provisionReturnsOnCall = make(map[int]struct {
			result1 *domain.ProvisionedServiceSpec
			result2 error
		})
	}
	fake.provisionReturnsOnCall[i] = struct {
		result1 *domain.ProvisionedServiceSpec
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) Unbind(arg1 context.Context, arg2 provider.UnbindData) (*domain.UnbindSpec, error) {
	fake.unbindMutex.Lock()
	ret, specificReturn := fake.unbindReturnsOnCall[len(fake.unbindArgsForCall)]
	fake.unbindArgsForCall = append(fake.unbindArgsForCall, struct {
		arg1 context.Context
		arg2 provider.UnbindData
	}{arg1, arg2})
	fake.recordInvocation("Unbind", []interface{}{arg1, arg2})
	fake.unbindMutex.Unlock()
	if fake.UnbindStub != nil {
		return fake.UnbindStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.unbindReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAsyncProvider) UnbindCallCount() int {
	fake.unbindMutex.RLock()
	defer fake.unbindMutex.RUnlock()
	return len(fake.unbindArgsForCall)
}

func (fake *FakeAsyncProvider) UnbindCalls(stub func(context.Context, provider.UnbindData) (*domain.UnbindSpec, error)) {
	fake.unbindMutex.Lock()
	defer fake.unbindMutex.Unlock()
	fake.UnbindStub = stub
}

func (fake *FakeAsyncProvider) UnbindArgsForCall(i int) (context.Context, provider.UnbindData) {
	fake.unbindMutex.RLock()
	defer fake.unbindMutex.RUnlock()
	argsForCall := fake.unbindArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeAsyncProvider) UnbindReturns(result1 *domain.UnbindSpec, result2 error) {
	fake.unbindMutex.Lock()
	defer fake.unbindMutex.Unlock()
	fake.UnbindStub = nil
	fake.unbindReturns = struct {
		result1 *domain.UnbindSpec
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) UnbindReturnsOnCall(i int, result1 *domain.UnbindSpec, result2 error) {
	fake.unbindMutex.Lock()
	defer fake.unbindMutex.Unlock()
	fake.UnbindStub = nil
	if fake.unbindReturnsOnCall == nil {
		fake.unbindReturnsOnCall = make(map[int]struct {
			result1 *domain.UnbindSpec
			result2 error
		})
	}
	fake.unbindReturnsOnCall[i] = struct {
		result1 *domain.UnbindSpec
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) Update(arg1 context.Context, arg2 provider.UpdateData) (*domain.UpdateServiceSpec, error) {
	fake.updateMutex.Lock()
	ret, specificReturn := fake.updateReturnsOnCall[len(fake.updateArgsForCall)]
	fake.updateArgsForCall = append(fake.updateArgsForCall, struct {
		arg1 context.Context
		arg2 provider.UpdateData
	}{arg1, arg2})
	fake.recordInvocation("Update", []interface{}{arg1, arg2})
	fake.updateMutex.Unlock()
	if fake.UpdateStub != nil {
		return fake.UpdateStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.updateReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeAsyncProvider) UpdateCallCount() int {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return len(fake.updateArgsForCall)
}

func (fake *FakeAsyncProvider) UpdateCalls(stub func(context.Context, provider.UpdateData) (*domain.UpdateServiceSpec, error)) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = stub
}

func (fake *FakeAsyncProvider) UpdateArgsForCall(i int) (context.Context, provider.UpdateData) {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	argsForCall := fake.updateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeAsyncProvider) UpdateReturns(result1 *domain.UpdateServiceSpec, result2 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	fake.updateReturns = struct {
		result1 *domain.UpdateServiceSpec
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) UpdateReturnsOnCall(i int, result1 *domain.UpdateServiceSpec, result2 error) {
	fake.updateMutex.Lock()
	defer fake.updateMutex.Unlock()
	fake.UpdateStub = nil
	if fake.updateReturnsOnCall == nil {
		fake.updateReturnsOnCall = make(map[int]struct {
			result1 *domain.UpdateServiceSpec
			result2 error
		})
	}
	fake.updateReturnsOnCall[i] = struct {
		result1 *domain.UpdateServiceSpec
		result2 error
	}{result1, result2}
}

func (fake *FakeAsyncProvider) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.bindMutex.RLock()
	defer fake.bindMutex.RUnlock()
	fake.deprovisionMutex.RLock()
	defer fake.deprovisionMutex.RUnlock()
	fake.getBindingMutex.RLock()
	defer fake.getBindingMutex.RUnlock()
	fake.lastBindingOperationMutex.RLock()
	defer fake.lastBindingOperationMutex.RUnlock()
	fake.lastOperationMutex.RLock()
	defer fake.lastOperationMutex.RUnlock()
	fake.provisionMutex.RLock()
	defer fake.provisionMutex.RUnlock()
	fake.unbindMutex.RLock()
	defer fake.unbindMutex.RUnlock()
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeAsyncProvider) recordInvocation(key string, args []interface{}) {
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

var _ provider.AsyncProvider = new(FakeAsyncProvider)
