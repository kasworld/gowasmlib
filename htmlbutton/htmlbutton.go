// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package htmlbutton handle html button(group)
package htmlbutton

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/kasworld/gowasmlib/jslog"
)

type HTMLButton struct {
	KeyCode    string
	IDBase     string
	ButtonText []string // state count
	ToolTip    string

	ClickFn func(obj interface{}, v *HTMLButton)
	State   int
}

func (v HTMLButton) JSID() string {
	return "jsid_" + v.IDBase
}
func (v HTMLButton) JSFnName() string {
	return "jsfn_" + v.IDBase
}
func (v *HTMLButton) MakeJSFn(obj interface{}) func(this js.Value, args []js.Value) interface{} {
	return func(this js.Value, args []js.Value) interface{} {
		v.State++
		v.State %= len(v.ButtonText)
		btnStr := v.ButtonText[v.State]
		js.Global().Get("document").Call("getElementById", v.JSID()).Set("innerHTML", btnStr)
		if v.ClickFn != nil {
			v.ClickFn(obj, v)
		}
		return nil
	}
}
func (v HTMLButton) MakeHTML() string {
	btnStr := v.ButtonText[v.State]
	return fmt.Sprintf(
		`<button class="button" id="%v" onclick="%v()">%s</button> `,
		v.JSID(), v.JSFnName(), btnStr,
	)
}

func (v *HTMLButton) Enable() {
	btn := js.Global().Get("document").Call("getElementById", v.JSID())
	btn.Call("removeAttribute", "disabled")
}

func (v *HTMLButton) Disable() {
	btn := js.Global().Get("document").Call("getElementById", v.JSID())
	btn.Set("disabled", true)
}

type HTMLButtonGroup []*HTMLButton

func (hbl HTMLButtonGroup) GetByIDBase(idb string) *HTMLButton {
	for _, v := range hbl {
		if v.IDBase == idb {
			return v
		}
	}
	jslog.Errorf("not found %v in %v", idb, hbl)
	return nil
}

func (hbl HTMLButtonGroup) MakeHTML(obj interface{}) string {
	var buf bytes.Buffer
	for _, v := range hbl {
		js.Global().Set(v.JSFnName(), js.FuncOf(v.MakeJSFn(obj)))
		buf.WriteString(v.MakeHTML())
	}
	return buf.String()
}
