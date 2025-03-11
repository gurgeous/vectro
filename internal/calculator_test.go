package internal

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCalculator(t *testing.T) {
	c := NewCalculator()

	// Clear
	c.PushInt(1, 2, 3)
	assert.Equal(t, 3, c.Len())
	c.Clear()
	assert.Equal(t, 0, c.Len())

	// Empty
	assert.True(t, c.Empty())
	c.PushInt(1)
	assert.False(t, c.Empty())
	c.Clear()

	// Push/Peek/Pop
	c.PushInt(1, 2, 3)
	assert.Equal(t, 3, c.Len())
	assert.Equal(t, 3, c.PopInt())
	assert.Equal(t, 2, c.PeekInt())
	assert.Equal(t, 2, c.Len())
}

func TestCalculatorHistory(t *testing.T) {
	c := NewCalculator()
	c.AddHistory("foo")
	c.AddHistory("bar")
	assert.Equal(t, []string{"foo", "bar"}, c.History())
	t.Run("trim", func(t *testing.T) {
		for range 123 {
			c.AddHistory("x")
		}
		assert.Equal(t, MaxArraySize, len(c.History()))
	})
}

func TestCalculatorStack(t *testing.T) {
	c := NewCalculator()

	t.Run("normalize", func(t *testing.T) {
		// pos
		c.PushFloat64(1.1)
		assert.Equal(t, 1.1, c.PopFloat64())
		c.PushFloat64(1.00000001)
		assert.Equal(t, 1.0, c.PopFloat64())
		// neg
		c.PushFloat64(-1.1)
		assert.Equal(t, -1.1, c.PopFloat64())
		c.PushFloat64(-1.00000001)
		assert.Equal(t, -1.0, c.PopFloat64())
	})
	t.Run("trim", func(t *testing.T) {
		for range 123 {
			c.PushInt(1)
		}
		assert.Equal(t, MaxArraySize, len(c.stack))
	})
}

func TestCalculatorRunFn(t *testing.T) {
	tests := []struct {
		cmd     string
		inputs  int
		outputs int
	}{
		{"CLEAR", 2, 0}, // fn()
		{"PI", 0, 1},    // fn() Num
		{"DROP", 1, 0},  // fn(x)
		{"NEG", 1, 1},   // fn(x) Num
		{"SWAP", 2, 2},  // fn(x, y)
		{"ADD", 2, 1},   // fn(x, y) Num
	}

	c := NewCalculator()
	for _, tc := range tests {
		undos := len(c.undo)
		c.Clear()
		for range tc.inputs {
			c.PushInt(123)
		}
		c.Run(tc.cmd)
		assert.Equal(t, tc.outputs, c.Len())
		assert.Equal(t, undos+1, len(c.undo)) // did the undo stack grow?
	}
}

func TestCalculatorRunValid(t *testing.T) {
	// add only works with 2 inputs
	c := NewCalculator()
	assert.Error(t, c.Run("ADD"))
	c.PushInt(123, 456)
	assert.NoError(t, c.Run("ADD"))
}

func TestCalculatorRunHistory(t *testing.T) {
	c := NewCalculator()
	c.PushInt(123, 456)
	c.Run("ADD")
	assert.Equal(t, "123 + 456 = 579", c.History()[0])
}

func TestEnter(t *testing.T) {
	c := NewCalculator()
	c.Enter(decimal.NewFromInt(123), false) // implicit (no undo)
	assert.Equal(t, 123, c.PeekInt())
	assert.Empty(t, c.undo)
	c.Enter(decimal.NewFromInt(456), true) // explicit (undo)
	assert.Equal(t, 456, c.PeekInt())
	assert.NotEmpty(t, c.undo)
}

func TestUndo(t *testing.T) {
	c := NewCalculator()
	c.PushInt(123)
	c.snapshotForUndo()
	c.PushInt(456)
	c.snapshotForUndo()
	c.PushInt(789)

	c.Undo()
	assert.Equal(t, 2, c.Len())
	assert.Equal(t, 456, c.PeekInt())
	c.Undo()
	assert.Equal(t, 1, c.Len())
	assert.Equal(t, 123, c.PeekInt())
}
