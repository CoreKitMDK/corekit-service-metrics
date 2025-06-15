package metrics

type Configuration struct {
	UseConsole   bool
	UseNATS      bool
	NatsURL      string
	NatsUsername string
	NatsPassword string
}

func NewConfiguration() *Configuration {
	return &Configuration{}
}

func (c *Configuration) Init() *MultiMetrics {
	var loggers []IMetrics

	if c.UseConsole {
		consoleLogger := NewMetricsConsole()
		loggers = append(loggers, consoleLogger)
	}

	if c.UseNATS {
		if c.NatsUsername != "" && c.NatsPassword != "" {
			if natsLogger, err := NewMetricsNATSWithAuth(c.NatsURL, c.NatsUsername, c.NatsPassword); err == nil {
				loggers = append(loggers, natsLogger)
			}
		} else {
			if natsLogger, err := NewMetricsNATS(c.NatsURL); err == nil {
				loggers = append(loggers, natsLogger)
			}
		}
	}

	multiLogger := NewMetrics(100, loggers...)
	return multiLogger
}
