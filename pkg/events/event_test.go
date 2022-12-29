package events

import (
	"errors"
	"testing"
	"time"
)

func TestSyncLocalEvents(t *testing.T) {
	handler := Handler()

	var counter1, counter2 int
	EVENT1 := EventType(100)
	EVENT2 := EventType(200)
	EVENT3 := EventType(300)

	handler.RegisterSyncCallBack(EVENT1, func(event LocalEvent) error {
		counter1 = counter1 + 1
		counter2 = counter2 + 1
		return nil
	})
	handler.RegisterSyncCallBack(EVENT2, func(event LocalEvent) error {
		counter1 = counter1 + 1
		return nil
	})

	if counter1 != 0 {
		t.Errorf("unexpected value for counter1, expected %d but got %d", 0, counter1)
	}
	if counter2 != 0 {
		t.Errorf("unexpected value for counter2, expected %d but got %d", 0, counter2)
	}

	handler.SendLocalEvent(LocalEvent{
		EventType: EVENT1,
	})
	if counter1 != 1 {
		t.Errorf("unexpected value for counter1, expected %d but got %d", 1, counter1)
	}
	if counter2 != 1 {
		t.Errorf("unexpected value for counter2, expected %d but got %d", 1, counter2)
	}

	handler.SendLocalEvent(LocalEvent{
		EventType: EVENT2,
	})
	if counter1 != 2 {
		t.Errorf("unexpected value for counter1, expected %d but got %d", 1, counter1)
	}
	if counter2 != 1 {
		t.Errorf("unexpected value for counter2, expected %d but got %d", 1, counter2)
	}

	handler.RegisterSyncCallBack(EVENT3, func(event LocalEvent) error {
		counter1 = counter1 + 1
		return errors.New("error")
	})
	// the next should be never called, because the previous failed
	handler.RegisterSyncCallBack(EVENT3, func(event LocalEvent) error {
		counter1 = counter1 + 1
		return nil
	})
	if counter1 != 2 {
		t.Errorf("unexpected value for counter1, expected %d but got %d", 2, counter1)
	}
}

func TestAsyncLocalEvents(t *testing.T) {
	handler := Handler()

	var counterFast, counterSlow int
	EVENT1 := EventType(100)
	EVENT2 := EventType(200)

	handler.RegisterAsyncCallBack(EVENT1, func(event LocalEvent) {
		counterFast = counterFast + 1
	})
	handler.RegisterAsyncCallBack(EVENT2, func(event LocalEvent) {
		time.Sleep(time.Second)
		counterSlow = counterSlow + 1
	})

	if counterFast != 0 {
		t.Errorf("unexpected value for counterFast, expected %d but got %d", 0, counterFast)
	}
	if counterSlow != 0 {
		t.Errorf("unexpected value for counterSlow, expected %d but got %d", 0, counterSlow)
	}

	handler.SendLocalEvent(LocalEvent{
		EventType: EVENT1,
	})
	handler.SendLocalEvent(LocalEvent{
		EventType: EVENT2,
	})

	time.Sleep(500 * time.Millisecond)

	if counterFast != 1 {
		t.Errorf("unexpected value for counterFast, expected %d but got %d", 1, counterFast)
	}
	if counterSlow != 0 {
		t.Errorf("unexpected value for counterSlow, expected %d but got %d", 0, counterSlow)
	}

	time.Sleep(1000 * time.Millisecond)

	if counterFast != 1 {
		t.Errorf("unexpected value for counterFast, expected %d but got %d", 1, counterFast)
	}
	if counterSlow != 1 {
		t.Errorf("unexpected value for counterSlow, expected %d but got %d", 1, counterSlow)
	}
}
