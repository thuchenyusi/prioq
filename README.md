# Priority Queue

This project implements a priority queue in Go, which allows tasks to be scheduled and executed based on their priority. Tasks can be paused and resumed, and the queue ensures that higher priority tasks are executed before lower priority ones.

## Installation

To install the project, use the following command:

```sh
go get -u github.com/yourusername/prioq
```

## Usage

Here is an example of how to use the priority queue:

```go
package main

import (
	"fmt"
	"prioq/pkg/prioq"
	"time"
)

func main() {
	pq := prioq.NewPriorityQueue()

	task1 := prioq.NewTask(1, false, func(p *prioq.PauseControl) {
		fmt.Println("Task 1 started")
		time.Sleep(1 * time.Second)
		fmt.Println("Task 1 completed")
	})
	task2 := prioq.NewTask(2, true, func(p *prioq.PauseControl) {
		fmt.Println("Task 2 started")
		for i := 0; i < 5; i++ {
			p.Check()
			time.Sleep(200 * time.Millisecond)
		}
		fmt.Println("Task 2 completed")
	})
	task3 := prioq.NewTask(3, false, func(p *prioq.PauseControl) {
		fmt.Println("Task 3 started")
		time.Sleep(3 * time.Second)
		fmt.Println("Task 3 completed")
	})

	pq.PushTask(task1)
	pq.PushTask(task2)

	pq.Start()
	pq.PushTask(task3)
	time.Sleep(50 * time.Second)
	pq.Stop()

	if pq.Len() != 0 {
		fmt.Printf("expected 0 tasks after scheduling, got %d\n", pq.Len())
	} else {
		fmt.Println("All tasks completed successfully")
	}
}
```

## Testing

To run the tests, use the following command:

```sh
go test ./...
```

## License

This project is licensed under the MIT License.
