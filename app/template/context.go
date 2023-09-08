package template

import (
	"context"
	"maps"
)

type templatectxkey string

const templateCtxKey templatectxkey = "csrfSession"

type M map[string]any

// add shit into the global template context that will get added into render map.
func WithContext(ctx context.Context, m M) context.Context {
	o := FromContext(ctx)
	maps.Copy(o, m)
	return context.WithValue(ctx, templateCtxKey, o)
}

func FromContext(ctx context.Context) M {
	m, ok := ctx.Value(templateCtxKey).(M)
	if !ok {
		return M{}
	}
	return m
}
