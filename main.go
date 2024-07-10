package main

import (
    "os"
    "image"
    "image/png"
    "image/color"
    "github.com/Knetic/govaluate"
    "fmt"
)

func Calc(expression *govaluate.EvaluableExpression, x, y int) uint8{
    
    parameters := make(map[string]interface{}, 8)
    parameters["x"] = x
    parameters["y"] = y

    result, err := expression.Evaluate(parameters)
    if err != nil{
        fmt.Println(err)
    }
    return uint8(result.(float64))
}

func MakeImage(expression *govaluate.EvaluableExpression, dx, dy int) [][]uint8{
    s := make([][]uint8, dy)
    for i:=0; i < dy; i++{
        s = append(s, make([]uint8, dx))
        for j:=0; j < dx; j++{
            s[i] = append(s[i], Calc(expression, j, i))
        }
    }

    return s
}

func WriteImage(s [][]uint8){
    dx, dy := len(s[0]), len(s)
    img := image.NewNRGBA(image.Rect(0, 0, dx, dy))
    for i:=0; i < dy; i++{
        for j:=0; j < dx; j++{
            v := uint8(uint32(s[i][j])*256/uint32(dx))
            img.Set(i, j, color.RGBA{v, v, 255, 255})
        }
    }

    outputFile, err := os.Create("img")
    if err != nil{
        
    }

    png.Encode(outputFile, img)
    outputFile.Close()
}

func main(){
    var x, y int
    fmt.Scanf("%v %v", &x, &y)
    var s string
    fmt.Scanf("%v", &s)
    expression, err := govaluate.NewEvaluableExpression(s)
    if err != nil{
        fmt.Println(err)
    } else {
        img := MakeImage(expression, x, y)
        WriteImage(img)
    }
}
