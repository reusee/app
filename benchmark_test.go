package app

import "testing"

func BenchmarkSignal(b *testing.B) {
	a := New()
	f := func() int {
		return 42
	}
	a.Load(func(loader Loader) {
		loader.Emit("foo", &f)
	})
	a.Load(func(loader Loader) {
		loader.Listen("foo", func(int) {})
	})
	a.FinishLoad()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f()
	}
}
