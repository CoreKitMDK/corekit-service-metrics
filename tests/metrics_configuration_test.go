package tests

import (
	"github.com/CoreKitMDK/corekit-service-metrics/v2/pkg/metrics"
	"testing"
	"time"
)

func TestMetricsConfiguration(t *testing.T) {

	config := metrics.NewConfiguration()

	config.UseConsole = true

	config.UseNATS = true
	config.NatsURL = "nats://localhost:4221"

	config.NatsPassword = "internal-metrics-broker"
	config.NatsUsername = "internal-metrics-broker"

	ogger := config.Init()
	defer ogger.Stop()

	_ = ogger.Log(metrics.NewMetric("test", 1))
	_ = ogger.Log(metrics.NewMetric("test", 2))
	_ = ogger.Log(metrics.NewMetric("test", 3))

	time.Sleep(4 * time.Second)
}
