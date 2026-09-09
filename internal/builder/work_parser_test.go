package builder

import (
	"os"
	"slices"
	"testing"

	"github.com/IridiumNan/project-todo/internal/models"
)

func isSameTomlConf(conf1 *mdTomlConfig, conf2 *mdTomlConfig) bool {
	if conf1.Title != conf2.Title {
		return false
	}

	if conf1.Energy != conf2.Energy {
		return false
	}

	if conf1.Viewer != conf2.Viewer {
		return false
	}

	if slices.Compare(conf1.DependenciesID, conf2.DependenciesID) != 0 {
		return false
	}

	return true
}

func TestParseFunc(t *testing.T) {
	expectedTomlConf := mdTomlConfig{
		Title:          "Test Work Build",
		Energy:         models.EnergyMedium,
		Viewer:         models.ViewerEditor,
		DependenciesID: []string{"abc", "cdf"},
	}

	mdParser := MDWorkParser{}

	data, err := os.ReadFile("new_work.md")
	if err != nil {
		t.Errorf("error when load the new_work.md, err: %s", err)
	}

	mdTomlConf, mdContext, err := mdParser.Parse(data)
	if err != nil {
		t.Errorf("encouter error when parse the markdown file, err: %s", err)
	}

	if !isSameTomlConf(mdTomlConf, &expectedTomlConf) {
		t.Errorf("toml config parse mismatch, want %v, got %v", expectedTomlConf, mdTomlConf)
	}

	expectedContext := []byte("**Write your context here**")

	if string(expectedContext) != string(mdContext) {
		t.Errorf("md context mismatch, want %s, got %s", expectedContext, mdContext)
	}
}
