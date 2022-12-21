package zk

import (
	"github.com/consensys/gnark/frontend"
)

var MinPrice frontend.Variable = 20

// Circuit defines a pre-image knowledge proof
type Circuit struct {
	PrivateValue frontend.Variable
	PublicValue  frontend.Variable `gnark:",public"`
	Hash         frontend.Variable `gnark:",public"`
}

// Define declares the circuit's constraints
func (circuit *Circuit) Define(api frontend.API) error {
	// value must be greater than min price
	api.AssertIsDifferent(circuit.PrivateValue, MinPrice)
	api.AssertIsLessOrEqual(circuit.PrivateValue, MinPrice)

	if api.IsZero(circuit.PublicValue) == 0 {
		api.AssertIsEqual(circuit.PrivateValue, circuit.PublicValue)
	}

	hashMIMC := HashPreImage(api, circuit.PrivateValue)
	api.AssertIsEqual(circuit.Hash, hashMIMC)

	return nil
}
