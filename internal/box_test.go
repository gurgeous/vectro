package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBox(t *testing.T) {
	box := NewBox(5, 16)

	// Width/Height
	assert.Equal(t, 123, box.Width(123).width)
	assert.Equal(t, 456, box.Height(456).height)

	// Rows and Cols
	a, b := box.Rows()
	assert.Equal(t, []int{5, 8}, []int{a.width, a.height})
	assert.Equal(t, []int{5, 8}, []int{b.width, b.height})
	a, b = box.Cols()
	assert.Equal(t, []int{2, 16}, []int{a.width, a.height})
	assert.Equal(t, []int{3, 16}, []int{b.width, b.height})

	// Cut
	a, b = box.CutLeft(1)
	assert.Equal(t, []int{1, 16}, []int{a.width, a.height})
	assert.Equal(t, []int{4, 16}, []int{b.width, b.height})
	a, b = box.CutRight(1)
	assert.Equal(t, []int{4, 16}, []int{a.width, a.height})
	assert.Equal(t, []int{1, 16}, []int{b.width, b.height})
	a, b = box.CutTop(1)
	assert.Equal(t, []int{5, 1}, []int{a.width, a.height})
	assert.Equal(t, []int{5, 15}, []int{b.width, b.height})
	a, b = box.CutBottom(1)
	assert.Equal(t, []int{5, 15}, []int{a.width, a.height})
	assert.Equal(t, []int{5, 1}, []int{b.width, b.height})
}
