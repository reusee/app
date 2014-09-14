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
	emits    map[string]interface{}
	listens  map[string][]interface{}
}

func New() *Application {
	app := &Application{
		provides: make(map[string]interface{}),
		requires: make(map[string][]interface{}),
		emits:    make(map[string]interface{}),
		listens:  make(map[string][]interface{}),
	}
	return app
}

type Loader struct {
	Provide func(name string, fn interface{})
	Require func(name string, fn interface{})
	Emit    func(name string, fn interface{})
	Listen  func(name string, fn interface{})
}

func (a *Application) Load(module interface{}) {
	modType := reflect.TypeOf(module)
	loader := Loader{
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
		Emit: func(name string, fn interface{}) {
			t := reflect.TypeOf(fn)
			if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Func {
				panic(sp("module %v: emitter %s is not a pointer to function", modType, name))
			}
			if _, in := a.emits[name]; in {
				panic(sp("module %v: multiple emitter %s", modType, name))
			}
			a.emits[name] = fn
		},
		Listen: func(name string, fn interface{}) {
			t := reflect.TypeOf(fn)
			if t.Kind() != reflect.Func {
				panic(sp("module %v: listener %s is not a function", modType, name))
			}
			a.listens[name] = append(a.listens[name], fn)
		},
	}
	if mod, ok := module.(interface {
		Load(Loader)
	}); ok {
		mod.Load(loader)
	} else if mod, ok := module.(func(Loader)); ok {
		mod(loader)
	} else {
		panic(sp("%v is not a module", modType))
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
		delete(a.requires, name)
	}
	for name, _ := range a.requires {
		panic(sp("%s not provided", name))
	}

	// handle signals
	for name, emit := range a.emits {
		emitType := reflect.TypeOf(emit).Elem()
		listens := a.listens[name]
		if len(listens) == 0 {
			panic(sp("%s not listened", name))
		}
		for _, listen := range listens {
			listenType := reflect.TypeOf(listen)
			if emitType.NumOut() != listenType.NumIn() {
				panic(sp("%s not match, emit %v, listen %v", name, emitType, listenType))
			}
			for i := 0; i < emitType.NumOut(); i++ {
				if emitType.Out(i) != listenType.In(i) {
					panic(sp("%s not match at arg #%d, emit %v, listen %v", name, i, emitType.Out(i), listenType.In(i)))
				}
			}
		}
		var wrapped reflect.Value
		switch emit := emit.(type) {
		// fast paths
		case *func() int:
			e := *emit
			wrapped = reflect.ValueOf(func() (ret int) {
				ret = e()
				for _, listen := range listens {
					listen.(func(int))(ret)
				}
				return
			})
		case *func() string:
			e := *emit
			wrapped = reflect.ValueOf(func() (ret string) {
				ret = e()
				for _, listen := range listens {
					listen.(func(string))(ret)
				}
				return
			})
		// generic with reflection
		default:
			listenValues := make([]reflect.Value, 0, len(listens))
			for _, listen := range listens {
				listenValues = append(listenValues, reflect.ValueOf(listen))
			}
			emitValue := reflect.ValueOf(reflect.ValueOf(emit).Elem().Interface())
			wrapped = reflect.MakeFunc(emitType, func(args []reflect.Value) (out []reflect.Value) {
				out = emitValue.Call(args)
				for _, listen := range listenValues {
					listen.Call(out)
				}
				return
			})
		}
		reflect.ValueOf(emit).Elem().Set(wrapped)
	}
}
