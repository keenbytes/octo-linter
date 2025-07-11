package naming

import (
	"testing"

	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

var DotGithub *dotgithub.DotGithub

func TestMain(m *testing.M) {
	DotGithub = &dotgithub.DotGithub{}
	DotGithub.ReadDir("../../../../tests/rules")
	m.Run()
}
