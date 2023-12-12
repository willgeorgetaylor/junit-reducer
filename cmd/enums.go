package cmd

// TestSuiteField

type TestSuiteField int

const (
	TestSuiteFieldName TestSuiteField = iota
	TestSuiteFieldFilepath
)

var testSuiteFieldKeys = map[TestSuiteField]string{
	TestSuiteFieldName:     "name",
	TestSuiteFieldFilepath: "filepath",
}

var testSuiteFieldValues = map[string]TestSuiteField{
	"name":     TestSuiteFieldName,
	"filepath": TestSuiteFieldFilepath,
}

// TestCaseField

type TestCaseField int

const (
	TestCaseFieldName TestCaseField = iota
	TestCaseFieldClassname
	TestCaseFieldFile
)

var testCaseFieldKeys = map[TestCaseField]string{
	TestCaseFieldName:      "name",
	TestCaseFieldClassname: "classname",
	TestCaseFieldFile:      "file",
}

var testCaseFieldValues = map[string]TestCaseField{
	"name":      TestCaseFieldName,
	"classname": TestCaseFieldClassname,
	"file":      TestCaseFieldFile,
}

// Aggregate operations

type AggregateOperation int

const (
	AggregateOperationMean AggregateOperation = iota
	AggregateOperationMode
	AggregateOperationMedian
	AgregateOperationMin
	AggregateOperationMax
	AggregateOperationSum
)

var aggregateOperationKeys = map[AggregateOperation]string{
	AggregateOperationMean:   "mean",
	AggregateOperationMode:   "mode",
	AggregateOperationMedian: "median",
	AgregateOperationMin:     "min",
	AggregateOperationMax:    "max",
	AggregateOperationSum:    "sum",
}

var aggregateOperationValues = map[string]AggregateOperation{
	"mean":   AggregateOperationMean,
	"mode":   AggregateOperationMode,
	"median": AggregateOperationMedian,
	"min":    AgregateOperationMin,
	"max":    AggregateOperationMax,
	"sum":    AggregateOperationSum,
}

// Rounding modes

type RoundingMode int

const (
	RoundingModeRound RoundingMode = iota
	RoundingModeCeil
	RoundingModeFloor
)

var roundingModeKeys = map[RoundingMode]string{
	RoundingModeRound: "round",
	RoundingModeCeil:  "ceil",
	RoundingModeFloor: "floor",
}

var roundingModeValues = map[string]RoundingMode{
	"round": RoundingModeRound,
	"ceil":  RoundingModeCeil,
	"floor": RoundingModeFloor,
}
