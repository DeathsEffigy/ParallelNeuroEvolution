## What is PNE?
Parallel Neuro-Evolution (PNE) is an experimental artificial intelligence that works very much unlike any other models currently out there. Basically, PNE does not use matrices at all. Instead, it emulates (with some degree of accuracy) a parallel activation paradigm where neurons either fire or don't fire, but they do not do so in sync. Consequently, the model will activate (and re-activate) a wide range of mechanical neurons (output neurons), that are ranked by frequency of activation and reactively inhibited. The neural model then evolves using forward propagation.

## What's next for PNE?
Right now, PNE works quite well. There is, however, a caveat: Since the model uses Goroutines to emulate a parallel activation paradigm, the model can get quite intense (and slow). Learning rate is perfectly fine and rather fast, but processing takes a while. This could be combatted by writing a `pne_brain.go` container that allows construction of several `pne_circuit.go` structs that perform quicker, much more focused tasks. As the PNE model should best be used in situations similar to contextual learning, this would also take the project one step further towards biological accuracy. For example, it seems ridiculous to assume that the brain dedicates a single circuit to the recognition of letters. Rather, the process likely involves multiple domain-general circuits (e.g., shape detection and recursion). An approximation of this will be implemented some time in the future that should allow for much more flexible models and *significantly* faster convergence.

## Should I use PNE?
If you're asking this, the answer is probably no. This isn't for software development (let alone production releases). This is mostly for people interested in computational neurosciences.

## How do I use PNE?
First off, initialise a circuit, like so:
```go
circuit := Circuit{}
circuit.Neurogenesis(sensory_neurons int, mechanical_neurons int)
```
`Neurogenesis()` takes two int arguments. The first one specifies your sensory neurons (inputs), whereas the second one dictates your mechanical neurons (outputs).
Next, you can expose the circuit to a stimulus as follows:
```go
response := circuit.ExposeTo([]float64)
```
You can then use the response struct to initialise forward propagation.
```go
circuit.CorrectFor(response []RankedResponse, correct_output int, stimulus []float64)
```
And that's it, really. You've just completed your first training cycle. A full one might look like this:
```go
package main

import (
  "fmt"
)

type Stimulus struct {
  Sensation []float64
  Correct int
}

func main() {
  circuit := Circuit{}
  circuit.Neurogenesis(256, 10)
  stimuli := LoadMyFancyStimuli() // as []Stimulus
  alpha_level := 0.9
  success_rate := float64(0)
  total := 0
  correct := 0
  error := 0
  
  for success_rate < alpha_level {
    for _, stimulus := range stimuli {
      total += 1
      result := circuit.ExposeTo(stimulus.Sensation)
      if len(result) > 0 {
        if result[0].outcome == stimulus.Correct {
          correct += 1
        } else {
          error += 1
        }
      } else {
        error += 1
      }
      circuit.CorrectFor(result, stimulus.Correct, stimulus.Sensation)
    }
    success_rate = float64(correct) / float64(total)
    fmt.Printf("success_rate=%f after %d trials.\n", success_rate, count)
  }
}
```