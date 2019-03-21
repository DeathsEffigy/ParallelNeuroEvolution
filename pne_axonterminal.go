package main

type axonTerminal struct {
    From *axon
    To *Dendrite
    SynapseIsExcitatory bool
}