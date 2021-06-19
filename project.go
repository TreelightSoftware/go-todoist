package todoist

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
)

// Project represents a project at Todoist, which holds the tasks, etc. Although the SKD doesn't interact with a DB, we provide an opinionated
// db name for the field that is the same as the JSON
type Project struct {
	ID           int64  `json:"id,omitempty" db:"id"`
	Name         string `json:"name" db:"name"`
	CommentCount int64  `json:"comment_count" db:"comment_count"`
	Order        int64  `json:"order" db:"order"`
	Color        Color  `json:"color" db:"color"`
	Shared       bool   `json:"shared" db:"shared"`
	SyncID       int64  `json:"sync_id" db:"sync_id"`
	Favorite     bool   `json:"favorite" db:"favorite"`
	InboxProject bool   `json:"inbox_project" db:"inbox_project"`
	URL          string `json:"url" db:"url"`
	TeamInbox    bool   `json:"team_inbox" db:"team_inbox"`
	ParentID     int64  `json:"parent_id" db:"parent_id"`
}

// ProjectParams are the fields set during creating or updating a project
type ProjectParams struct {
	Name     *string `json:"name,omitempty" db:"name"`
	ParentID *int64  `json:"parent_id,omitempty" db:"parent_id"`
	Color    Color   `json:"color,omitempty" db:"color"`
	Favorite *bool   `json:"favorite,omitempty" db:"favorite"`
}

// GetAllProjects returns all of the project for a user's token. https://developer.todoist.com/rest/v1/#get-all-projects
func GetAllProjects(token string) ([]Project, error) {
	projects := []Project{}
	resp, err := makeCall(token, EndpointNameGetProjects, map[string]string{}, nil)
	if err != nil {
		return projects, err
	}
	err = json.Unmarshal(resp.Body, &projects)
	return projects, err
}

// CreateProject creates a new project for the user. Tasks belong to projects and require, at a minimum, a name. https://developer.todoist.com/rest/v1/#create-a-new-project
func CreateProject(token string, input *ProjectParams) (*Project, error) {
	if input == nil {
		return nil, errors.New("you must provide a valid project struct with at least a name field")
	}
	if input.Name == nil {
		return nil, errors.New("name is required")
	}
	resp, err := makeCall(token, EndpointNameCreateProject, map[string]string{}, input)
	if err != nil {
		return nil, err
	}
	created := &Project{}
	err = json.Unmarshal(resp.Body, &created)
	return created, err
}

// CreateTestProject creates a simple test project to be used in tests
func CreateTestProject(token string) (*Project, error) {
	r := rand.Int63n((999999999))
	input := ProjectParams{
		Name: String(fmt.Sprintf("Test Project %d", r)),
	}
	return CreateProject(token, &input)
}

// GetProject gets a single project by its id. https://developer.todoist.com/rest/v1/#get-a-project
func GetProject(token string, projectID int64) (*Project, error) {
	resp, err := makeCall(token, EndpointNameGetProject, map[string]string{
		"id": fmt.Sprintf("%d", projectID),
	}, nil)
	if err != nil {
		return nil, err
	}
	found := &Project{}
	err = json.Unmarshal(resp.Body, &found)
	return found, err
}

// UpdateProject updates a project. Currently, only name, color, and favorit are supported. https://developer.todoist.com/rest/v1/#update-a-project
func UpdateProject(token string, projectID int64, params *ProjectParams) (*Project, error) {
	if params == nil {
		return nil, errors.New("you must pass in valid update data")
	}
	_, err := makeCall(token, EndpointNameUpdateProject, map[string]string{
		"id": fmt.Sprintf("%d", projectID),
	}, params)
	if err != nil {
		return nil, err
	}
	// the update itself returns nothing, so we need to
	// get it again if we want the updated information
	return GetProject(token, projectID)
}

// DeleteProject deletes a project. https://developer.todoist.com/rest/v1/#delete-a-project
func DeleteProject(token string, projectID int64) error {
	resp, err := makeCall(token, EndpointNameDeleteProject, map[string]string{
		"id": fmt.Sprintf("%d", projectID),
	}, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}
	return nil
}
