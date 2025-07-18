package slf4go_logrus_provider

import (
	"github.com/MariusSchmidt/slf4go/slf4go_api"
	"github.com/sirupsen/logrus"
)

func ExampleNew() {
	// default instantiation
	logger := New(logrus.StandardLogger())
	logger.Infof("Server starting...")
}

func ExampleSlf4GoLogrusLogger_ForComponent() {
	// with custom Component-Label
	logger := New(logrus.StandardLogger()).ForComponent("service")
	logger.Infof("Server starting...")
}

func ExampleSlf4GoLogrusLogger_WithAppComponentLabel() {
	// with custom Component-Label
	logger := New(logrus.StandardLogger()).WithAppComponentLabel("appLabel")
	logger.Infof("Server starting...")
}

func ExampleSlf4GoLogrusLogger_WithStaticTags() {
	logger := New(logrus.StandardLogger()).WithStaticTags(slf4go_api.LogTags{
		"requestID": "abc-123",
		"userID":    "user-456",
	})

	logger.Infof("Processing request...")
}
