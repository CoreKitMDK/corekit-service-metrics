package metrics

import "fmt"

func MetricToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int, int64, float64, float32:
		return fmt.Sprintf("%v", val)
	default:
		return "internal conversion error"
	}
}
