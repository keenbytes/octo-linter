package dotgithub

// File represents both GitHub Actions action and workflow.
type File interface {
	Unmarshal(fromRaw bool) error
	GetType() int
}
