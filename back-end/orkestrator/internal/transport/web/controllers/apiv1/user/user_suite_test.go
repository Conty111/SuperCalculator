package user_test

import (
	"testing"

	. "github.com/onsi/ginkgo" //nolint:revive
	. "github.com/onsi/gomega" //nolint:revive
)

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Suite")
}
