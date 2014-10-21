package app

import "testing"

func TestFunc(t *testing.T) {
	a := New()
	var foo func(int)
	var bar func(string)
	var baz func(bool)
	var qux func(float64)
	var quux func()
	a.Load(func(loader Loader) {
		loader.Define("foo", &foo)
		loader.Define("bar", &bar)
		loader.Define("baz", &baz)
		loader.Define("qux", &qux)
		loader.Define("quux", &quux)
	})

	fooCalled := false
	barCalled := false
	bazCalled := false
	quxCalled := false
	a.Load(func(loader Loader) {
		loader.Implement("foo", func(n int) {
			if n != 42 {
				t.Fatal("foo() is not 42")
			}
			fooCalled = true
		})
		loader.Implement("bar", func(s string) {
			if s != "bar" {
				t.Fatal("bar() is not bar")
			}
			barCalled = true
		})
		loader.Implement("baz", func(b bool) {
			if !b {
				t.Fatal("baz() is not true")
			}
			bazCalled = true
		})
		loader.Implement("qux", func(f float64) {
			if f != 42.0 {
				t.Fatal("qux() is not 42.0")
			}
			quxCalled = true
		})
		loader.Implement("quux", func() {})
	})
	a.FinishLoad()

	foo(42)
	if !fooCalled {
		t.Fatal("foo not called")
	}
	bar("bar")
	if !barCalled {
		t.Fatal("bar not called")
	}
	baz(true)
	if !bazCalled {
		t.Fatal("baz not called")
	}
	qux(42.0)
	if !quxCalled {
		t.Fatal("qux no called")
	}
	quux()
}

func TestFunc2(t *testing.T) {
	a := New()
	a.Load(func(loader Loader) {
		func() {
			defer func() {
				err := recover()
				if err == nil {
					t.Fatal("should panic")
				}
				if err.(string) != "module func(app.Loader): argument is not a pointer to function" {
					panic(err)
				}
			}()
			loader.Define("foo", 42)
		}()
	})

	a.Load(func(loader Loader) {
		var f1 func(int)
		loader.Define("foo", &f1)
		f2 := func() {}
		func() {
			defer func() {
				err := recover()
				if err == nil {
					t.Fatal("should panic")
				}
				if err.(string) != "module func(app.Loader): multiple definition of foo" {
					panic(err)
				}
			}()
			loader.Define("foo", &f2)
		}()
	})

	a.Load(func(loader Loader) {
		func() {
			defer func() {
				err := recover()
				if err == nil {
					t.Fatal("should panic")
				}
				if err.(string) != "module func(app.Loader): implementation of foo is not a function" {
					panic(err)
				}
			}()
			loader.Implement("foo", 42)
		}()
	})

	a.Load(func(loader Loader) {
		loader.Implement("foo", func(int) {})
	})
	a.FinishLoad()
}

func TestFunc3(t *testing.T) {
	a := New()
	a.Load(func(loader Loader) {
		var f func(int)
		loader.Define("foo", &f)
	})
	a.Load(func(loader Loader) {
		loader.Implement("foo", func() {})
	})
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal("should panic")
			}
			if err.(string) != "defined func(int), implemented func()" {
				panic(err)
			}
		}()
		a.FinishLoad()
	}()
}

func TestFunc4(t *testing.T) {
	a := New()
	a.Load(func(loader Loader) {
		var f func(int)
		loader.Define("foo", &f)
	})
	a.Load(func(loader Loader) {
		loader.Implement("foo", func(string) {})
	})
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal("should panic")
			}
			if err.(string) != "defined func(int), implemented func(string)" {
				panic(err)
			}
		}()
		a.FinishLoad()
	}()
}

func TestFunc5(t *testing.T) {
	a := New()
	a.Load(func(loader Loader) {
		var f func(int)
		loader.Define("foo", &f)
	})
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal("should panic")
			}
			if err.(string) != "no implementation for foo" {
				panic(err)
			}
		}()
		a.FinishLoad()
	}()
}

func TestDefineProvide(t *testing.T) {
	a := New()
	a.Load(func(loader Loader) {
		var foo func(int)
		loader.DefineProvide("foo", &foo)
	})
	var foo func(int)
	called := false
	a.Load(func(loader Loader) {
		loader.Implement("foo", func(i int) {
			called = true
			if i != 42 {
				t.Fatal()
			}
		})
		loader.Require("foo", &foo)
	})
	a.FinishLoad()
	foo(42)
	if !called {
		t.Fatal()
	}
}
