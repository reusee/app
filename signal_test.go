package app

import "testing"

func TestSignal1(t *testing.T) {
	a := New()
	foo := func() int {
		return 42
	}
	bar := func() string {
		return "bar"
	}
	baz := func() []int {
		return []int{1, 2, 3}
	}
	a.Load(func(loader Loader) {
		loader.Emit("foo", &foo)
		loader.Emit("bar", &bar)
		loader.Emit("baz", &baz)
	})

	fooEmitted := false
	barEmitted := false
	a.Load(func(loader Loader) {
		loader.Listen("foo", func(n int) {
			if n != 42 {
				t.Fatal()
			}
			fooEmitted = true
		})
		loader.Listen("bar", func(s string) {
			if s != "bar" {
				t.Fatal()
			}
			barEmitted = true
		})
		loader.Listen("baz", func(is []int) {})
	})
	a.FinishLoad()

	foo()
	if !fooEmitted {
		t.Fatal()
	}
	bar()
	if !barEmitted {
		t.Fatal()
	}
	baz()
}

func TestSignal2(t *testing.T) {
	a := New()
	a.Load(func(loader Loader) {
		func() {
			defer func() {
				err := recover()
				if err == nil {
					t.Fatal()
				}
				if err.(string) != "module func(app.Loader): emitter foo is not a pointer to function" {
					t.Fatal()
				}
			}()
			loader.Emit("foo", 42)
		}()
	})

	a.Load(func(loader Loader) {
		f1 := func() int { return 42 }
		loader.Emit("foo", &f1)
		f2 := func() {}
		func() {
			defer func() {
				err := recover()
				if err == nil {
					t.Fatal()
				}
				if err.(string) != "module func(app.Loader): multiple emitter foo" {
					t.Fatal()
				}
			}()
			loader.Emit("foo", &f2)
		}()
	})

	a.Load(func(loader Loader) {
		func() {
			defer func() {
				err := recover()
				if err == nil {
					t.Fatal()
				}
				if err.(string) != "module func(app.Loader): listener foo is not a function" {
					t.Fatal()
				}
			}()
			loader.Listen("foo", 42)
		}()
	})

	a.Load(func(loader Loader) {
		loader.Listen("foo", func(int) {})
	})
	a.FinishLoad()
}

func TestSignal3(t *testing.T) {
	a := New()
	a.Load(func(loader Loader) {
		f := func() int { return 42 }
		loader.Emit("foo", &f)
	})
	a.Load(func(loader Loader) {
		loader.Listen("foo", func() {})
	})
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal()
			}
			if err.(string) != "foo not match, emit func() int, listen func()" {
				t.Fatal()
			}
		}()
		a.FinishLoad()
	}()
}

func TestSignal4(t *testing.T) {
	a := New()
	a.Load(func(loader Loader) {
		f := func() int { return 42 }
		loader.Emit("foo", &f)
	})
	a.Load(func(loader Loader) {
		loader.Listen("foo", func(string) {})
	})
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal()
			}
			if err.(string) != "foo not match at arg #0, emit int, listen string" {
				t.Fatal()
			}
		}()
		a.FinishLoad()
	}()
}

func TestSignal5(t *testing.T) {
	a := New()
	a.Load(func(loader Loader) {
		f := func() int { return 42 }
		loader.Emit("foo", &f)
	})
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal()
			}
			if err.(string) != "foo not listened" {
				t.Fatal()
			}
		}()
		a.FinishLoad()
	}()
}
