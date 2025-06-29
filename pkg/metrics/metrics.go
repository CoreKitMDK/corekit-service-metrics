package metrics

import (
	"fmt"
	"github.com/CoreKitMDK/corekit-service-metrics/v2/internal/metrics"
	"os"
	"strings"
	"time"
)

type Metric struct {
	Timestamp string            `json:"timestamp"`
	Data      string            `json:"data"`
	Key       string            `json:"key"`
	Tags      map[string]string `json:"tags"`
}

func NewMetric(key string, data interface{}) Metric {
	return Metric{
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      metrics.MetricToString(data),
		Tags:      make(map[string]string),
		Key:       key,
	}
}

type IMultiMetrics interface {
	Record(key string, data interface{}) error
	Stop()
}

type IMetrics interface {
	Record(mm Metric) error
}

type MultiMetrics struct {
	metrics   []IMetrics
	bufferLen int
	tags      map[string]string
	logCh     chan Metric
	quitLogCh chan struct{}
	stopped   bool
}

func (l *MultiMetrics) Record(key string, data interface{}) error {
	m := NewMetric(key, data)
	l.log(m)
	return nil
}

func (l *MultiMetrics) Stop() {
	l.stopped = true
	close(l.quitLogCh)
}

func NewMetrics(bufferLen int, metrics ...IMetrics) *MultiMetrics {

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
		hostname = "unknown"
	}

	tags := make(map[string]string)
	tags["hostname"] = hostname

	mm := &MultiMetrics{
		metrics:   metrics,
		bufferLen: bufferLen,
		tags:      tags,
		logCh:     make(chan Metric, bufferLen),
		quitLogCh: make(chan struct{}),
		stopped:   false,
	}

	go mm.startWorker()
	return mm
}

func (l *MultiMetrics) processMetric(metric Metric) {
	var didLog = false

	for _, logger := range l.metrics {

		metric.Tags = l.tags

		if err := logger.Record(metric); err != nil {
			fallbackErrorLog(fmt.Sprintln("Error logging metric: ", err))
		} else {
			didLog = true
		}
	}

	if !didLog {
		fallbackLog(metric)
	}
}

func (l *MultiMetrics) log(m Metric) {

	select {
	case l.logCh <- m:
		break
	default:
		fallbackErrorLog("Channel overflow detected: " + m.Key)
		go func() {
			l.processMetric(m)
		}()
	}
}

func (l *MultiMetrics) startWorker() {
	defer close(l.logCh)

	for {
		select {
		case entry := <-l.logCh:
			l.processMetric(entry)
		case <-l.quitLogCh:
			for entry := range l.logCh {
				l.processMetric(entry)
			}
			return
		}
	}
}

func (l *Metric) formatTags() string {
	if len(l.Tags) == 0 {
		return ""
	}

	var builder strings.Builder

	for key, value := range l.Tags {
		builder.WriteString(key)
		builder.WriteString(":")
		builder.WriteString(value)
		builder.WriteString(",")
	}

	result := builder.String()
	if len(result) > 0 {
		result = result[:len(result)-1]
	}

	result += ";"

	return result
}

func fallbackErrorLog(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf(fmt.Sprintf("%s - [FALLBACK] : %s\n", timestamp, message))
}

func fallbackLog(m Metric) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf(fmt.Sprintf("%s - [FALLBACK] [%s] : %s = %s\n", timestamp, m.formatTags(), m.Key, m.Data))
}
