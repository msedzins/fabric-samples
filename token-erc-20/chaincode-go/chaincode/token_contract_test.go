package chaincode

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-samples/token-erc-20/chaincode-go/chaincode/tests/testsfakes"
	"github.com/stretchr/testify/assert"
)

var _SaveKeyPair = []struct {
	name                   string
	identity               func() cid.ClientIdentity
	transientdata          func() (map[string][]byte, error)
	expectedError          string
	expectedPutState       func(key string, value []byte) error
	expectedPutPrivateData func(collection string, key string, value []byte) error
}{
	{
		name: "Wrong organisation",
		identity: func() cid.ClientIdentity {
			identity := &testsfakes.FakeTestClientIdentity{}
			identity.GetMSPIDStub = func() (string, error) {
				return "WRONG_ORG", nil
			}
			return identity
		},
		expectedError: "client is not authorized to call SaveKeyPair",
	},
	{
		name: "No transient data",
		identity: func() cid.ClientIdentity {
			identity := &testsfakes.FakeTestClientIdentity{}
			identity.GetMSPIDStub = func() (string, error) {
				return BANK_ORG, nil
			}
			return identity
		},
		expectedError: "public key not found",
	},
	{
		name: "No private key",
		identity: func() cid.ClientIdentity {
			identity := &testsfakes.FakeTestClientIdentity{}
			identity.GetMSPIDStub = func() (string, error) {
				return BANK_ORG, nil
			}
			return identity
		},
		transientdata: func() (map[string][]byte, error) {
			return map[string][]byte{
				"public": []byte("SOME KEY"),
			}, nil
		},
		expectedError: "private key not found",
	},
	{
		name: "OK",
		identity: func() cid.ClientIdentity {
			identity := &testsfakes.FakeTestClientIdentity{}
			identity.GetMSPIDStub = func() (string, error) {
				return BANK_ORG, nil
			}
			return identity
		},
		transientdata: func() (map[string][]byte, error) {
			return map[string][]byte{
				"public":  []byte("SOME PUBLIC KEY"),
				"private": []byte("SOME PRIVATE KEY"),
			}, nil
		},
		expectedError: "",
		expectedPutState: func(key string, value []byte) error {
			if key != BANK_ORG {
				return fmt.Errorf("expected: %v, got: %v", BANK_ORG, key)
			}
			if string(value) != "SOME PUBLIC KEY" {
				return fmt.Errorf("expected: %v, got: %v", "SOME PUBLIC KEY", string(value))
			}
			return nil
		},
		expectedPutPrivateData: func(collection string, key string, value []byte) error {
			if collection != BANK_PDC {
				return fmt.Errorf("expected: %v, got: %v", BANK_PDC, collection)
			}
			if key != BANK_ORG {
				return fmt.Errorf("expected: %v, got: %v", BANK_ORG, key)
			}
			if string(value) != "SOME PRIVATE KEY" {
				return fmt.Errorf("expected: %v, got: %v", "SOME PRIVATE KEY", string(value))
			}
			return nil
		},
	},
}

func TestSaveKeyPair(t *testing.T) {

	//Prepare fixed data
	sc := SmartContract{}
	stub := &testsfakes.FakeTestChaincodeStubInterface{}
	tc := &testsfakes.FakeTestTransactionContextInterface{}
	tc.GetStubStub = func() shim.ChaincodeStubInterface {
		return stub
	}

	for _, tt := range _SaveKeyPair {
		t.Run(tt.name, func(t *testing.T) {

			//Prepare dynamic data
			tc.GetClientIdentityStub = tt.identity
			stub.GetTransientStub = tt.transientdata
			stub.PutStateStub = tt.expectedPutState
			stub.PutPrivateDataStub = tt.expectedPutPrivateData

			err := sc.SaveKeyPair(tc)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
