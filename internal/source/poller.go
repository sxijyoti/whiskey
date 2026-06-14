package source

import "time"

type Poller struct {
	Interval time.Duration
}

func (p Poller) Start(
	fn func(),
) {

	ticker := time.NewTicker(
		p.Interval,
	)

	go func() {

		for range ticker.C {
			fn()
		}
	}()
}