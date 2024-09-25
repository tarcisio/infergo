package example

import "github.com/tarcisio/infergo"

type Payload struct {
	Age          int    `json:"age"`
	State        string `json:"state"`
	PanicsOnWhen bool   `json:"panics_on_when"`
	PanicsOnThen bool   `json:"panics_on_then"`

	AgeCheck   bool `json:"age_check"`
	StateCheck bool `json:"state_check"`
}

func AgeRule() *infergo.Rule[*Payload] {
	return infergo.NewRule(
		"Age > 18",
		func(payload *Payload) bool {
			return !payload.AgeCheck && payload.Age > 18
		},
		func(payload *Payload) {
			payload.AgeCheck = true
		},
	)
}

func StateRule() *infergo.Rule[*Payload] {
	return infergo.NewRule(
		"State is not empty",
		func(payload *Payload) bool {
			return !payload.StateCheck && payload.State != ""
		},
		func(payload *Payload) {
			payload.StateCheck = true
		},
	)
}

func RulePanicsOnWhen() *infergo.Rule[*Payload] {
	return infergo.NewRule(
		"RulePanicsOnWhen",
		func(payload *Payload) bool {
			if payload.PanicsOnWhen {
				panic("RulePanicsOnWhen")
			}

			return false
		},
		func(_ *Payload) {
		},
	)
}

func RulePanicOnThen() *infergo.Rule[*Payload] {
	return infergo.NewRule(
		"RulePanicOnThen",
		func(payload *Payload) bool {
			return payload.PanicsOnThen
		},
		func(_ *Payload) {
			panic("RulePanicOnThen")
		},
	)
}

func RuleWithNoResolutions() *infergo.Rule[*Payload] {
	return infergo.NewRule(
		"RuleWithNoResolutions",
		func(_ *Payload) bool {
			return true
		},
		func(_ *Payload) {
		},
	)
}
