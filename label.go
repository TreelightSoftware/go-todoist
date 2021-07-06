package todoist

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Label represent a label that can be assigned to a task
type Label struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Color    int64  `json:"color" db:"color"`
	Order    int64  `json:"order" db:"order"`
	Favorite bool   `json:"favorite" db:"favorite"`
}

// LabelParams are used when creating or updating a label
type LabelParams struct {
	Name     string `json:"name" db:"name"`
	Color    *int64 `json:"color" db:"color"`
	Order    *int64 `json:"order" db:"order"`
	Favorite *bool  `json:"favorite" db:"favorite"`
}

// GetAllLabels returns all of the labels for a user's token. https://developer.todoist.com/rest/v1/#get-all-labels
func GetAllLabels(token string) ([]Label, error) {
	labels := []Label{}
	resp, err := makeCall(token, EndpointNameGetAllLabels, map[string]string{}, nil)
	if err != nil {
		return labels, err
	}
	err = json.Unmarshal(resp.Body, &labels)
	return labels, err
}

// CreateLabel creates a label and requires at least a name. https://developer.todoist.com/rest/v1/#create-a-new-label
func CreateLabel(token string, input *LabelParams) (*Label, error) {
	if input == nil {
		return nil, errors.New("you must provide a valid input with at least a name field")
	}
	if input.Name == "" {
		return nil, errors.New("name is required")
	}
	resp, err := makeCall(token, EndpointNameCreateLabel, map[string]string{}, input)
	if err != nil {
		return nil, err
	}
	created := &Label{}
	err = json.Unmarshal(resp.Body, &created)
	return created, err
}

// GetLabel gets a single label. https://developer.todoist.com/rest/v1/#get-a-label
func GetLabel(token string, labelID int64) (*Label, error) {
	resp, err := makeCall(token, EndpointNameGetLabel, map[string]string{
		"id": fmt.Sprintf("%d", labelID),
	}, nil)
	if err != nil {
		return nil, err
	}
	found := &Label{}
	err = json.Unmarshal(resp.Body, &found)
	return found, err
}

// UpdateLabel updates a label. https://developer.todoist.com/rest/v1/#update-a-label
func UpdateLabel(token string, labelID int64, input *LabelParams) (*Label, error) {
	if input == nil {
		return nil, errors.New("you must provide a valid input")
	}
	_, err := makeCall(token, EndpointNameUpdateLabel, map[string]string{
		"id": fmt.Sprintf("%d", labelID),
	}, input)
	if err != nil {
		return nil, err
	}
	// the update itself returns nothing, so we need to
	// get it again if we want the updated information
	return GetLabel(token, labelID)
}

// DeleteLabel deletes a label. https://developer.todoist.com/rest/v1/#delete-a-label
func DeleteLabel(token string, labelID int64) error {
	resp, err := makeCall(token, EndpointNameDeleteLabel, map[string]string{
		"id": fmt.Sprintf("%d", labelID),
	}, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}
	return nil
}
