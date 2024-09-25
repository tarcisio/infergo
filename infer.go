package infergo

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"time"
)

// Engine is a generic rule engine that evaluates and executes a collection of rules on a given payload.
//
// Each rule is defined by a condition (when) and an action (then), and the engine processes these rules
// in the order of their priority. Rules are added to the engine using the `AddRule` method, and the engine
// executes them on the provided payload with the `Execute` method.
//
// The `Engine` tracks the number of execution cycles and safeguards against infinite loops using a maximum
// cycle limit (`maxCycle`).
// If the engine runs more cycles than allowed without making progress, it terminates and returns an error.
//
// Type Parameters:
//   - P: The type of payload that the engine will process. The engine is generic, allowing it to handle various
//     payload types.
type Engine[P any] struct {
	maxCycle uint64
	rules    []*Rule[P]
	logger   *slog.Logger
}

// NewEngine initializes and returns a new instance of the rule engine.
func NewEngine[P any](maxCycle uint64) *Engine[P] {
	return &Engine[P]{
		maxCycle: maxCycle,
		rules:    []*Rule[P]{},
		logger:   slog.Default(),
	}
}

// AddRule adds a rule to the engine and sorts all the rules by priority.
//
// This method sorts the rules after each addition so that rules with
// higher priority are evaluated first during execution.
//
// Example Usage:
//
//	rule := infergo.NewRule("CheckBalance", checkBalanceWhenFunc, adjustBalanceThenFunc)
//	engine.AddRule(rule, 5) // Add the rule with a priority of 5
func (e *Engine[P]) AddRule(rule *Rule[P], priority int) {
	rule.priority = priority
	e.logger.Info("adding rule", slog.String("name", rule.Name))
	e.rules = append(e.rules, rule)

	// sort the runnable rules by priority
	sort.SliceStable(e.rules, func(i, j int) bool {
		return e.rules[i].priority > e.rules[j].priority
	})
}

// Execute runs the rule engine on the provided payload. It evaluates each rule's condition (when function),
// and if the condition returns true, the rule's action (then function) is executed.
//
// The engine evaluates rules in priority order, and it will continue executing rules until no more rules can be run,
// or until the `maxCycle` limit is reached, in which case an error is returned.
//
// Example Usage:
//
//	err := engine.Execute(myPayload)
//	if err != nil {
//	    log.Println("Execution failed:", err)
//	}
func (e *Engine[P]) Execute(payload P) error {
	e.logger.Info("starting rule execution", slog.Int("rule_count", len(e.rules)))

	startTime := time.Now()

	var cycle uint64

	for {
		e.logger.Debug("inside cycle", slog.Uint64("cycle", cycle))

		runnable := e.runnable(payload)
		e.logger.Info("selected rules", slog.Int("count", len(runnable)))

		if len(runnable) > 0 {
			e.logger.Debug("len(runnable rules) > 0")

			cycle++

			if cycle > e.maxCycle {
				e.logger.Error("max cycle reached", slog.Uint64("max_cycle", e.maxCycle))

				return fmt.Errorf("max cycle of %d reached. error: %w", e.maxCycle, ErrMaxCycleReached)
			}

			// get the first rule with the highest priority and execute it
			runner := runnable[0]
			e.logger.Info("executing rule", slog.String("name", runner.Name))
			runner.then(payload)

		} else {
			e.logger.Info("No more rule to run")

			break
		}
	}

	e.logger.Info(
		"finished rules execution",
		slog.Uint64("cycles", cycle),
		slog.Int64("duration_ms", time.Since(startTime).Milliseconds()))

	return nil
}

func (e *Engine[P]) runnable(payload P) []*Rule[P] {
	runnable := make([]*Rule[P], 0)

	for _, rule := range e.rules {
		can := rule.when(payload)
		if can {
			runnable = append(runnable, rule)
			e.logger.Debug("rule is runnable", slog.String("name", rule.Name))
		} else {
			e.logger.Debug("rule is not runnable", slog.String("name", rule.Name))
		}
	}

	return runnable
}

// WhenFunc defines the condition function for a rule. It receives the current payload and
// returns a boolean value indicating whether the rule can be executed.
//
// If the condition evaluates to true, the corresponding action (then function) will be executed.
type WhenFunc[P any] func(payload P) bool

// ThenFunc defines the action function for a rule. It receives the current payload and performs
// the desired operation when the rule's condition (when function) evaluates to true.
type ThenFunc[P any] func(payload P)

// Rule represents a single rule in the engine. A rule consists of a condition (when function),
// an action (then function), and a priority that determines the order in which it is executed.
type Rule[P any] struct {
	Name     string      // The name of the rule, which can be used for logging or debugging purposes.
	when     WhenFunc[P] // The condition function that determines whether the rule should be executed.
	then     ThenFunc[P] // The action function that is executed when the condition is true.
	priority int         // The priority of the rule. Higher priority rules are executed first.
}

// NewRule creates and returns a new rule. A rule consists of a name, a condition (when function),
// and an action (then function).
// The rule's priority defaults to 0, but it can be set when adding the rule to the engine with `AddRule`.
func NewRule[P any](name string, when WhenFunc[P], then ThenFunc[P]) *Rule[P] {
	return &Rule[P]{
		Name:     name,
		when:     when,
		then:     then,
		priority: 0,
	}
}

var ErrMaxCycleReached = errors.New("max cycle reached")
