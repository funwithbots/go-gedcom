package gedcom

const (
	LevelDebug Level = iota - 1
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

type Level int

type Warnings []warning

type warning struct {
	level   Level
	warning Warning
}

type Warning interface {
	String() string
}

func (w *Warnings) AddWarning(l Level, warn Warning) {
	*w = append(*w, warning{l, warn})
}

func (w *Warnings) GetWarnings(minLevel Level) []Warning {
	list := make([]Warning, 0)
	for _, ww := range *w {
		if ww.level >= minLevel {
			list = append(list, ww.warning)
		}
	}

	return list
}
