# go-example

This repository demonstrates how to implement and use Stately.ai state machines in Go. It features a simple traffic light state machine example that showcases:

- finite states
- transitions & transition actions
- guarded transitions
- entry, exit, and transition actions

This serves as a starting point for developers looking to integrate Stately.ai state machines into their Go projects. More advanced examples are coming in the future, and contributions are welcome to help expand the functionality and showcase additional state machine/statechart features in Go.

## Running the Example

To run this example, follow these steps:

1. Ensure you have Go installed on your system. You can download it from [https://golang.org/](https://golang.org/).

2. Clone this repository

3. Run the example:

```bash
go run main.go
```

## Overview

The algorithm for transitioning the state machine in `main.go` can be summarized in the following steps:

1. Check if the current state exists in the state machine definition
2. Look for a transition matching the given event type in the current state
3. If a guard is specified for the transition:
   b. Execute the guard function with the provided parameters
   c. If the guard returns false, return the current state without transitioning
4. Determine the next state based on the transition's target
5. Collect actions to be executed in the following order:
   a. Exit actions of the current state
   b. Transition actions
   c. Entry actions of the next state
6. Return a tuple of the next state and the collected actions to execute.

This algorithm is implemented in the `TransitionStateMachine` function.
