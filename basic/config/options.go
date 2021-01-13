package config

import "github.com/micro/go-micro/v2/config/source"

/*
Functional option is a pattern for flexible API argument management design.
Being different from traditional approaches, where a function reads args
and set itself according to some pre-defined logic, the new pattern passes
functions, the Option(s), to the API (public function)
with user-defined arguments (implemented with closures),
and the function only has one line of code to cope with that -
running all the Option(s) and let THEM configure the API, contrarily.
*/
type Options struct {
	App     map[string]interface{}
	AppName string
	Sources []source.Source
}

type Option func(o *Options)

func WithSource(src source.Source) Option {
	return func(ops *Options) {
		ops.Sources = append(ops.Sources, src)
	}
}

func WithApp(appName string) Option {
	return func(o *Options) {
		o.AppName = appName
	}
}
