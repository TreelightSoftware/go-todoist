package todoist

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskCRUD(t *testing.T) {
	setup()
	existingToken := config.AuthToken
	config.AuthToken = ""
	tasks, err := GetActiveTasks("")
	assert.NotNil(t, err)
	assert.Zero(t, len(tasks))
	assert.Equal(t, "Empty token", err.Error())
	config.AuthToken = existingToken
	if config.AuthToken == "" {
		// everything below will fail, so exit out
		fmt.Println("Skipping the rest of the Task tests")
		return
	}
	// create a test task, find it, get it, update it, close it, reopen it, delete it
	tasks, err = GetActiveTasks("")
	fmt.Printf("\n%+v\n%+v\n", tasks, err)
}
