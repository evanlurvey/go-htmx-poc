package identity

import "context"

type sessionctxkey string

const sessionCtxKey sessionctxkey = "session"

func SessionWithContext(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, sessionCtxKey, session)
}

func SessionFromContext(ctx context.Context) Session {
	session, ok := ctx.Value(sessionCtxKey).(Session)
	if !ok {
		return NewAnonymousSession()
	}
	return session
}
