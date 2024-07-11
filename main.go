package main

import (
    "os"
    "image"
    "image/png"
    "image/color"
    "github.com/Knetic/govaluate"
    "errors"
    "fmt"
    "net/http"
    "github.com/labstack/echo/v4"
)

func Calc(expression *govaluate.EvaluableExpression, x, y int) (uint8, error){
    
    parameters := make(map[string]interface{}, 8)
    parameters["x"] = x
    parameters["y"] = y

    result, err := expression.Evaluate(parameters)
    if err != nil{
        return 0, errors.New("error in calculation")
    }
    return uint8(result.(float64)), nil
}

func MakeImage(expression *govaluate.EvaluableExpression, dx, dy int, name string) error{
    s := make([][]uint8, dy)
    for i:=0; i < dy; i++{
        s = append(s, make([]uint8, dx))
        for j:=0; j < dx; j++{
            v, err := Calc(expression, j, i)
            if err != nil{
                return errors.New("error in calculation")
            } 
            s[i] = append(s[i], v)
        }
    }

    img := image.NewNRGBA(image.Rect(0, 0, dx, dy))
    for i:=0; i < dy; i++{
        for j:=0; j < dx; j++{
            v := uint8(uint32(s[i][j])*256/uint32(dx))
            img.Set(i, j, color.RGBA{v, v, 255, 255})
        }
    }

    outputFile, err := os.Create(name)
    if err != nil{
        return errors.New("error in writing to the file")
    }
    png.Encode(outputFile, img)
    outputFile.Close()
    return nil
}

func gen(c echo.Context) error {
    s := c.Param("exp")
    name := c.Param("name")
    expression, err := govaluate.NewEvaluableExpression(s)
    if err != nil{
        return c.String(http.StatusOK, "There is an error in expression")
    } else {
        err := MakeImage(expression, 256, 256, name)
        if err != nil{
            return c.String(http.StatusOK, fmt.Sprint(err))
        }
        return c.File(name)
    }
}

func main(){
    e := echo.New()
    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Just a fucking root")
    })
    e.GET("/gen/:exp/:name", gen)
    e.Logger.Fatal(e.Start(":1323"))
}
