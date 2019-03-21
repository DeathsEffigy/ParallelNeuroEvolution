package main

import (
    "math"
    "math/rand"
    "time"
    //"fmt"
)

type axon struct {
    From *Neuron
    Terminals []*axonTerminal
}

func (a *axon) Genesis(neuron *Neuron) {
    a.From = neuron
}

func (a *axon) GrowTerminals() {
    rand.Seed(time.Now().UnixNano())
    
    if a.From.Type == neurontype.Sensory {
        // connect to deep neurons
        deep := len(a.From.circuit.Cluster) - a.From.circuit.In - a.From.circuit.Out
        per := int(math.Ceil(float64(a.From.circuit.In) / float64(deep)))
        this := int(math.Floor(float64(a.From.Index) / float64(per))) + a.From.circuit.In
        a.Terminals = append(a.Terminals, &axonTerminal{a, a.From.circuit.Cluster[this].GetVacantDendrite(), true})
        
        
        /*for i := a.From.circuit.In; i < len(a.From.circuit.Cluster) - a.From.circuit.Out; i++ {
            a.Terminals = append(a.Terminals, &axonTerminal{a, a.From.circuit.Cluster[i].GetVacantDendrite(), MakeBool(rand.Float64())})
        }*/
    } else if a.From.Type == neurontype.Deep {
        // connect to mechanical neurons
        for i := (len(a.From.circuit.Cluster) - a.From.circuit.Out); i < len(a.From.circuit.Cluster); i++ {
            a.Terminals = append(a.Terminals, &axonTerminal{a, a.From.circuit.Cluster[i].GetVacantDendrite(), MakeBool(rand.Float64())})
        }
        /*
        // connect to deep neurons (but not themselves, obv)
        rem := a.From.circuit.MaxConn - a.From.circuit.Out
        fmt.Println(rem)
        if rem < 1 {
            return
        }
        offset := a.From.Index - (rem / 2)
        if offset < a.From.circuit.In {
            offset = a.From.circuit.In
        }
        for i := 0; i < rem; i++ {
            // don't connect to self
            if (offset + i) == a.From.Index {
                continue
            }
            // don't connect to mechanicals
            if (offset + i) >= (len(a.From.circuit.Cluster) - a.From.circuit.Out) {
                continue
            }
            a.Terminals = append(a.Terminals, &axonTerminal{a, a.From.circuit.Cluster[offset + i].GetVacantDendrite(), MakeBool(rand.Float64())})
        }
        */
    } else {
        // mechanicals don't get any connections
    }
}

func (a *axon) HasTerminalTo(ptr *Neuron) *axonTerminal {
    if ptr == nil {
        return nil
    }
    for _, at := range a.Terminals {
        if at == nil {
            continue
        }
        if at.To == nil {
            continue
        }
        if at.To.PartOf == nil {
            continue
        }
        
        if &(*at.To.PartOf) == &(*ptr) {
            return at
        }
    }
    return nil
}

func MakeBool(p float64) bool {
    /*if p < 0.5 {
        return false
    }*/
    return true
}

func (a *axon) GrowSingleTerminal(to int, exc bool) {
    a.Terminals = append(a.Terminals, &axonTerminal{a, a.From.circuit.Cluster[to].GetVacantDendrite(), exc})
}