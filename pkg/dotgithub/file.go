package dotgithub

type File interface {
	Unmarshal(logLevel int, fromRaw bool) error
	GetType() int
}
