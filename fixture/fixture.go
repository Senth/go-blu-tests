package fixture

import (
	"os"
	"path/filepath"
	"testing"
)

const testDataDir = "testdata"

func ReadFixture(t *testing.T, filename string) []byte {
	t.Helper()

	filePath := filepath.Join(testDataDir, filename)

	// Read data from file
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file '%s': %s", filePath, err)
	}

	return data
}

func WriteFixture(t *testing.T, filename string, data []byte) {
	t.Helper()

	filePath := filepath.Join(testDataDir, filename)

	// Write data to file
	if err := os.WriteFile(filePath, data, 0600); err != nil {
		t.Fatalf("failed to write file '%s': %s", filePath, err)
	}
}
