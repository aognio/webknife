package events_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aognio/webknife/internal/events"
)

func TestEventBus_PublishSubscribe(t *testing.T) {
	bus := events.NewEventBus()
	var received events.Event
	bus.Subscribe(func(e events.Event) {
		received = e
	})
	event := events.RequestCompleted{
		Method:     "GET",
		Path:       "/test",
		StatusCode: 200,
		Timestamp:  time.Now(),
	}
	bus.Publish(event)
	if received == nil {
		t.Fatal("expected to receive event")
	}
	rc, ok := received.(events.RequestCompleted)
	if !ok {
		t.Fatal("expected RequestCompleted event")
	}
	if rc.Method != "GET" {
		t.Fatalf("expected GET, got %s", rc.Method)
	}
	if rc.Path != "/test" {
		t.Fatalf("expected /test, got %s", rc.Path)
	}
}

func TestEventBus_MultipleSubscribers(t *testing.T) {
	bus := events.NewEventBus()
	count := 0
	bus.Subscribe(func(e events.Event) {
		count++
	})
	bus.Subscribe(func(e events.Event) {
		count++
	})
	bus.Publish(events.ServerStarted{Addr: ":8080"})
	if count != 2 {
		t.Fatalf("expected 2 subscribers called, got %d", count)
	}
}

func TestRequestMiddleware_PublishesEvents(t *testing.T) {
	bus := events.NewEventBus()
	var published []events.Event
	bus.Subscribe(func(e events.Event) {
		published = append(published, e)
	})
	pub := &testPublisher{bus: bus}
	handler := events.RequestMiddleware(pub)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if len(published) != 2 {
		t.Fatalf("expected 2 events (received + completed), got %d", len(published))
	}
	if _, ok := published[0].(events.RequestReceived); !ok {
		t.Fatal("first event should be RequestReceived")
	}
	if _, ok := published[1].(events.RequestCompleted); !ok {
		t.Fatal("second event should be RequestCompleted")
	}
	rc := published[1].(events.RequestCompleted)
	if rc.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rc.StatusCode)
	}
}

type testPublisher struct {
	bus *events.EventBus
}

func (p *testPublisher) Publish(event events.Event) {
	p.bus.Publish(event)
}
