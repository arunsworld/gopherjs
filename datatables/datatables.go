/*
Package datatables provides gopherjs binding for the datatables library: https://datatables.net/.

Usage:
	elem := dom.GetWindow().Document().QuerySelector("#myTable")
	options := datatables.Options{
		LengthMenu: []datatables.MenuLengthDefn{
			datatables.MenuLengthDefn{Count: 5, Display: "5"},
			datatables.MenuLengthDefn{Count: 10, Display: "10"},
			datatables.MenuLengthDefn{Count: -1, Display: "All"},
		},
	}
	rf1 := func(data *js.Object, renderType string) interface{} {
		if renderType == "display" {
			if data != nil {
				return "$ " + data.String()
			}
		}
		return data
	}
	columns := []datatables.Column{
		datatables.Column{ID: "a", Title: "A", Visible: true},
		datatables.Column{ID: "b", Title: "B", Visible: true, Render: rf1},
		datatables.Column{ID: "c", Title: "C", Visible: true},
		datatables.Column{ID: "d", Title: "D", Visible: true},
	}
	data := `[
		{"a": "Arun", "b": 1234, "c": null, "d": 0},
		{"a": "Barua", "b": 99, "c": null, "d": 0},
		{"a": null, "b": 0, "c": "Arun", "d": 1234},
		{"a": null, "b": 0, "c": "Barua", "d": 9999}
	]`
	dt := datatables.NewDatatable(elem, options, columns, js.Global.Get("JSON").Call("parse", "[]"))
	clicked := make(chan *js.Object)
	go func() {
		for {
			r := <-clicked
			console.Log(r)
		}
	}()
	dt.RowClicked(clicked)
	go func() {
		time.Sleep(time.Second)
		dt.Update(js.Global.Get("JSON").Call("parse", data))
	}()
*/
package datatables

import (
	"github.com/arunsworld/go-js-dom"
	"github.com/gopherjs/gopherjs/js"
)

// Column is a type that defines the columns for the datatable
type Column struct {
	ID             string
	Title          string
	Visible        bool
	CellType       string
	Render         Renderer
	DefaultContent string
}

// Renderer is a type that provides renderers for column identified by the column's id
type Renderer func(data *js.Object, renderType string) interface{}

// MenuLengthDefn defines a menu length
type MenuLengthDefn struct {
	Count   int
	Display string
}

// Options define options on the datatable
type Options struct {
	LengthMenu []MenuLengthDefn
}

// Datatable wrapper around the datatable object
type Datatable struct {
	dt *js.Object
}

// NewDatatable creates a new datatable
func NewDatatable(elem dom.Element, options Options, columns []Column, data *js.Object) Datatable {
	result := Datatable{}
	result.dt = js.Global.Call("jQuery", elem)
	opts := js.Global.Get("Object").New(js.M{"columns": columnOptions(columns), "data": data,
		"columnDefs": columnDefnOptions(columns)})
	updateOptions(opts, options)
	result.dt.Call("dataTable", opts)
	return result
}

// NewDatatableWithArray creates a new datatable with a flat 2D array
func NewDatatableWithArray(elem dom.Element, options Options, columns []Column, data *js.Object) Datatable {
	result := Datatable{}
	result.dt = js.Global.Call("jQuery", elem)
	opts := js.Global.Get("Object").New(js.M{"columns": columnOptionsJustTitle(columns), "data": data,
		"columnDefs": columnDefnOptions(columns)})
	updateOptions(opts, options)
	result.dt.Call("dataTable", opts)
	return result
}

// Destroy removes the existing table in prep for a refresh
func (dt *Datatable) Destroy() {
	api := dt.dt.Call("api")
	api.Call("destroy")
}

// Update updates the datatable with the passed in data
func (dt *Datatable) Update(data *js.Object) {
	api := dt.dt.Call("api")
	api.Call("clear")
	api.Get("rows").Call("add", data)
	api.Call("draw")
}

// RowClicked is the callback when a row is clicked
func (dt *Datatable) RowClicked(clicked chan *js.Object) {
	tbody := dt.dt.Call("children", "tbody")
	tbody.Call("on", "click", "tr", func(e *js.Object) {
		row := e.Get("currentTarget")
		dataTableRow := dt.dt.Call("api").Call("row", row).Call("data")
		clicked <- dataTableRow
	})
}

func updateOptions(opts *js.Object, options Options) {
	if len(options.LengthMenu) > 0 {
		lenMenuCount := make([]interface{}, len(options.LengthMenu))
		lenMenuDisplay := make([]interface{}, len(options.LengthMenu))
		for i, o := range options.LengthMenu {
			lenMenuCount[i] = o.Count
			lenMenuDisplay[i] = o.Display
		}
		opts.Set("lengthMenu", js.S{lenMenuCount, lenMenuDisplay})
	}
}

func columnOptions(columns []Column) *js.Object {
	result := js.S{}
	for _, c := range columns {
		result = append(result, js.M{"data": c.ID, "name": c.ID, "title": c.Title})
	}
	return js.Global.Get("Object").New(result)
}

func columnOptionsJustTitle(columns []Column) *js.Object {
	result := js.S{}
	for _, c := range columns {
		result = append(result, js.M{"title": c.Title})
	}
	return js.Global.Get("Object").New(result)
}

func columnDefnOptions(columns []Column) *js.Object {
	result := js.S{}
	for i, c := range columns {
		d := js.Global.Get("Object").New()
		d.Set("targets", i)
		d.Set("visible", c.Visible)
		if c.CellType != "" {
			d.Set("cellType", c.CellType)
		}
		if c.DefaultContent != "" {
			d.Set("defaultContent", c.DefaultContent)
		}
		if c.Render != nil {
			d.Set("render", c.Render)
		}
		result = append(result, d)
	}
	return js.Global.Get("Object").New(result)
}
