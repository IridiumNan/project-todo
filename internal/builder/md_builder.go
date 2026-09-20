package builder

import "github.com/IridiumNan/project-todo/internal/models"

type MDBuilder func([]*models.Work, string) []byte

// TODO: User the With function to wrap table create
