package app

import "testing"

func TestLoad1(t *testing.T) {
	a := New()

	// module load
	foo := new(moduleFoo)
	a.Load(foo)
	a.Load(new(moduleBar))

	// multiple provides
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal()
			}
			if err.(string) != "module *app.moduleBaz: multiple provides of bar" {
				t.Fatal()
			}
		}()
		a.Load(new(moduleBaz))
	}()

	// not providing function
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal()
			}
			if err.(string) != "module *app.moduleQux: provided qux is not a function" {
				t.Fatal()
			}
		}()
		a.Load(new(moduleQux))
	}()

	// not requiring function
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal()
			}
			if err.(string) != "module *app.moduleQuux: required quux is not a pointer to function" {
				t.Fatal()
			}
		}()
		a.Load(new(moduleQuux))
	}()

	a.FinishLoad()

	if foo.bar() != 42 {
		t.Fatal()
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

func TestLoad2(t *testing.T) {
	a := New()
	a.Load(new(moduleBar))
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal()
			}
			if err.(string) != "bar is not required by any module" {
				t.Fatal()
			}
		}()
		a.FinishLoad()
	}()
}

func TestLoad3(t *testing.T) {
	a := New()
	a.Load(new(moduleBar))
	a.Load(new(moduleA))
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatal()
			}
			if err.(string) != "bar not match, func() int provided, func() required" {
				t.Fatal()
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
