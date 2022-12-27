package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/go-playground/validator/v10"
	"gnark-bid/zk"
	"math/big"
	"syscall/js"
)

var validation = validator.New()

func hash() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "{'error': 'Invalid no of arguments passed'}"
		}
		fmt.Println("args", args)

		if err := validation.Var(args[0].String(), "required,hexadecimal"); err != nil {
			//return "{'error': 'Invalid argument input passed'}"
			return fmt.Sprintf("{'error': '%s'}", err.Error())
		}

		value := common.FromHex(args[0].String())
		return fmt.Sprintf("{'data': '%s'}", zk.HashMIMC(value).String())
	})
}

func proof() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "{'error': 'Invalid no of arguments passed'}"
		}
		if err := validation.Var(args[0].String(), "required,hexadecimal"); err != nil {
			//return "{'error': 'Invalid argument passed'}"
			return fmt.Sprintf("{'error': '%s'}", err.Error())
		}

		inputBytes := common.FromHex(args[0].String())
		privateValue := new(big.Int).SetBytes(inputBytes)
		fmt.Println("privateValue", privateValue.String())
		assignment := zk.Circuit{
			PrivateValue: privateValue.String(),
			Hash:         zk.HashMIMC(inputBytes).String(),
		}

		if err := validation.Var(args[0].String(), "required,hexadecimal"); err != nil {
			//return "{'error': 'Invalid argument input passed'}"
			return fmt.Sprintf("{'error': '%s'}", err.Error())
		}

		vkKey, err := ReadJsonVPKey()
		if err != nil {
			return fmt.Sprintf("{'error': '%s'}", err.Error())
		}

		g16, err := zk.NewGnarkGroth16(vkKey)
		if err != nil {
			return fmt.Sprintf("{'error': '%s'}", err.Error())
		}
		inputProof := [1]*big.Int{zk.HashMIMC(inputBytes)}
		proofGenerated, err := g16.GenerateProof(assignment, inputProof)
		if err != nil {
			return fmt.Sprintf("{'error': '%s'}", err.Error())
		}
		proofJSON, _ := json.Marshal(proofGenerated)

		return string(proofJSON)
	})
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("hash", hash())
	js.Global().Set("generateProof", proof())
	<-make(chan bool)
}
