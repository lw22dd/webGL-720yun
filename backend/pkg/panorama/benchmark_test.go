package panorama

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkConvert(b *testing.B) {
	tempDir, _ := os.MkdirTemp("", "bench_test_")
	defer os.RemoveAll(tempDir)

	inputPath := filepath.Join(tempDir, "test_panorama.jpg")
	createTestPanorama(inputPath, 4096, 2048)

	benchmarks := []struct {
		name  string
		opts  *Options
	}{
		{"None_256", &Options{Size: 256, Transform: TransformCubemap, Layout: LayoutNone}},
		{"Row_256", &Options{Size: 256, Transform: TransformCubemap, Layout: LayoutRow}},
		{"CrossL_256", &Options{Size: 256, Transform: TransformCubemap, Layout: LayoutCrossL}},
		{"None_512", &Options{Size: 512, Transform: TransformCubemap, Layout: LayoutNone}},
		{"Row_512", &Options{Size: 512, Transform: TransformCubemap, Layout: LayoutRow}},
		{"CrossL_512", &Options{Size: 512, Transform: TransformCubemap, Layout: LayoutCrossL}},
		{"None_1024", &Options{Size: 1024, Transform: TransformCubemap, Layout: LayoutNone}},
		{"Row_1024", &Options{Size: 1024, Transform: TransformCubemap, Layout: LayoutRow}},
		{"CrossL_1024", &Options{Size: 1024, Transform: TransformCubemap, Layout: LayoutCrossL}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				outputPath := filepath.Join(tempDir, "output")
				converter := NewConverter(bm.opts)
				converter.Convert(inputPath, outputPath)
			}
		})
	}
}
