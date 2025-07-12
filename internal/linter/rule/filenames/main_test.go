package filenames

import (
	"testing"

	"github.com/keenbytes/octo-linter/v2/pkg/dotgithub"
)

var DotGithub *dotgithub.DotGithub

func TestMain(m *testing.M) {
	DotGithub = &dotgithub.DotGithub{}
	_ = DotGithub.ReadDir("../../../../tests/rules")

	m.Run()
}
