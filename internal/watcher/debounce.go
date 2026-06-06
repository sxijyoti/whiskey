package watcher

import (
	"sync"
	"time"
)

type Debouncer struct {
	mu    sync.Mutex
	timer *time.Timer
	delay time.Duration
	fn    func()
}

func NewDebouncer(
	delay time.Duration,
) *Debouncer {
	return &Debouncer{
		delay: delay,
	}
}

func (d *Debouncer) Run(
	fn func(),
) {

	d.mu.Lock()
	defer d.mu.Unlock()

	d.fn = fn

	if d.timer == nil {

		d.timer = time.AfterFunc(
			d.delay,
			func() {
				d.fire()
			},
		)

		return
	}

	d.timer.Reset(
		d.delay,
	)
}

func (d *Debouncer) fire() {

	d.mu.Lock()
	fn := d.fn
	d.mu.Unlock()

	if fn != nil {
		fn()
	}
}