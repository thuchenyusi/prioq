package prioq

import (
	"container/heap"
	"sync"
)

type PriorityQueue struct {
	tasks   []*Task
	mu      sync.Mutex
	running map[*Task]struct{}
	notify  chan struct{}
	finish  chan *Task
	wg      *sync.WaitGroup
	stop    chan struct{}
	active  bool
}

func (pq *PriorityQueue) Len() int { return len(pq.tasks) }

func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.tasks[i].Priority > pq.tasks[j].Priority
}

func (pq *PriorityQueue) Swap(i, j int) {
	pq.tasks[i], pq.tasks[j] = pq.tasks[j], pq.tasks[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	pq.tasks = append(pq.tasks, x.(*Task))
	select {
	case pq.notify <- struct{}{}:
	default:
	}
}

func (pq *PriorityQueue) Pop() interface{} {
	old := pq.tasks
	n := len(old)
	item := old[n-1]
	pq.tasks = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) PushTask(task *Task) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	heap.Push(pq, task)
}

func (pq *PriorityQueue) PopTask() *Task {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	return heap.Pop(pq).(*Task)
}

func (pq *PriorityQueue) PeekTask() *Task {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if len(pq.tasks) > 0 {
		return pq.tasks[0]
	}
	return nil
}

func (pq *PriorityQueue) Schedule() {
	const maxRunningTasks = 1
	for {
		select {
		case t := <-pq.finish:
			delete(pq.running, t)
			pq.notify <- struct{}{}
		case <-pq.notify:
			task := pq.PeekTask()
			if task == nil {
				continue
			}
			if len(pq.running) < maxRunningTasks {
				pq.running[task] = struct{}{}
				pq.PopTask().Run(pq)
			} else {
				for t := range pq.running {
					if t.CanPause() && t.Priority < task.Priority {
						oldTask := t.Pause(pq)
						if oldTask == nil {
							continue
						}
						delete(pq.running, t)
						pq.Push(oldTask)
						pq.running = map[*Task]struct{}{task: {}}
						pq.PopTask().Run(pq)
						break
					}
				}
			}
		case <-pq.stop:
			for t := range pq.running {
				if t.CanPause() {
					oldTask := t.Pause(pq)
					if oldTask == nil {
						continue
					}
					delete(pq.running, t)
					pq.Push(oldTask)
				}
			}
			pq.wg.Wait()
			return
		}
	}
}

func (pq *PriorityQueue) Start() {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if pq.active {
		return
	}
	pq.active = true
	go pq.Schedule()
}

func (pq *PriorityQueue) Stop() {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	if !pq.active {
		return
	}
	pq.active = false
	close(pq.stop)
}

func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{
		tasks:   []*Task{},
		running: make(map[*Task]struct{}),
		notify:  make(chan struct{}, 1),
		finish:  make(chan *Task),
		wg:      &sync.WaitGroup{},
		stop:    make(chan struct{}),
	}
	heap.Init(pq)
	return pq
}
