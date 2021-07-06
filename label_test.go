package todoist

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabelCRUD(t *testing.T) {
	setup()
	tokenToUse := "" // if you didn't want to pass it to the env for a one off test, set this (but don't commit it!)
	if tokenToUse == "" {
		tokenToUse = config.AuthToken
	}
	existingToken := config.AuthToken
	config.AuthToken = ""
	sections, err := GetAllLabels("")
	assert.NotNil(t, err)
	assert.Zero(t, len(sections))
	assert.Equal(t, "Empty token", err.Error())
	config.AuthToken = existingToken
	if tokenToUse == "" {
		// all further tests will fail, so we return from here
		fmt.Println("Skipping the rest of the Label tests")
		return
	}

	project, err := CreateTestProject(tokenToUse)
	assert.Nil(t, err)
	require.NotNil(t, project)
	defer DeleteProject(tokenToUse, project.ID)
	r := rand.Int63n(999999)

	// for this test, we simply work with one label; other tests will handle the task assignment
	params := &LabelParams{}
	created, err := CreateLabel(tokenToUse, params)
	assert.Nil(t, created)
	assert.NotNil(t, err)
	name := fmt.Sprintf("label_%d", r)
	params.Name = name
	created, err = CreateLabel(tokenToUse, params)
	assert.NotNil(t, created)
	assert.Nil(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, name, created.Name)
	defer DeleteLabel(tokenToUse, created.ID)

	// get it and make sure everything matches
	foundLabel := false
	allLabels, err := GetAllLabels(tokenToUse)
	assert.Nil(t, err)
	for _, s := range allLabels {
		if s.ID == created.ID {
			foundLabel = true
		}
	}
	assert.True(t, foundLabel)

	label, err := GetLabel(tokenToUse, created.ID)
	require.NotNil(t, label)
	assert.Nil(t, err)
	assert.NotZero(t, label.ID)
	assert.Equal(t, created.ID, label.ID)
	assert.Equal(t, name, label.Name)

	// update it, make sure it sticks
	newName := "updated_label"
	updated, err := UpdateLabel(tokenToUse, label.ID, &LabelParams{
		Name: newName,
	})
	assert.Nil(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, newName, updated.Name)

	label, err = GetLabel(tokenToUse, created.ID)
	require.NotNil(t, label)
	assert.Nil(t, err)
	assert.NotZero(t, label.ID)
	assert.Equal(t, created.ID, label.ID)
	assert.Equal(t, newName, label.Name)

	// delete it and make sure it is gone
	err = DeleteLabel(tokenToUse, label.ID)
	assert.Nil(t, err)
	shouldBeGone, err := GetLabel(tokenToUse, label.ID)
	assert.Nil(t, shouldBeGone)
	assert.NotNil(t, err)
}

func TestTaskLabelAssignment(t *testing.T) {
	setup()
	tokenToUse := ""
	if tokenToUse == "" {
		// all further tests will fail, so we return from here
		fmt.Println("Skipping the rest of the Label tests")
		return
	}

	project, err := CreateTestProject(tokenToUse)
	assert.Nil(t, err)
	require.NotNil(t, project)
	defer DeleteProject(tokenToUse, project.ID)
	r := rand.Int63n(999999)

	name := fmt.Sprintf("label_%d", r)
	params := &LabelParams{
		Name: name,
	}
	createdLabel, err := CreateLabel(tokenToUse, params)
	assert.NotNil(t, createdLabel)
	assert.Nil(t, err)
	assert.NotZero(t, createdLabel.ID)
	assert.Equal(t, name, createdLabel.Name)
	defer DeleteLabel(tokenToUse, createdLabel.ID)

	// create a task with that section set
	createTaskInput := &TaskParams{
		Content:     String(fmt.Sprintf("My New Task %d", r)),
		Description: String("Created from a unit test"),
		ProjectID:   Int64(project.ID),
		Priority:    PriorityUrgent,
		LabelIDs:    &[]int64{createdLabel.ID},
	}
	createdTask, err := CreateTask(tokenToUse, createTaskInput)
	assert.Nil(t, err)
	require.NotNil(t, createdTask)
	defer DeleteTask(tokenToUse, createdTask.ID)
	assert.NotZero(t, createdTask.ID)
	assert.NotZero(t, len(createdTask.LabelIDs))
	foundLabel := false
	for _, l := range createdTask.LabelIDs {
		if l == createdLabel.ID {
			foundLabel = true
		}
	}
	assert.True(t, foundLabel)
}
