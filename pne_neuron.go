package main

import (
    //"time"
    "math"
)

type NeuronType struct {
    Undetermined, Sensory, Deep, Mechanical int
}

var neurontype = NeuronType{-1, 0, 1, 2}

type Neuron struct {
    circuit *Circuit
    Index int
    Type int
    Axon *axon
    MembranePotential float64
    ThresholdPotential float64
    InRefractoryPeriod bool
    Dendrites []Dendrite
}

func (n *Neuron) Genesis(index int, t int, circuit *Circuit) {
    n.Index = index
    n.Type = t
    n.circuit = circuit
    n.Axon = &axon{}
    n.Axon.Genesis(n)
    
    for i := 0; i < n.circuit.MaxConn; i++ {
        n.Dendrites = append(n.Dendrites, Dendrite{nil, n})
    }
    
    n.AssumeRestingPotential()
}

func (n *Neuron) GetVacantDendrite() *Dendrite{
    for i := 0; i < len(n.Dendrites); i++ {
        if n.Dendrites[i].ReceptiveTo == nil {
            return &n.Dendrites[i]
        }
    }
    return nil
}

func (n *Neuron) Hyperpolarization() {
    n.MembranePotential = -0.90
    n.InRefractoryPeriod = true
    //time.Sleep(1 * time.Millisecond)
    n.AssumeRestingPotential()
}

func (n *Neuron) AssumeRestingPotential() {
    n.MembranePotential = -0.70
    n.ThresholdPotential = -0.55
    n.InRefractoryPeriod = false
}

func (n *Neuron) Inhibit() {
    if (*n).InRefractoryPeriod {
        return 
    }
    
    n.MembranePotential -= 0.075
}

func (n *Neuron) Excite(in ... float64) {
    if (*n).InRefractoryPeriod {
        return
    }
    
    if len(in) > 0 {
        n.MembranePotential += in[0]
    } else {
        n.MembranePotential += 0.075
    }
        
    if n.MembranePotential >= n.ThresholdPotential {
        go n.Activate()
    }
}

func (n *Neuron) Activate() {
    defer func() {
        n.Hyperpolarization()
    }()
    
    if n.Type == neurontype.Mechanical {
        inilen := int(math.Ceil((float64(n.circuit.In) - float64(n.circuit.Out)) / float64(2))) + n.circuit.Out + n.circuit.In + n.circuit.Out
        out := n.Index - n.circuit.In - (inilen - n.circuit.In - n.circuit.Out)
        n.circuit.Results = append(n.circuit.Results, Percept{out})
    } else {
        for i := 0; i < len(n.Axon.Terminals); i++ {
            if n.Axon.Terminals[i].To != nil {
                if n.Axon.Terminals[i].SynapseIsExcitatory {
                    go n.Axon.Terminals[i].To.PartOf.Excite()
                } else {
                    go n.Axon.Terminals[i].To.PartOf.Inhibit()
                }
            }
        }
    }
}