package hash

import (
	"testing"
)

func BenchmarkSha256(b *testing.B) {
	b.ResetTimer()
	str := "lucifer"
	for i := 0; i < b.N; i++ {
		s := []byte(str)
		_ = ByteMd5(s)
	}
}
