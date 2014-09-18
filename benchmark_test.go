package app

import "testing"

func BenchmarkIntSignal(b *testing.B) {
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

func BenchmarkStringSignal(b *testing.B) {
	a := New()
	f := func() string {
		return "foobar"
	}
	a.Load(func(loader Loader) {
		loader.Emit("foo", &f)
	})
	a.Load(func(loader Loader) {
		loader.Listen("foo", func(string) {})
	})
	a.FinishLoad()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f()
	}
}

func BenchmarkStructSignal(b *testing.B) {
	AddSignalType((*func() struct{ int })(nil), func(emit interface{}, listens []interface{}) interface{} {
		return func() (ret struct{ int }) {
			ret = emit.(func() struct{ int })()
			for _, l := range listens {
				l.(func(struct{ int }))(ret)
			}
			return
		}
	})

	a := New()
	f := func() struct{ int } {
		return struct{ int }{42}
	}
	a.Load(func(loader Loader) {
		loader.Emit("foo", &f)
	})
	a.Load(func(loader Loader) {
		loader.Listen("foo", func(struct{ int }) {})
	})
	a.FinishLoad()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f()
	}
}

func BenchmarkBoolSignal(b *testing.B) {
	a := New()
	f := func() bool {
		return true
	}
	a.Load(func(loader Loader) {
		loader.Emit("foo", &f)
		loader.Listen("foo", func(b bool) {})
	})
	a.FinishLoad()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f()
	}
}
