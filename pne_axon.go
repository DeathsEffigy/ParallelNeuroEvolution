package main

type axon struct {
    From *Neuron
    Terminals []*axonTerminal
}

func (a *axon) Genesis(neuron *Neuron) {
    a.From = neuron
}

func (a *axon) GrowTerminals() {
    if a.From.Type == neurontype.Sensory {
        // connect to deep neurons
        for i := a.From.circuit.In; i < len(a.From.circuit.Cluster) - a.From.circuit.Out; i++ {
            a.Terminals = append(a.Terminals, &axonTerminal{a, a.From.circuit.Cluster[i].GetVacantDendrite()})
        }
    } else if a.From.Type == neurontype.Deep {
        // connect to mechanical neurons
        for i := (len(a.From.circuit.Cluster) - a.From.circuit.Out); i < len(a.From.circuit.Cluster); i++ {
            a.Terminals = append(a.Terminals, &axonTerminal{a, a.From.circuit.Cluster[i].GetVacantDendrite()})
        }
        // connect to deep neurons (but not themselves, obv)
        rem := a.From.circuit.MaxConn - a.From.circuit.Out
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
            a.Terminals = append(a.Terminals, &axonTerminal{a, a.From.circuit.Cluster[offset + i].GetVacantDendrite()})
        }
    } else {
        // mechanicals don't get any connections
    }
}