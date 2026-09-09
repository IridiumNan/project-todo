package builder

import (
	_ "embed"
	"log/slog"
	"testing"

	testutils "github.com/IridiumNan/project-todo/internal/test_utils"
)

//go:embed inject_test/expected_work.md
var expectedWork string

func TestMDTableInject(t *testing.T) {
	mdInjector := MDInjector{
		KeyContentMap: map[string]string{},
	}

	idTitleMap := map[string]string{
		"fnosigaif":    "Hello world",
		"fsodifaisfd":  "Hello golang",
		"ifjasodfijad": "Complete the injection",
		"jfoaidfjaofd": "Haha",
	}

	mdInjector.buildWorksTable(idTitleMap)

	if got := mdInjector.Inject(newWorkMDTemplate); got != expectedWork {
		t.Error("not match expected content")
		if err := testutils.CompareWithDiff(got, expectedWork); err != nil {
			slog.Error("while calling diff command", "err", err)
		}
	}
}
