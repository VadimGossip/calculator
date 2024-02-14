package workerctrl

type service struct {
	lockChan chan struct{}
}

func NewService(n int) *service {
	return &service{lockChan: make(chan struct{}, n)}
}

func (s *service) Acquire(n int) {
	for i := 0; i < n; i++ {
		s.lockChan <- struct{}{}
	}
}

func (s *service) Release(n int) {
	for i := 0; i < n; i++ {
		<-s.lockChan
	}
}
