package websocket

import (
	"github.com/arunsworld/go-js-dom"
	"github.com/gopherjs/gopherjs/js"
)

// ReadyState represents the state that a WebSocket is in. For more information
// about the available states, see
// http://dev.w3.org/html5/websockets/#dom-websocket-readystate
type ReadyState uint16

func (rs ReadyState) String() string {
	switch rs {
	case Connecting:
		return "Connecting"
	case Open:
		return "Open"
	case Closing:
		return "Closing"
	case Closed:
		return "Closed"
	default:
		return "Unknown"
	}
}

const (
	// Connecting means that the connection has not yet been established.
	Connecting ReadyState = 0
	// Open means that the WebSocket connection is established and communication
	// is possible.
	Open ReadyState = 1
	// Closing means that the connection is going through the closing handshake,
	// or the Close() method has been invoked.
	Closing ReadyState = 2
	// Closed means that the connection has been closed or could not be opened.
	Closed ReadyState = 3
)

// WSURL returns the WS URL given the current document
func WSURL(doc dom.HTMLDocument, relativeLocation string) string {
	p := doc.Location().Protocol
	h := doc.Location().Host
	newP := "wss://"
	if p == "http:" {
		newP = "ws://"
	}
	return newP + h + relativeLocation
}

// New creates a new WebSocket.
func New(url string) (ws *WebSocket, err error) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if jsErr, ok := e.(*js.Error); ok && jsErr != nil {
			ws = nil
			err = jsErr
		} else {
			panic(e)
		}
	}()

	object := js.Global.Get("WebSocket").New(url)

	ws = &WebSocket{
		Object: object,
	}

	return
}

// WebSocket is a wrapper around JS WebSocket object
type WebSocket struct {
	*js.Object
	URL string `js:"url"`
	// ready state
	ReadyState     ReadyState `js:"readyState"`
	BufferedAmount uint32     `js:"bufferedAmount"`
	// networking
	Extensions string `js:"extensions"`
	Protocol   string `js:"protocol"`
	// messaging
	BinaryType string `js:"binaryType"`
}

// AddEventListener provides the ability to bind callback
// functions to the following available events:
// open, error, close, message
func (ws *WebSocket) AddEventListener(typ string, useCapture bool, listener func(*js.Object)) {
	ws.Call("addEventListener", typ, listener, useCapture)
}

// RemoveEventListener removes a previously bound callback function
func (ws *WebSocket) RemoveEventListener(typ string, useCapture bool, listener func(*js.Object)) {
	ws.Call("removeEventListener", typ, listener, useCapture)
}

// Send sends a message on the WebSocket. The data argument can be a string or a
// *js.Object fulfilling the ArrayBufferView definition.
//
// See: http://dev.w3.org/html5/websockets/#dom-websocket-send
func (ws *WebSocket) Send(data interface{}) (err error) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if jsErr, ok := e.(*js.Error); ok && jsErr != nil {
			err = jsErr
		} else {
			panic(e)
		}
	}()
	ws.Object.Call("send", data)
	return
}

// Close closes the underlying WebSocket.
//
// See: http://dev.w3.org/html5/websockets/#dom-websocket-close
func (ws *WebSocket) Close() (err error) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if jsErr, ok := e.(*js.Error); ok && jsErr != nil {
			err = jsErr
		} else {
			panic(e)
		}
	}()

	// Use close code closeNormalClosure to indicate that the purpose
	// for which the connection was established has been fulfilled.
	// See https://tools.ietf.org/html/rfc6455#section-7.4.
	ws.Object.Call("close", closeNormalClosure)
	return
}

// Close codes defined in RFC 6455, section 11.7.
const (
	// 1000 indicates a normal closure, meaning that the purpose for
	// which the connection was established has been fulfilled.
	closeNormalClosure = 1000
)
