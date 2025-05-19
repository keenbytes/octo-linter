package dotgithub

type File interface {
	Unmarshal(fromRaw bool) error
	GetType() int
}
