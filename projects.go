package todoist

import (
	"encoding/json"
	"errors"

	"github.com/mitchellh/mapstructure"
)

// Project represents a project at Todoist, which holds the tasks, etc. Although the SKD doesn't interact with a DB, we provide an opinionated
// db name for the field that is the same as the JSON
type Project struct {
	ID           int64  `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	CommentCount int64  `json:"comment_count" db:"comment_count"`
	Order        int64  `json:"order" db:"order"`
	Color        int64  `json:"color" db:"color"`
	Shared       bool   `json:"shared" db:"shared"`
	SyncID       int64  `json:"sync_id" db:"sync_id"`
	Favorite     bool   `json:"favorite" db:"favorite"`
	InboxProject bool   `json:"inbox_project" db:"inbox_project"`
	URL          string `json:"url" db:"url"`
	TeamInbox    bool   `json:"team_inbox" db:"team_inbox"`
	ParentID     int64  `json:"parent_id" db:"parent_id"`
}

// GetAllProjects returns all of the project for a user's token
func GetAllProjects(token string) ([]Project, error) {
	projects := []Project{}
	resp, err := makeCall(token, EndpointNameGetProjects, map[string]string{}, nil)
	if err != nil {
		// parse the error string
		return projects, err
	}
	ps := []interface{}{}
	err = json.Unmarshal(resp.Body, &ps)
	if err != nil {
		return projects, errors.New("could not parse the response")
	}
	err = mapstructure.Decode(ps, &projects)
	return projects, err
}
