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
}

func New() *Application {
	app := &Application{
		provides: make(map[string]interface{}),
		requires: make(map[string][]interface{}),
	}
	return app
}

type Loader struct {
	Provide func(name string, fn interface{})
	Require func(name string, fn interface{})
}

func (a *Application) Load(module interface{}) {
	modType := reflect.TypeOf(module)
	if mod, ok := module.(interface {
		Load(Loader)
	}); ok {
		mod.Load(Loader{
			Provide: func(name string, fn interface{}) {
				if _, in := a.provides[name]; in {
					panic(sp("module %v: multiple provides of %s", modType, name))
				}
				t := reflect.TypeOf(fn)
				if t.Kind() != reflect.Func {
					panic(sp("module %v: provided %s is not a function", modType, name))
				}
				a.provides[name] = fn
			},
			Require: func(name string, fn interface{}) {
				t := reflect.TypeOf(fn)
				if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Func {
					panic(sp("module %v: required %s is not a pointer to function", modType, name))
				}
				a.requires[name] = append(a.requires[name], fn)
			},
		})
	}
}

func (a *Application) FinishLoad() {
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
	}
}

//TODO emit and listen
