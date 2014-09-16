package app

func init() {
	AddSignalType((*func() int)(nil), sigHandlerInt)
	AddSignalType((*func() string)(nil), sigHandlerString)
}

type _TypeRet interface{}

/* gogen: [
	{ "name": "sigHandlerInt",
		"_TypeRet": "int"
	},
	{ "name": "sigHandlerString",
		"_TypeRet": "string"
	}
]
*/
func sigTypeTemplate(emit interface{}, listens []interface{}) {
	emitPtr := emit.(*func() _TypeRet)
	e := *emitPtr
	*emitPtr = func() (ret _TypeRet) {
		ret = e()
		for _, listen := range listens {
			listen.(func(_TypeRet))(ret)
		}
		return
	}
}

// sigHandlerInt is generated from sigTypeTemplate by gogen
func sigHandlerInt(emit interface{}, listens []interface{}) {
	emitPtr := emit.(*func() int)
	e := *emitPtr
	*emitPtr = func() (ret int) {
		ret = e()
		for _, listen := range listens {
			listen.(func(int))(ret)
		}
		return
	}
}

// sigHandlerString is generated from sigTypeTemplate by gogen
func sigHandlerString(emit interface{}, listens []interface{}) {
	emitPtr := emit.(*func() string)
	e := *emitPtr
	*emitPtr = func() (ret string) {
		ret = e()
		for _, listen := range listens {
			listen.(func(string))(ret)
		}
		return
	}
}
