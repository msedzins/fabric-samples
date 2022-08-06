package chaincode

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cryptoballot/rsablind"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//Anonymous payments extension
//Functionality responsible for payment

const DEBIT_PROOF = "debit_" + BANK_ORG
const CREDIT_PROOF = "credit_" + BANK_ORG

//STEP 0 - Payer hides the messate to be sign
//To not reveal the data the request must go to the peer that is trusted to the payer + no blockchain transaction can be generated
func (s *SmartContract) BlindToken(ctx contractapi.TransactionContextInterface, uuid string) (string, error) {

	pubKey, err := ctx.GetStub().GetState(BANK_ORG)
	if err != nil {
		return "", fmt.Errorf("failed to get pub key: %v", err)
	}

	raw, err := base64.StdEncoding.DecodeString(string(pubKey))
	if err != nil {
		return "", fmt.Errorf("failed to decode pubkey : %v", err)
	}

	var key rsa.PublicKey
	if err = json.Unmarshal(raw, &key); err != nil {
		return "", fmt.Errorf("failed to unmarshal pubkey : %v", err)
	}

	// Blind the hashed message
	blinded, unblinder, err := rsablind.Blind(&key, []byte(uuid))
	if err != nil {
		return "", fmt.Errorf("failed to blind the message : %v", err)
	}

	resp := struct {
		Blinded   string
		Unblinder string
	}{
		Blinded:   base64.StdEncoding.EncodeToString(blinded),
		Unblinder: base64.StdEncoding.EncodeToString(unblinder),
	}

	response, err := json.Marshal(&resp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %v", err)
	}

	return string(response), nil
}

//STEP 1 - Payer debits his account
func (s *SmartContract) DebitMyAccount(ctx contractapi.TransactionContextInterface, blinded string) error {

	account, err := ctx.GetStub().GetState(BANK_ACCOUNT)
	if err != nil {
		return fmt.Errorf("failed to get bank account: %v", err)
	}

	debit, err := ctx.GetStub().GetState(DEBIT_PROOF + blinded)
	if err != nil {
		return fmt.Errorf("failed to get debit proof: %v", err)
	}
	if len(debit) > 0 {
		return errors.New("debit operation can be only done once for one blinded token")
	}

	if err = ctx.GetStub().PutState(DEBIT_PROOF+blinded, []byte(blinded)); err != nil {
		return fmt.Errorf("failed to put message: %v", err)
	}

	return s.Transfer(ctx, string(account), 1)
}

// //STEP 2 - Payer asks bank to blindsign the token. Bank verifies if STEP 1 took place.
// //This request should go only to peer(s) that belong to Org1MSP (otherwise it will fail due to lack of private data)
func (s *SmartContract) BlindSignToken(ctx contractapi.TransactionContextInterface, blinded string) (string, error) {

	debit, err := ctx.GetStub().GetState(DEBIT_PROOF + blinded)
	if err != nil {
		return "", fmt.Errorf("failed to get debit proof: %v", err)
	}
	if len(debit) == 0 {
		return "", errors.New("token not paid. please call DebitMyAccount first")
	}

	k, err := ctx.GetStub().GetPrivateData(BANK_PDC, BANK_ORG)
	if err != nil {
		return "", fmt.Errorf("failed to get private data: %v", err)
	}

	var key rsa.PrivateKey
	if err := json.Unmarshal(k, &key); err != nil {
		return "", fmt.Errorf("failed to unmarshal key: %v", err)
	}

	raw, err := base64.StdEncoding.DecodeString(blinded)
	if err != nil {
		return "", fmt.Errorf("failed to decode blinded message: %v", err)
	}

	// Blind sign the blinded message
	sig, err := rsablind.BlindSign(&key, raw)
	if err != nil {
		return "", fmt.Errorf("failed to blindsign: %v", err)
	}

	return base64.StdEncoding.EncodeToString(sig), nil
}

//This request should go only to peer(s) that belongs to the Payer
func (s *SmartContract) UnblindSignature(ctx contractapi.TransactionContextInterface, sig string, unblinder string) (string, error) {

	pubkey, err := ctx.GetStub().GetState(BANK_ORG)
	if err != nil {
		return "", fmt.Errorf("failed to get pubkey: %v", err)
	}
	raw, err := base64.StdEncoding.DecodeString(string(pubkey))
	if err != nil {
		return "", fmt.Errorf("failed to decode pubkey: %v", err)
	}

	var key rsa.PublicKey
	if err := json.Unmarshal(raw, &key); err != nil {
		return "", fmt.Errorf("failed to unmarshal key: %v", err)
	}

	unblinderBytes, err := base64.StdEncoding.DecodeString(unblinder)
	if err != nil {
		return "", fmt.Errorf("failed to decode unblinder: %v", err)
	}
	sigBytes, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return "", fmt.Errorf("failed to decode sig: %v", err)
	}

	// Unblind the signature
	unblindedSig := rsablind.Unblind(&key, sigBytes, unblinderBytes)

	return base64.StdEncoding.EncodeToString(unblindedSig), nil
}

//This is done by Payee
func (s *SmartContract) CreditMyAccount(ctx contractapi.TransactionContextInterface, unblindedSig string, uuid string) error {

	pubkey, err := ctx.GetStub().GetState(BANK_ORG)
	if err != nil {
		return fmt.Errorf("failed to get pubkey: %v", err)
	}
	raw, err := base64.StdEncoding.DecodeString(string(pubkey))
	if err != nil {
		return fmt.Errorf("failed to decode pubkey: %v", err)
	}

	var key rsa.PublicKey
	if err := json.Unmarshal(raw, &key); err != nil {
		return fmt.Errorf("failed to unmarshal key: %v", err)
	}

	unblindedSigBytes, err := base64.StdEncoding.DecodeString(unblindedSig)
	if err != nil {
		return fmt.Errorf("failed to decode unblindedSig: %v", err)
	}

	if err := rsablind.VerifyBlindSignature(&key, []byte(uuid), unblindedSigBytes); err != nil {
		return fmt.Errorf("failed to verify signature: %v", err)
	}

	//validate for double payments
	credit, err := ctx.GetStub().GetState(CREDIT_PROOF + uuid)
	if err != nil {
		return fmt.Errorf("failed to get credit proof: %v", err)
	}
	if len(credit) > 0 {
		return errors.New("credit operation can be only done once for one blinded token")
	}

	if err = ctx.GetStub().PutState(CREDIT_PROOF+uuid, []byte(uuid)); err != nil {
		return fmt.Errorf("failed to put message: %v", err)
	}

	// Get ID of submitting client identity
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return fmt.Errorf("failed to get client id: %v", err)
	}

	account, err := ctx.GetStub().GetState(BANK_ACCOUNT)
	if err != nil {
		return fmt.Errorf("failed to get bank account: %v", err)
	}
	return transferHelper(ctx, string(account), clientID, 1)
}
