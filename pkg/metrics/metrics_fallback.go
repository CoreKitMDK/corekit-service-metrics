package metrics

import (
	"fmt"
	"time"
)

type Fallback struct {
}

func NewMetricsFallback() *Fallback {
	return &Fallback{}
}

func (lf *Fallback) Record(key string, data interface{}) error {
	fmt.Printf("%s - [%s] : %s\n", time.Now().Format("2006-01-02 15:04:05"), key, data)
	return nil
}
