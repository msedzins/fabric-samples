# ERC-20 token scenario - untraceable payments extension

Extension builds on top of already existing "ERC-20 token smart contract" implementation ([golang implementation](https://github.com/msedzins/fabric-samples/tree/main/token-erc-20/chaincode-go)).
It demonstrates how to implement payment system as described in ["Blind signatures for untraceable payments"](http://www.hit.bme.hu/~buttyan/courses/BMEVIHIM219/2009/Chaum.BlindSigForPayment.1982.PDF)

**NOTE:**
The implementation has not undergone a security review or audit and should not be used in production code.

## Prerequisites 

Please execute all steps described in [ERC-20 README](README.md) file first. It includes setting up the network, deployment of the chaincode (the updated one) + initialisation, minting and transfering tokens to Org2.

## Configuration \[as Org1MSP]

### Generate RSA key pair (public and private) 
```
KEY=$(peer chaincode query -C mychannel -n token_erc20 -c '{"function":"GenerateKeyPair","Args":[]}')
```

Call should go to the peer(s) owned by Org1 (which represents Bank).
Additionally, no transaction must be generated (otherwise the response will be stored on-chain).

### Save private key in Private Data Collection
```
export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt")

peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n token_erc20 -c '{"function":"SavePrivateKey","Args":[]}' --transient "{\"key\":\"$KEY\"}" --waitForEvent
```

We are saving data in PDC not on-chain to keep it private. The request should go only to the peers owned by Org1. To that end, we modify TARGET_TLS_OPTIONS appropriately. 
Please note that we are using implicit PDC which requires only one endorsement.

### Save public key on-chain
```
#Restore original configuration (TARGET_TLS_OPTIONS)
source env_org1.sh 

#jq doesn't preserve  BigInt values. It can't be used.
PUBKEY=$(echo $KEY | base64 -d | sed 's/,"D".*/}\n/' | base64) 

peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n token_erc20 -c '{"function":"SavePublicKey","Args":["'"$PUBKEY"'"]}'  --waitForEvent
```

### Set bank account 
```
BANK_ACCOUNT=$(peer chaincode query -C mychannel -n token_erc20 -c '{"function":"ClientAccountID","Args":[]}')

peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n token_erc20 -c '{"function":"SetBankAccount","Args":["'"$BANK_ACCOUNT"'"]}'  --waitForEvent
```

Bank account is used as a placeholder for blinded tokens (described in the "Payment" section)
In the example, minter account is used for this purpose. Propobably it's better to have a dedicated account.

## Payment 

### Generate token and blind it \[Payer;Org2MSP]
```
#UUID represents our token
uuid=$(uuidgen)

RESPONSE=$(peer chaincode query -C mychannel -n token_erc20 -c '{"function":"BlindToken","Args":["'"$uuid"'"]}') 
```

To not reveal the data the request must go to the peer that is trusted to the payer (belongs to Org2MSP in our case) + no blockchain transaction can be generated

### Debit account \[Payer;Org2MSP]
```
BLINDED=$(echo $RESPONSE | jq -r '.Blinded')

peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n token_erc20 -c '{"function":"DebitMyAccount","Args":["'"$BLINDED"'"]}' --waitForEvent
```

Originally, there is one step -> bank blind signs the token + debits the account of the client. It won't work for HLF because we can't prevent situation in which client calls the function, gets the signature, but doesn't generate the transaction. Hence, we split the  process into two steps:
1. debit the account (it must generate the transaction)
2. ask for a signature

After the call, Payer account should be debited by 1 token. Can be  verified by running:
```
peer chaincode query -C mychannel -n token_erc20 -c '{"function":"ClientAccountBalance","Args":[]}'
```

Minter account is credited by 1 token respectively.

### Blind sign token \[Payer;Org2MSP]
```
BLINDED=$(echo $RESPONSE | jq -r '.Blinded')

#call is made as Org2MSP
#but goes to the peer(s) owned by Org1MSP  (otherwise it will fail due to lack of private data) 
export CORE_PEER_ADDRESS=localhost:7051 
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt

SIG=$(peer chaincode query -C mychannel -n token_erc20 -c '{"function":"BlindSignToken","Args":["'"$BLINDED"'"]}')

#restore original configuration
source env_org2.sh 
```

The call to BlindSignToken will only succeed if call to DebitMyAccount was made first. It must go to the peer which has an access to private key, which is stored in Org1MSP private data collection.

### Unblind the signature \[Payer;Org2MSP]
```
UNBLINDER=$(echo $RESPONSE | jq -r '.Unblinder')

UNBLIND_SIG=$(peer chaincode query -C mychannel -n token_erc20 -c '{"function":"UnblindSignature","Args":["'"$SIG"'","'"$UNBLINDER"'"]}')
```

To not reveal the data - the call should go to the peer(s) owned by the Payer (Org2MSP). 

### Use the token \[Payee;Org3MSP]
```
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n token_erc20 -c '{"function":"CreditMyAccount","Args":["'"$UNBLIND_SIG"'","'"$uuid"'"]}') --waitForEvent
```

**NOTE:**
1. In our example, for simplification, the call is made by Org2MSP.
2. The call can be made only one to avoid double spending.

After the call, Payee account should be credited by 1 token. Can be  verified by running:
```
peer chaincode query -C mychannel -n token_erc20 -c '{"function":"ClientAccountBalance","Args":[]}'
```
