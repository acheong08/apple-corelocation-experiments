//go:build js && wasm

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"wloc/lib"
)

func main() {
	done := make(chan struct{}, 0)
	wlocFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			fmt.Println("No arguments provided")
			return nil
		}
		bssids := make([]string, len(args))
		for i, arg := range args {
			if arg.Type() != js.TypeString {
				fmt.Println("Invalid type provided ", arg.Type())
				return nil
			}
			bssids[i] = arg.String()
		}
		fmt.Println(bssids)
		handler := js.FuncOf(func(this js.Value, args []js.Value) any {
			resolve := args[0]

			go func() {
				devices, err := lib.QueryBssid(bssids, true)
				if err != nil {
					resolve.Invoke(fmt.Sprintf("an error occured: %v", err))
					return
				}
				b, _ := json.Marshal(devices.GetWifiDevices())
				fmt.Println(string(b))
				resolve.Invoke(string(b))
			}()
			return nil
		})
		return js.Global().Get("Promise").New(handler)
	})
	js.Global().Set("wloc", wlocFunc)
	js.Global().Set("wlocDone", js.FuncOf(func(this js.Value, args []js.Value) any {
		done <- struct{}{}
		return nil
	}))
	<-done
	wlocFunc.Release()
}
