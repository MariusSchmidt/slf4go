package test_mocks

import (
	"github.com/MariusSchmidt/slf4go/slf4go_api"
)

// This file contains compile-time assertions to ensure that the mock implementations of
// slf4go_api.Slf4GoLogger stay in sync with the interface. If any method is added, removed
// or changed in the interface, the compiler will throw an error if the mock implementation
// is not updated accordingly. This helps catch interface compatibility issues early during
// development.
var _ slf4go_api.Slf4GoLogger = &MockSlf4GoLogger{}
