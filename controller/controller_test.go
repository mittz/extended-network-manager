package controller

import (
	"testing"
)

func TestCompareMaps(t *testing.T) {
	a := map[string]interface{}{"key1": "value1", "key2": "value2", "key3": "value3"}
	b := map[string]interface{}{"key2": "value2"}

	actual := CompareMaps(a, b)
	expected := true

	if actual != expected {
		t.Errorf("actual: %v / expected: %v", actual, expected)
	}
}

func TestCallPipework(t *testing.T) {
	var event Event
	CallPipework(event)
}
