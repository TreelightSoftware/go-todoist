package todoist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProjects(t *testing.T) {
	setup()
	projects, err := GetAllProjects("")
	assert.NotNil(t, err)
	assert.Zero(t, len(projects))
	assert.Equal(t, "Empty token", err.Error())
}
