package builder

import "context"

type ArchContext struct {
	context.Context
	Stat *BuildStat
}

type TaskContext struct {
	ArchContext
	GOOS, GOARCH    string
	NameSuffix, Env []string
}
