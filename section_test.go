package todoist

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSectionCRUD(t *testing.T) {
	setup()
	tokenToUse := "" // if you didn't want to pass it to the env for a one off test, set this (but don't commit it!)
	if tokenToUse == "" {
		tokenToUse = config.AuthToken
	}
	existingToken := config.AuthToken
	config.AuthToken = ""
	sections, err := GetAllSections("", 1)
	assert.NotNil(t, err)
	assert.Zero(t, len(sections))
	assert.Equal(t, "Empty token", err.Error())
	config.AuthToken = existingToken
	if tokenToUse == "" {
		// all further tests will fail, so we return from here
		fmt.Println("Skipping the rest of the Section tests")
		return
	}

	project, err := CreateTestProject(tokenToUse)
	assert.Nil(t, err)
	require.NotNil(t, project)
	defer DeleteProject(tokenToUse, project.ID)
	r := rand.Int63n(999999)

	// for this test, we simply work with one section; other tests will handle the task assignment
	params := &SectionParams{}
	created, err := CreateSection(tokenToUse, params)
	assert.Nil(t, created)
	assert.NotNil(t, err)
	name := fmt.Sprintf("Section %d", r)
	params.Name = name
	created, err = CreateSection(tokenToUse, params)
	assert.Nil(t, created)
	assert.NotNil(t, err)
	params.ProjectID = Int64(project.ID)
	created, err = CreateSection(tokenToUse, params)
	assert.NotNil(t, created)
	assert.Nil(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, name, created.Name)
	assert.Equal(t, project.ID, created.ProjectID)
	defer DeleteSection(tokenToUse, created.ID)

	// get it and make sure everything matches
	foundSection := false
	allSections, err := GetAllSections(tokenToUse, project.ID)
	assert.Nil(t, err)
	for _, s := range allSections {
		if s.ID == created.ID {
			foundSection = true
		}
	}
	assert.True(t, foundSection)

	section, err := GetSection(tokenToUse, created.ID)
	require.NotNil(t, section)
	assert.Nil(t, err)
	assert.NotZero(t, section.ID)
	assert.Equal(t, created.ID, section.ID)
	assert.Equal(t, name, section.Name)

	// update it, make sure it sticks
	newName := "Updated section"
	updated, err := UpdateSection(tokenToUse, section.ID, &SectionParams{
		ProjectID: Int64(1), // should be ignored
		Name:      newName,
	})
	assert.Nil(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, newName, updated.Name)
	assert.Equal(t, project.ID, updated.ProjectID)

	section, err = GetSection(tokenToUse, created.ID)
	require.NotNil(t, section)
	assert.Nil(t, err)
	assert.NotZero(t, section.ID)
	assert.Equal(t, created.ID, section.ID)
	assert.Equal(t, newName, section.Name)

	// delete it and make sure it is gone
	err = DeleteSection(tokenToUse, section.ID)
	assert.Nil(t, err)
	shouldBeGone, err := GetSection(tokenToUse, section.ID)
	assert.Nil(t, shouldBeGone)
	assert.NotNil(t, err)
}

func TestTaskSectionAssignment(t *testing.T) {
	setup()
	tokenToUse := ""
	if tokenToUse == "" {
		// all further tests will fail, so we return from here
		fmt.Println("Skipping the rest of the Section tests")
		return
	}

	project, err := CreateTestProject(tokenToUse)
	assert.Nil(t, err)
	require.NotNil(t, project)
	defer DeleteProject(tokenToUse, project.ID)
	r := rand.Int63n(999999)

	name := fmt.Sprintf("Section %d", r)
	params := &SectionParams{
		ProjectID: Int64(project.ID),
		Name:      name,
	}
	createdSection, err := CreateSection(tokenToUse, params)
	assert.NotNil(t, createdSection)
	assert.Nil(t, err)
	assert.NotZero(t, createdSection.ID)
	assert.Equal(t, name, createdSection.Name)
	assert.Equal(t, project.ID, createdSection.ProjectID)
	defer DeleteSection(tokenToUse, createdSection.ID)

	// create a task with that section set
	createTaskInput := &TaskParams{
		Content:     String(fmt.Sprintf("My New Task %d", r)),
		Description: String("Created from a unit test"),
		ProjectID:   Int64(project.ID),
		Priority:    PriorityUrgent,
		SectionID:   Int64(createdSection.ID),
	}
	createdTask, err := CreateTask(tokenToUse, createTaskInput)
	assert.Nil(t, err)
	require.NotNil(t, createdTask)
	defer DeleteTask(tokenToUse, createdTask.ID)
	assert.NotZero(t, createdTask.ID)
	assert.Equal(t, createdSection.ID, createdTask.SectionID)
}
