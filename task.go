package todoist

import (
	"encoding/json"
	"errors"
)

// Task represents a single todo item to track
type Task struct {
	ID           int64       `json:"id" db:"id"`
	ProjectID    int64       `json:"project_id" db:"project_id"`
	SectionID    int64       `json:"section_id" db:"section_id"`
	Content      string      `json:"content" db:"content"`
	Description  string      `json:"description" db:"description"`
	Completed    bool        `json:"completed" db:"completed"`
	LabelIDs     []string    `json:"label_ids" db:"label_ids"`
	ParentID     int64       `json:"parent_id" db:"parent_id"`
	Order        int64       `json:"order" db:"order"`
	Priority     int64       `json:"priority" db:"priority"`
	Due          TaskDueInfo `json:"due" db:"due"`
	URL          string      `json:"url" db:"url"`
	CommentCount int64       `json:"comment_count" db:"comment_count"`
	Assignee     int64       `json:"assignee" db:"assignee"`
	Assigner     int64       `json:"assigner" db:"assigner"`
}

// TaskDueInfo is the date/time information for a task
type TaskDueInfo struct {
	Date      string `json:"date" db:"date"`
	DateTime  string `json:"datetime" db:"datetime"`
	Recurring bool   `json:"recurring" db:"recurring"`
	String    string `json:"string" db:"string"`
	Timezone  string `json:"timezone" db:"timezone"`
}

// GetActiveTasks gets the active tasks for a user
func GetActiveTasks(token string) ([]Task, error) {
	tasks := []Task{}
	resp, err := makeCall(token, EndpointNameGetAllActiveTasks, map[string]string{}, nil)
	if err != nil {
		return tasks, err
	}
	err = json.Unmarshal(resp.Body, &tasks)
	return tasks, err
}

// CreateTask creates a returns a new task
func CreateTask(token string, input *Task) (*Task, error) {
	if input == nil {
		return nil, errors.New("you must provide a valid task struct with at least a content field")
	}
	if input.Content == "" {
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
