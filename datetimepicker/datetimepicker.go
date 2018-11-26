/*
Package datetimepicker provides gopherjs binding for the bootstrap4-datetimepicker library.
Library: https://www.npmjs.com/package/bootstrap4-datetimepicker.
Version at the time of writing: 5.2.3

Usage:
	elem := dom.GetWindow().Document().GetElementByID("dateTimePicker")
	dtp := datetimepicker.NewDatetimepicker(elem, datetimepicker.Options{})
	timeSubscription := make(chan time.Time)
	go func() {
		e := dom.GetWindow().Document().GetElementByID("chosenDate").(*dom.HTMLDivElement)
		s := dom.GetWindow().Document().GetElementByID("chosenDateValue")
		for {
			t := <-timeSubscription
			if t.IsZero() {
				e.Style().SetProperty("display", "none", "")
			} else {
				s.SetInnerHTML(t.In(time.UTC).String())
				e.Style().SetProperty("display", "block", "")
			}
		}
	}()
	dtp.Subscribe(timeSubscription)
	dtp.SetDate("24-Nov-2018 12:19 AM")
	btn := dom.GetWindow().Document().GetElementByID("clearButton").(*dom.HTMLButtonElement)
	btn.AddEventListener("click", false, func(e dom.Event) {
		dtp.ClearDate()
	})
*/
package datetimepicker

import (
	"time"

	"github.com/arunsworld/go-js-dom"
	"github.com/gopherjs/gopherjs/js"
)

// Datetimepicker represents a Datetimepicker JS entity
type Datetimepicker struct {
	dtp *js.Object
}

// Options allow configuration of Datetimepicker
type Options struct {
	Format             string
	SideBySide         bool
	DaysOfWeekDisabled []int
}

// NewDatetimepicker creates a new Datetimepicker
func NewDatetimepicker(elem dom.Element, options Options) Datetimepicker {
	if js.Global.Get("jQuery") == js.Undefined {
		js.Global.Call("alert", "jQuery is required to use the datetimepicker.")
		return Datetimepicker{}
	}
	dtp := Datetimepicker{}
	dtp.dtp = js.Global.Call("jQuery", elem)
	if dtp.dtp.Get("datetimepicker") == js.Undefined {
		js.Global.Call("alert", "bootstrap-datetimepicker.min.js library is not loaded.")
		return Datetimepicker{}
	}
	opts := defaultOptions()
	updateOptions(opts, options)
	dtp.dtp.Call("datetimepicker", opts)
	return dtp
}

// GetDate gets the date on the Datetimepicker
func (dtp *Datetimepicker) GetDate() time.Time {
	currentDate := dtp.dtp.Call("data", "DateTimePicker").Call("date")
	if currentDate == nil {
		return time.Time{}
	}
	return currentDate.Call("toDate").Interface().(time.Time)
}

// SetDate sets the date on the Datetimepicker
func (dtp *Datetimepicker) SetDate(date string) {
	dtp.dtp.Call("data", "DateTimePicker").Call("date", date)
}

// ClearDate clears out the date
func (dtp *Datetimepicker) ClearDate() {
	dtp.dtp.Call("data", "DateTimePicker").Call("date", nil)
}

// Subscribe subscribes to changes to date
func (dtp *Datetimepicker) Subscribe(subscription chan time.Time) {
	dtp.dtp.Call("on", "dp.change", func(e *js.Object) {
		newDate := e.Get("date")
		if !newDate.Bool() {
			subscription <- time.Time{}
			return
		}
		subscription <- newDate.Call("toDate").Interface().(time.Time)
	})
}

func updateOptions(opts *js.Object, options Options) {
	if options.Format != "" {
		opts.Set("format", options.Format)
	}
	opts.Set("sideBySide", options.SideBySide)
	if len(options.DaysOfWeekDisabled) > 0 {
		daysOfWeek := js.Global.Get("Array").Call("from", options.DaysOfWeekDisabled)
		opts.Set("daysOfWeekDisabled", daysOfWeek)
	}
}

func defaultOptions() *js.Object {
	options := js.M{
		"format":  "D-MMM-YYYY h:mm A",
		"minDate": false,
		"maxDate": false,
		"icons": js.M{
			"time":     "fa fa-clock-o",
			"date":     "fa fa-calendar",
			"up":       "fa fa-arrow-up",
			"down":     "fa fa-arrow-down",
			"previous": "fa fa-arrow-left",
			"next":     "fa fa-arrow-right",
			"today":    "fa fa-bullseye",
			"clear":    "fa fa-trash",
			"close":    "fa fa-times",
		}}
	return js.Global.Get("Object").New(options)
}
