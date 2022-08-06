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

/*
Private + Public key used in tests (base64 encoded)
export KEY="eyJOIjoyNDc4NzE5NTI3NjkzMDY0OTIzMDI4NzI1ODIyNDM0MDkzNzgxNzEzNDY2NzU0ODEyMjk5MjU3MTY4NzkyNjcwMDUyMzc5MTkxODk5NTAyMjM3MTM5OTY4MDQyNDYwMzE4NjcwNTYzMjkyNjI4MzQwMDAzMDE0MjIyOTU1NTI5ODU4NzcxNzYyMjI0NTc1ODAxNzAwOTYxMjUzMTA2NDI4MDk5ODI1NDc1NjgxMTAyMzQxNTMwMzk3OTg1NjcxMDQyMzE1OTgwNzQ3ODI0Nzg5NTYzODM3MTM1Nzg0NTg0MDE2ODc4MTAwMTk5NjE5NjY0MTk0MTI0NTE2ODY4NTE4MzgwMTk2Njc2Mzk4NjQxMDU4NDk1Mzc1MzkzNTQ5MzUzODgwODgyNzAwNDY0Njg3ODk4NDcyNDc2NDU3ODM5ODMxMjg4NzUzODA0Mjg3MzI3NDQ1Mjc5Njg1Mjk2NTA1MjcxNDY4NzI5NDExNzUwMDM2MTYwMjAxMjczMjEzODMzNzY5OTQ5NDAzOTc2ODc5MTgwOTE0MDQ0ODE0MTg1NzgxNzQxNjUxNjA4Nzk5MzgyNTkyMDc2MjU0ODE3NTQ0ODQyNzQ5NDgzNTY1ODc4ODc5MDU5ODIwMjUwODI0MTI0MjU3NDM1ODIwMTM5NzUwNzg4ODgxOTUwODQzODkzMzg4MTg3MzYzNzk3OTk1NzAyNjYzODkxMTc5MDU2OTYyODA4NDk4MjU0MDQ4NDI3MjA4NjY5MDQyNzk0MzIzMjk4OTg1MjMyNTI0Mzg1NTk1OTc2MjQ1ODIwMDM0MzM3NDU5MjM0NDg4OSwiRSI6NjU1MzcsIkQiOjE1MTI4NjcyNTIyMDQ0NDMyNDQ1OTY5MzA0ODA0NTE3MTA1MDM1MTAwNTc5ODU0NTA4NDQxMDc3MDYzNTk4NzAwMjkwNzAxMDgxMjE4MjU2Nzg0MDQ1NDU4NTQ5ODg5Nzk0NTkyNzUzOTcxODIyNTczNTI2NDkxNzQwMjcyMDg5NzEyODE0MjEwMTM4MjQ3NDEyNjEyNDg2MDk1NTI3Mjc4MzMwMjU5MDk0NDQxOTMyNTk3MDM5MjE3NDUzNzgwNTY1NTc5NjMxMzQwMDA1MTI0MjEyODA4ODYwODU0Njg2NzY4NTYxODgyMjkzMTc1ODQ5OTQyMDE0NTM3NTQ5ODkzMjIxODIzODYzNzYzMzE1NzM2MTIwODE5NjUwMjY5OTQ3NTQ2NDcxMTU2NTM1NjIzNTMxMzg1MjE5OTA2NjY5Mzg1NzEzMjM2MTk3OTUzOTYyMzAxMDI1MDY5NjQ5MTI1Njk1NTA1MTg1Njg0NTM4MjU5MDc0MDM5MjI4NzgwODIyNjE0MDY5NzAzNjQ3MzMzNjgyMDA3NzAwOTkyMTE1MjAwNjUwOTU3NTgyMDkyMjI2Mzg0MzY1MDg0NTM5NDM5OTUwMzI5NTI4MzA0OTgyNDc4MjI1MTUzMTA5NzQzODcxMDcyMzM4NjA2ODYzOTg4MjMzMjIyODM1NzM2MDEwNzEzNzU0ODM2NjA5OTUyODg1NjY1MDg0MjM0NTIyNzU1MjU0NzU2NjA0ODg5MTI5NTU3MjIyMjc4NDQwOTU3MjgzMTU4MzkwOTkzNDQwMDk1NjU3NTQ5MDQzMjMwODkzNjcyOTkzLCJQcmltZXMiOlsxNjA0NTQ2NDM0OTg1MTM1MDg1NzE0MDc5MzY1NTk2OTQ0NzgzNDMxNjcwMjM0Mjg4OTExNDkzMTUxMTg0MTIxOTgyNTQ5NjgwMjk4MTM0MzcwMTM2Mzk4MDkwNDE3MTcwMTM4NzE0NjQ1OTU4NTA0NDg5ODY0MTQ3NDY3MzAzMTk1ODU0MDIyMzkxNzUyMDUxNzY5MjI4NzUwODE2MDcxMjE4NjMzODMwNDc2NTA3NzU5MjcwODkxMzk0Mzc5NzMxNzA3MzM3Njc0OTUyNjU2NjQwNzg3Njc2MDk5Nzc2NDcxNzM5OTkzODc5NDM0OTkyODM2MDM1NTUwNTUxNTQyNTM5MzcxNDIwODM1MDY2NzQ5NzI4MTIwNzM0OTQyNTE3NTg1NTc3NTI5ODY3MDgwMDYzNjcsMTU0NDgxMDA5MzE1MDExMDkxNzU3NjE3MzAzOTk0MDc1MjYxNTEwMDU4MTYwNTE4MTM4MDYwOTQ3MTkyMzUyMjgwMzc1MDQwNjkzMzA0OTcxMDIzMjQ1NzYxMzY3NDA0MTgzMTA5MTc5MjU1NjE5MzQwOTg0NDk2ODA3NjM5NTg3MzczNzgyMjk5NzQxMzIxNjk5MjU1MzcwNjE2MjY2NTkwMjkyNjQxMzYxNDg5ODM2NzM0NzMwMzA2OTM5OTk0NjMwNjUwMDMwNjY1MjYyNTI1NzkxMDM5Njg0OTkyMjk1OTIyMTU0MjMzNzM5MzA3NDI0MzQ2NTYzNzQ0ODM3NjQ0NTc5MjQ2NTI0MDQ2NTE4NTY1NjY4MzU1Mzc2NzgyOTc5ODUxNDM0ODEwOTE4MTY0OTY3XSwiUHJlY29tcHV0ZWQiOnsiRHAiOjkwMTQ2NjM0MzIyODI5MzU0MTkwNzUwODc2OTcyMjEzNDE2NzM1NTE0NDM5MjEyMjI3NzgxNTI0Njc1NTI1ODQyNDk3MzM2MjA0ODU3MjY3NjYzMTg1OTUyNDk4ODMzMDMyNDk2ODY5MzQ3Mzc5NTQ5NDQwNDY1NTUzNDIxODg5NDIzMjk1Mzk3MjAyMDU0NjM1MDA0NjU3ODM0NTc4NTQ2ODgyMDYzMDE1MDA2ODE0MDA3ODk1MTE0NDI1NTMzMjQyOTk4ODc1NzM2OTk4NjY4NzEyMDI4Njc2OTIxMDgyMjczMzIxMjkxMjQxMjc4MTExOTQxMDg1MTQ0NzM5MzAxOTIwNTI2OTAxMDU5MTg0NDY4Mjk5NDQyMjM5NjI1NzAzODAyMzc4Mjc0NDE4ODU5NDksIkRxIjo3MDY1MzQyODk1MTcwODg0MzMxNTExNjk3MzE1ODk1NDY2NTEyODEzMzQ3MTIyMDM4OTU1NDU4NDkwODQyNjgwMDkyMTAyODg4MDQ5Njg2NzQ0MDU0MTUwMjUyODc0ODIzMzU4NTgyOTk3NDAyODkzMjE1NTQxMzA4NDM5Nzk1ODI4MjIxODQ1MTQ0NjQ1NTg0MzQ2OTgwMjk5NDUyMTc5MzQ1MTUxMDMyNDk5NjUxMTM3OTYyOTkyMjMzNzI5NjQ3NDY0OTQ5NTk5NzA3OTE5MTExNDM0MTg3NTg1NTc2MjA1NzYyNjI0MTY5NzA4NzE1Mjg2NTc2ODk4MDY5NDMxODYwODA4OTA5OTQ2NzAyMzkzMjg1MjMzMjAyNzE1NTU0NjMwMzEwOTgwMDkxMzM5MzYwNSwiUWludiI6MTI4Nzg2NTUwMjc2MzgwMjc1NTkxNTc2NDcwNjQ5MjAyOTA1NTQ4NjkzNzM4Nzk3MTAzMDQ1MDMyNjg2NDIxNjE5NTA4MzQ1NDY2NTQwMTg0MDM4MjQ1MjExMzE1NDc2ODc4Mjk5NjQ1NjkwMzU2NDEzNzAwNDY2NDQ4ODE0MjY2ODgyMzA3NTI2NTQxMTU5NTA5MjE3NjQ1MjkyOTQ4NjAzMTk4MDI5MDQ4NTUzNDIwMzQxNDkwOTI5ODE0MjA5MjQxMzE2MzQxMDc1NTgyMzI0Mjg5OTA5NjIzNzA2ODEwOTEwNDgyNjk5MzEwMzA2OTEzNzg2Mjc0MTAxNDM3NzM2NDkxMzY5ODkzNTg2Mjg5NTc3NDgwNzIzOTI3OTg3MjA3NTk3NzUwNjg4OTM1OTI5MDI4LCJDUlRWYWx1ZXMiOltdfX0="
*/

var _UnblindSignature = []struct {
	name             string
	sig              string
	unblinder        string
	expectedError    string
	expectedGetState func(key string) ([]byte, error)
}{
	{
		name: "OK",
		//taken from BlindSignToken
		sig: "nW3KYizdcSlSRPxdH0rn7FPGYdOGd+ygoYJkKxC73d9D5FS/Bs3p3cG/lB+RhPrYMmjRh3DLJyatpnLdGHdYQjf4IOiD1mYeYiC5vYrt0pssjMOwf+I9GwaVY+JnG2YTbjlFqCwJxGyErHLhLq1iP4dAgGdLa2Kqr0wekwYbi1O+Bj+Ne8ho9GF/lv6Hvtf3Ht1YJTL5gxMYGJp9jCB5w2uxi/VqL36km2/57i5ogdIGfgwU97MdPqRNaU1Aoo3r4qvVlVr1XSfFt0bMzbA3re7eNKvb40X4n6Gh2rQ4z0j9ZDa5hA8vjB6qXTluq06vwXDamVgNglWswQ/0PRvN7w==",
		//taken from BlindToken
		unblinder:     "CwOdEOGwzLxu0dE3paIukhZR0dhu54yHLfkm/xkFHVfYB9ml+XcF/p74YGBIilSQEkAUWRq9/UvRQB7egjUSIpbVIDw55pYcE867c0EmcqykBr350ZMmvmVfAVufUAW2hrBuwyAKD0e8TvaHck4TiHX6age6sz0W7Hk0Pe4/E8GDh34v9dPgcNlrJmfHuSTFsCNUQMEVQznX4z2MGp2ftqoiTX3h534TAspNkFld50CLx8kyl82E1UaSYZWK7O0AZY8pBsPEB2PV7SeZM2pU3951dc303DHPZsBpT4zpxyPbxsgbHE3YAect5xD/UlDnY6K5OtihyL7CRT75eqN0vA==",
		expectedError: "",
		expectedGetState: func(key string) ([]byte, error) {
			//{"N":24787195276930649230287258224340937817134667548122992571687926700523791918995022371399680424603186705632926283400030142229555298587717622245758017009612531064280998254756811023415303979856710423159807478247895638371357845840168781001996196641941245168685183801966763986410584953753935493538808827004646878984724764578398312887538042873274452796852965052714687294117500361602012732138337699494039768791809140448141857817416516087993825920762548175448427494835658788790598202508241242574358201397507888819508438933881873637979957026638911790569628084982540484272086690427943232989852325243855959762458200343374592344889,"E":65537}
			return []byte("eyJOIjoyNDc4NzE5NTI3NjkzMDY0OTIzMDI4NzI1ODIyNDM0MDkzNzgxNzEzNDY2NzU0ODEyMjk5MjU3MTY4NzkyNjcwMDUyMzc5MTkxODk5NTAyMjM3MTM5OTY4MDQyNDYwMzE4NjcwNTYzMjkyNjI4MzQwMDAzMDE0MjIyOTU1NTI5ODU4NzcxNzYyMjI0NTc1ODAxNzAwOTYxMjUzMTA2NDI4MDk5ODI1NDc1NjgxMTAyMzQxNTMwMzk3OTg1NjcxMDQyMzE1OTgwNzQ3ODI0Nzg5NTYzODM3MTM1Nzg0NTg0MDE2ODc4MTAwMTk5NjE5NjY0MTk0MTI0NTE2ODY4NTE4MzgwMTk2Njc2Mzk4NjQxMDU4NDk1Mzc1MzkzNTQ5MzUzODgwODgyNzAwNDY0Njg3ODk4NDcyNDc2NDU3ODM5ODMxMjg4NzUzODA0Mjg3MzI3NDQ1Mjc5Njg1Mjk2NTA1MjcxNDY4NzI5NDExNzUwMDM2MTYwMjAxMjczMjEzODMzNzY5OTQ5NDAzOTc2ODc5MTgwOTE0MDQ0ODE0MTg1NzgxNzQxNjUxNjA4Nzk5MzgyNTkyMDc2MjU0ODE3NTQ0ODQyNzQ5NDgzNTY1ODc4ODc5MDU5ODIwMjUwODI0MTI0MjU3NDM1ODIwMTM5NzUwNzg4ODgxOTUwODQzODkzMzg4MTg3MzYzNzk3OTk1NzAyNjYzODkxMTc5MDU2OTYyODA4NDk4MjU0MDQ4NDI3MjA4NjY5MDQyNzk0MzIzMjk4OTg1MjMyNTI0Mzg1NTk1OTc2MjQ1ODIwMDM0MzM3NDU5MjM0NDg4OSwiRSI6NjU1Mzd9Cg=="),
				nil
		},
	},
}

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
			return raw, nil
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
			return raw, nil
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
			return raw, nil
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

func TestUnblindSignature(t *testing.T) {

	//Prepare fixed data
	sc := SmartContract{}
	stub := &testsfakes.FakeTestChaincodeStubInterface{}
	tc := &testsfakes.FakeTestTransactionContextInterface{}
	tc.GetStubStub = func() shim.ChaincodeStubInterface {
		return stub
	}

	for _, tt := range _UnblindSignature {
		t.Run(tt.name, func(t *testing.T) {

			//Prepare dynamic data
			stub.GetStateStub = tt.expectedGetState

			r, err := sc.UnblindSignature(tc, tt.sig, tt.unblinder)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, r)
			}
		})
	}

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
