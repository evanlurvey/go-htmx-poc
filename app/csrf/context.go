package csrf

import (
	"context"
)

type csrftokenkey string

const csrfTokenCtxKey csrftokenkey = "csrfSession"

func WithContext(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, csrfTokenCtxKey, token)
}

func FromContext(ctx context.Context) string {
	sid, ok := ctx.Value(csrfTokenCtxKey).(string)
	if !ok {
		panic("missing csrf session")
	}
	return sid
}
