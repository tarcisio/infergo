package main

import (
	"fmt"

	"github.com/tarcisio/infergo"
	"github.com/tarcisio/infergo/example"
)

const (
	MaxCycles       = 100
	DefaultPriority = 100
)

func main() {
	eng := infergo.NewEngine[*example.Payload](MaxCycles)
	eng.AddRule(example.AgeRule(), DefaultPriority)

	var age int
	fmt.Print("Enter your age: ") //nolint

	if _, err := fmt.Scanf("%d", &age); err != nil {
		panic(err)
	}

	payload := example.Payload{
		Age:          age,
		State:        "SP",
		PanicsOnWhen: false,
		PanicsOnThen: false,
		AgeCheck:     false,
		StateCheck:   false,
	}

	if err := eng.Execute(&payload); err != nil {
		panic(err)
	}

	if payload.AgeCheck {
		fmt.Println("You are old enough") //nolint
	} else {
		fmt.Println("You are not old enough") //nolint
	}
}
