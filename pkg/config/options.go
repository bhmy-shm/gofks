package config

type Options struct {
	Encoding map[string]Encoder
}

type Option func(o *Options)

func newOptions(opts ...Option) Options {
	options := Options{
		Encoding: map[string]Encoder{
			"json": NewJsonEncoder(),
			"yaml": NewYamlEncoder(),
		},
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
