package app

import "testing"

func BenchmarkIntSignal(b *testing.B) {
	a := New()
	var f func(int)
	a.Load(func(loader Loader) {
		loader.Define("foo", &f)
	})
	a.Load(func(loader Loader) {
		loader.Implement("foo", func(int) {})
	})
	a.FinishLoad()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f(42)
	}
}

func BenchmarkStringSignal(b *testing.B) {
	a := New()
	var f func(string)
	a.Load(func(loader Loader) {
		loader.Define("foo", &f)
	})
	a.Load(func(loader Loader) {
		loader.Implement("foo", func(string) {})
	})
	a.FinishLoad()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f("foobar")
	}
}

func BenchmarkBoolSignal(b *testing.B) {
	a := New()
	var f func(bool)
	a.Load(func(loader Loader) {
		loader.Define("foo", &f)
		loader.Implement("foo", func(b bool) {})
	})
	a.FinishLoad()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f(true)
	}
}

func BenchmarkStructSignal(b *testing.B) {
	a := New()
	var f func(struct{ int })
	a.Load(func(loader Loader) {
		loader.Define("foo", &f)
	})
	a.Load(func(loader Loader) {
		loader.Implement("foo", func(struct{ int }) {})
	})
	a.FinishLoad()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f(struct{ int }{42})
	}
}
