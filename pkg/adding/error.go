package adding

import (
	"errors"
)

// ErrAlreadyExists is used item in repository already exists
var ErrAlreadyExists = errors.New("Already exists")

//ErrRepositoryUnknown is used when an unknown issue occurs with repository retrieval
var ErrRepositoryUnknown = errors.New("An unknown issue occurred when attempting to get from repository")

//ErrHashGenerationFailed is used when an attempt to generate a unique hash fails multiple times
var ErrHashGenerationFailed = errors.New("Attempts to generate a unique hash failed more times than the max permitted value")
