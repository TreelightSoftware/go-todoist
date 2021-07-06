package todoist

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Task represents a single todo item to track
type Task struct {
	ID           int64       `json:"id" db:"id"`
	ProjectID    int64       `json:"project_id" db:"project_id"`
	SectionID    int64       `json:"section_id" db:"section_id"`
	Content      string      `json:"content" db:"content"`
	Description  string      `json:"description" db:"description"`
	Completed    bool        `json:"completed" db:"completed"`
	LabelIDs     []int64     `json:"label_ids" db:"label_ids"`
	ParentID     int64       `json:"parent_id" db:"parent_id"`
	Order        int64       `json:"order" db:"order"`
	Priority     Priority    `json:"priority" db:"priority"`
	Due          TaskDueInfo `json:"due" db:"due"`
	URL          string      `json:"url" db:"url"`
	CommentCount int64       `json:"comment_count" db:"comment_count"`
	Assignee     int64       `json:"assignee" db:"assignee"`
	Assigner     int64       `json:"assigner" db:"assigner"`
}

// TaskDueInfo is the date/time information for a task
type TaskDueInfo struct {
	Date      string `json:"date" db:"date"`
	Datetime  string `json:"datetime" db:"datetime"`
	Recurring bool   `json:"recurring" db:"recurring"`
	String    string `json:"string" db:"string"`
	Timezone  string `json:"timezone" db:"timezone"`
}

// TaskParams are the fields you can set when creating or updating
type TaskParams struct {
	ProjectID    *int64   `json:"project_id,omitempty" db:"project_id"`
	SectionID    *int64   `json:"section_id,omitempty" db:"section_id"`
	Content      *string  `json:"content,omitempty" db:"content"`
	Description  *string  `json:"description,omitempty" db:"description"`
	Completed    *bool    `json:"completed,omitempty" db:"completed"`
	LabelIDs     *[]int64 `json:"label_ids,omitempty" db:"label_ids"`
	ParentID     *int64   `json:"parent_id,omitempty" db:"parent_id"`
	Order        *int64   `json:"order,omitempty" db:"order"`
	Priority     Priority `json:"priority,omitempty" db:"priority"`
	URL          *string  `json:"url,omitempty" db:"url"`
	CommentCount *int64   `json:"comment_count,omitempty" db:"comment_count"`
	Assignee     *int64   `json:"assignee,omitempty" db:"assignee"`
	Assigner     *int64   `json:"assigner,omitempty" db:"assigner"`

	// we lift the following up for the update and create calls
	DueLang     *string `json:"due_lang" db:"due_lang"`
	DueString   *string `json:"due_string" db:"due_string"`
	DueDate     *string `json:"due_date" db:"due_date"`
	DueDatetime *string `json:"due_datetime" db:"due_datetime"`
}

// GetActiveTasks gets the active tasks for a user. https://developer.todoist.com/rest/v1/#get-active-tasks
func GetActiveTasks(token string) ([]Task, error) {
	tasks := []Task{}
	resp, err := makeCall(token, EndpointNameGetAllActiveTasks, map[string]string{}, nil)
	if err != nil {
		return tasks, err
	}
	err = json.Unmarshal(resp.Body, &tasks)
	return tasks, err
}

// CreateTask creates a returns a new task. The only required field is the content field. https://developer.todoist.com/rest/v1/#create-a-new-task
func CreateTask(token string, input *TaskParams) (*Task, error) {
	if input == nil {
		return nil, errors.New("you must provide a valid input with at least a content field")
	}
	if input.Content == nil || StringValue(input.Content) == "" {
		return nil, errors.New("content is required")
	}
	resp, err := makeCall(token, EndpointNameCreateTask, map[string]string{}, input)
	if err != nil {
		return nil, err
	}
	created := &Task{}
	err = json.Unmarshal(resp.Body, &created)
	return created, err
}

// GetActiveTask gets a single task by its id. https://developer.todoist.com/rest/v1/#get-an-active-task
func GetActiveTask(token string, taskID int64) (*Task, error) {
	resp, err := makeCall(token, EndpointNameGetTask, map[string]string{
		"id": fmt.Sprintf("%d", taskID),
	}, nil)
	if err != nil {
		return nil, err
	}
	found := &Task{}
	err = json.Unmarshal(resp.Body, &found)
	return found, err
}

// UpdateTask updates a task. https://developer.todoist.com/rest/v1/#update-a-task
func UpdateTask(token string, taskID int64, newData *TaskParams) (*Task, error) {
	if newData == nil {
		return nil, errors.New("you must pass in a valid input")
	}
	_, err := makeCall(token, EndpointNameUpdateTask, map[string]string{
		"id": fmt.Sprintf("%d", taskID),
	}, newData)
	if err != nil {
		return nil, err
	}
	// the update itself returns nothing, so we need to
	// get it again if we want the updated information
	return GetActiveTask(token, taskID)
}

// DeleteTask deletes a task. You probably want to close it instead? https://developer.todoist.com/rest/v1/#delete-a-task
func DeleteTask(token string, taskID int64) error {
	resp, err := makeCall(token, EndpointNameDeleteTask, map[string]string{
		"id": fmt.Sprintf("%d", taskID),
	}, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}
	return nil
}

// CloseTask closes a task. According to the docs, this will cause root tasks to be marked complete and moved to the history. https://developer.todoist.com/rest/v1/#close-a-task
func CloseTask(token string, taskID int64) error {
	resp, err := makeCall(token, EndpointNameCloseTask, map[string]string{
		"id": fmt.Sprintf("%d", taskID),
	}, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}
	return nil
}

// ReopenTask reopens a closed task. https://developer.todoist.com/rest/v1/#reopen-a-task
func ReopenTask(token string, taskID int64) error {
	resp, err := makeCall(token, EndpointNameReopenTask, map[string]string{
		"id": fmt.Sprintf("%d", taskID),
	}, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("received status code %d", resp.StatusCode)
	}
	return nil
}
