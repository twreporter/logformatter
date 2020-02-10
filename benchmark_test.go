package logformatter

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
)

func genErrorFrames(at, depth int) error {
	if at >= depth {
		return errors.Errorf("Synetic errors")
	}
	return genErrorFrames(at+1, depth)
}

func BenchmarkFormatStack(b *testing.B) {
	for _, c := range [...]struct {
		frame int
	}{
		{frame: 10},
		{frame: 20},
		{frame: 30},
		{frame: 100},
	} {
		name := fmt.Sprint(c.frame) + " frame"
		b.Run(name, func(b *testing.B) {
			err := genErrorFrames(0, c.frame)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				FormatStack(err)
			}
			b.StopTimer()
		})
	}
}
