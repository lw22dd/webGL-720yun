package panorama

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"runtime"
	"sync"

	"github.com/disintegration/imaging"
	"golang.org/x/sync/errgroup"
)

type Converter struct {
	options *Options
}

func NewConverter(opts *Options) *Converter {
	if opts == nil {
		opts = DefaultOptions()
	}
	return &Converter{options: opts}
}

func (c *Converter) Convert(inputPath, outputPath string) error {
	img, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	size := c.options.Size
	if size <= 0 {
		size = img.Bounds().Dx() / 4
	}

	indices, err := GenerateFaceIndices(size, c.options.Transform)
	if err != nil {
		return fmt.Errorf("failed to generate indices: %w", err)
	}

	if c.options.Order != nil {
		indices = c.reorderIndices(indices)
	}

	if c.options.Rotate != nil {
		indices = c.rotateIndices(indices, size)
	}

	switch c.options.Layout {
	case LayoutNone:
		return c.saveSeparateFaces(img, indices, size, outputPath)
	case LayoutRow, LayoutColumn:
		return c.saveJoinedFaces(img, indices, size, outputPath)
	case LayoutCrossL, LayoutCrossR, LayoutCrossH:
		return c.saveCrossLayout(img, indices, size, outputPath)
	default:
		return c.saveSeparateFaces(img, indices, size, outputPath)
	}
}

func (c *Converter) saveSeparateFaces(img image.Image, indices FaceIndices, size int, outputPath string) error {
	faceNames := c.options.FaceNames
	if len(faceNames) != 6 {
		faceNames = []string{"0", "1", "2", "3", "4", "5"}
	}

	var eg errgroup.Group

	for i := 0; i < 6; i++ {
		i := i
		eg.Go(func() error {
			faceImg := c.extractFace(img, indices[i], size)

			if c.options.Inverse != "" {
				faceImg = c.applyInverse(faceImg, c.options.Inverse)
			}

			outputFile := fmt.Sprintf("%s_%s.jpg", outputPath, faceNames[i])
			return c.saveImage(faceImg, outputFile)
		})
	}

	return eg.Wait()
}

func (c *Converter) saveJoinedFaces(img image.Image, indices FaceIndices, size int, outputPath string) error {
	faces := make([]image.Image, 6)

	for i := 0; i < 6; i++ {
		face := c.extractFace(img, indices[i], size)

		if c.options.Inverse != "" {
			face = c.applyInverse(face, c.options.Inverse)
		}
		faces[i] = face
	}

	var result image.Image
	if c.options.Layout == LayoutRow {
		result = c.joinHorizontal(faces)
	} else {
		result = c.joinVertical(faces)
	}

	return c.saveImage(result, outputPath+".jpg")
}

func (c *Converter) saveCrossLayout(img image.Image, indices FaceIndices, size int, outputPath string) error {
	faces := make([]image.Image, 6)

	for i := 0; i < 6; i++ {
		face := c.extractFace(img, indices[i], size)
		faces[i] = face
	}

	var result image.Image

	s1 := size
	s2 := size * 2
	s3 := size * 3
	s4 := size * 4

	switch c.options.Layout {
	case LayoutCrossL:
		result = c.createCrossL(faces, s1, s2, s3, s4)
	case LayoutCrossR:
		result = c.createCrossR(faces, s1, s2, s3, s4)
	case LayoutCrossH:
		result = c.createCrossH(faces, s1, s2, s3, s4)
	}

	if c.options.Inverse != "" {
		result = c.applyInverse(result, c.options.Inverse)
	}

	return c.saveImage(result, outputPath+".jpg")
}

func (c *Converter) createCrossL(faces []image.Image, s1, s2, s3, s4 int) image.Image {
	result := image.NewRGBA(image.Rect(0, 0, s4, s3))

	c.insertFace(result, faces[FaceNX], 0, s1)
	c.insertFace(result, faces[FacePZ], s1, s1)
	c.insertFace(result, faces[FacePX], s2, s1)
	c.insertFace(result, faces[FaceNZ], s3, s1)
	c.insertFace(result, faces[FacePY], s1, 0)
	c.insertFace(result, faces[FaceNY], s1, s2)

	return result
}

func (c *Converter) createCrossR(faces []image.Image, s1, s2, s3, s4 int) image.Image {
	result := image.NewRGBA(image.Rect(0, 0, s4, s3))

	c.insertFace(result, faces[FaceNZ], 0, s1)
	c.insertFace(result, faces[FaceNX], s1, s1)
	c.insertFace(result, faces[FacePZ], s2, s1)
	c.insertFace(result, faces[FacePX], s3, s1)
	c.insertFace(result, faces[FacePY], s2, 0)
	c.insertFace(result, faces[FaceNY], s2, s2)

	return result
}

func (c *Converter) createCrossH(faces []image.Image, s1, s2, s3, s4 int) image.Image {
	result := image.NewRGBA(image.Rect(0, 0, s3, s4))

	c.insertFace(result, faces[FaceNX], 0, s1)
	c.insertFace(result, faces[FacePZ], s1, s1)
	c.insertFace(result, faces[FacePX], s2, s1)
	c.insertFace(result, imaging.Rotate180(faces[FaceNZ]), s1, s3)
	c.insertFace(result, faces[FacePY], s1, 0)
	c.insertFace(result, faces[FaceNY], s1, s2)

	return result
}

func (c *Converter) insertFace(dst *image.RGBA, face image.Image, x, y int) {
	bounds := face.Bounds()
	for i := bounds.Min.Y; i < bounds.Max.Y; i++ {
		for j := bounds.Min.X; j < bounds.Max.X; j++ {
			dst.Set(j+x, i+y, face.At(j, i))
		}
	}
}

func (c *Converter) extractFace(img image.Image, idx *Index, size int) image.Image {
	result := image.NewRGBA(image.Rect(0, 0, size, size))
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	fac := float64(imgWidth) / 4.0

	var rgbaImg *image.RGBA
	if rgba, ok := img.(*image.RGBA); ok {
		rgbaImg = rgba
	} else {
		rgbaImg = image.NewRGBA(bounds)
		for y := 0; y < imgHeight; y++ {
			for x := 0; x < imgWidth; x++ {
				rgbaImg.Set(x, y, img.At(x, y))
			}
		}
	}

	srcPix := rgbaImg.Pix
	srcStride := rgbaImg.Stride
	dstPix := result.Pix
	dstStride := result.Stride

	numCPU := runtime.NumCPU()
	if numCPU < 1 {
		numCPU = 1
	}

	rowsPerGoroutine := (size + numCPU - 1) / numCPU
	if rowsPerGoroutine < 1 {
		rowsPerGoroutine = 1
	}

	var wg sync.WaitGroup

	for startRow := 0; startRow < size; startRow += rowsPerGoroutine {
		endRow := startRow + rowsPerGoroutine
		if endRow > size {
			endRow = size
		}

		wg.Add(1)
		go func(yStart, yEnd int) {
			defer wg.Done()

			for y := yStart; y < yEnd; y++ {
				for x := 0; x < size; x++ {
					i := y*size + x
					srcX := int(float64(idx.X[i]) * fac)
					srcY := int(float64(idx.Y[i]) * fac)

					srcX = srcX % imgWidth
					if srcX < 0 {
						srcX += imgWidth
					}
					srcY = srcY % imgHeight
					if srcY < 0 {
						srcY += imgHeight
					}

					if srcX >= 0 && srcX < imgWidth && srcY >= 0 && srcY < imgHeight {
						srcOffset := srcY*srcStride + srcX*4
						dstOffset := y*dstStride + x*4

						dstPix[dstOffset] = srcPix[srcOffset]
						dstPix[dstOffset+1] = srcPix[srcOffset+1]
						dstPix[dstOffset+2] = srcPix[srcOffset+2]
						dstPix[dstOffset+3] = srcPix[srcOffset+3]
					}
				}
			}
		}(startRow, endRow)
	}

	wg.Wait()
	return result
}

func (c *Converter) applyInverse(img image.Image, inverse string) image.Image {
	switch inverse {
	case "horizontal":
		return imaging.FlipH(img)
	case "vertical":
		return imaging.FlipV(img)
	case "both":
		return imaging.FlipH(imaging.FlipV(img))
	default:
		return img
	}
}

func (c *Converter) joinHorizontal(images []image.Image) image.Image {
	if len(images) == 0 {
		return nil
	}

	height := images[0].Bounds().Dy()
	totalWidth := 0
	for _, img := range images {
		totalWidth += img.Bounds().Dx()
	}

	result := image.NewRGBA(image.Rect(0, 0, totalWidth, height))
	x := 0
	for _, img := range images {
		c.insertFace(result, img, x, 0)
		x += img.Bounds().Dx()
	}

	return result
}

func (c *Converter) joinVertical(images []image.Image) image.Image {
	if len(images) == 0 {
		return nil
	}

	width := images[0].Bounds().Dx()
	totalHeight := 0
	for _, img := range images {
		totalHeight += img.Bounds().Dy()
	}

	result := image.NewRGBA(image.Rect(0, 0, width, totalHeight))
	y := 0
	for _, img := range images {
		c.insertFace(result, img, 0, y)
		y += img.Bounds().Dy()
	}

	return result
}

func (c *Converter) reorderIndices(indices FaceIndices) FaceIndices {
	if len(c.options.Order) != 6 {
		return indices
	}

	reordered := FaceIndices{}
	for i, order := range c.options.Order {
		if order >= 0 && order < 6 {
			reordered[i] = indices[order]
		}
	}
	return reordered
}

func (c *Converter) rotateIndices(indices FaceIndices, size int) FaceIndices {
	if len(c.options.Rotate) != 6 {
		return indices
	}

	rotated := FaceIndices{}
	for i, angle := range c.options.Rotate {
		if indices[i] == nil {
			continue
		}

		rotated[i] = &Index{
			X: make([]float32, size*size),
			Y: make([]float32, size*size),
		}

		for j := 0; j < size*size; j++ {
			rotated[i].X[j] = indices[i].X[j]
			rotated[i].Y[j] = indices[i].Y[j]
		}

		switch angle {
		case 90:
			for j := 0; j < size*size; j++ {
				row := j / size
				col := j % size
				newIdx := col*size + (size - 1 - row)
				rotated[i].X[newIdx] = indices[i].X[j]
				rotated[i].Y[newIdx] = indices[i].Y[j]
			}
		case 180:
			for j := 0; j < size*size; j++ {
				newIdx := size*size - 1 - j
				rotated[i].X[newIdx] = indices[i].X[j]
				rotated[i].Y[newIdx] = indices[i].Y[j]
			}
		case 270:
			for j := 0; j < size*size; j++ {
				row := j / size
				col := j % size
				newIdx := (size-1-col)*size + row
				rotated[i].X[newIdx] = indices[i].X[j]
				rotated[i].Y[newIdx] = indices[i].Y[j]
			}
		}
	}

	return rotated
}

func (c *Converter) saveImage(img image.Image, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
}

func createTestPanorama(path string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	stripWidth := width / 3

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var col color.Color
			if x < stripWidth {
				col = color.RGBA{R: 255, G: 0, B: 0, A: 255}
			} else if x < stripWidth*2 {
				col = color.RGBA{R: 0, G: 255, B: 0, A: 255}
			} else {
				col = color.RGBA{R: 0, G: 0, B: 255, A: 255}
			}
			img.Set(x, y, col)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
}
