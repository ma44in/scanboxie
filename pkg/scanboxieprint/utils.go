package scanboxieprint

import (
	"bytes"
	"fmt"
	"html/template"
	"image/color"

	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/code128"
)

// GetBarcodeSvg returns Barcode Image as SVG
// Thx to: https://github.com/boombuler/barcode/issues/57
func GetBarcodeSvg(barcode string) template.HTML {
	bc, err := code128.Encode(barcode)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	canvas := svg.New(buf)

	height := 100
	bounds := bc.Bounds()
	//canvas.Startview(bounds.Dx(), height, 0, 0, 100, 100)
	//canvas.Startpercent(bounds.Dx(), 100, "viewBox=\"0 0 100 100\"")

	// TODO heigth 50  problem ... still scales with aspect ratio
	canvas.Startraw("preserveAspectRatio=\"none\"", fmt.Sprintf("viewBox=\"0 0 %d 50\"", bounds.Dx()))

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		if bc.At(x, bounds.Min.Y) == color.Black {
			start := x
			x++

			for x < bounds.Max.X && bc.At(x, bounds.Min.Y) == color.Black {
				x++
			}

			canvas.Rect(start, 0, x-start, height, "fill:black")
		}
	}

	canvas.End()
	return template.HTML(buf.String())
}
