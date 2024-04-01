//go:build integration
// +build integration

package agent_integration_test

import "testing"

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Agents integration")
}
