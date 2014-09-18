package app

func init() {
	AddSignalType((*func() int)(nil), func(emit interface{}, listens []interface{}) interface{} {
		return func() (ret int) {
			ret = emit.(func() int)()
			for _, listen := range listens {
				listen.(func(int))(ret)
			}
			return
		}
	})
	AddSignalType((*func() string)(nil), func(emit interface{}, listens []interface{}) interface{} {
		return func() (ret string) {
			ret = emit.(func() string)()
			for _, listen := range listens {
				listen.(func(string))(ret)
			}
			return
		}
	})
	AddSignalType((*func() bool)(nil), func(emit interface{}, listens []interface{}) interface{} {
		return func() (ret bool) {
			ret = emit.(func() bool)()
			for _, listen := range listens {
				listen.(func(bool))(ret)
			}
			return
		}
	})
}
