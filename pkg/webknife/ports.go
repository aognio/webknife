// Package webknife defines shared types and ports for the Webknife HTTP diagnostic laboratory.
//
// Feature packages (static, proxy, echo, respond, redirect) and middleware packages
// (auth, headers, observe) are independently reusable. They communicate through
// standard http.Handler and small port interfaces defined here.
package webknife

import "net/http"

// Middleware is a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Event represents an observable fact that happened during request processing.
// Feature packages publish events through the Publisher port.
type Event interface{}

// Publisher is the port through which feature packages publish events.
// Implementations may fan out to subscribers, log to stdout, send to
// metrics backends, or do nothing (noopPublisher).
type Publisher interface {
	Publish(event Event)
}

// NoopPublisher discards all events. Use when observation is not needed.
type NoopPublisher struct{}

func (NoopPublisher) Publish(Event) {}

// Chain applies middlewares in order (first listed = outermost wrapper).
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
