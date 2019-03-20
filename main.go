package main

import (
    "fmt"
    "path/filepath"
    "os"
    "strings"
    "regexp"
    "image"
    _"image/png"
)

const path_stim = "/users/fabianschneider/desktop/programming/go/PNE/data"

func main() {
    c := Circuit{}
    c.Neurogenesis(256, 4)
    
    stimuli := LoadStimuli()
    
    for _, stimulus := range stimuli {
        res := c.ExposeTo(stimulus.GreyScale)
        fmt.Printf("NN thinks %s(%s) is %d.\n", stimulus.Type, stimulus.Variant, res[0].outcome)
    }
    c.ExposeTo(stimuli[0].GreyScale)
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

func LoadStimuli () []ImgStimulus {
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
    return float64(gs) / float64(255)
}