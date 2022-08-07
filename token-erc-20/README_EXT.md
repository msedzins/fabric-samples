# ERC-20 token scenario - untraceable payments extension

Extension builds on top of already existing "ERC-20 token smart contract" implementation.
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

**NOTE:**
Call should go to the peer(s) owned by Org1 (which represents Bank).
Additionally, no transaction must be generated (otherwise the response will be stored on-chain).

### Save private key in Private Data Collection
```
export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt")

peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n token_erc20 -c '{"function":"SavePrivateKey","Args":[]}' --transient "{\"key\":\"$KEY\"}" --waitForEvent
```

**NOTE:**
We are saving data in PDC not on-chain to keep it private. The request should go only to the peers owned by Org1. To that end, we modify TARGET_TLS_OPTIONS appropriately. 

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

### Generate UUID and blind it \[Payer;Org2MSP]
```

```

