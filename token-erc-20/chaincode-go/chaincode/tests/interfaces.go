package tests

import (
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . TestTransactionContextInterface
type TestTransactionContextInterface interface {
	contractapi.TransactionContextInterface
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . TestChaincodeStubInterface
type TestChaincodeStubInterface interface {
	shim.ChaincodeStubInterface
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . TestClientIdentity
type TestClientIdentity interface {
	cid.ClientIdentity
}
