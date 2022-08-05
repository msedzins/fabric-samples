package chaincode

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//Anonymous payments extension
//Functionality responsible for configuration

const BANK_ORG = "Org1MSP"
const BANK_PDC = "_implicit_org_" + BANK_ORG
const BANK_ACCOUNT = "account_" + BANK_ORG
const DEBIT_PROOF = "debit_" + BANK_ORG

//NOTE: Call to this function must not generate blockchain transaction ("query", not "invoke")
//Otherwise private key will be stored on-chain and revealed to everyone
func (s *SmartContract) GenerateKeyPair(ctx contractapi.TransactionContextInterface) (string, error) {

	//Org1MSP act as a bank and is the only one entitled to generate signing keys
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != BANK_ORG {
		return "", errors.New("client is not authorized to call GenerateKeyPair")
	}

	// Generate a key
	keysize := 2048
	//Line below is not deterministic, not a problem because we use only one peer
	key, _ := rsa.GenerateKey(rand.Reader, keysize)

	raw, err := json.Marshal(key)
	if err != nil {
		return "", fmt.Errorf("failed to marshal key: %v", err)
	}

	return base64.StdEncoding.EncodeToString(raw), nil
}

// //This request should go only to peer(s) that belong to Org1MSP (to avoid revealing the data)
func (s *SmartContract) SavePrivateKey(ctx contractapi.TransactionContextInterface) error {

	// Org1MSP act as a bank and is the only one entitled to store signing keys
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != BANK_ORG {
		return fmt.Errorf("client is not authorized to call SavePrivateKey")
	}

	//ContractAPI doesn't support transient map....
	//We must use transient map so that private key is not revealed
	tr, err := ctx.GetStub().GetTransient()
	if err != nil {
		return fmt.Errorf("failed to get Transient field: %v", err)
	}
	key, ok := tr["key"]
	if !ok {
		return errors.New("key not found")
	}

	// 	//private key goes to implicit private data collection
	// 	//access control must be implemented in the chaincode!
	if err = ctx.GetStub().PutPrivateData(BANK_PDC, BANK_ORG, key); err != nil {
		return fmt.Errorf("failed to put private key: %v", err)
	}

	return nil
}

func (s *SmartContract) SavePublicKey(ctx contractapi.TransactionContextInterface, public string) error {
	// Org1MSP act as a bank and is the only one entitled to store signing keys
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != BANK_ORG {
		return fmt.Errorf("client is not authorized to call SavePublicKey")
	}

	if err = ctx.GetStub().PutState(BANK_ORG, []byte(public)); err != nil {
		return fmt.Errorf("failed to put public key: %v", err)
	}

	return nil
}

func (s *SmartContract) GetPrivateKey(ctx contractapi.TransactionContextInterface) (string, error) {

	// Org1MSP act as a bank and is the only one entitled to generate signing keys
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != BANK_ORG {
		return "", fmt.Errorf("client is not authorized to call GetPrivateKey")
	}

	key, err := ctx.GetStub().GetPrivateData(BANK_PDC, BANK_ORG)
	if err != nil {
		return "", fmt.Errorf("failed to get private data: %v", err)
	}

	return string(key), nil
}

func (s *SmartContract) SetBankAccount(ctx contractapi.TransactionContextInterface, account string) error {

	// Org1MSP act as a bank and is the only one entitled to set up bank Account
	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSPID: %v", err)
	}
	if clientMSPID != BANK_ORG {
		return errors.New("client is not authorized to call SetBankAccount")
	}

	if err = ctx.GetStub().PutState(BANK_ACCOUNT, []byte(account)); err != nil {
		return fmt.Errorf("failed to put public key: %v", err)
	}

	return nil
}
