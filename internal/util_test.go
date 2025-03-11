package internal

import (
	"os"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestIsInt(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected bool
	}{
		{"Integer", 42, true},
		{"Intish", 1.00000001, true},
		{"Pi", 3.14159, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsInt(decimal.NewFromFloat(tc.input))
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"Integer", 42, 42},
		{"Intify", 1.00000001, 1},
		{"-Intify", -1.00000001, -1},
		{"Pi", 3.14159, 3.14159},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Normalize(decimal.NewFromFloat(tc.input)).InexactFloat64()
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFileExists(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp(t.TempDir(), "test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Existing File", tmpFile.Name(), true},
		{"Non-Existing File", "/path/that/does/not/exist", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FileExists(tc.path)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestReversed(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []int{5, 4, 3, 2, 1}, Reversed(numbers))
}

func TestStyleBetweenStars(t *testing.T) {
	boldStyle := lipgloss.NewStyle().Bold(true)
	result := StyleBetweenStars("**Hello** **world**!", boldStyle)
	assert.Equal(t, boldStyle.Render("Hello")+" "+boldStyle.Render("world")+"!", result)
}
