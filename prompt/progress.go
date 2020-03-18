package prompt

import (
	"fmt"
	"io"
	"os"
	"time"
)

type ProgressBar struct {
	w              io.Writer
	description    string
	cur            uint64
	max            uint64
	barSegments    int
	tickInterval   time.Duration
	close          chan struct{}
	update         chan struct{}
	running        chan struct{}
	additionalInfo []func() string
}

type ProgressOption func(*ProgressBar)

func ProgressOptionDescription(description string) ProgressOption {
	return func(bar *ProgressBar) {
		bar.description = description
	}
}

func ProgressOptionBarSegments(segments int) ProgressOption {
	return func(bar *ProgressBar) {
		bar.barSegments = segments
	}
}

func ProgressOptionTickInterval(interval time.Duration) ProgressOption {
	return func(bar *ProgressBar) {
		bar.tickInterval = interval
	}
}

func ProgressOptionAppendAdditionalInfo(f func() string) ProgressOption {
	return func(bar *ProgressBar) {
		bar.additionalInfo = append(bar.additionalInfo, f)
	}
}

func NewProgressBar(opts ...ProgressOption) *ProgressBar {
	bar := &ProgressBar{
		max:         100,
		barSegments: 30,
		close:       make(chan struct{}),
		update:      make(chan struct{}),
		running:     make(chan struct{}),
	}

	for _, o := range opts {
		o(bar)
	}

	return bar
}

func (p *ProgressBar) percentage() float64 {
	if p.max == 0 {
		return 0
	}

	return float64(p.cur) / float64(p.max)
}

func (p *ProgressBar) bar() string {
	segments := int(float64(p.barSegments) * p.percentage())

	if p.percentage() > 0 && segments == 0 {
		segments = 1
	}

	bar := make([]byte, p.barSegments+2)
	bar[0] = '['
	bar[p.barSegments+1] = ']'
	b := bar[1 : len(bar)-1]
	for i := range b {
		switch {
		case i == segments-1:
			b[i] = '>'
		case i < segments:
			b[i] = '='
		default:
			b[i] = ' '
		}
	}

	return string(bar)
}

func (ProgressBar) clearLine() {
	fmt.Fprint(os.Stderr, "\r")
}

func (p *ProgressBar) print() {
	d := p.description
	if len(d) > 0 {
		d = "  " + d
	}
	fmt.Fprintf(os.Stderr, "%s  %.1f %% %s", d, p.percentage()*100, p.bar())
	for _, f := range p.additionalInfo {
		fmt.Fprintf(os.Stderr, "  %s", f())
	}
}

func (p *ProgressBar) run() {
	defer close(p.running)
	for {
		p.clearLine()
		p.print()
		var tick <-chan time.Time
		if p.tickInterval > 0 {
			tick = time.After(p.tickInterval)
		}
		select {
		case <-p.close:
			return
		case <-tick:
		case <-p.update:
		}
	}
}

func (p *ProgressBar) Update(cur uint64) {
	p.Set(cur, p.max)
}

func (p *ProgressBar) Set(cur, max uint64) {
	p.cur, p.max = cur, max
	select {
	case p.update <- struct{}{}:
	default:
	}
}

func (p *ProgressBar) Finish() {
	p.Set(p.max, p.max)
	p.Stop()
	<-p.running
	fmt.Fprintln(os.Stderr)
}

func (p *ProgressBar) Start() {
	go p.run()
}

func (p *ProgressBar) Stop() {
	close(p.close)
}
