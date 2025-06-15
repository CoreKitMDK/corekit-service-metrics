package tests

import (
	"github.com/CoreKitMDK/corekit-service-metrics/v2/pkg/metrics"
	"testing"
	"time"
)

func TestMetric(t *testing.T) {

	mnats, err := metrics.NewMetricsNATSWithAuth("127.0.0.1:4221", "internal-metrics-broker", "internal-metrics-broker")

	if err != nil {
		t.Error(err)
	}

	nm := metrics.NewMetricsConsole()

	m := metrics.NewMetrics(10, mnats, nm)
	_ = m.Log(metrics.NewMetric("test", 123))

	time.Sleep(4 * time.Second)

	m.Stop()
}
