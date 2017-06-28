package workers

type WorkingGroup struct {
	main    chan func()
	allDone chan bool
}

func New(nWorkers int) WorkingGroup {
	wg := WorkingGroup{
		main:    make(chan func()),
		allDone: make(chan bool),
	}

	procDone := make(chan bool)

	for i := 0; i < nWorkers; i++ {
		go func() {
			for f := range wg.main {
				f()
			}
			procDone <- true
		}()
	}

	go func() {
		for i := 0; i < nWorkers; i++ {
			<-procDone
		}
		wg.allDone <- true
	}()

	return wg
}

func (wg WorkingGroup) Add(f func()) {
	wg.main <- f
}

func (wg WorkingGroup) Wait() {
	close(wg.main)
	<-wg.allDone
}
