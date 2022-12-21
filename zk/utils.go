package zk

import (
	"encoding/json"
	"fmt"
	"math/big"
)

const fpSize = 4 * 8

func ParserProof(proofBytes []byte) *Proof {
	proof := &Proof{}
	proof.A[0] = new(big.Int).SetBytes(proofBytes[fpSize*0 : fpSize*1])
	proof.A[1] = new(big.Int).SetBytes(proofBytes[fpSize*1 : fpSize*2])
	proof.B[0][0] = new(big.Int).SetBytes(proofBytes[fpSize*2 : fpSize*3])
	proof.B[0][1] = new(big.Int).SetBytes(proofBytes[fpSize*3 : fpSize*4])
	proof.B[1][0] = new(big.Int).SetBytes(proofBytes[fpSize*4 : fpSize*5])
	proof.B[1][1] = new(big.Int).SetBytes(proofBytes[fpSize*5 : fpSize*6])
	proof.C[0] = new(big.Int).SetBytes(proofBytes[fpSize*6 : fpSize*7])
	proof.C[1] = new(big.Int).SetBytes(proofBytes[fpSize*7 : fpSize*8])
	fmt.Println("a", proof.A)
	fmt.Println("b", proof.B)
	fmt.Println("c", proof.C)

	return proof
}

func GetVPKey(jsonBytes []byte) (*VPKey, error) {
	var vpKey *VPKey
	err := json.Unmarshal(jsonBytes, &vpKey)
	if err != nil {
		return nil, err
	}
	return vpKey, nil
}
