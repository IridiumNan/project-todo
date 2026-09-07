package builder

import (
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"testing"
)

//go:embed inject_test/expected_work.md
var expectedWork string

func compareWithDiff(got string, want string) (err error) {
	clearFile := func(f *os.File) {
		f.Close()

		err := os.Remove(f.Name())
		if err != nil {
			slog.Error("while removing test file", "err", err)
		}
	}

	file1, err := os.CreateTemp("/tmp/", "go-inject-test-*")
	if err != nil {
		return fmt.Errorf("error when create a tmp file, err: %s, tmp file path: %s", err.Error(), file1.Name())
	}
	defer clearFile(file1)

	file2, err := os.CreateTemp("/tmp/", "go-inject-test-*")
	if err != nil {
		return fmt.Errorf("error when create a tmp file, err: %s, tmp file path: %s", err.Error(), file2.Name())
	}

	defer clearFile(file2)

	_, err = file1.WriteString(got)
	if err != nil {
		return fmt.Errorf("error when write string into file, fileName: %s, err: %w", file1.Name(), err)
	}

	_, err = file2.WriteString(want)
	if err != nil {
		return fmt.Errorf("error when write string into file, fileName: %s, err: %w", file2.Name(), err)
	}

	slog.Info("show file name", "file1", file1.Name(), "file2", file2.Name())
	cmd := exec.Command("diff", file1.Name(), file2.Name())
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error when exec diff command, err : %w", err)
	}

	return nil
}

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

	mdInjector.BuildWorksTable(idTitleMap)

	if got := mdInjector.Inject(newWorkMDTemplate); got != expectedWork {
		t.Error("not match expected content")
		if err := compareWithDiff(got, expectedWork); err != nil {
			slog.Error("while calling diff command", "err", err)
		}
	}
}
