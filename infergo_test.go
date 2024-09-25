package infergo_test

import (
	"testing"

	"github.com/tarcisio/infergo"
	"github.com/tarcisio/infergo/example"
)

func TestOneGoodRule(t *testing.T) {
	t.Parallel()

	eng := infergo.NewEngine[*example.Payload](100)
	eng.AddRule(example.AgeRule(), 100)

	payload := example.Payload{
		Age:          25,
		State:        "CA",
		PanicsOnWhen: false,
		PanicsOnThen: false,
		AgeCheck:     false,
		StateCheck:   false,
	}
	if err := eng.Execute(&payload); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !payload.AgeCheck {
		t.Fatal("expected AgeCheck to be true")
	}
}

func TestTwoGoodRule(t *testing.T) {
	t.Parallel()

	eng := infergo.NewEngine[*example.Payload](100)
	eng.AddRule(example.AgeRule(), 100)
	eng.AddRule(example.StateRule(), 101)

	payload := example.Payload{
		Age:          25,
		State:        "CA",
		PanicsOnWhen: false,
		PanicsOnThen: false,
		AgeCheck:     false,
		StateCheck:   false,
	}

	if err := eng.Execute(&payload); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !payload.AgeCheck {
		t.Fatal("expected AgeCheck to be true")
	}

	if !payload.StateCheck {
		t.Fatal("expected StateCheck to be true")
	}
}

func TestRuleWithNoResolutions(t *testing.T) {
	t.Parallel()

	eng := infergo.NewEngine[*example.Payload](100)
	eng.AddRule(example.AgeRule(), 100)
	eng.AddRule(example.RuleWithNoResolutions(), 101)

	payload := example.Payload{
		Age:          25,
		State:        "CA",
		PanicsOnWhen: false,
		PanicsOnThen: false,
		AgeCheck:     false,
		StateCheck:   false,
	}

	if err := eng.Execute(&payload); err == nil {
		t.Fatal("expected error, got nil")
	}
}
