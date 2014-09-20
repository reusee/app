package main

import (
	"fmt"

	app ".."
)

func init() {
	application.Load(new(ModFoo))
}

type ModFoo struct {
}

func (m *ModFoo) Load(loader app.Loader) {
	fmt.Printf("load foo\n")
	loader.Implement("run", func() {
		fmt.Printf("run foo\n")
	})
}
