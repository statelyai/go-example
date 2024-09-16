package main

import (
	"encoding/json"
	"fmt"
)

// StateMachine represents the structure of the XState JSON format
type StateMachine struct {
	Initial string                 `json:"initial"`
	States  map[string]StateConfig `json:"states"`
	Guards  map[string]GuardImplementation
}

// StateConfig represents the configuration for each state
type StateConfig struct {
	On    map[string]Transition `json:"on"`
	Entry []Action              `json:"entry,omitempty"`
	Exit  []Action              `json:"exit,omitempty"`
}

// State represents the current state
type State struct {
	Value string `json:"value"`
}

// Event represents an event that can trigger a state transition
type Event struct {
	Type string `json:"type"`
}

// Action represents an action to be executed
type Action struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params,omitempty"`
}

// Transition represents a transition between states
type Transition struct {
	Target  string   `json:"target"`
	Guard   *Guard   `json:"guard,omitempty"`
	Actions []Action `json:"actions,omitempty"`
}

// Guard represents a condition for a transition
type Guard struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params,omitempty"`
}

// GuardImplementation is a function type for guard implementations
type GuardImplementation func(params map[string]interface{}) bool

// CreateMachine creates a new StateMachine with the given JSON and guard implementations
func CreateMachine(machineJSON string, guards map[string]GuardImplementation) (*StateMachine, error) {
	var machine StateMachine
	err := json.Unmarshal([]byte(machineJSON), &machine)
	if err != nil {
		return nil, err
	}
	machine.Guards = guards
	return &machine, nil
}

// TransitionStateMachine takes a state machine definition, current state, and event,
// and returns the next state and actions to execute
func (machine *StateMachine) TransitionStateMachine(currentState State, event Event) (State, []Action) {
	stateConfig, exists := machine.States[currentState.Value]
	if !exists {
		return currentState, []Action{}
	}

	transition, exists := stateConfig.On[event.Type]
	if !exists {
		return currentState, []Action{}
	}

	if transition.Guard != nil {
		guardImpl, exists := machine.Guards[transition.Guard.Type]
		if !exists || !guardImpl(transition.Guard.Params) {
			return currentState, []Action{}
		}
	}

	nextState := State{Value: transition.Target}
	nextStateConfig, exists := machine.States[transition.Target]
	if !exists {
		return nextState, []Action{}
	}

	actions := []Action{}
	if stateConfig.Exit != nil {
		actions = append(actions, stateConfig.Exit...)
	}
	if transition.Actions != nil {
		actions = append(actions, transition.Actions...)
	}
	if nextStateConfig.Entry != nil {
		actions = append(actions, nextStateConfig.Entry...)
	}

	return nextState, actions
}

func main() {
	// Define guard implementations
	guards := map[string]GuardImplementation{
		"isValid": func(params map[string]interface{}) bool {
			id, ok := params["someParam"].(string)
			return ok && id == "accepted"
		},
	}

	// Example usage
	machineJSON := `{
		"initial": "green",
		"states": {
			"green": {
				"on": {
					"timer": {
						"target": "yellow",
						"guard": {
							"type": "isValid",
							"params": {
								"someParam": "accepted"
							}
						},
						"actions": [
							{
								"type": "logTransition",
								"params": {
									"message": "Transitioning from green to yellow"
								}
							}
						]
					}
				},
				"entry": [
					{
						"type": "turnOnLight",
						"params": {
							"color": "green"
						}
					}
				],
				"exit": [
					{
						"type": "turnOffLight",
						"params": {
							"color": "green"
						}
					}
				]
			},
			"yellow": {
				"on": {
					"timer": {
						"target": "red"
					}
				},
				"entry": [
					{
						"type": "startTimer",
						"params": {
							"duration": 5000
						}
					},
					{
						"type": "turnOnLight",
						"params": {
							"color": "yellow"
						}
					}
				]
			},
			"red": {
				"on": {
					"timer": {
						"target": "green"
					}
				},
				"entry": [
					{
						"type": "turnOnLight",
						"params": {
							"color": "red"
						}
					}
				]
			}
		}
	}`

	machine, err := CreateMachine(machineJSON, guards)
	if err != nil {
		fmt.Println("Error creating state machine:", err)
		return
	}

	currentState := State{Value: "green"}
	event := Event{Type: "timer"}

	nextState, actions := machine.TransitionStateMachine(currentState, event)
	fmt.Printf("Current state: %s\n", currentState.Value)
	fmt.Printf("Event: %s\n", event.Type)
	fmt.Printf("Next state: %s\n", nextState.Value)
	fmt.Println("Actions to execute:")
	for _, action := range actions {
		fmt.Printf("  - Type: %s, Params: %v\n", action.Type, action.Params)
	}
}
