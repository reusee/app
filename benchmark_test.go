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
	SignalHandler((*func() struct{ int })(nil), func(emit interface{}, listens []interface{}) {
		emitPtr := emit.(*func() struct{ int })
		e := *emitPtr
		*emitPtr = func() (ret struct{ int }) {
			ret = e()
			for _, listen := range listens {
				listen.(func(struct{ int }))(ret)
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
