/*
Package dropzone provides gopherjs binding for the Dropzone library: https://www.dropzonejs.com/.
Curent version of dropzone at the time of writing: 5.2.0

Usage:
	dropzone.AutoDiscover(false)
	el := dom.GetWindow().Document().QuerySelector(".dropzone")
	dz := dropzone.NewDropzone(el, dropzone.Props{
		URL: "http://localhost:3000/debug/",
		OnError: func(f dropzone.File, s int, m string) {
			js.Global.Call("alert", "Upload failed.")
		},
		OnSuccess: func(f dropzone.File, r string) {
			msg := f.Name + " successfully uploaded."
			js.Global.Get("console").Call("log", msg)
		},
		AcceptedFiles: ".jpg, .jpeg",
		ImageResize: [...]float64{300, 300},
	})
	btn := dom.GetWindow().Document().QuerySelector("#clearButton")
	btn.AddEventListener("click", false, func(e dom.Event) {
		dz.RemoveFiles()
	})
*/
package dropzone

import (
	"github.com/arunsworld/go-js-dom"
	"github.com/gopherjs/gopherjs/js"
)

var dz = js.Global.Get("Dropzone")

// AutoDiscover sets whether we should auto discover dropzone components
func AutoDiscover(flag bool) {
	dz.Set("autoDiscover", flag)
}

// Dropzone wraps the Dropzone js object
type Dropzone struct {
	*js.Object
}

// File wraps the Dropzone file object
// Per https://github.com/gopherjs/gopherjs/wiki/JavaScript-Tips-and-Gotchas it should always be
// referenced as a pointer and never directly.
type File struct {
	*js.Object
	Name string `js:"name"`
	Size int    `js:"size"`
	Type string `js:"type"`
}

// FormData wraps the FormData object
type FormData struct {
	*js.Object
}

// Append appends parameters and values to the FormData object
func (f *FormData) Append(param string, value interface{}) {
	f.Call("append", param, value)
}

// RemoveFiles removes files already loaded
func (d *Dropzone) RemoveFiles() {
	d.Call("removeAllFiles")
}

// Props are properties to create a new dropzone object
type Props struct {
	URL           string
	Params        map[string]string
	Headers       map[string]string
	ImageResize   [2]float64
	AcceptedFiles string
	OnSuccess     func(file *File, response string)
	OnError       func(file *File, status int, errorMessage string)
	OnSend        func(file *File, formData *FormData)
}

// NewDropzone creates a new dropzone object
func NewDropzone(elem dom.Element, props Props) *Dropzone {
	p := js.Global.Get("Object").New()
	p.Set("url", props.URL)
	p.Set("params", props.Params)
	p.Set("headers", props.Headers)
	if props.ImageResize[0] > 0 && props.ImageResize[1] > 0 {
		p.Set("resizeWidth", props.ImageResize[0])
		p.Set("resizeHeight", props.ImageResize[1])
	}
	p.Set("acceptedFiles", props.AcceptedFiles)
	d := dz.New(elem, p)
	if props.OnError != nil {
		d.Call("on", "error", func(file *File, errorMessage string, xhr *js.Object) {
			props.OnError(file, xhr.Get("status").Int(), errorMessage)
		})
	}
	if props.OnSuccess != nil {
		d.Call("on", "success", func(file *File, resp string) {
			props.OnSuccess(file, resp)
		})
	}
	if props.OnSend != nil {
		d.Call("on", "sending", func(file *File, xhr *js.Object, formData *FormData) {
			props.OnSend(file, formData)
		})
	}
	return &Dropzone{Object: d}
}
