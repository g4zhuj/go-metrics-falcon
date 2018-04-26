package falconmetrics

import (
	"fmt"
	"testing"
	"time"

	metrics "github.com/rcrowley/go-metrics"
)

func TestRegistry(t *testing.T) {
	c := metrics.NewCounter()
	metrics.Register("counter_name", c)
	c.Inc(47)

	g := metrics.NewGauge()
	err := metrics.Register("gauge_name", g)
	if err != nil {
		fmt.Printf("Register Gauge %v\n", err)
		return
	}
	g.Update(47)

	s := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	metrics.Register("histogram_name", h)
	h.Update(47)
	for i := 0; i < 100; i++ {
		h.Update(int64(i))
	}

	m := metrics.NewMeter()
	metrics.Register("meter_name", m)
	m.Mark(47)
	for i := 0; i < 10; i++ {
		m.Mark(int64(1))
		time.Sleep(time.Millisecond * 1)
	}

	timer := metrics.NewTimer()
	metrics.Register("timer_name", t)
	timer.Time(func() { time.Sleep(time.Millisecond * 10) })
	timer.Update(47)

	//report to open falcon
	cfg := DefaultFalconConfig
	cfg.Debug = true
	//report for each 5 seconds
	cfg.Step = 5
	falcon := NewFalcon(&cfg)
	falcon.ReportRegistry(metrics.DefaultRegistry)
}
