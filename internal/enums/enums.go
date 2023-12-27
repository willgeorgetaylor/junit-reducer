package enums

// TestSuiteField

type TestSuiteField int

const (
	TestSuiteFieldName TestSuiteField = iota
	TestSuiteFieldFilepath
	TestSuiteFieldNameFilepath
)

var TestSuiteFieldKeys = map[TestSuiteField]string{
	TestSuiteFieldName:         "name",
	TestSuiteFieldFilepath:     "filepath",
	TestSuiteFieldNameFilepath: "name+filepath",
}

var TestSuiteFieldValues = map[string]TestSuiteField{
	"name":          TestSuiteFieldName,
	"filepath":      TestSuiteFieldFilepath,
	"name+filepath": TestSuiteFieldNameFilepath,
}

func GetTestSuiteFields() []string {
	TestSuiteFieldInputs := make([]string, len(TestSuiteFieldValues))
	i := 0
	for key := range TestSuiteFieldValues {
		TestSuiteFieldInputs[i] = key
		i++
	}
	return TestSuiteFieldInputs
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

func GetTestCaseFields() []string {
	TestCaseFieldInputs := make([]string, len(TestCaseFieldValues))
	i := 0
	for key := range TestCaseFieldValues {
		TestCaseFieldInputs[i] = key
		i++
	}
	return TestCaseFieldInputs
}

// Aggregate operations

type AggregateOperation int

const (
	AggregateOperationMean AggregateOperation = iota
	AggregateOperationMode
	AggregateOperationMedian
	AggregateOperationMin
	AggregateOperationMax
	AggregateOperationSum
)

var AggregateOperationKeys = map[AggregateOperation]string{
	AggregateOperationMean:   "mean",
	AggregateOperationMode:   "mode",
	AggregateOperationMedian: "median",
	AggregateOperationMin:    "min",
	AggregateOperationMax:    "max",
	AggregateOperationSum:    "sum",
}

var AggregateOperationValues = map[string]AggregateOperation{
	"mean":   AggregateOperationMean,
	"mode":   AggregateOperationMode,
	"median": AggregateOperationMedian,
	"min":    AggregateOperationMin,
	"max":    AggregateOperationMax,
	"sum":    AggregateOperationSum,
}

func GetAggregateOperations() []string {
	AggregateOperationInputs := make([]string, len(AggregateOperationValues))
	i := 0
	for key := range AggregateOperationValues {
		AggregateOperationInputs[i] = key
		i++
	}
	return AggregateOperationInputs
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

func GetRoundingModes() []string {
	RoundingModeInputs := make([]string, len(RoundingModeValues))
	i := 0
	for key := range RoundingModeValues {
		RoundingModeInputs[i] = key
		i++
	}
	return RoundingModeInputs
}
