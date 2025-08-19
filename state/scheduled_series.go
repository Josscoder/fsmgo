package state

import "time"

type ScheduledStateSeries struct {
	*Series
	ticker   *time.Ticker
	quitCh   chan struct{}
	interval time.Duration
}

func NewScheduledStateSeries(states []State, interval time.Duration) *ScheduledStateSeries {
	s := &ScheduledStateSeries{
		Series:   NewStateSeries(states),
		quitCh:   make(chan struct{}),
		interval: interval,
	}
	s.BaseState = NewBaseState(s)
	s.ticker = time.NewTicker(interval)
	return s
}

func (s *ScheduledStateSeries) OnStart() {
	s.Series.OnStart()

	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.Update()
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
