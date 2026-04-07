package panorama

import (
	"math"
)

type TransformType string

const (
	TransformCubemap TransformType = "cubemap"
	TransformEAC     TransformType = "eac"
	TransformOTC     TransformType = "otc"
)

type LayoutType string

const (
	LayoutNone   LayoutType = "none"
	LayoutRow    LayoutType = "row"
	LayoutColumn LayoutType = "column"
	LayoutCrossL LayoutType = "crossL"
	LayoutCrossR LayoutType = "crossR"
	LayoutCrossH LayoutType = "crossH"
)

type FaceIndex int

const (
	FacePX FaceIndex = iota
	FaceNX
	FacePY
	FaceNY
	FacePZ
	FaceNZ
)

type Index struct {
	X []float32
	Y []float32
}

type FaceIndices [6]*Index

type Options struct {
	Size       int
	Transform  TransformType
	Layout     LayoutType
	Inverse    string
	FaceNames  []string
	Order      []int
	Rotate     []int
	Resample   string
}

func DefaultOptions() *Options {
	return &Options{
		Transform: TransformCubemap,
		Layout:    LayoutNone,
		Resample:  "bilinear",
		FaceNames: []string{"0", "1", "2", "3", "4", "5"},
	}
}

func GenerateFaceIndices(size int, transform TransformType) (FaceIndices, error) {
	indices := FaceIndices{}

	ls := make([]float32, size)
	for i := 0; i < size; i++ {
		ls[i] = -1.0 + float32(i)*2.0/float32(size)
	}

	switch transform {
	case TransformEAC:
		for i := range ls {
			ls[i] = float32(math.Tan(float64(ls[i]) / (4.0 / math.Pi)))
		}
	case TransformOTC:
		for i := range ls {
			ls[i] = float32(math.Tan(float64(ls[i])*0.8687) / math.Tan(0.8687))
		}
	}

	x0 := make([]float32, size*size)
	y0 := make([]float32, size*size)
	x1 := make([]float32, size*size)
	y1 := make([]float32, size*size)

	for j := 0; j < size; j++ {
		for i := 0; i < size; i++ {
			idx := j*size + i
			xv := ls[i]
			yv := ls[j]

			x0[idx] = float32(math.Atan(float64(xv)))
			y0[idx] = float32(math.Atan2(float64(yv), math.Hypot(1.0, float64(xv))))
			x1[idx] = float32(math.Atan2(float64(xv), float64(yv)))
			y1[idx] = float32(math.Atan(math.Hypot(float64(yv), float64(xv))))
		}
	}

	pio2 := float32(math.Pi / 2.0)

	for i := range x0 {
		x0[i] /= pio2
		y0[i] /= pio2
		x1[i] /= pio2
		y1[i] /= pio2
	}

	indices[FacePX] = &Index{X: make([]float32, size*size), Y: make([]float32, size*size)}
	indices[FaceNX] = &Index{X: make([]float32, size*size), Y: make([]float32, size*size)}
	indices[FacePY] = &Index{X: make([]float32, size*size), Y: make([]float32, size*size)}
	indices[FaceNY] = &Index{X: make([]float32, size*size), Y: make([]float32, size*size)}
	indices[FacePZ] = &Index{X: make([]float32, size*size), Y: make([]float32, size*size)}
	indices[FaceNZ] = &Index{X: make([]float32, size*size), Y: make([]float32, size*size)}

	for i := 0; i < size*size; i++ {
		indices[FacePX].X[i] = x0[i] + 3
		indices[FacePX].Y[i] = y0[i] + 1

		indices[FaceNX].X[i] = x0[i] + 1
		indices[FaceNX].Y[i] = y0[i] + 1

		indices[FacePY].X[i] = mod4(x1[i] - 2)
		indices[FacePY].Y[i] = y1[i]

		indices[FaceNY].X[i] = mod4(4 - x1[i])
		indices[FaceNY].Y[i] = 2 - y1[i]

		indices[FacePZ].X[i] = x0[i] + 2
		indices[FacePZ].Y[i] = y0[i] + 1

		indices[FaceNZ].X[i] = mod4(x0[i])
		indices[FaceNZ].Y[i] = y0[i] + 1
	}

	return indices, nil
}

func mod4(x float32) float32 {
	result := float32(math.Mod(float64(x), 4.0))
	if result < 0 {
		result += 4
	}
	return result
}
