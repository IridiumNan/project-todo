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

// MDInjector the injector for markdown file
// It replace key to value on markdown template file
// To add new key-value pair, use [MDInjector.Push]
//
// Specially, init it with [NewMDInjector] with idTitleMap
type MDInjector struct {
	keyContentMap map[string]string
}

// buildWorksTable Build the work table content
// The idTitleMap comes from the [WorkBuilder.DataDir]
func (mi *MDInjector) buildWorksTable(idTitleMap map[string]string) {
	if idTitleMap == nil {
		return
	}
	emptyTable := "There is no exist works"

	if len(idTitleMap) == 0 {
		mi.keyContentMap[worksTableName] = emptyTable
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

	mi.keyContentMap[worksTableName] = tableHead + strings.ReplaceAll(tableBody, "\n\n", "\n")
}

// Inject replace all key with built content
// Just for count 1
func (mi *MDInjector) Inject(rawStr string) (injectedStr string) {
	for key, content := range mi.keyContentMap {
		rawStr = strings.Replace(rawStr, key, content, 1)
	}

	injectedStr = rawStr
	return
}

// Push add new key-value pair to the injector
// When [MDInjector.Inject] function is called, the key string will be replaced by value string
func (mi *MDInjector) Push(key string, value string) {
	mi.keyContentMap[key] = value
}

// NewMDInjector create a new markdown injector with idTitleMap
// This idTitleMap will be used to build a id title table which is specially injected into new built markdown context configuration file
// If the idTitleMap == nil, it will not build it
func NewMDInjector(idTitleMap map[string]string) *MDInjector {
	mi := MDInjector{keyContentMap: map[string]string{}}

	mi.buildWorksTable(idTitleMap)

	return &mi
}
