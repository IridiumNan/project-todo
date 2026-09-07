package builder

import (
	"fmt"
	"strings"
)

type Injector interface {
	Inject(rawStr string) (injectedStr string)
}

const (
	worksTableName = "{{WORK TABLE}}"
)

type MDInjector struct {
	KeyContentMap map[string]string
}

// BuildWorksTable Build the work table content
func (mi *MDInjector) BuildWorksTable(idTitleMap map[string]string) {
	emptyTable := "There is no exist works"

	if len(idTitleMap) == 0 {
		mi.KeyContentMap[worksTableName] = emptyTable
		return
	}

	tableHead := `| ID | Title |
| -------------- | --------------- |
`

	tableBody := ""
	for id, title := range idTitleMap {
		newLine := fmt.Sprintf("| %s | %s |\n", id, title)
		tableBody = tableBody + newLine
	}

	mi.KeyContentMap[worksTableName] = tableHead + strings.ReplaceAll(tableBody, "\n\n", "\n")
}

// Inject replace all key with built content
func (mi *MDInjector) Inject(rawStr string) (injectedStr string) {
	for key, content := range mi.KeyContentMap {
		rawStr = strings.Replace(rawStr, key, content, 1)
	}

	injectedStr = rawStr
	return
}
