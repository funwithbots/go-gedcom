package gedcom7

type DocOptions func(opts docOptions) docOptions

type docOptions struct {
	// Allow tags if they were deprecated after this version.
	// e.g. if set to '5.5.1', then tags deprecated on or after v5.5.1 will be allowed.
	maxDeprecatedTagVersion string
	// base path for gedcom documents
	docPath string
}

func (d *docOptions) withOpts(opts []DocOptions) {
	for _, opt := range opts {
		*d = opt(*d)
	}
}

func WithDocPath(path string) DocOptions {
	if path == "" {
		path = "./"
	}
	if path[len(path)-1] != '/' {
		path += "/"
	}
	return func(opts docOptions) docOptions {
		opts.docPath = path
		return opts
	}
}

func WithMaxDeprecatedTags(ver string) DocOptions {
	return func(opts docOptions) docOptions {
		opts.maxDeprecatedTagVersion = ver
		return opts
	}
}
