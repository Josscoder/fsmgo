package state

import "time"

type ScheduledStateSeries struct {
	*Series
	ticker *time.Ticker
	quitCh chan struct{}
}

func NewScheduledStateSeries(states []State) *ScheduledStateSeries {
	s := &ScheduledStateSeries{
		Series: NewStateSeries(states),
		quitCh: make(chan struct{}),
	}
	s.BaseState = NewBaseState(s)
	s.ticker = time.NewTicker(time.Second)
	return s
}

func (s *ScheduledStateSeries) OnStart() {
	s.Series.OnStart()

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.OnUpdate()
			case <-s.quitCh:
				s.ticker.Stop()
				return
			}
		}
	}()
}

func (s *ScheduledStateSeries) OnEnd() {
	close(s.quitCh)
	s.Series.OnEnd()
}
