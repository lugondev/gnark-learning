package zk

import (
	"fmt"
	"github.com/consensys/gnark/frontend"
	"gnark-bid/circuits"
)

// Circuit defines a pre-image knowledge proof
type Circuit struct {
	PrivateValue frontend.Variable
	Hash         frontend.Variable `gnark:",public"`
}

// Define declares the circuit's constraints
func (circuit *Circuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuits.IsZero(api, circuit.PrivateValue), 0)
	fmt.Println("PrivateValue:", circuit.PrivateValue)

	hashMIMC := HashPreImage(api, circuit.PrivateValue)
	api.AssertIsEqual(circuit.Hash, hashMIMC)

	return nil
}
