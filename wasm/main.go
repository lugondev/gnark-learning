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
			return "{'error': 'Invalid argument passed'}"
		}

		value := common.FromHex(args[0].String())
		return fmt.Sprintf("{'data': '%s'}", zk.HashMIMC(value).String())
	})
}

func proof() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 2 {
			return "{'error': 'Invalid no of arguments passed'}"
		}
		if err := validation.Var(args[0].String(), "required,hexadecimal"); err != nil {
			//return "{'error': 'Invalid argument passed'}"
			return fmt.Sprintf("{'error': '%s'}", err.Error())
		}

		isReveal := args[1].Bool()
		vkKey, err := ReadJsonVPKey()
		if err != nil {
			return fmt.Sprintf("{'error': '%s'}", err.Error())
		}

		g16, err := zk.NewGnarkGroth16(vkKey)
		if err != nil {
			return fmt.Sprintf("{'error': '%s'}", err.Error())
		}
		inputBytes := common.FromHex(args[0].String())
		inputValue := new(big.Int).SetBytes(inputBytes)
		privateValue := big.NewInt(0)
		if isReveal {
			privateValue = inputValue
		}
		assignment := zk.Circuit{
			PrivateValue: privateValue,
			PublicValue:  inputValue,
			Hash:         zk.HashMIMC(inputBytes),
		}

		if err := validation.Var(args[0].String(), "required,hexadecimal"); err != nil {
			//return "{'error': 'Invalid argument input passed'}"
			return fmt.Sprintf("{'error': '%s'}", err.Error())
		}
		inputProof := [2]*big.Int{privateValue, zk.HashMIMC(inputBytes)}
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
