package chaincode

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
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
		expectedError: "key not found",
	},
	{
		name: "No key",
		identity: func() cid.ClientIdentity {
			identity := &testsfakes.FakeTestClientIdentity{}
			identity.GetMSPIDStub = func() (string, error) {
				return BANK_ORG, nil
			}
			return identity
		},
		transientdata: func() (map[string][]byte, error) {
			return map[string][]byte{
				"some_field": []byte("SOME KEY"),
			}, nil
		},
		expectedError: "key not found",
	},
	{
		name: "Wrong key format",
		identity: func() cid.ClientIdentity {
			identity := &testsfakes.FakeTestClientIdentity{}
			identity.GetMSPIDStub = func() (string, error) {
				return BANK_ORG, nil
			}
			return identity
		},
		transientdata: func() (map[string][]byte, error) {
			return map[string][]byte{
				"key": []byte("SOME KEY"),
			}, nil
		},
		expectedError: "failed to unmarshal key: invalid character 'S' looking for beginning of value",
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
				"key": []byte(
					func() []byte {
						sc := &SmartContract{}
						key, _ := sc.GenerateKeyPair(nil)
						raw, _ := base64.StdEncoding.DecodeString(key)
						return raw
					}()),
			}, nil
		},
		expectedError: "",
		expectedPutState: func(key string, value []byte) error {
			if key != BANK_ORG {
				return fmt.Errorf("expected: %v, got: %v", BANK_ORG, key)
			}
			var pk rsa.PublicKey
			if err := json.Unmarshal(value, &pk); err != nil {
				return fmt.Errorf("expectedPutState: failed to unmarshal key: %v", err)
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

			var pk rsa.PrivateKey
			if err := json.Unmarshal(value, &pk); err != nil {
				return fmt.Errorf("expectedPutPrivateData: failed to unmarshal key: %v", err)
			}

			return nil
		},
	},
}

func TestGenerateKeyPair(t *testing.T) {

	sc := SmartContract{}
	r, _ := sc.GenerateKeyPair(nil)
	raw, _ := base64.StdEncoding.DecodeString(r)

	var key rsa.PrivateKey
	assert.NoError(t, json.Unmarshal(raw, &key))
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
