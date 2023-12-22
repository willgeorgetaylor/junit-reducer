package enums

// TestSuiteField

type TestSuiteField int

const (
	TestSuiteFieldName TestSuiteField = iota
	TestSuiteFieldFilepath
)

var TestSuiteFieldKeys = map[TestSuiteField]string{
	TestSuiteFieldName:     "name",
	TestSuiteFieldFilepath: "filepath",
}

var TestSuiteFieldValues = map[string]TestSuiteField{
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

var TestCaseFieldKeys = map[TestCaseField]string{
	TestCaseFieldName:      "name",
	TestCaseFieldClassname: "classname",
	TestCaseFieldFile:      "file",
}

var TestCaseFieldValues = map[string]TestCaseField{
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

var AggregateOperationKeys = map[AggregateOperation]string{
	AggregateOperationMean:   "mean",
	AggregateOperationMode:   "mode",
	AggregateOperationMedian: "median",
	AgregateOperationMin:     "min",
	AggregateOperationMax:    "max",
	AggregateOperationSum:    "sum",
}

var AggregateOperationValues = map[string]AggregateOperation{
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

var RoundingModeKeys = map[RoundingMode]string{
	RoundingModeRound: "round",
	RoundingModeCeil:  "ceil",
	RoundingModeFloor: "floor",
}

var RoundingModeValues = map[string]RoundingMode{
	"round": RoundingModeRound,
	"ceil":  RoundingModeCeil,
	"floor": RoundingModeFloor,
}