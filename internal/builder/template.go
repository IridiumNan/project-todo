package builder

import _ "embed"

type WorkTemplateProvider struct {
	RawTemplate string
	Injector    Injector
}

//go:embed new_work.md
var newWorkMDTemplate string
