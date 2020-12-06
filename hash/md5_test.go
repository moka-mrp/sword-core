package hash

import (
	"testing"
)

func TestMd5(t *testing.T) {
	name := "lucifer"
	res := Md5(name)
	t.Log(res)
	t.Log(Md5Check(name, res))
}

func BenchmarkMd5(b *testing.B) {
	b.ResetTimer()
	name := "lucifer"
	for i := 0; i < b.N; i++ {
		_ = Md5(name)
	}
}

func BenchmarkByteMd5(b *testing.B) {
	b.ResetTimer()
	name := "lucifer"
	for i := 0; i < b.N; i++ {
		s := []byte(name)
		_ = ByteMd5(s)
	}
}
