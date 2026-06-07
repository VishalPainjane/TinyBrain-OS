package events_test

import (
	"sync"
	"testing"
	"time"

	"github.com/VishalPainjane/TinyBrain-OS/internal/events"
)

func TestChannelBus_PublishDeliversToSubscriber(t *testing.T) {
	var bus events.EventBus = events.NewChannelBus(8)

	var wg sync.WaitGroup
	wg.Add(1)

	var received events.Event
	unsub := bus.Subscribe(events.TypeTaskCreated, func(ev events.Event) {
		received = ev
		wg.Done()
	})
	defer unsub()

	payload := events.TaskCreatedPayload{TaskID: "task-42"}
	ev := events.NewEvent(events.TypeTaskCreated, payload, time.Now())
	bus.Publish(ev)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("subscriber did not receive event in time")
	}

	if received.Type != events.TypeTaskCreated {
		t.Errorf("Type = %q, want TaskCreated", received.Type)
	}
	got, ok := received.Payload.(events.TaskCreatedPayload)
	if !ok {
		t.Fatalf("Payload type = %T, want TaskCreatedPayload", received.Payload)
	}
	if got.TaskID != "task-42" {
		t.Errorf("TaskID = %q, want task-42", got.TaskID)
	}
}

func TestChannelBus_MultipleSubscribers(t *testing.T) {
	var bus events.EventBus = events.NewChannelBus(8)

	var count int
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	handler := func(events.Event) {
		mu.Lock()
		count++
		mu.Unlock()
		wg.Done()
	}

	unsub1 := bus.Subscribe(events.TypeTaskCreated, handler)
	unsub2 := bus.Subscribe(events.TypeTaskCreated, handler)
	defer unsub1()
	defer unsub2()

	bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{TaskID: "t1"}, time.Now()))

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("subscribers did not receive event in time")
	}

	mu.Lock()
	defer mu.Unlock()
	if count != 2 {
		t.Errorf("delivery count = %d, want 2", count)
	}
}

func TestChannelBus_UnsubscribeStopsDelivery(t *testing.T) {
	var bus events.EventBus = events.NewChannelBus(8)

	var count int
	var mu sync.Mutex

	unsub := bus.Subscribe(events.TypeTaskCreated, func(events.Event) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{TaskID: "t1"}, time.Now()))
	time.Sleep(100 * time.Millisecond)

	unsub()

	bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{TaskID: "t2"}, time.Now()))
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Errorf("delivery count = %d, want 1 after unsubscribe", count)
	}
}

func TestChannelBus_DifferentEventTypesIsolated(t *testing.T) {
	var bus events.EventBus = events.NewChannelBus(8)

	var count int
	var mu sync.Mutex

	unsub := bus.Subscribe(events.TypeTaskCreated, func(ev events.Event) {
		mu.Lock()
		count++
		mu.Unlock()
	})
	defer unsub()

	bus.Publish(events.NewEvent(events.TypeTaskAssigned, events.TaskAssignedPayload{}, time.Now()))
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if count != 0 {
		mu.Unlock()
		t.Fatalf("TaskCreated subscriber count = %d, want 0 after TaskAssigned publish", count)
	}
	mu.Unlock()

	bus.Publish(events.NewEvent(events.TypeTaskCreated, events.TaskCreatedPayload{TaskID: "t1"}, time.Now()))
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Errorf("TaskCreated subscriber count = %d, want 1", count)
	}
}
