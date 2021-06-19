package todoist

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskCRUD(t *testing.T) {
	setup()
	tokenToUse := "" // if you didn't want to pass it to the env for a one off test, set this (but don't commit it!)
	if tokenToUse == "" {
		tokenToUse = config.AuthToken
	}
	existingToken := config.AuthToken
	config.AuthToken = ""
	tasks, err := GetActiveTasks("")
	assert.NotNil(t, err)
	assert.Zero(t, len(tasks))
	assert.Equal(t, "Empty token", err.Error())
	config.AuthToken = existingToken
	if tokenToUse == "" {
		// everything below will fail, so exit out
		fmt.Println("Skipping the rest of the Task tests")
		return
	}

	// start with a project we can use to catch all of these
	project, err := CreateTestProject(tokenToUse)
	assert.Nil(t, err)
	require.NotNil(t, project)
	defer DeleteProject(tokenToUse, project.ID)
	r := rand.Int63n(999999)

	// create a test task, find it, get it, update it, close it, reopen it, delete it
	// first, a bad create without content
	bad, err := CreateTask(tokenToUse, &TaskParams{})
	assert.NotNil(t, err)
	assert.Nil(t, bad)
	bad, err = CreateTask(tokenToUse, nil)
	assert.NotNil(t, err)
	assert.Nil(t, bad)

	createInput := &TaskParams{
		Content:     String(fmt.Sprintf("My New Task %d", r)),
		Description: String("Created from a unit test"),
		ProjectID:   Int64(project.ID),
		Priority:    PriorityUrgent,
	}
	created, err := CreateTask(tokenToUse, createInput)
	assert.Nil(t, err)
	require.NotNil(t, created)
	assert.NotZero(t, created.ID)
	assert.Equal(t, StringValue(createInput.Description), created.Description)
	assert.Equal(t, StringValue(createInput.Content), created.Content)
	assert.Equal(t, project.ID, created.ProjectID)
	assert.Equal(t, PriorityUrgent, created.Priority)
	defer DeleteTask(tokenToUse, created.ID)

	foundInSlice := false
	tasks, err = GetActiveTasks(tokenToUse)
	assert.Nil(t, err)
	assert.NotZero(t, len(tasks))
	for i := range tasks {
		if tasks[i].ID == created.ID {
			foundInSlice = true
		}
	}
	assert.True(t, foundInSlice)

	found, err := GetActiveTask(tokenToUse, -1)
	assert.NotNil(t, err)
	require.Nil(t, found)
	found, err = GetActiveTask(tokenToUse, created.ID)
	assert.Nil(t, err)
	require.NotNil(t, found)
	assert.Equal(t, StringValue(createInput.Content), found.Content)
	assert.Equal(t, StringValue(createInput.Description), found.Description)
	assert.Equal(t, Int64Value(createInput.ProjectID), found.ProjectID)
	assert.Equal(t, PriorityUrgent, found.Priority)

	// change some of the data, get it to make sure
	_, err = UpdateTask(tokenToUse, -1, &TaskParams{})
	assert.NotNil(t, err)
	_, err = UpdateTask(tokenToUse, -1, nil)
	assert.NotNil(t, err)

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02T15:04:05Z")
	updateData := &TaskParams{
		Content:     String(fmt.Sprintf("Updated tasks %d", r)),
		DueDatetime: String(tomorrow),
	}
	updated, err := UpdateTask(tokenToUse, created.ID, updateData)
	assert.Nil(t, err)
	assert.Equal(t, created.ID, updated.ID)

	found, err = GetActiveTask(tokenToUse, created.ID)
	assert.Nil(t, err)
	require.NotNil(t, found)
	assert.Equal(t, fmt.Sprintf("Updated tasks %d", r), found.Content)
	assert.Equal(t, StringValue(createInput.Description), found.Description)
	assert.Equal(t, Int64Value(createInput.ProjectID), found.ProjectID)
	assert.Equal(t, PriorityUrgent, found.Priority)

	// close it
	err = CloseTask(tokenToUse, -1)
	assert.NotNil(t, err)
	err = CloseTask(tokenToUse, created.ID)
	assert.Nil(t, err)
	found, err = GetActiveTask(tokenToUse, created.ID)
	assert.NotNil(t, err)
	assert.Nil(t, found)

	// reopen it and get it
	err = ReopenTask(tokenToUse, -1)
	assert.NotNil(t, err)
	err = ReopenTask(tokenToUse, created.ID)
	assert.Nil(t, err)
	found, err = GetActiveTask(tokenToUse, created.ID)
	assert.Nil(t, err)
	assert.NotNil(t, found)

	// delete it, make sure it is gone
	err = DeleteTask(tokenToUse, -1)
	assert.NotNil(t, err)
	err = DeleteTask(tokenToUse, created.ID)
	assert.Nil(t, err)

	found, err = GetActiveTask(tokenToUse, created.ID)
	assert.NotNil(t, err)
	assert.Nil(t, found)
}
