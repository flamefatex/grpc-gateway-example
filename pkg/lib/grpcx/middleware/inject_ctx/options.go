package inject_ctx

var (
	defaultOptions = &options{
		injectFromOrigCtxFuncs: nil,
	}
)

type options struct {
	injectFromOrigCtxFuncs []InjectCtxFunc
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type Option func(*options)

func WithInjectCtxFunc(fns ...InjectCtxFunc) Option {
	return func(o *options) {
		for _, fn := range fns {
			o.injectFromOrigCtxFuncs = append(o.injectFromOrigCtxFuncs, fn)
		}
	}
}
