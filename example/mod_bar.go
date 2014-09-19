package main

import (
	"fmt"

	app ".."
)

func init() {
	application.Load(new(ModBar))
}

type ModBar struct {
}

func (m *ModBar) Load(loader app.Loader) {
	fmt.Printf("load bar\n")
}
