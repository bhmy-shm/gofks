package plugin

import "context"

type appKey struct {
}

// WithContext returns a new Context that carries value.
func WithContext(ctx context.Context, base PluginBase) context.Context {
	return context.WithValue(ctx, appKey{}, base)
}
