package select2

import (
	"errors"

	"github.com/arunsworld/go-js-dom"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
)

// Select2 wraps the select2 library JS object
type Select2 struct {
	*js.Object
	placeholder string
}

// Options to configure Select2
type Options struct {
	Placeholder string
	Data        *js.Object
}

// Selection is a selection from select2
type Selection struct {
	*js.Object
	ID   string `js:"id"`
	Text string `js:"text"`
}

// NewSelect2 creates a new Select2 object
func NewSelect2(elem dom.Element, options Options) *Select2 {
	if js.Global.Get("jQuery") == js.Undefined {
		js.Global.Call("alert", "jQuery is required to use select2.")
		return &Select2{}
	}
	s := &Select2{Object: js.Global.Call("jQuery", elem)}
	if s.Get("select2") == js.Undefined {
		js.Global.Call("alert", "select2.min.js library is not loaded.")
		return &Select2{}
	}
	if (options.Data != nil) && (options.Data.Get("constructor") != js.Global.Get("Array")) {
		js.Global.Call("alert", "Bad configuration for Select2. Data is not an array.")
		return &Select2{}
	}
	s.Call("select2", js.M{
		"placeholder": options.Placeholder,
		"theme":       "bootstrap4",
		"data":        options.Data,
	})
	s.placeholder = options.Placeholder
	return s
}

// GetSelection returns the values that have been selected
func (s *Select2) GetSelection() []*Selection {
	selection := s.Call("select2", "data")
	result := make([]*Selection, selection.Length())
	for i := 0; i < selection.Length(); i++ {
		result[i] = &Selection{Object: selection.Index(i)}
	}
	return result
}

// SetSelection sets select2 selection
func (s *Select2) SetSelection(selections []string) {
	s.Call("val", selections).Call("trigger", "change")
}

// ResetSelection removes the select2 selection
func (s *Select2) ResetSelection() {
	s.SetSelection([]string{})
}

// ModifyOptions changes the options available in select2
func (s *Select2) ModifyOptions(options *js.Object) error {
	if options.Get("constructor") != js.Global.Get("Array") {
		console.Log(options, " is not an Array. Cannot ModifyOptions with it.")
		return errors.New("call to ModifyOptions was without an Array")
	}
	s.Call("empty")
	s.Call("append", js.Global.Get("Option").New())
	s.Call("select2", js.M{
		"placeholder": s.placeholder,
		"theme":       "bootstrap4",
		"data":        options,
	})
	return nil
}

// Subscribe subscribes to new selections
func (s *Select2) Subscribe(subscription chan *Selection) {
	s.Call("on", "select2:select", func(e *js.Object) {
		selection := &Selection{Object: e.Get("params").Get("data")}
		subscription <- selection
	})
}
