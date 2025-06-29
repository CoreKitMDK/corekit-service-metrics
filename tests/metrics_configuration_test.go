package tests

import (
	"testing"
	"time"

	"github.com/CoreKitMDK/corekit-service-metrics/v2/pkg/metrics"
)

func TestMetricsConfiguration(t *testing.T) {

	config := metrics.NewConfiguration()

	config.UseConsole = true

	config.UseNATS = true
	config.NatsURL = "internal-metrics-broker-nats-client"

	config.NatsPassword = "internal-metrics-broker"
	config.NatsUsername = "internal-metrics-broker"

	ogger := config.Init()
	defer ogger.Stop()

	_ = ogger.Record("test", 1)
	_ = ogger.Record("test", 2)
	_ = ogger.Record("test", 3)

	time.Sleep(2 * time.Second)
}
