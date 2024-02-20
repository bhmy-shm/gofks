package confs

type (
	TraceConfig struct {
		Trace *trace `yaml:"trace"`
	}
	trace struct {
		Endpoint     string            `yaml:"endpoint"`
		Exporter     string            `yaml:"exporter"` //default=jaeger,options=jaeger|otlpgrpc|otlphttp|file
		OtlpHeaders  map[string]string `yaml:"headers"`
		OtlpHttpPath string            `yaml:"httpPath"`
		Name         string            `yaml:"name"`
		Namespace    string            `yaml:"namespace"`
		Version      string            `yaml:"version"`
		Sampler      float64           `yaml:"sampler"` //default=1.0
	}
)

func (c *TraceConfig) IsEnable() bool {
	if c.Trace != nil {
		return len(c.Trace.Exporter) > 0 && len(c.Trace.Endpoint) > 0
	}
	return false
}
