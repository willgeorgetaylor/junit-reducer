package enums

import (
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
