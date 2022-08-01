// Code generated by counterfeiter. DO NOT EDIT.
package testsfakes

import (
	"crypto/x509"
	"sync"

	"github.com/hyperledger/fabric-samples/token-erc-20/chaincode-go/chaincode/tests"
)

type FakeTestClientIdentity struct {
	AssertAttributeValueStub        func(string, string) error
	assertAttributeValueMutex       sync.RWMutex
	assertAttributeValueArgsForCall []struct {
		arg1 string
		arg2 string
	}
	assertAttributeValueReturns struct {
		result1 error
	}
	assertAttributeValueReturnsOnCall map[int]struct {
		result1 error
	}
	GetAttributeValueStub        func(string) (string, bool, error)
	getAttributeValueMutex       sync.RWMutex
	getAttributeValueArgsForCall []struct {
		arg1 string
	}
	getAttributeValueReturns struct {
		result1 string
		result2 bool
		result3 error
	}
	getAttributeValueReturnsOnCall map[int]struct {
		result1 string
		result2 bool
		result3 error
	}
	GetIDStub        func() (string, error)
	getIDMutex       sync.RWMutex
	getIDArgsForCall []struct {
	}
	getIDReturns struct {
		result1 string
		result2 error
	}
	getIDReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	GetMSPIDStub        func() (string, error)
	getMSPIDMutex       sync.RWMutex
	getMSPIDArgsForCall []struct {
	}
	getMSPIDReturns struct {
		result1 string
		result2 error
	}
	getMSPIDReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	GetX509CertificateStub        func() (*x509.Certificate, error)
	getX509CertificateMutex       sync.RWMutex
	getX509CertificateArgsForCall []struct {
	}
	getX509CertificateReturns struct {
		result1 *x509.Certificate
		result2 error
	}
	getX509CertificateReturnsOnCall map[int]struct {
		result1 *x509.Certificate
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTestClientIdentity) AssertAttributeValue(arg1 string, arg2 string) error {
	fake.assertAttributeValueMutex.Lock()
	ret, specificReturn := fake.assertAttributeValueReturnsOnCall[len(fake.assertAttributeValueArgsForCall)]
	fake.assertAttributeValueArgsForCall = append(fake.assertAttributeValueArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.AssertAttributeValueStub
	fakeReturns := fake.assertAttributeValueReturns
	fake.recordInvocation("AssertAttributeValue", []interface{}{arg1, arg2})
	fake.assertAttributeValueMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeTestClientIdentity) AssertAttributeValueCallCount() int {
	fake.assertAttributeValueMutex.RLock()
	defer fake.assertAttributeValueMutex.RUnlock()
	return len(fake.assertAttributeValueArgsForCall)
}

func (fake *FakeTestClientIdentity) AssertAttributeValueCalls(stub func(string, string) error) {
	fake.assertAttributeValueMutex.Lock()
	defer fake.assertAttributeValueMutex.Unlock()
	fake.AssertAttributeValueStub = stub
}

func (fake *FakeTestClientIdentity) AssertAttributeValueArgsForCall(i int) (string, string) {
	fake.assertAttributeValueMutex.RLock()
	defer fake.assertAttributeValueMutex.RUnlock()
	argsForCall := fake.assertAttributeValueArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeTestClientIdentity) AssertAttributeValueReturns(result1 error) {
	fake.assertAttributeValueMutex.Lock()
	defer fake.assertAttributeValueMutex.Unlock()
	fake.AssertAttributeValueStub = nil
	fake.assertAttributeValueReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTestClientIdentity) AssertAttributeValueReturnsOnCall(i int, result1 error) {
	fake.assertAttributeValueMutex.Lock()
	defer fake.assertAttributeValueMutex.Unlock()
	fake.AssertAttributeValueStub = nil
	if fake.assertAttributeValueReturnsOnCall == nil {
		fake.assertAttributeValueReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.assertAttributeValueReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTestClientIdentity) GetAttributeValue(arg1 string) (string, bool, error) {
	fake.getAttributeValueMutex.Lock()
	ret, specificReturn := fake.getAttributeValueReturnsOnCall[len(fake.getAttributeValueArgsForCall)]
	fake.getAttributeValueArgsForCall = append(fake.getAttributeValueArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetAttributeValueStub
	fakeReturns := fake.getAttributeValueReturns
	fake.recordInvocation("GetAttributeValue", []interface{}{arg1})
	fake.getAttributeValueMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeTestClientIdentity) GetAttributeValueCallCount() int {
	fake.getAttributeValueMutex.RLock()
	defer fake.getAttributeValueMutex.RUnlock()
	return len(fake.getAttributeValueArgsForCall)
}

func (fake *FakeTestClientIdentity) GetAttributeValueCalls(stub func(string) (string, bool, error)) {
	fake.getAttributeValueMutex.Lock()
	defer fake.getAttributeValueMutex.Unlock()
	fake.GetAttributeValueStub = stub
}

func (fake *FakeTestClientIdentity) GetAttributeValueArgsForCall(i int) string {
	fake.getAttributeValueMutex.RLock()
	defer fake.getAttributeValueMutex.RUnlock()
	argsForCall := fake.getAttributeValueArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTestClientIdentity) GetAttributeValueReturns(result1 string, result2 bool, result3 error) {
	fake.getAttributeValueMutex.Lock()
	defer fake.getAttributeValueMutex.Unlock()
	fake.GetAttributeValueStub = nil
	fake.getAttributeValueReturns = struct {
		result1 string
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeTestClientIdentity) GetAttributeValueReturnsOnCall(i int, result1 string, result2 bool, result3 error) {
	fake.getAttributeValueMutex.Lock()
	defer fake.getAttributeValueMutex.Unlock()
	fake.GetAttributeValueStub = nil
	if fake.getAttributeValueReturnsOnCall == nil {
		fake.getAttributeValueReturnsOnCall = make(map[int]struct {
			result1 string
			result2 bool
			result3 error
		})
	}
	fake.getAttributeValueReturnsOnCall[i] = struct {
		result1 string
		result2 bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeTestClientIdentity) GetID() (string, error) {
	fake.getIDMutex.Lock()
	ret, specificReturn := fake.getIDReturnsOnCall[len(fake.getIDArgsForCall)]
	fake.getIDArgsForCall = append(fake.getIDArgsForCall, struct {
	}{})
	stub := fake.GetIDStub
	fakeReturns := fake.getIDReturns
	fake.recordInvocation("GetID", []interface{}{})
	fake.getIDMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeTestClientIdentity) GetIDCallCount() int {
	fake.getIDMutex.RLock()
	defer fake.getIDMutex.RUnlock()
	return len(fake.getIDArgsForCall)
}

func (fake *FakeTestClientIdentity) GetIDCalls(stub func() (string, error)) {
	fake.getIDMutex.Lock()
	defer fake.getIDMutex.Unlock()
	fake.GetIDStub = stub
}

func (fake *FakeTestClientIdentity) GetIDReturns(result1 string, result2 error) {
	fake.getIDMutex.Lock()
	defer fake.getIDMutex.Unlock()
	fake.GetIDStub = nil
	fake.getIDReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeTestClientIdentity) GetIDReturnsOnCall(i int, result1 string, result2 error) {
	fake.getIDMutex.Lock()
	defer fake.getIDMutex.Unlock()
	fake.GetIDStub = nil
	if fake.getIDReturnsOnCall == nil {
		fake.getIDReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getIDReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeTestClientIdentity) GetMSPID() (string, error) {
	fake.getMSPIDMutex.Lock()
	ret, specificReturn := fake.getMSPIDReturnsOnCall[len(fake.getMSPIDArgsForCall)]
	fake.getMSPIDArgsForCall = append(fake.getMSPIDArgsForCall, struct {
	}{})
	stub := fake.GetMSPIDStub
	fakeReturns := fake.getMSPIDReturns
	fake.recordInvocation("GetMSPID", []interface{}{})
	fake.getMSPIDMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeTestClientIdentity) GetMSPIDCallCount() int {
	fake.getMSPIDMutex.RLock()
	defer fake.getMSPIDMutex.RUnlock()
	return len(fake.getMSPIDArgsForCall)
}

func (fake *FakeTestClientIdentity) GetMSPIDCalls(stub func() (string, error)) {
	fake.getMSPIDMutex.Lock()
	defer fake.getMSPIDMutex.Unlock()
	fake.GetMSPIDStub = stub
}

func (fake *FakeTestClientIdentity) GetMSPIDReturns(result1 string, result2 error) {
	fake.getMSPIDMutex.Lock()
	defer fake.getMSPIDMutex.Unlock()
	fake.GetMSPIDStub = nil
	fake.getMSPIDReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeTestClientIdentity) GetMSPIDReturnsOnCall(i int, result1 string, result2 error) {
	fake.getMSPIDMutex.Lock()
	defer fake.getMSPIDMutex.Unlock()
	fake.GetMSPIDStub = nil
	if fake.getMSPIDReturnsOnCall == nil {
		fake.getMSPIDReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getMSPIDReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeTestClientIdentity) GetX509Certificate() (*x509.Certificate, error) {
	fake.getX509CertificateMutex.Lock()
	ret, specificReturn := fake.getX509CertificateReturnsOnCall[len(fake.getX509CertificateArgsForCall)]
	fake.getX509CertificateArgsForCall = append(fake.getX509CertificateArgsForCall, struct {
	}{})
	stub := fake.GetX509CertificateStub
	fakeReturns := fake.getX509CertificateReturns
	fake.recordInvocation("GetX509Certificate", []interface{}{})
	fake.getX509CertificateMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeTestClientIdentity) GetX509CertificateCallCount() int {
	fake.getX509CertificateMutex.RLock()
	defer fake.getX509CertificateMutex.RUnlock()
	return len(fake.getX509CertificateArgsForCall)
}

func (fake *FakeTestClientIdentity) GetX509CertificateCalls(stub func() (*x509.Certificate, error)) {
	fake.getX509CertificateMutex.Lock()
	defer fake.getX509CertificateMutex.Unlock()
	fake.GetX509CertificateStub = stub
}

func (fake *FakeTestClientIdentity) GetX509CertificateReturns(result1 *x509.Certificate, result2 error) {
	fake.getX509CertificateMutex.Lock()
	defer fake.getX509CertificateMutex.Unlock()
	fake.GetX509CertificateStub = nil
	fake.getX509CertificateReturns = struct {
		result1 *x509.Certificate
		result2 error
	}{result1, result2}
}

func (fake *FakeTestClientIdentity) GetX509CertificateReturnsOnCall(i int, result1 *x509.Certificate, result2 error) {
	fake.getX509CertificateMutex.Lock()
	defer fake.getX509CertificateMutex.Unlock()
	fake.GetX509CertificateStub = nil
	if fake.getX509CertificateReturnsOnCall == nil {
		fake.getX509CertificateReturnsOnCall = make(map[int]struct {
			result1 *x509.Certificate
			result2 error
		})
	}
	fake.getX509CertificateReturnsOnCall[i] = struct {
		result1 *x509.Certificate
		result2 error
	}{result1, result2}
}

func (fake *FakeTestClientIdentity) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.assertAttributeValueMutex.RLock()
	defer fake.assertAttributeValueMutex.RUnlock()
	fake.getAttributeValueMutex.RLock()
	defer fake.getAttributeValueMutex.RUnlock()
	fake.getIDMutex.RLock()
	defer fake.getIDMutex.RUnlock()
	fake.getMSPIDMutex.RLock()
	defer fake.getMSPIDMutex.RUnlock()
	fake.getX509CertificateMutex.RLock()
	defer fake.getX509CertificateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTestClientIdentity) recordInvocation(key string, args []interface{}) {
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

var _ tests.TestClientIdentity = new(FakeTestClientIdentity)
