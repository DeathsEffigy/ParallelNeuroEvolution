package main

import (
    "fmt"
    "path/filepath"
    "os"
    "strings"
    "regexp"
    "image"
    _"image/png"
    "math"
)

//const path_stim = "/users/fabianschneider/desktop/programming/go/PNE/data/letters"
const path_stim = "/users/fabianschneider/desktop/programming/go/PNE/data/numbers"

func main() {
    // circuit method
    c := Circuit{}
    c.Neurogenesis(256, 10)
    size := len(c.Cluster)
    
    //constim := map[string]int{"A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5}
    //stimcon := map[int]string{0: "A", 1: "B", 2: "C", 3: "D", 4: "E", 5: "F"}
    constim := map[string]int{"ZERO": 0, "ONE": 1, "TWO": 2, "THREE": 3, "FOUR": 4, "FIVE": 5, "SIX": 6, "SEVEN": 7, "EIGHT": 8, "NINE": 9}
    stimcon := map[int]string{0: "ZERO", 1: "ONE", 2: "TWO", 3: "THREE", 4: "FOUR", 5: "FIVE", 6: "SIX", 7: "SEVEN", 8: "EIGHT", 9: "NINE"}
    stimuli := LoadStimuli()
    
    count := 0
    right := 0
    error := 0
    
    for i := 0; i < 100; i++ {
        countThis := 0
        countRight := 0
        countError := 0
        for _, stimulus := range stimuli {
            res := c.ExposeTo(stimulus.GreyScale)
            if len(res) > 0 {
                count += 1
                countThis += 1
                if stimulus.Type == stimcon[res[0].outcome] {
                    right += 1
                    countRight += 1
                } else {
                    error += 1
                    countError += 1
                }
                c.CorrectFor(res, constim[stimulus.Type], stimulus.GreyScale)
            }
        }
        fmt.Printf("success_rate_overall=%f after trials=%d. success_rate this trial=%f.\n", (float64(right) / float64(count)), count, (float64(countRight) / float64(countThis)))
    }
    
    fmt.Printf("Size %d -> %d.\n", size, len(c.Cluster))
    obs := float64(right)
    exp := float64(1) / float64(len(constim)) * float64(count)
    x2 := ((obs - exp) * (obs - exp)) / exp
    p := chi2p(2, x2)
    fmt.Printf("Stats: X^2=%f, p=%f\n", x2, p)
    
    
    /*
    // LSBN method
    stimuli := LoadStimuli()
    constim := map[string]int{"ZERO": 0, "ONE": 1, "TWO": 2, "THREE": 3, "FOUR": 4, "FIVE": 5, "SIX": 6, "SEVEN": 7, "EIGHT": 8, "NINE": 9}
    var sensations []float64
    for _, stimulus := range stimuli {
        for _, sensation := range stimulus.GreyScale {
            sensations = append(sensations, sensation)
        }
    }
    lsbn := LargeScaleBrainNetwork{}
    fmt.Printf("Growing ShapesCircuit...\n")
    types, _, _ := lsbn.GrowCircuit("shapes", 64, 0.6, 256, sensations, true, 0.9)
    fmt.Printf("ShapesCircuit grown with types=%d.\n", types)
    lsbn.Grow("numbers", "shapes", 10)
    fmt.Printf("NumbersCircuit grown.\n")
    
    count := 0
    succs := 0
    for i := 0; i < 100; i++ {
        total := 0
        success := 0
        
        for _, stimulus := range stimuli {
            total += 1
            count += 1
            stim, res := lsbn.Expose("numbers", stimulus.GreyScale)
            if len(res) > 0 {
                if res[0].outcome == constim[stimulus.Type] {
                    success += 1
                    succs += 1
                }
                lsbn.Correct("numbers", res, constim[stimulus.Type], stim)
            }
        }
        
        fmt.Printf("Trials no.%d done with overall_success_rate=%.4f success_rate=%.4f.\n", i, float64(succs) / float64(count), float64(success) / float64(total))
    }
    
    obs := float64(succs)
    exp := float64(1) / float64(len(constim)) * float64(count)
    x2 := ((obs - exp) * (obs - exp)) / exp
    p := chi2p(2, x2)
    fmt.Printf("Stats: X^2=%f, p=%f\n", x2, p)*/
}

func chi2p(dof int, distance float64) float64 {
    return gammaIncQ(.5*float64(dof), .5*distance)
}

type ifctn func(float64) float64

func gammaIncQ(a, x float64) float64 {
    aa1 := a - 1
    var f ifctn = func(t float64) float64 {
        return math.Pow(t, aa1) * math.Exp(-t)
    }
    y := aa1
    h := 1.5e-2
    for f(y)*(x-y) > 2e-8 && y < x {
        y += .4
    }
    if y > x {
        y = x
    }
    return 1 - simpson38(f, 0, y, int(y/h/math.Gamma(a)))
}

func simpson38(f ifctn, a, b float64, n int) float64 {
    h := (b - a) / float64(n)
    h1 := h / 3
    sum := f(a) + f(b)
    for j := 3*n - 1; j > 0; j-- {
        if j%3 == 0 {
            sum += 2 * f(a+h1*float64(j))
        } else {
            sum += 3 * f(a+h1*float64(j))
        }
    }
    return h * sum / 8
}

type LoadStimulus struct {
    Type string
    Variant string
    Path string
}

type ImgStimulus struct {
    Type string
    Variant string
    Path string
    GreyScale []float64
}

type Pixel struct {
    R, G, B, A int
}

func LoadStimuli() []ImgStimulus {
    files := []LoadStimulus{}
    
    filepath.Walk(path_stim, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            if path[len(path)-4:] == ".png" {
                s := strings.Split(path, "/")
                name := s[len(s)-1]
                name = name[:len(name)-4]
                re := regexp.MustCompile(`[^a-zA-Z]+`)
                Type := re.ReplaceAllString(name, "")
                re2 := regexp.MustCompile(`[^0-9]+`)
                Variant := re2.ReplaceAllString(name, "")
                files = append(files, LoadStimulus{Type, Variant, path})
            }
        }
        return nil
    })
    
    images := []ImgStimulus{}
    
    for _, stimulus := range files {
        reader, err := os.Open(stimulus.Path)
        if err != nil {
            continue
        }
        defer reader.Close()
        
        img, _, err2 := image.Decode(reader)
        if err2 != nil {
            continue
        }
        
        bounds := img.Bounds()
        width, height := bounds.Max.X, bounds.Max.Y
        
        var Greyscale []float64
        for y := 0; y < height; y++ {
            for x := 0; x < width; x++ {
                RGBA := RGBAToPixel(img.At(x, y).RGBA())
                GS := GSToGR(PixelToGS(RGBA))
                Greyscale = append(Greyscale, GS)
            }
        }
        
        images = append(images, ImgStimulus{stimulus.Type, stimulus.Variant, stimulus.Path, Greyscale})
    }
    
    return images
}

func RGBAToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
    return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func PixelToGS(pixel Pixel) int {
    return int(float64(pixel.R) * 0.299 + float64(pixel.G) * 0.587 + float64(pixel.B) * 0.114)
}

func GSToGR(gs int) float64 {
    // this makes white the priority. we want black to be the priority.
    //return (float64(gs) / float64(255)) * 0.2
    bs := 255 - gs
    return (float64(bs) / float64(255)) * 0.2
}