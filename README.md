## What is PNE?
Parallel Neuro-Evolution (PNE) is an experimental artificial intelligence that works very much unlike any other models currently out there. Basically, PNE does not use matrices at all. Instead, it emulates (with some degree of accuracy) a parallel activation paradigm where neurons either fire or don't fire, but they do not do so in sync. Consequently, the model will activate (and re-activate) a wide range of mechanical neurons (output neurons), that are ranked by frequency of activation and reactively inhibited. The neural model then evolves using forward propagation.

## What's next for PNE?
Right now, PNE works quite well. There is, however, a caveat: Since the model uses Goroutines to emulate a parallel activation paradigm, the `pne_circuit.go`-model can get quite intense (and slow). Learning rate is perfectly fine and rather fast, but processing takes a while. This could be combatted by writing a `pne_lsbn.go` container that allows construction of several `pne_circuit.go` structs that perform quicker, much more focused tasks. As the PNE model should best be used in situations similar to contextual learning, this would also take the project one step further towards biological accuracy. For example, it seems ridiculous to assume that the brain dedicates a single circuit to the recognition of letters. Rather, the process likely involves multiple domain-general circuits (e.g., shape detection and recursion). An approximation of this has now been implemented as `pne_lsbn.go`. However, this is very rudimentary and currently not advised. A better solution to this issue would be an implementation of circuits whose mechanical neurons can fire at sensory neurons of adjacent circuits. This will be implemented some time in the future, as it should allow for much more flexible models and *significantly* faster convergence.

## Should I use PNE?
If you're asking this, the answer is probably no. This isn't for software development (let alone production releases). This is mostly for people interested in computational neurosciences.

## How do I use PNE?
There are currently two ways to separate models supported by PNE. One utilises `pne_circuit.go` only to form a single circuit struct, while the other uses `pne_lsbn.go` to create multiple circuit structs within a large scale brain network. Currently, it is advisable to make use of the circuit model, as it is faster and far more accurate.

### Method 1: PNE.Circuit
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
        circuit.CorrectFor(result, stimulus.Correct, stimulus.Sensation)
      } else {
        error += 1
      }
    }
    success_rate = float64(correct) / float64(total)
    fmt.Printf("success_rate=%f after %d trials.\n", success_rate, count)
  }
}
```

### Method 2: PNE.LSBN
In order to use the LSBN method, first initialise your large scale brain network as follows:
```go
lsbn := LargeScaleBrainNetwork{}
```
Next, setup a circuit that will handle chunked stimuli data. This can be achieved by calling `lsbn.GrowCircuit(circuit_name string, chunk_length int, chunk_beta float64, stimulus_length int, stimuli []float64, do_train bool, train_until_alpha_level float64)`. For example, if we wanted to slice a 16x16 (i.e. `stimulus_length=256`) greyscale matrix of handwritten numbers into four chunks (i.e. `chunk_length=64`) that are handled by a "shapes"-circuit (i.e. `circuit_name="shapes"`) that distinguishes shapes at `β=0.6` and is trained for `α=0.9`, this could be achieved like so:
```go
lsbn.GrowCircuit("shapes", 64, 0.6, 256, stimuli, true, 0.9)
```
However, `lsbn.GrowCircuit()` also returns three values, that is `(types int, chunks []LSBNChunk, success bool)` where `types` denotes the number of different shapes this struct has learned to differentiate at our α-level, `chunks` represents the re-ordered and differentially analysed struct of chunks that were used in the process, and `success` can be used for catching errors.
Once we've grown (and trained) this circuit, we can grow recursion for our LSBN like so:
```go
lsbn.Grow("numbers", "shapes", 10)
```
As may become apparent, `lsbn.Grow()` takes a `recursion_name string` parameter, a `use_circuit string` parameter and an `outputs int` parameter (where outputs will denote the final outputs of the NN). Unlike `lsbn.GrowCircuit()`, `lsbn.Grow()` does not return any values, as it is untrained and serves only to finalise our initialisation of the `LargeScaleBrainNetwork{}` struct. In order to train the LSBN, we must first expose it to some stimulus, like so:
```go
stim_chunked, res := lsbn.Expose("numbers", stimulus.Sensation)
```
Again, `lsbn.Expose()` takes `recursion_circuit string` and `stimulus []float64` as inputs and outputs `stim_chunked []float64` and `results []RankedResult`. These can now be used to train the model:
```go
if len(res) > 0 {
  lsbn.Correct("numbers", res, stimulus.Correct, stim_chunked)
}
```
And, like so, we've just completed our first trial and trained the model. A full cycle of training a large scale brain network, then, might look as follows:
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
  stimuliOne := LoadMyFancyStimuliAsFloat() // as one big []float64
  stimuli := LoadMyFancyStimuli() // as []Stimulus
  
  lsbn := LargeScaleBrainNetwork{}
  types, _, _ := lsbn.GrowCircuit("shapes", 64, 0.6, 256, stimuliOne, true, 0.9)
  lsbn.Grow("numbers", "shapes", 10)
  
  alpha_level := 0.9
  success_rate := 0
  count := 0
  success := 0
  
  for success_rate < alpha_level {
    for _, stimulus := range stimuli {
      count += 1
      stim, res := lsbn.Expose("numbers", stimulus.Sensation)
      if len(res) > 0 {
        if res[0].outcome == stimulus.Correct {
          success += 1
        }
        lsbn.Correct("numbers", res, stimulus.Correct, stim)
      }
    }
    success_rate = float64(success) / float64(count)
    fmt.Printf("success_rate=%f after %d trials.\n", success_rate, count)
  }
}
```