package internal

import "github.com/rotisserie/eris"

// InvalidArgumentError Sentinel error for invalid argument.
//nolint:gochecknoglobals // Intends to be a sentinel value
var InvalidArgumentError error = eris.New("Invalid argument")

// OutOfRangeArgumentError Sentinel error for out of range argument.
//nolint:gochecknoglobals // Intends to be a sentinel value
var OutOfRangeArgumentError error = eris.New("Out of range argument")

// AlreadyExistsError Sentinel error for element already existing.
//nolint:gochecknoglobals // Intends to be a sentinel value
var AlreadyExistsError error = eris.New("Element already exists")

// NotExistsError Sentinel error for non existent argument.
//nolint:gochecknoglobals // Intends to be a sentinel value
var NotExistsError error = eris.New("Element does not exist")

// InternalError Sentinel error for all purpose argument.
//nolint:gochecknoglobals // Intends to be a sentinel value
var InternalError error = eris.New("Internal error")

// ConfigurationError Sentinel error for configuration argument.
//nolint:gochecknoglobals // Intends to be a sentinel value
var ConfigurationError error = eris.New("Error handling configuration configuration")
