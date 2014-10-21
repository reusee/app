package app

import (
	"fmt"
	"reflect"
)

var (
	sp = fmt.Sprintf
)

type defineInfo struct {
	ptr      interface{}
	provides bool
}

type Application struct {
	provides map[string]interface{}
	requires map[string][]interface{}
	defs     map[string]defineInfo
	impls    map[string][]interface{}
}

func New() *Application {
	app := &Application{
		provides: make(map[string]interface{}),
		requires: make(map[string][]interface{}),
		defs:     make(map[string]defineInfo),
		impls:    make(map[string][]interface{}),
	}
	return app
}

type Loader struct {
	Provide       func(name string, fn interface{})
	Require       func(name string, fn interface{})
	Define        func(name string, fnPtr interface{})
	Implement     func(name string, fn interface{})
	DefineProvide func(name string, fnPtr interface{})
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
			a.defs[name] = defineInfo{fnPtr, false}
		},
		Implement: func(name string, fn interface{}) {
			t := reflect.TypeOf(fn)
			if t.Kind() != reflect.Func {
				panic(sp("module %v: implementation of %s is not a function", modType, name))
			}
			a.impls[name] = append(a.impls[name], fn)
		},
		DefineProvide: func(name string, fnPtr interface{}) {
			t := reflect.TypeOf(fnPtr)
			if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Func {
				panic(sp("module %v: argument is not a pointer to function", modType))
			}
			if _, in := a.defs[name]; in {
				panic(sp("module %v: multiple definition of %s", modType, name))
			}
			a.defs[name] = defineInfo{fnPtr, true}
		},
	}
	if mod, ok := module.(interface {
		Load(Loader)
	}); ok {
		mod.Load(loader)
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

func (a *Application) FinishLoad() {
	// handle defs / impls
	for name, info := range a.defs {
		fnPtr := info.ptr
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
		if info.provides {
			a.provides[name] = reflect.ValueOf(fnPtr).Elem().Interface()
		}
	}

	// match provides and requires
	for name, provide := range a.provides {
		requires := a.requires[name]
		provideValue := reflect.ValueOf(provide)
		for _, require := range requires {
			requireValue := reflect.ValueOf(require).Elem()
			if provideValue.Type() != requireValue.Type() {
				panic(sp("%s not match, %v provided, %v required", name, provideValue.Type(), requireValue.Type()))
			}
			requireValue.Set(provideValue)
		}
		delete(a.requires, name)
	}
	for name, _ := range a.requires {
		panic(sp("%s not provided", name))
	}

}
