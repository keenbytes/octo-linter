package rule

import (
	"github.com/keenbytes/octo-linter/pkg/dotgithub"
)

const (
	DotGithubFileTypeAction   = 1
	DotGithubFileTypeWorkflow = 2
)

type Rule interface {
	Validate(conf interface{}) error
	Lint(config interface{}, f dotgithub.File, d *dotgithub.DotGithub, chErrors chan<- string) (bool, error)
	ConfigName() string
	FileType() int
}
