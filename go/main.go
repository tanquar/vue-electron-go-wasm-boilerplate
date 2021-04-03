package main

import (
	"syscall/js"
  "local.packages/gossi"
)

func run(this js.Value, inputs []js.Value) interface{} {
  var program = inputs[0].String()
  var data = inputs[1].String()
  var result = gossi.Run(program, data)
  return result
}

func main() {
	js.Global().Set("run", js.FuncOf(run))

	ch := make(chan struct{}, 0)
	<-ch
}
