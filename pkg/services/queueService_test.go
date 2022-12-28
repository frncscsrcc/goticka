package services

import (
	"goticka/pkg/domain/queue"
	"goticka/testUtils"
	"testing"
)

func TestCreateQueue(t *testing.T) {
	testUtils.ResetTestDependencies()

	createdQueue, err := NewQueueService().Create(
		queue.Queue{
			Name: "Queue1",
		},
	)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if createdQueue.ID == 0 {
		t.Error("queue ID should be initialized")
	}

	// TODO Call UserGet to check if the returned value is correct
}
