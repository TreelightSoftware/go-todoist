package todoist

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Section divides a project into logical sections
type Section struct {
	ID        int64  `json:"id" db:"id"`
	ProjectID int64  `json:"project_id" db:"project_id"`
	Name      string `json:"name" db:"name"`
	Order     int64  `json:"order" db:"order"`
}

// SectionParams are the fields used when creating or editing sections
type SectionParams struct {
	ProjectId *int64 `json:"project_id" db:"project_id"`
	Name      string `json:"name" db:"name"`
	Order     *int64 `json:"order" db:"order"`
}

// GetAllSections returns all of the sections for a user's token. If provided a non-zero project ID, it will get only sections for that project. https://developer.todoist.com/rest/v1/#get-all-sections
func GetAllSections(token string, projectID int64) ([]Section, error) {
	sections := []Section{}
	data := map[string]string{}
	if projectID != 0 {
		data["project_id"] = fmt.Sprintf("%d", projectID)
	}
	resp, err := makeCall(token, EndpointNameGetAllSections, map[string]string{}, data)
	if err != nil {
		return sections, err
	}
	err = json.Unmarshal(resp.Body, &sections)
	return sections, err
}

// CreateSection creates a section and requires at least a name and project_id. https://developer.todoist.com/rest/v1/#create-a-new-section
func CreateSection(token string, input *SectionParams) (*Section, error) {
	if input == nil {
		return nil, errors.New("you must provide a valid input with at least a name field and a project_id field")
	}
	if input.Name == "" || input.ProjectId == nil || Int64Value(input.ProjectId) == 0 {
		return nil, errors.New("name and project_id are required")
	}
	resp, err := makeCall(token, EndpointNameCreateSection, map[string]string{}, input)
	if err != nil {
		return nil, err
	}
	created := &Section{}
	err = json.Unmarshal(resp.Body, &created)
	return created, err
}

// GetSection gets a single section. https://developer.todoist.com/rest/v1/#get-a-single-section
func GetSection(token string, sectionID int64) (*Section, error) {
	resp, err := makeCall(token, EndpointNameGetSection, map[string]string{
		"id": fmt.Sprintf("%d", sectionID),
	}, nil)
	if err != nil {
		return nil, err
	}
	found := &Section{}
	err = json.Unmarshal(resp.Body, &found)
	return found, err
}

// UpdateSection updates a section. Currently, only the name may change. https://developer.todoist.com/rest/v1/#update-a-section
func UpdateSection(token string, sectionID int64, input *SectionParams) (*Section, error) {
	if input == nil {
		return nil, errors.New("you must provide a valid input with at least a name field")
	}
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	_, err := makeCall(token, EndpointNameUpdateSection, map[string]string{
		"id": fmt.Sprintf("%d", sectionID),
	}, input)
	if err != nil {
		return nil, err
	}
	// the update itself returns nothing, so we need to
	// get it again if we want the updated information
	return GetSection(token, sectionID)
}

// DeleteSection deletes a section. https://developer.todoist.com/rest/v1/#delete-a-section
func DeleteSection(token string, sectionID int64) error {
	resp, err := makeCall(token, EndpointNameDeleteSection, map[string]string{
		"id": fmt.Sprintf("%d", sectionID),
	}, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}
	return nil
}
