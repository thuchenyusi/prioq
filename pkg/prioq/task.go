package prioq

type PauseControl struct {
	pauseCh    chan struct{}
	continueCh chan struct{}
	finishCh   chan struct{}
}

func (c *PauseControl) Check() {
	select {
	case <-c.pauseCh:
		// Wait for continue signal
		<-c.continueCh
	default:
	}
}

func (c *PauseControl) Finish() {
	close(c.finishCh)
}

type Task struct {
	Priority   int
	Pausable   bool
	IsPaused   bool
	Execute    func(*PauseControl)
	pauseCh    chan struct{}
	continueCh chan struct{}
	finishCh   chan struct{}
}

func NewTask(priority int, pausable bool, execute func(*PauseControl)) *Task {
	return &Task{
		Priority:   priority,
		Pausable:   pausable,
		Execute:    execute,
		pauseCh:    make(chan struct{}),
		continueCh: make(chan struct{}),
		finishCh:   make(chan struct{}),
	}
}

func (t *Task) Run(pq *PriorityQueue) {
	pq.wg.Add(1)
	if t.IsPaused {
		t.IsPaused = false
		t.continueCh <- struct{}{}
		return
	}
	pauseControl := &PauseControl{
		pauseCh:    t.pauseCh,
		continueCh: t.continueCh,
		finishCh:   t.finishCh,
	}
	go func() {
		defer func() {
			pq.wg.Done()
		}()
		if t.Execute != nil {
			func() {
				defer func() {
					pauseControl.Finish()
					pq.finish <- t
				}()
				t.Execute(pauseControl)
			}()
		}
	}()
}

func (t *Task) CanPause() bool {
	return t.Pausable
}

func (t *Task) Pause(pq *PriorityQueue) *Task {
	if t.Pausable {
		select {
		case t.pauseCh <- struct{}{}:
			defer pq.wg.Done()
			t.IsPaused = true
			return t
		case <-t.finishCh:
			return nil
		}
	}
	return nil
}
