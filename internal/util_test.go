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
	tmpFile, err := os.CreateTemp("", "test")
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

func TestFilter(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result := Filter(numbers, func(n int) bool {
		return n%2 == 0
	})
	assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
}

func TestMap(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	result := Map(numbers, func(n int) string {
		return "Item " + string(rune('0'+n))
	})
	assert.Equal(t, []string{"Item 1", "Item 2", "Item 3", "Item 4", "Item 5"}, result)
}

func TestMapWithIndex(t *testing.T) {
	numbers := []int{10, 20, 30, 40, 50}
	result := MapWithIndex(numbers, func(i int, n int) int {
		return n * i
	})
	assert.Equal(t, []int{0, 20, 60, 120, 200}, result)
}

func TestReverse(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []int{5, 4, 3, 2, 1}, Reverse(numbers))
}

func TestStyleBetweenStars(t *testing.T) {
	boldStyle := lipgloss.NewStyle().Bold(true)
	result := StyleBetweenStars("**Hello** **world**!", boldStyle)
	assert.Equal(t, boldStyle.Render("Hello")+" "+boldStyle.Render("world")+"!", result)
}
