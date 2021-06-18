package todoist

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProjects(t *testing.T) {
	setup()
	tokenToUse := "" // if you didn't want to pass it to the env for a one off test, set this (but don't commit it!)
	if tokenToUse == "" {
		tokenToUse = config.AuthToken
	}
	existingToken := config.AuthToken
	config.AuthToken = ""
	projects, err := GetAllProjects("")
	assert.NotNil(t, err)
	assert.Zero(t, len(projects))
	assert.Equal(t, "Empty token", err.Error())
	config.AuthToken = existingToken
	if tokenToUse == "" {
		// all further tests will fail, so we return from here
		fmt.Println("Skipping the rest of the Project tests")
		return
	}
	r := rand.Int63n(9999999)
	projectName := fmt.Sprintf("My Test Project %d", r)
	created, err := CreateProject(tokenToUse, &ProjectParams{
		Name:  String(projectName),
		Color: ColorLightBlue,
	})
	assert.Nil(t, err)
	fmt.Printf("\n%+v\n", created)
	assert.NotZero(t, created.ID)
	assert.Equal(t, projectName, created.Name)
	assert.Equal(t, ColorLightBlue, created.Color)
	defer DeleteProject(tokenToUse, created.ID)

	// get it in the list AND singularly
	allProjects, err := GetAllProjects(tokenToUse)
	fmt.Printf("\nALL\n%+v\n", allProjects)
	assert.Nil(t, err)
	foundInSlice := false
	for i := range allProjects {
		if allProjects[i].ID == created.ID {
			foundInSlice = true
			break
		}
	}
	assert.True(t, foundInSlice)
	found, err := GetProject(tokenToUse, created.ID)
	assert.Nil(t, err)
	require.NotNil(t, found)
	assert.Equal(t, created.Name, found.Name)

	// change it a bit; since the Update call then chains to a Get, we don't need to re-get at this time
	updateData := &ProjectParams{
		Name:     String(fmt.Sprintf("Updated %d", r)),
		Color:    ColorTaupe,
		Favorite: Bool(true),
	}
	updated, err := UpdateProject(tokenToUse, created.ID, updateData)
	assert.Nil(t, err)
	require.NotNil(t, updated)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, updateData.Color, updated.Color)
	assert.Equal(t, fmt.Sprintf("Updated %d", r), updated.Name)
	assert.True(t, updated.Favorite)

	err = DeleteProject("", created.ID)
	assert.Nil(t, err)

	// make sure it is gone

	allProjects, err = GetAllProjects(tokenToUse)
	assert.Nil(t, err)
	foundInSlice = false
	for i := range allProjects {
		if allProjects[i].ID == created.ID {
			foundInSlice = true
			break
		}
	}
	assert.False(t, foundInSlice)
	notFound, err := GetProject(tokenToUse, created.ID)
	assert.NotNil(t, err)
	assert.Nil(t, notFound)

}
