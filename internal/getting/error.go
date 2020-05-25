package getting

import (
	"errors"
)

// ErrNotFound is used when item in a repository could not be found
var ErrNotFound = errors.New("Not found")

//ErrRepositoryUnknown is used when an unknown issue occurs with repository retrieval
var ErrRepositoryUnknown = errors.New("An unknown issue occurred when attempting to get from repository")
