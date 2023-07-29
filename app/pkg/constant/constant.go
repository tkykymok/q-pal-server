package constant

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var ContextExecutorKey = contextKey("ContextExecutor")
