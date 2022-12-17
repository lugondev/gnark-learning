package main

import (
	"log"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/examples/cubic"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

func main() {
	err := generateGroth16()
	if err != nil {
		log.Fatal("groth16 error:", err)
	}
}

func generateGroth16() error {
	var circuit cubic.Circuit

	r1csCompiled, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return err
	}

	pk, vk, err := groth16.Setup(r1csCompiled)
	if err != nil {
		return err
	}
	{
		f, err := os.Create("cubic.g16.vk")
		if err != nil {
			return err
		}
		_, err = vk.WriteRawTo(f)
		if err != nil {
			return err
		}
	}

	{
		f, err := os.Create("cubic.g16.pk")
		if err != nil {
			return err
		}
		_, err = pk.WriteRawTo(f)
		if err != nil {
			return err
		}
	}

	{
		f, err := os.Create("contract_g16.sol")
		if err != nil {
			return err
		}
		err = vk.ExportSolidity(f)
		if err != nil {
			return err
		}
	}
	return nil
}
