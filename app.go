package app

import (
	"fmt"
	"reflect"
)

var (
	sp = fmt.Sprintf
)

type Application struct {
	provides map[string]interface{}
	requires map[string][]interface{}
	defs     map[string]interface{}
	impls    map[string][]interface{}
	runs     []func()
	clears   []func()
}

func New() *Application {
	app := &Application{
		provides: make(map[string]interface{}),
		requires: make(map[string][]interface{}),
		defs:     make(map[string]interface{}),
		impls:    make(map[string][]interface{}),
	}
	return app
}

type Loader struct {
	Provide   func(name string, fn interface{})
	Require   func(name string, fn interface{})
	Define    func(name string, fn interface{})
	Implement func(name string, fn interface{})
}

func (a *Application) Load(module interface{}) {
	modType := reflect.TypeOf(module)
	loader := Loader{
		Provide: func(name string, fn interface{}) {
			if _, in := a.provides[name]; in {
				panic(sp("module %v: multiple provides of %s", modType, name))
			}
			a.provides[name] = fn
		},
		Require: func(name string, fn interface{}) {
			t := reflect.TypeOf(fn)
			if t.Kind() != reflect.Ptr {
				panic(sp("module %v: required %s is not a pointer", modType, name))
			}
			a.requires[name] = append(a.requires[name], fn)
		},
		Define: func(name string, fnPtr interface{}) {
			t := reflect.TypeOf(fnPtr)
			if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Func {
				panic(sp("module %v: argument is not a pointer to function", modType))
			}
			if _, in := a.defs[name]; in {
				panic(sp("module %v: multiple definition of %s", modType, name))
			}
			a.defs[name] = fnPtr
		},
		Implement: func(name string, fn interface{}) {
			t := reflect.TypeOf(fn)
			if t.Kind() != reflect.Func {
				panic(sp("module %v: implementation of %s is not a function", modType, name))
			}
			a.impls[name] = append(a.impls[name], fn)
		},
	}
	if mod, ok := module.(interface {
		Load(Loader)
	}); ok {
		mod.Load(loader)
	}
	if mod, ok := module.(interface {
		Run()
	}); ok {
		a.runs = append(a.runs, mod.Run)
	}
	if mod, ok := module.(func(Loader)); ok {
		mod(loader)
	}
}

var fnHandlers = make(map[reflect.Type]func(interface{}, []interface{}))

func AddFuncType(fnNilPtr interface{}, handler func(impls []interface{}) interface{}) {
	fnHandlers[reflect.TypeOf(fnNilPtr).Elem()] = func(fnPtr interface{}, impls []interface{}) {
		fn := handler(impls)
		reflect.ValueOf(fnPtr).Elem().Set(reflect.ValueOf(fn))
	}
}

func (a *Application) Run() {
	// match provides and requires
	for name, provide := range a.provides {
		requires, ok := a.requires[name]
		if !ok {
			panic(sp("%s is not required by any module", name))
		}
		provideValue := reflect.ValueOf(provide)
		for _, require := range requires {
			requireValue := reflect.ValueOf(require).Elem()
			if provideValue.Type() != requireValue.Type() {
				panic(sp("%s not match, %v provided, %v required", name, provideValue.Type(), requireValue.Type()))
			}
			reflect.ValueOf(require).Elem().Set(provideValue)
		}
		delete(a.requires, name)
	}
	for name, _ := range a.requires {
		panic(sp("%s not provided", name))
	}

	// handle defs / impls
	for name, fnPtr := range a.defs {
		impls := a.impls[name]
		if len(impls) == 0 {
			panic(sp("no implementation for %s", name))
		}
		fnType := reflect.TypeOf(fnPtr).Elem()
		for _, impl := range impls {
			if t := reflect.TypeOf(impl); t != fnType {
				panic(sp("defined %v, implemented %v", fnType, t))
			}
		}
		handler, ok := fnHandlers[fnType]
		if ok {
			handler(fnPtr, impls)
		} else {
			implValues := make([]reflect.Value, 0, len(impls))
			for _, impl := range impls {
				implValues = append(implValues, reflect.ValueOf(impl))
			}
			reflect.ValueOf(fnPtr).Elem().Set(reflect.MakeFunc(fnType,
				func(args []reflect.Value) (ret []reflect.Value) {
					for _, impl := range implValues {
						impl.Call(args)
					}
					return
				}))
		}
	}

	// run
	for _, fn := range a.runs {
		fn()
	}
	// clear
	for _, fn := range a.clears {
		fn()
	}
}

func (a *Application) OnRun(fn func()) {
	a.runs = append(a.runs, fn)
}

func (a *Application) OnClear(fn func()) {
	a.clears = append(a.clears, fn)
}
