package main

import (
    //"fmt"
    //"gonum.org/v1/gonum/mat"
    //"gonum.org/v1/gonum/stat"
    "math"
)

/**
 * (1) allow for initialisation of circuit types that can be trained independently
 * (2) allow for parallel activation of circuit layers with specific tasks
 * (3) add gamma layer (recursion) as a defer() that'll put the pieces together
 *  --> how do we convert results from circuit layers to inputs for recursive layer?
 *  ====> SOLUTION: Increase input layer count for outputs? e.g., {0 0 0 1}
 *        that way we could use a sort of matrix multiplication for the final products?
 *          => might be slow due to go routines required
 * [later: (4) add LSTM capacities/attention to the mix; might be relatively easy to import in this setting]
 */

type LSBNCircuit struct {
    Identifier string
    Circuit *Circuit
    Data []LSBNChunk
    StimLength int
    ChunkLength int
    ChunkBeta float64
    RawStimuli *[]float64
    Types int
    ConnectsTo string
}

type LargeScaleBrainNetwork struct {
    IsSetup bool
    Circuits map[string]*LSBNCircuit
}

func (lsbn *LargeScaleBrainNetwork) ForceGenesis () {
    if lsbn.IsSetup != true {
        lsbn.Circuits = make(map[string]*LSBNCircuit)
        lsbn.IsSetup = true
    }
}

type LSBNChunk struct {
    Mean float64
    Sensations []float64
    PeakMap map[int]bool
    Type int
}

func (c *LSBNChunk) Setup(s []float64) {
    c.Sensations = s
    c.PeakMap = make(map[int]bool, len(s))
    c.Type = -1
    
    total := float64(0)
    num := 0
    for _, p := range s {
        total += p
        num += 1
    }
    c.Mean = total / float64(num)
    for index, p := range s {
        meancentered := p - c.Mean
        if meancentered > 0 {
            c.PeakMap[index] = true
        } else {
            c.PeakMap[index] = false
        }
    }
}

func (c *LSBNChunk) GetOverlapWith(o LSBNChunk) float64 {
    cxy := int(math.Sqrt(float64(len(c.Sensations))))
    
    features := 0
    overlap := 0
    for index, isPeak := range c.PeakMap {
        if isPeak {
            features += 1
            
            for i := 0; i < 9; i++ {
                // adjust index
                inc := index
                if i < 3 {
                    inc = inc - cxy - 1 + i
                } else if i < 6 {
                    inc = inc - 1 + i
                } else {
                    inc = inc + cxy - 1 + i
                }
                
                // no overflows
                if inc < 0 || inc > len(c.Sensations) {
                    continue
                }
                
                if o.PeakMap[inc] == isPeak {
                    overlap += 1
                    break
                }
            }
        }
    }
    return float64(overlap) / float64(features)
}

func (lsbn *LargeScaleBrainNetwork) GrowCircuit (identifier string, chunk_length int, chunk_beta float64, stim_length int, stimuli []float64, train bool, train_alpha float64) (int, []LSBNChunk, bool) {
    // make sure lsbn is set up properly
    lsbn.ForceGenesis()
    
    // make sure the In:Length ratio works
    if stim_length % chunk_length != 0 {
        return -1, nil, false
    }
    
    // re-structure stimuli input into []chunks
    types, chunks := lsbn.MakeChunks(chunk_length, chunk_beta, stim_length, stimuli)
    
    // setup circuit
    (*lsbn).Circuits[identifier] = &LSBNCircuit{identifier, &Circuit{}, chunks, stim_length, chunk_length, chunk_beta, &stimuli, types, ""}
    (*lsbn).Circuits[identifier].Circuit.Neurogenesis(chunk_length, types)
    
    if train {
        (*lsbn).TrainCircuit(identifier, train_alpha)
    }
    
    return types, chunks, true
}

func (lsbn *LargeScaleBrainNetwork) MakeChunks (chunk_length int, chunk_beta float64, stim_length int, stimuli []float64) (int, []LSBNChunk) {
    chunks_per_stim := stim_length / chunk_length
    chunks_per_xy := int(math.Sqrt(float64(chunks_per_stim)))
    chunk_xy := int(math.Sqrt(float64(chunk_length)))
    chunkx := make([][][]float64, (len(stimuli) / stim_length) * chunks_per_stim)
    for i := 0; i < len(stimuli); i += chunk_xy {
        c := i / chunk_xy
        col_x := int(c % chunks_per_xy)
        row_x := int(math.Floor(float64(c) / float64(chunks_per_xy)))
        cc := int(math.Floor(float64(row_x) / float64(chunk_xy)))
        if cc > 0 {
            cc += col_x + cc
        } else {
            cc += col_x
        }
        chunkx[cc] = append(chunkx[cc], stimuli[i:i+chunk_xy])
    }
    
    // setup chunks
    var chunks []LSBNChunk
    for _, mchunks := range chunkx {
        // get whole chunk
        var realchunk []float64
        for _, chunk := range mchunks {
            for _, p := range chunk {
                realchunk = append(realchunk, p)
            }
        }
        
        lsc := LSBNChunk{}
        lsc.Setup(realchunk)
        chunks = append(chunks, lsc)
    }
    
    // calculate similarities and find types
    types := 0
    for i := 0; i < len(chunks); i++ {
        for n := 0; n < len(chunks); n++ {
            if i == n {
                continue
            }
            ol := chunks[i].GetOverlapWith(chunks[n])
            if ol > chunk_beta {
                if chunks[n].Type >= 0 {
                    chunks[i].Type = chunks[n].Type
                }
            }
        }
        
        if chunks[i].Type < 0 {
            chunks[i].Type = types
            types += 1
        }
    }
    
    return types, chunks
}

func (lsbn *LargeScaleBrainNetwork) TrainCircuit (identifier string, alpha float64) bool {
    success_rate := float64(0)
    
    for success_rate < alpha {
        total := 0
        correct := 0
                
        for i := 0; i < len((*lsbn).Circuits[identifier].Data); i++ {
            total += 1
            res := (*lsbn).Circuits[identifier].Circuit.ExposeTo((*lsbn).Circuits[identifier].Data[i].Sensations)
            if len(res) > 0 {
                if res[0].outcome == (*lsbn).Circuits[identifier].Data[i].Type {
                    correct += 1
                }
                (*lsbn).Circuits[identifier].Circuit.CorrectFor(res, (*lsbn).Circuits[identifier].Data[i].Type, (*lsbn).Circuits[identifier].Data[i].Sensations)
            }
        }
        
        success_rate = float64(correct) / float64(total)
    }
    
    return true
}

func (lsbn *LargeScaleBrainNetwork) Grow (recursion string, circuit string, outs int) {
    (*lsbn).Circuits[recursion] = &LSBNCircuit{}
    (*lsbn).Circuits[recursion].Identifier = recursion
    (*lsbn).Circuits[recursion].Circuit = &Circuit{}
    (*lsbn).Circuits[recursion].StimLength = (*lsbn).Circuits[circuit].Types * ((*lsbn).Circuits[circuit].StimLength / (*lsbn).Circuits[circuit].ChunkLength)
    (*lsbn).Circuits[recursion].ChunkLength = (*lsbn).Circuits[recursion].StimLength
    (*lsbn).Circuits[recursion].RawStimuli = (*lsbn).Circuits[circuit].RawStimuli
    (*lsbn).Circuits[recursion].ConnectsTo = circuit
    (*lsbn).Circuits[recursion].Circuit.Neurogenesis((*lsbn).Circuits[recursion].StimLength, outs)
}

func (lsbn *LargeScaleBrainNetwork) Expose (circuit string, stimulus []float64) ([]float64, []RankedResult) {
    r := (*lsbn).Circuits[circuit]
    c := (*lsbn).Circuits[(*lsbn).Circuits[circuit].ConnectsTo]
    _, stim := lsbn.MakeChunks(c.ChunkLength, c.ChunkBeta, c.StimLength, stimulus)
    
    var ins []float64
    for _, s := range stim {
        in := make([]float64, c.Types)
        
        res := c.Circuit.ExposeTo(s.Sensations)
        if len(res) > 0 {
            in[res[0].outcome] = 1
        }
        for _, i := range in {
            ins = append(ins, i)
        }
    }
    
    res := r.Circuit.ExposeTo(ins)
    
    return ins, res
}

func (lsbn *LargeScaleBrainNetwork) Correct (circuit string, res []RankedResult, v int, stimulus []float64) {
    (*lsbn).Circuits[circuit].Circuit.CorrectFor(res, v, stimulus)
}