package panorama

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/disintegration/imaging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFaceIndices(t *testing.T) {
	tests := []struct {
		name      string
		size      int
		transform TransformType
	}{
		{"Cubemap_256", 256, TransformCubemap},
		{"EAC_256", 256, TransformEAC},
		{"OTC_256", 256, TransformOTC},
		{"Cubemap_512", 512, TransformCubemap},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indices, err := GenerateFaceIndices(tt.size, tt.transform)
			require.NoError(t, err, "GenerateFaceIndices should not return error")
			assert.Len(t, indices, 6, "Should generate 6 face indices")

			for i, idx := range indices {
				assert.NotNil(t, idx, "Face index %d should not be nil", i)
				assert.Len(t, idx.X, tt.size*tt.size, "X coordinates length should be size^2")
				assert.Len(t, idx.Y, tt.size*tt.size, "Y coordinates length should be size^2")

				for j := 0; j < len(idx.X); j++ {
					assert.GreaterOrEqual(t, idx.X[j], float32(-1.0), "X coordinate should be >= -1")
					assert.LessOrEqual(t, idx.X[j], float32(5.0), "X coordinate should be <= 5")
					assert.GreaterOrEqual(t, idx.Y[j], float32(-1.0), "Y coordinate should be >= -1")
					assert.LessOrEqual(t, idx.Y[j], float32(3.0), "Y coordinate should be <= 3")
				}
			}
		})
	}
}

func TestConverter_Convert(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "panorama_test_*")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tempDir)

	inputPath := filepath.Join(tempDir, "test_panorama.jpg")
	err = createTestPanorama(inputPath, 4096, 2048)
	require.NoError(t, err, "Failed to create test panorama")

	tests := []struct {
		name   string
		opts   *Options
		layout LayoutType
	}{
		{
			name: "None_Layout",
			opts: &Options{
				Size:      512,
				Transform: TransformCubemap,
				Layout:    LayoutNone,
				FaceNames: []string{"r", "l", "u", "d", "f", "b"},
			},
			layout: LayoutNone,
		},
		{
			name: "Row_Layout",
			opts: &Options{
				Size:      512,
				Transform: TransformCubemap,
				Layout:    LayoutRow,
			},
			layout: LayoutRow,
		},
		{
			name: "Column_Layout",
			opts: &Options{
				Size:      512,
				Transform: TransformCubemap,
				Layout:    LayoutColumn,
			},
			layout: LayoutColumn,
		},
		{
			name: "CrossL_Layout",
			opts: &Options{
				Size:      512,
				Transform: TransformCubemap,
				Layout:    LayoutCrossL,
			},
			layout: LayoutCrossL,
		},
		{
			name: "CrossR_Layout",
			opts: &Options{
				Size:      512,
				Transform: TransformCubemap,
				Layout:    LayoutCrossR,
			},
			layout: LayoutCrossR,
		},
		{
			name: "CrossH_Layout",
			opts: &Options{
				Size:      512,
				Transform: TransformCubemap,
				Layout:    LayoutCrossH,
			},
			layout: LayoutCrossH,
		},
		{
			name: "EAC_Transform",
			opts: &Options{
				Size:      512,
				Transform: TransformEAC,
				Layout:    LayoutRow,
			},
			layout: LayoutRow,
		},
		{
			name: "OTC_Transform",
			opts: &Options{
				Size:      512,
				Transform: TransformOTC,
				Layout:    LayoutRow,
			},
			layout: LayoutRow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDir := filepath.Join(tempDir, tt.name)
			err := os.MkdirAll(outputDir, 0755)
			require.NoError(t, err, "Failed to create output directory")

			outputPath := filepath.Join(outputDir, "output")

			converter := NewConverter(tt.opts)
			err = converter.Convert(inputPath, outputPath)
			require.NoError(t, err, "Convert should not return error")

			verifyOutput(t, outputDir, tt.layout, tt.opts.Size)
		})
	}
}

func TestConverter_WithInverse(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "panorama_test_*")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tempDir)

	inputPath := filepath.Join(tempDir, "test_panorama.jpg")
	err = createTestPanorama(inputPath, 4096, 2048)
	require.NoError(t, err, "Failed to create test panorama")

	tests := []struct {
		name    string
		inverse string
	}{
		{"Horizontal", "horizontal"},
		{"Vertical", "vertical"},
		{"Both", "both"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDir := filepath.Join(tempDir, "inverse_"+tt.name)
			err := os.MkdirAll(outputDir, 0755)
			require.NoError(t, err, "Failed to create output directory")

			outputPath := filepath.Join(outputDir, "output")

			opts := &Options{
				Size:      256,
				Transform: TransformCubemap,
				Layout:    LayoutRow,
				Inverse:   tt.inverse,
			}

			converter := NewConverter(opts)
			err = converter.Convert(inputPath, outputPath)
			require.NoError(t, err, "Convert should not return error")

			outputFile := outputPath + ".jpg"
			assert.FileExists(t, outputFile, "Output file should exist")
		})
	}
}

func TestConverter_WithOrderAndRotate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "panorama_test_*")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tempDir)

	inputPath := filepath.Join(tempDir, "test_panorama.jpg")
	err = createTestPanorama(inputPath, 4096, 2048)
	require.NoError(t, err, "Failed to create test panorama")

	outputDir := filepath.Join(tempDir, "order_rotate")
	err = os.MkdirAll(outputDir, 0755)
	require.NoError(t, err, "Failed to create output directory")

	outputPath := filepath.Join(outputDir, "output")

	opts := &Options{
		Size:      256,
		Transform: TransformCubemap,
		Layout:    LayoutRow,
		Order:     []int{0, 1, 5, 4, 2, 3},
		Rotate:    []int{270, 90, 180, 0, 0, 180},
	}

	converter := NewConverter(opts)
	err = converter.Convert(inputPath, outputPath)
	require.NoError(t, err, "Convert should not return error")

	outputFile := outputPath + ".jpg"
	assert.FileExists(t, outputFile, "Output file should exist")
}

func TestConverter_DefaultOptions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "panorama_test_*")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tempDir)

	inputPath := filepath.Join(tempDir, "test_panorama.jpg")
	err = createTestPanorama(inputPath, 4096, 2048)
	require.NoError(t, err, "Failed to create test panorama")

	outputDir := filepath.Join(tempDir, "default")
	err = os.MkdirAll(outputDir, 0755)
	require.NoError(t, err, "Failed to create output directory")

	outputPath := filepath.Join(outputDir, "output")

	converter := NewConverter(nil)
	err = converter.Convert(inputPath, outputPath)
	require.NoError(t, err, "Convert with default options should not return error")

	for i := 0; i < 6; i++ {
		outputFile := filepath.Join(outputDir, "output_"+string(rune('0'+i))+".jpg")
		assert.FileExists(t, outputFile, "Output file %d should exist", i)
	}
}

func verifyOutput(t *testing.T, outputDir string, layout LayoutType, size int) {
	switch layout {
	case LayoutNone:
		files, err := filepath.Glob(filepath.Join(outputDir, "*.jpg"))
		require.NoError(t, err, "Failed to list output files")
		assert.Len(t, files, 6, "Should generate 6 separate face files")

		for _, file := range files {
			img, err := imaging.Open(file)
			require.NoError(t, err, "Failed to load output image")

			assert.Equal(t, size, img.Bounds().Dx(), "Face width should match size")
			assert.Equal(t, size, img.Bounds().Dy(), "Face height should match size")
		}

	case LayoutRow:
		files, err := filepath.Glob(filepath.Join(outputDir, "*.jpg"))
		require.NoError(t, err, "Failed to list output files")
		require.Len(t, files, 1, "Should generate 1 row layout file")

		img, err := imaging.Open(files[0])
		require.NoError(t, err, "Failed to load output image")

		assert.Equal(t, size*6, img.Bounds().Dx(), "Row layout width should be 6*size")
		assert.Equal(t, size, img.Bounds().Dy(), "Row layout height should be size")

	case LayoutColumn:
		files, err := filepath.Glob(filepath.Join(outputDir, "*.jpg"))
		require.NoError(t, err, "Failed to list output files")
		require.Len(t, files, 1, "Should generate 1 column layout file")

		img, err := imaging.Open(files[0])
		require.NoError(t, err, "Failed to load output image")

		assert.Equal(t, size, img.Bounds().Dx(), "Column layout width should be size")
		assert.Equal(t, size*6, img.Bounds().Dy(), "Column layout height should be 6*size")

	case LayoutCrossL, LayoutCrossR:
		files, err := filepath.Glob(filepath.Join(outputDir, "*.jpg"))
		require.NoError(t, err, "Failed to list output files")
		require.Len(t, files, 1, "Should generate 1 cross layout file")

		img, err := imaging.Open(files[0])
		require.NoError(t, err, "Failed to load output image")

		assert.Equal(t, size*4, img.Bounds().Dx(), "Cross layout width should be 4*size")
		assert.Equal(t, size*3, img.Bounds().Dy(), "Cross layout height should be 3*size")

	case LayoutCrossH:
		files, err := filepath.Glob(filepath.Join(outputDir, "*.jpg"))
		require.NoError(t, err, "Failed to list output files")
		require.Len(t, files, 1, "Should generate 1 cross layout file")

		img, err := imaging.Open(files[0])
		require.NoError(t, err, "Failed to load output image")

		assert.Equal(t, size*3, img.Bounds().Dx(), "CrossH layout width should be 3*size")
		assert.Equal(t, size*4, img.Bounds().Dy(), "CrossH layout height should be 4*size")
	}
}
