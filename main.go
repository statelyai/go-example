package main

import (
	"encoding/json"
	"fmt"
)

// StateMachine represents the structure of the XState JSON format
type StateMachine struct {
	Initial string                 `json:"initial"`
	States  map[string]StateConfig `json:"states"`
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
	Actions []Action `json:"actions,omitempty"`
}

// TransitionStateMachine takes a state machine definition, current state, and event,
// and returns the next state and actions to execute
func TransitionStateMachine(machine StateMachine, currentState State, event Event) (State, []Action) {
	stateConfig, exists := machine.States[currentState.Value]
	if !exists {
		return currentState, []Action{}
	}

	transition, exists := stateConfig.On[event.Type]
	if !exists {
		return currentState, []Action{}
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
	// Example usage
	machineJSON := `{
		"initial": "green",
		"states": {
			"green": {
				"on": {
					"timer": {
						"target": "yellow",
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

	var machine StateMachine
	err := json.Unmarshal([]byte(machineJSON), &machine)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	currentState := State{Value: "green"}
	event := Event{Type: "timer"}

	nextState, actions := TransitionStateMachine(machine, currentState, event)
	fmt.Printf("Current state: %s\n", currentState.Value)
	fmt.Printf("Event: %s\n", event.Type)
	fmt.Printf("Next state: %s\n", nextState.Value)
	fmt.Println("Actions to execute:")
	for _, action := range actions {
		fmt.Printf("  - Type: %s, Params: %v\n", action.Type, action.Params)
	}
}
