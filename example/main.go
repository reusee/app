package main

import (
	"github.com/reusee/app"
)

var application = app.New()

func main() {
	var run func()
	application.Load(func(loader app.Loader) {
		loader.Define("run", &run)
	})
	application.FinishLoad()
	run()
}
