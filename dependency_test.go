package app

import "testing"

func TestDep1(t *testing.T) {
	a := New()

	// module load
	foo := new(moduleFoo)
	a.Load(foo)
	a.Load(new(moduleBar))

	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal("should panic")
			}
			if err.(string) != "int is not a module" {
				t.Fatal(err)
			}
		}()
		a.Load(42)
	}()

	// multiple provides
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal("should panic")
			}
			if err.(string) != "module *app.moduleBaz: multiple provides of bar" {
				t.Fatal(err)
			}
		}()
		a.Load(new(moduleBaz))
	}()

	// not requiring function
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal("should panic")
			}
			if err.(string) != "module *app.moduleQuux: required quux is not a pointer" {
				t.Fatal(err)
			}
		}()
		a.Load(new(moduleQuux))
	}()

	a.FinishLoad()

	if foo.bar() != 42 {
		t.Fatal("foo.bar() is not 42")
	}
}

type moduleFoo struct {
	bar func() int
}

func (m *moduleFoo) Load(loader Loader) {
	loader.Require("bar", &m.bar)
}

type moduleBar struct{}

func (m *moduleBar) Load(loader Loader) {
	loader.Provide("bar", m.bar)
}

func (m *moduleBar) bar() int {
	return 42
}

type moduleBaz struct{}

func (m *moduleBaz) Load(loader Loader) {
	loader.Provide("bar", func() int {
		return 24
	})
}

type moduleQux struct{}

func (m *moduleQux) Load(loader Loader) {
	loader.Provide("qux", 24)
}

type moduleQuux struct{}

func (m *moduleQuux) Load(loader Loader) {
	loader.Require("quux", 24)
}

func TestDep2(t *testing.T) {
	a := New()
	a.Load(new(moduleBar))
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal("should panic")
			}
			if err.(string) != "bar is not required by any module" {
				t.Fatal(err)
			}
		}()
		a.FinishLoad()
	}()
}

func TestDep3(t *testing.T) {
	a := New()
	a.Load(new(moduleBar))
	a.Load(new(moduleA))
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal("should panic")
			}
			if err.(string) != "bar not match, func() int provided, func() required" {
				t.Fatal(err)
			}
		}()
		a.FinishLoad()
	}()
}

type moduleA struct{}

func (m *moduleA) Load(loader Loader) {
	var f func()
	loader.Require("bar", &f)
}

func TestDep4(t *testing.T) {
	a := New()
	a.Load(func(loader Loader) {
		var f func()
		loader.Require("bar", &f)
	})
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal("should panic")
			}
			if err.(string) != "bar not provided" {
				t.Fatal(err)
			}
		}()
		a.FinishLoad()
	}()
}

func TestDep5(t *testing.T) {
	a := New()
	a.Load(func(loader Loader) {
		i := 42
		loader.Provide("foo", i)
	})
	var i int
	a.Load(func(loader Loader) {
		loader.Require("foo", &i)
	})
	a.FinishLoad()
	if i != 42 {
		t.Fatal("i is not 42")
	}
}
