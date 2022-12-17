SHELL=/bin/bash

.PHONY: generator-sol clean-abi-generated solc abigen

all: generator-sol abigen go-test

generator-sol:
	cd solidity && go run contract/main.go

clean-abi-generated:
	cd solidity && rm -fr ./abi/*

solc: clean-abi-generated
	cd solidity && solc --bin --abi -o ./abi contract_g16.sol

abigen: solc
	cd solidity && abigen --bin ./abi/Verifier.bin --abi abi/Verifier.abi --pkg solidity --out solidity_groth16.go --type Verifier

go-test:
	cd solidity && go test
