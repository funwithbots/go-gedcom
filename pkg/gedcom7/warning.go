package gedcom7

type warning struct {
	Node    *Line
	Line    string
	Message string
}

func (d *document) Warnings() []string {
	warnings := make([]string, len(d.warnings))
	for i, w := range d.warnings {
		warnings[i] = w.String()
	}

	return warnings
}

func (w *warning) String() string {
	return w.Line + " " + w.Message
}
