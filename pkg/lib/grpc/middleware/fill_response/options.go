package fill_response

var (
	defaultOptions = &options{
		fillResponseFunc: nil,
	}
)

type options struct {
	fillResponseFunc FillResponseFunc
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

func WithFillResponseFunc(fn FillResponseFunc) Option {
	return func(o *options) {
		o.fillResponseFunc = fn
	}
}
