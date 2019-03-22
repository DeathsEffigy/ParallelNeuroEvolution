package main

import (
    "time"
    "math"
    "runtime"
)

type Circuit struct {
    In int
    Out int
    Cluster []*Neuron
    Results []Percept
    MaxConn int
    Inhibitors int
}

type Percept struct {
    outcome int
}

type RankedResult struct {
    outcome int
    amplitude int
    confidence float64
}

func (circuit *Circuit) Neurogenesis(in int, out int) {
    // setup counts
    circuit.In = in
    circuit.Out = out
    n := int(math.Ceil((float64(circuit.In) - float64(circuit.Out)) / float64(2))) + circuit.Out
    if n > circuit.In {
        n = circuit.In
    }
    n = n + circuit.In + circuit.Out
    circuit.Inhibitors = 0
    circuit.MaxConn = 15
    
    for i := 0; i < n; i++ {
        // determine current type
        t := neurontype.Undetermined
        if i < circuit.In {
            t = neurontype.Sensory
        } else if i >= circuit.In && i < (n - circuit.Out) {
            t = neurontype.Deep
        } else if i >= (n - circuit.Out) {
            t = neurontype.Mechanical
        }
        // grow neuron
        circuit.GrowNeuron(t)
    }
    
    // grow axon terminals
    for i := 0; i < len(circuit.Cluster); i++ {
        circuit.Cluster[i].Axon.GrowTerminals()
    }
}

func (circuit *Circuit) GrowNeuron(t int) {
    circuit.Cluster = append(circuit.Cluster, &Neuron{})
    circuit.Cluster[len(circuit.Cluster)-1].Genesis(len(circuit.Cluster)-1, t, circuit)
}

func (circuit *Circuit) ExposeTo(stimulus []float64) []RankedResult {
    defer func() {
        circuit.Results = nil
    }()
    
    // make sure stimulus isn't bigger than inputs
    if len(stimulus) > circuit.In {
        return nil
    }
    
    for index, stim := range stimulus {
        go circuit.Cluster[index].Excite(stim)
    }
    
    for {
        time.Sleep(100 * time.Millisecond)
        if runtime.NumGoroutine() == 1 {
            break
        }
    }
    
    mechano := 0
    count := make(map[int]int)
    for _, res := range circuit.Results {
        count[res.outcome] += 1
        mechano += 1
    }
    
    t := len(count)
    ranked := []RankedResult{}
    for i := 0; i < t; i++ {
        highest := -1
        highestI := -1
        for n := 0; n < t; n++ {
            if count[n] > highest {
                highest = count[n]
                highestI = n
            }
        }
        count[highestI] = -1
        ranked = append(ranked, RankedResult{highestI, highest, float64(highest) / float64(mechano)})
    }
    
    return ranked
}

func (c *Circuit) CorrectFor(r []RankedResult, v int, stimulus []float64) {
    if r[0].outcome != v {
        deep := int(math.Ceil((float64(c.In) - float64(c.Out)) / float64(2))) + c.Out
        per := int(math.Ceil(float64(c.In) / float64(deep)))
        
        designatedOut := deep + (*c).In + (*c).Out - (*c).Out + r[0].outcome
        realOut := deep + (*c).In + (*c).Out - (*c).Out + v
        
        deepPotentials := make([]float64, deep + c.In)
        for in, stim := range stimulus {
            this := int(math.Floor(float64(in) / float64(per))) + c.In
            if c.Cluster[in].MembranePotential + stim > c.Cluster[in].ThresholdPotential {
                deepPotentials[this] += 0.075
            }
        }
        
        for index, pot := range deepPotentials {
            if c.Cluster[index].MembranePotential + pot > c.Cluster[index].ThresholdPotential {
                at := (*c).Cluster[index].Axon.HasTerminalTo((*c).Cluster[designatedOut])
                if at == nil {
                    continue
                }
                
                if (*at).SynapseIsExcitatory {
                    (*c).GrowNeuron(neurontype.Deep)
                    (*c).Cluster[len((*c).Cluster)-1].Axon.GrowSingleTerminal(realOut, true)
                    if (*c).Inhibitors < deep * ((*c).Out-1) {
                        (*c).Cluster[index].Axon.GrowSingleTerminal(designatedOut, false)
                        (*c).Inhibitors += 1
                    }
                    (*c).Cluster[index].Axon.GrowSingleTerminal(len((*c).Cluster)-1, true)
                }
            }
        }
        
        
    }
}