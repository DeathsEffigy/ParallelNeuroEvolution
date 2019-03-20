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
}

type Percept struct {
    outcome int
}

type RankedResult struct {
    outcome int
    amplitude int
}

func (circuit *Circuit) Neurogenesis(in int, out int) {
    // setup counts
    circuit.In = in
    circuit.Out = out
    n := int(math.Ceil((float64(circuit.In) - float64(circuit.Out)) / float64(2))) + circuit.Out + circuit.In + circuit.Out
    circuit.MaxConn = 5
    
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
        circuit.Results = circuit.Results[:0]
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
    
    count := make(map[int]int)
    for _, res := range circuit.Results {
        count[res.outcome] += 1
    }
    
    ranked := []RankedResult{}
    for i := 0; i < len(count); i++ {
        highest := -1
        highestI := -1
        for n := 0; n < len(count); n++ {
            if count[n] > highest {
                highest = count[n]
                highestI = n
            }
        }
        count[highestI] = -1
        ranked = append(ranked, RankedResult{highestI, highest})
    }
    
    return ranked
}