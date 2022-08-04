package chaincode

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-samples/token-erc-20/chaincode-go/chaincode/tests/testsfakes"
	"github.com/stretchr/testify/assert"
)

var _GenerateKeyPair = []struct {
	name          string
	identity      func() cid.ClientIdentity
	expectedError string
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
		expectedError: "client is not authorized to call GenerateKeyPair",
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
		expectedError: "",
	},
}

var _SavePrivateKey = []struct {
	name                   string
	identity               func() cid.ClientIdentity
	transientdata          func() (map[string][]byte, error)
	expectedError          string
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
		expectedError: "client is not authorized to call SavePrivateKey",
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
						keysize := 2048
						key, _ := rsa.GenerateKey(rand.Reader, keysize)
						raw, _ := json.Marshal(key)
						return raw
					}()),
			}, nil
		},
		expectedError: "",
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

var _SavePublicKey = []struct {
	name             string
	identity         func() cid.ClientIdentity
	publicKey        string
	expectedError    string
	expectedPutState func(key string, value []byte) error
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
		expectedError: "client is not authorized to call SavePublicKey",
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
		publicKey:     "base64key",
		expectedError: "",
		expectedPutState: func(key string, value []byte) error {
			if key != BANK_ORG {
				return fmt.Errorf("expected: %v, got: %v", BANK_ORG, key)
			}

			expected := "base64key"
			if expected != string(value) {
				return errors.New("expectedPutState: Wrong input parameter")
			}

			return nil
		},
	},
}

var _GetPrivateKey = []struct {
	name                   string
	identity               func() cid.ClientIdentity
	expectedError          string
	expectedOut            string
	expectedGetPrivateData func(collection, key string) ([]byte, error)
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
		expectedError: "client is not authorized to call GetPrivateKey",
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
		expectedError: "",
		expectedOut:   "PRIVATE_KEY",
		expectedGetPrivateData: func(collection, key string) ([]byte, error) {
			if collection != BANK_PDC {
				return nil, fmt.Errorf("expected: %v, got: %v", BANK_PDC, collection)
			}
			if key != BANK_ORG {
				return nil, fmt.Errorf("expected: %v, got: %v", BANK_ORG, key)
			}

			return []byte("PRIVATE_KEY"), nil
		},
	},
}

func TestGenerateKeyPair(t *testing.T) {

	//Prepare fixed data
	sc := SmartContract{}
	stub := &testsfakes.FakeTestChaincodeStubInterface{}
	tc := &testsfakes.FakeTestTransactionContextInterface{}
	tc.GetStubStub = func() shim.ChaincodeStubInterface {
		return stub
	}

	for _, tt := range _GenerateKeyPair {
		t.Run(tt.name, func(t *testing.T) {

			//Prepare dynamic data
			tc.GetClientIdentityStub = tt.identity

			r, err := sc.GenerateKeyPair(tc)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)

				raw, _ := base64.StdEncoding.DecodeString(r)
				var key rsa.PrivateKey
				assert.NoError(t, json.Unmarshal(raw, &key))
			}
		})
	}

}

func TestSavePrivateKey(t *testing.T) {

	//Prepare fixed data
	sc := SmartContract{}
	stub := &testsfakes.FakeTestChaincodeStubInterface{}
	tc := &testsfakes.FakeTestTransactionContextInterface{}
	tc.GetStubStub = func() shim.ChaincodeStubInterface {
		return stub
	}

	for _, tt := range _SavePrivateKey {
		t.Run(tt.name, func(t *testing.T) {

			//Prepare dynamic data
			tc.GetClientIdentityStub = tt.identity
			stub.GetTransientStub = tt.transientdata
			stub.PutPrivateDataStub = tt.expectedPutPrivateData

			err := sc.SavePrivateKey(tc)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestSavePublicKey(t *testing.T) {

	//Prepare fixed data
	sc := SmartContract{}
	stub := &testsfakes.FakeTestChaincodeStubInterface{}
	tc := &testsfakes.FakeTestTransactionContextInterface{}
	tc.GetStubStub = func() shim.ChaincodeStubInterface {
		return stub
	}

	for _, tt := range _SavePublicKey {
		t.Run(tt.name, func(t *testing.T) {

			//Prepare dynamic data
			tc.GetClientIdentityStub = tt.identity
			stub.PutStateStub = tt.expectedPutState

			err := sc.SavePublicKey(tc, tt.publicKey)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestGetPrivateKey(t *testing.T) {

	//Prepare fixed data
	sc := SmartContract{}
	stub := &testsfakes.FakeTestChaincodeStubInterface{}
	tc := &testsfakes.FakeTestTransactionContextInterface{}
	tc.GetStubStub = func() shim.ChaincodeStubInterface {
		return stub
	}

	for _, tt := range _GetPrivateKey {
		t.Run(tt.name, func(t *testing.T) {

			//Prepare dynamic data
			tc.GetClientIdentityStub = tt.identity
			stub.GetPrivateDataStub = tt.expectedGetPrivateData

			key, err := sc.GetPrivateKey(tc)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOut, key)
			}

		})
	}
}
