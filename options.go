package bitcask

// Options for bitcask
type Options struct {
	Dir         string
	MaxFileSize int64
}

type Option func(*Options)

func WithDir(dir string) Option {
	return func(o *Options) {
		o.Dir = dir
	}
}

func WithMaxFileSize(size int64) Option {
	return func(o *Options) {
		o.MaxFileSize = size
	}
}
