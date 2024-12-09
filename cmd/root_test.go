package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd(t *testing.T) {
	rootCmd := rootCmd
	statusCmd := statusCmd

	rootCmd.AddCommand(statusCmd)

	err := statusCmd.Execute()
	assert.NoError(t, err, "Cmd.Execute should not return an error")

	t.Log("hi", tasks)
}
