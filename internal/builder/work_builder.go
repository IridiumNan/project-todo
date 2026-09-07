package builder

import (
	_ "embed"
)

//go:embed new_work.md
var newWorkMDTemplate string

type WorkBuilder struct{}
