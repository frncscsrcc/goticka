package events

import (
	"errors"
	"testing"
)

func TestSyncLocalEvents(t *testing.T) {
	handler := Handler()

	var counter1, counter2 int
	EVENT1 := EventType(100)
	EVENT2 := EventType(200)
	EVENT3 := EventType(300)

	handler.RegisterCallBack(EVENT1, func(event LocalEvent) error {
		counter1 = counter1 + 1
		counter2 = counter2 + 1
		return nil
	})
	handler.RegisterCallBack(EVENT2, func(event LocalEvent) error {
		counter1 = counter1 + 1
		return nil
	})

	if counter1 != 0 {
		t.Errorf("unexpected value for counter1, expected %d but got %d", 0, counter1)
	}
	if counter2 != 0 {
		t.Errorf("unexpected value for counter2, expected %d but got %d", 0, counter2)
	}

	handler.SendSyncLocalEvent(LocalEvent{
		EventType: EVENT1,
	})
	if counter1 != 1 {
		t.Errorf("unexpected value for counter1, expected %d but got %d", 1, counter1)
	}
	if counter2 != 1 {
		t.Errorf("unexpected value for counter2, expected %d but got %d", 1, counter2)
	}

	handler.SendSyncLocalEvent(LocalEvent{
		EventType: EVENT2,
	})
	if counter1 != 2 {
		t.Errorf("unexpected value for counter1, expected %d but got %d", 1, counter1)
	}
	if counter2 != 1 {
		t.Errorf("unexpected value for counter2, expected %d but got %d", 1, counter2)
	}

	handler.RegisterCallBack(EVENT3, func(event LocalEvent) error {
		counter1 = counter1 + 1
		return errors.New("error")
	})
	// the next should be never called, because the previous failed
	handler.RegisterCallBack(EVENT3, func(event LocalEvent) error {
		counter1 = counter1 + 1
		return nil
	})
	if counter1 != 2 {
		t.Errorf("unexpected value for counter1, expected %d but got %d", 2, counter1)
	}

}
