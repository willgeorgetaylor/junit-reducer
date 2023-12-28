package enums

import (
	"reflect"
	"testing"
)

func TestAggregateOperationValues(t *testing.T) {
	expectedValues := map[string]AggregateOperation{
		"mean":   AggregateOperationMean,
		"mode":   AggregateOperationMode,
		"median": AggregateOperationMedian,
		"min":    AggregateOperationMin,
		"max":    AggregateOperationMax,
		"sum":    AggregateOperationSum,
	}

	for key, expectedValue := range expectedValues {
		actualValue, ok := AggregateOperationValues[key]
		if !ok {
			t.Errorf("Expected key '%s' not found in AggregateOperationValues", key)
		}

		if actualValue != expectedValue {
			t.Errorf("Expected value '%d' for key '%s', got '%d'", expectedValue, key, actualValue)
		}
	}
}

func TestGetTestSuiteFields(t *testing.T) {
	expectedFields := []string{"filepath", "name", "name+filepath"}

	actualFields := GetTestSuiteFields()

	if !reflect.DeepEqual(actualFields, expectedFields) {
		t.Errorf("Expected fields %v, but got %v", expectedFields, actualFields)
	}
}

func TestGetTestCaseFields(t *testing.T) {
	expectedFields := []string{"classname", "file", "name"}

	actualFields := GetTestCaseFields()

	if !reflect.DeepEqual(actualFields, expectedFields) {
		t.Errorf("Expected fields %v, but got %v", expectedFields, actualFields)
	}
}

func TestGetAggregateOperations(t *testing.T) {
	expectedOperations := []string{"max", "mean", "median", "min", "mode", "sum"}

	actualOperations := GetAggregateOperations()

	if !reflect.DeepEqual(actualOperations, expectedOperations) {
		t.Errorf("Expected operations %v, but got %v", expectedOperations, actualOperations)
	}
}

func TestGetRoundingModes(t *testing.T) {
	expectedModes := []string{"ceil", "floor", "round"}

	actualModes := GetRoundingModes()

	if !reflect.DeepEqual(actualModes, expectedModes) {
		t.Errorf("Expected modes %v, but got %v", expectedModes, actualModes)
	}
}
