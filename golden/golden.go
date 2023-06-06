package golden

import (
	"encoding/json"
	"flag"
	"github.com/Senth/go-blu-tests/fixture"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// Golden tests are used to compare the actual output of a function with the expected output.
type Golden interface {
	// AssertGolden compares the actual output of a function with the expected output.
	// Update the golden by passing the -update flag to the test.
	AssertGolden(t *testing.T, goldenFile string, actual interface{}, msgsAndArgs ...interface{})
}

type goldenImpl struct {
	updateFlag bool
}

// NewGolden creates a new goldenImpl implementation used for tests
func NewGolden(t *testing.T) Golden {
	updateFlag := flag.Lookup("update") != nil
	return &goldenImpl{
		updateFlag: updateFlag,
	}
}

// AssertGolden compares the actual output of a function with the expected output.
func (g *goldenImpl) AssertGolden(t *testing.T, goldenFile string, actual interface{}, msgsAndArgs ...interface{}) {
	t.Helper()

	// Update the golden file
	if g.updateFlag {
		g.updateGolden(t, goldenFile, actual)
	}

	expected := string(fixture.ReadFixture(t, goldenFile))

	// Use JSON to compare the actual and expected output
	if isJSONFile(goldenFile) {
		g.assertGoldenJSON(t, expected, actual, msgsAndArgs...)
		return
	}
}

func (g *goldenImpl) assertGoldenJSON(t *testing.T, expected string, actual interface{}, msgsAndArgs ...interface{}) {
	t.Helper()

	// Convert the struct to JSON
	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatalf("failed to marshal actual output to JSON: %s", err)
	}

	assert.Equal(t, expected, string(actualJSON), msgsAndArgs...)
}

func (g *goldenImpl) updateGolden(t *testing.T, goldenFile string, actual interface{}) {
	t.Helper()

	// Update as JSON
	if isJSONFile(goldenFile) {
		g.updateGoldenJSON(t, goldenFile, actual)
		return
	}

	switch actual := actual.(type) {
	case []byte:
		fixture.WriteFixture(t, goldenFile, actual)
	case string:
		fixture.WriteFixture(t, goldenFile, []byte(actual))
	default:
		t.Fatalf("unknown type for golden update: %T", actual)
	}
}

func (g *goldenImpl) updateGoldenJSON(t *testing.T, goldenFile string, actual interface{}) {
	t.Helper()

	// Convert the struct to JSON
	actualJSON, err := json.MarshalIndent(actual, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal actual output to JSON: %s", err)
	}

	fixture.WriteFixture(t, goldenFile, actualJSON)
}

func isJSONFile(filename string) bool {
	return strings.HasSuffix(filename, ".json")
}
