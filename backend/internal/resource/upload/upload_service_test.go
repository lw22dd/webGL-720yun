package upload

import (
	"testing"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/stretchr/testify/assert"
)

func TestNoopWriter(t *testing.T) {
	writer := &noopWriter{}

	n, err := writer.Write([]byte("test data"))

	assert.NoError(t, err)
	assert.Equal(t, 9, n)
}

func TestNoopWriterEmptySlice(t *testing.T) {
	writer := &noopWriter{}

	n, err := writer.Write([]byte{})

	assert.NoError(t, err)
	assert.Equal(t, 0, n)
}

func TestNoopWriterLargeData(t *testing.T) {
	writer := &noopWriter{}
	largeData := make([]byte, 1024*1024) // 1MB

	n, err := writer.Write(largeData)

	assert.NoError(t, err)
	assert.Equal(t, 1024*1024, n)
}

func TestProgressbarWithNoopWriter(t *testing.T) {
	writer := &noopWriter{}

	bar := progressbar.NewOptions64(100,
		progressbar.OptionSetWriter(writer),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
	)

	assert.NotNil(t, bar)

	for i := 0; i <= 10; i++ {
		bar.Set(i * 10)
		time.Sleep(10 * time.Millisecond)
	}

	err := bar.Close()
	assert.NoError(t, err)
}

func TestProgressbarOnCompletion(t *testing.T) {
	writer := &noopWriter{}
	completed := false

	bar := progressbar.NewOptions64(100,
		progressbar.OptionSetWriter(writer),
		progressbar.OptionOnCompletion(func() {
			completed = true
		}),
	)

	assert.NotNil(t, bar)

	bar.Set(100)
	bar.Close()

	assert.True(t, completed)
}
