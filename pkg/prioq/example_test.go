package prioq_test

import (
	"fmt"
	"prioq/pkg/prioq"
	"time"
)

func ExamplePriorityQueue() {
	pq := prioq.NewPriorityQueue()

	task1 := prioq.NewTask(1, false, func(p *prioq.PauseControl) {
		time.Sleep(1 * time.Second)
	})
	task2 := prioq.NewTask(2, true, func(p *prioq.PauseControl) {
		for i := 0; i < 5; i++ {
			p.Check()
			time.Sleep(200 * time.Millisecond)
		}
	})
	task3 := prioq.NewTask(3, false, func(p *prioq.PauseControl) {
		time.Sleep(3 * time.Second)
	})

	pq.PushTask(task1)
	pq.PushTask(task2)

	pq.Start()
	pq.PushTask(task3)
	time.Sleep(8 * time.Second)
	pq.Stop()

	if pq.Len() != 0 {
		fmt.Printf("expected 0 tasks after scheduling, got %d\n", pq.Len())
	} else {
		fmt.Println("All tasks completed successfully")
	}
	// Output:
	// All tasks completed successfully
}
