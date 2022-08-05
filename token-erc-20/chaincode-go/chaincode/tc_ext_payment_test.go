package chaincode

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-samples/token-erc-20/chaincode-go/chaincode/tests/testsfakes"
	"github.com/stretchr/testify/assert"
)

var _BlindSignToken = []struct {
	name                   string
	inBlindedMessage       string
	expectedError          string
	expectedGetPrivateData func(collection, key string) ([]byte, error)
	expectedGetState       func(key string) ([]byte, error)
}{
	{
		name:          "Not paid",
		expectedError: "token not paid. please call DebitMyAccount first",
		expectedGetPrivateData: func(collection, key string) ([]byte, error) {

			keysize := 2048
			k, _ := rsa.GenerateKey(rand.Reader, keysize)
			raw, _ := json.Marshal(k)
			return []byte(base64.StdEncoding.EncodeToString(raw)), nil
		},
		expectedGetState: func(key string) ([]byte, error) {
			return []byte(""), nil
		},
	},
	{
		name:             "Bad input parameter",
		inBlindedMessage: "BAD_MESSAGE",
		expectedError:    "failed to decode blinded message: illegal base64 data at input byte 3",
		expectedGetPrivateData: func(collection, key string) ([]byte, error) {

			keysize := 2048
			k, _ := rsa.GenerateKey(rand.Reader, keysize)
			raw, _ := json.Marshal(k)
			return []byte(base64.StdEncoding.EncodeToString(raw)), nil
		},
		expectedGetState: func(key string) ([]byte, error) {
			return []byte("DEBIT_PROOF"), nil
		},
	},
	{
		name: "OK",
		//response from BlindToken
		inBlindedMessage: "HvJ06atFNIHNY8emzzrAOHH2lr5yJqqi60OgzWc4SUEndCbI59NBPYkEOSgWwc8SjaqSRm3uB+i9qsmCy7Y/I6b1xOHaRxxFdVzx8hCmz5tKD1yxTcPQoPJCKM6q38IMvL7YwMOH+FFqA0G7fL1GZC9OVUWaGBiaAE7K64QyQpW+h9kfqK+k5faAz34hSXwmLtRCGRp0rSRS1Y8ctU6AymFQoLVnD+Sd1LPO0nIEWlwtGJJs8X7tOvfEM0Oz8JpIW/4grPfu7kFQlEgT+oRPauWn/1/D/qBTxLR+GvqNvYSfB9itth36irV9WMaSOWno+C6B8agHClqjJqz94Wiigw==",
		expectedError:    "",
		expectedGetPrivateData: func(collection, key string) ([]byte, error) {

			keysize := 2048
			k, _ := rsa.GenerateKey(rand.Reader, keysize)
			raw, _ := json.Marshal(k)
			return []byte(base64.StdEncoding.EncodeToString(raw)), nil
		},
		expectedGetState: func(key string) ([]byte, error) {
			return []byte("DEBIT_PROOF"), nil
		},
	},
}

var _BlindToken = []struct {
	name             string
	expectedError    string
	expectedGetState func(key string) ([]byte, error)
}{
	{
		name:          "Wrong public key",
		expectedError: "failed to decode pubkey : illegal base64 data at input byte 3",
		expectedGetState: func(key string) ([]byte, error) {
			if key != BANK_ORG {
				return nil, fmt.Errorf("expected: %v, got: %v", BANK_ORG, key)
			}

			return []byte("BAD_KEY"), nil
		},
	},
	{
		name:          "Wrong public key",
		expectedError: "failed to unmarshal pubkey : invalid character 'B' looking for beginning of value",
		expectedGetState: func(key string) ([]byte, error) {
			if key != BANK_ORG {
				return nil, fmt.Errorf("expected: %v, got: %v", BANK_ORG, key)
			}

			return []byte(base64.StdEncoding.EncodeToString([]byte("BAD_KEY"))), nil
		},
	},
	{
		name:          "OK",
		expectedError: "",
		expectedGetState: func(key string) ([]byte, error) {
			if key != BANK_ORG {
				return nil, fmt.Errorf("expected: %v, got: %v", BANK_ORG, key)
			}

			//{"N":24787195276930649230287258224340937817134667548122992571687926700523791918995022371399680424603186705632926283400030142229555298587717622245758017009612531064280998254756811023415303979856710423159807478247895638371357845840168781001996196641941245168685183801966763986410584953753935493538808827004646878984724764578398312887538042873274452796852965052714687294117500361602012732138337699494039768791809140448141857817416516087993825920762548175448427494835658788790598202508241242574358201397507888819508438933881873637979957026638911790569628084982540484272086690427943232989852325243855959762458200343374592344889,"E":65537}
			return []byte("eyJOIjoyNDc4NzE5NTI3NjkzMDY0OTIzMDI4NzI1ODIyNDM0MDkzNzgxNzEzNDY2NzU0ODEyMjk5MjU3MTY4NzkyNjcwMDUyMzc5MTkxODk5NTAyMjM3MTM5OTY4MDQyNDYwMzE4NjcwNTYzMjkyNjI4MzQwMDAzMDE0MjIyOTU1NTI5ODU4NzcxNzYyMjI0NTc1ODAxNzAwOTYxMjUzMTA2NDI4MDk5ODI1NDc1NjgxMTAyMzQxNTMwMzk3OTg1NjcxMDQyMzE1OTgwNzQ3ODI0Nzg5NTYzODM3MTM1Nzg0NTg0MDE2ODc4MTAwMTk5NjE5NjY0MTk0MTI0NTE2ODY4NTE4MzgwMTk2Njc2Mzk4NjQxMDU4NDk1Mzc1MzkzNTQ5MzUzODgwODgyNzAwNDY0Njg3ODk4NDcyNDc2NDU3ODM5ODMxMjg4NzUzODA0Mjg3MzI3NDQ1Mjc5Njg1Mjk2NTA1MjcxNDY4NzI5NDExNzUwMDM2MTYwMjAxMjczMjEzODMzNzY5OTQ5NDAzOTc2ODc5MTgwOTE0MDQ0ODE0MTg1NzgxNzQxNjUxNjA4Nzk5MzgyNTkyMDc2MjU0ODE3NTQ0ODQyNzQ5NDgzNTY1ODc4ODc5MDU5ODIwMjUwODI0MTI0MjU3NDM1ODIwMTM5NzUwNzg4ODgxOTUwODQzODkzMzg4MTg3MzYzNzk3OTk1NzAyNjYzODkxMTc5MDU2OTYyODA4NDk4MjU0MDQ4NDI3MjA4NjY5MDQyNzk0MzIzMjk4OTg1MjMyNTI0Mzg1NTk1OTc2MjQ1ODIwMDM0MzM3NDU5MjM0NDg4OSwiRSI6NjU1Mzd9Cg=="),
				nil
		},
	},
}

var _DebitMyAccount = []struct {
	name             string
	expectedError    string
	expectedGetState func(key string) ([]byte, error)
}{
	{
		name:          "Can't debit twice",
		expectedError: "debit operation can be only done once for one blinded token",
		expectedGetState: func(key string) ([]byte, error) {
			if key == BANK_ACCOUNT {
				return []byte("BANK_ACCOUNT"), nil
			}

			return []byte("PROOF_FOUND"), nil
		},
	},
	{
		name:          "Transfer not initialised",
		expectedError: "Contract options need to be set before calling any function, call Initialize() to initialize contract",
		expectedGetState: func(key string) ([]byte, error) {
			if key == BANK_ACCOUNT {
				return []byte("BANK_ACCOUNT"), nil
			}

			return nil, nil
		},
	},
}

func TestBlindSignToken(t *testing.T) {

	//Prepare fixed data
	sc := SmartContract{}
	stub := &testsfakes.FakeTestChaincodeStubInterface{}
	tc := &testsfakes.FakeTestTransactionContextInterface{}
	tc.GetStubStub = func() shim.ChaincodeStubInterface {
		return stub
	}

	for _, tt := range _BlindSignToken {
		t.Run(tt.name, func(t *testing.T) {

			//Prepare dynamic data
			stub.GetPrivateDataStub = tt.expectedGetPrivateData
			stub.GetStateStub = tt.expectedGetState

			r, err := sc.BlindSignToken(tc, tt.inBlindedMessage)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, r)
			}
		})
	}

}

func TestDebitMyAccount(t *testing.T) {

	//Prepare fixed data
	sc := SmartContract{}
	stub := &testsfakes.FakeTestChaincodeStubInterface{}
	tc := &testsfakes.FakeTestTransactionContextInterface{}
	tc.GetStubStub = func() shim.ChaincodeStubInterface {
		return stub
	}

	for _, tt := range _DebitMyAccount {
		t.Run(tt.name, func(t *testing.T) {

			//Prepare dynamic data
			stub.GetStateStub = tt.expectedGetState

			err := sc.DebitMyAccount(tc, "BLINDED")
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestBlindToken(t *testing.T) {

	//Prepare fixed data
	sc := SmartContract{}
	stub := &testsfakes.FakeTestChaincodeStubInterface{}
	tc := &testsfakes.FakeTestTransactionContextInterface{}
	tc.GetStubStub = func() shim.ChaincodeStubInterface {
		return stub
	}

	for _, tt := range _BlindToken {
		t.Run(tt.name, func(t *testing.T) {

			//Prepare dynamic data
			stub.GetStateStub = tt.expectedGetState

			r, err := sc.BlindToken(tc, "UUID")
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)

				var response map[string]interface{}
				json.Unmarshal([]byte(r), &response)

				assert.NotNil(t, response["Blinded"])
				assert.NotNil(t, response["Unblinder"])

				assert.NotEqual(t, "", response["Blinded"])
				assert.NotEqual(t, "", response["Unblinder"])
			}
		})
	}

}
