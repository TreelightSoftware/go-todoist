package todoist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientPointerHelpers(t *testing.T) {

	int64Val := int64(42)
	boolVal := true
	stringVal := "platypus"

	assert.Equal(t, int64Val, Int64Value(&int64Val))
	assert.NotNil(t, Int64(int64Val))
	assert.Equal(t, boolVal, BoolValue(&boolVal))
	assert.NotNil(t, Bool(boolVal))
	assert.Equal(t, stringVal, StringValue(&stringVal))
	assert.NotNil(t, String(stringVal))
	assert.Equal(t, int64(0), Int64Value(nil))
	assert.Equal(t, false, BoolValue(nil))
	assert.Equal(t, "", StringValue(nil))
}

func TestNoEndpointCall(t *testing.T) {
	resp, err := makeCall("test", "EndpointNameDoesNotExist", map[string]string{}, nil)
	assert.NotNil(t, err)
	assert.Zero(t, len(resp.Body))
}
