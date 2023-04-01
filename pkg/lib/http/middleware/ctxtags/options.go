package ctxtags

var (
	defaultOptions = &options{
		requestFieldsExtractorFunc: nil,
	}
)

type options struct {
	requestFieldsExtractorFunc RequestFieldExtractorFunc
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

// WithFieldExtractor 添加从请求注入信息到ctx tags
func WithFieldExtractor(f RequestFieldExtractorFunc) Option {
	return func(o *options) {
		o.requestFieldsExtractorFunc = f
	}
}
