package web

type Middleware func(Handler) Handler

func wrapMiddleware(mw []Middleware, handler Handler) Handler {

	// looping backwards to ensure middleware is applied in the correct order
	// the last middleware will be applied first
	// ensuring the first middleware is executed first by request
	for i := len(mw) - 1; i >= 0; i-- {
		mwFunc := mw[i]
		if mwFunc != nil {
			handler = mwFunc(handler)
		}
	}

	return handler
}
