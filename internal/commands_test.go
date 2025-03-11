package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEachCommand(t *testing.T) {
	tests := []struct {
		cmd     string
		inputs  []float64
		outputs []float64
	}{
		{"ADD", []float64{3, 5}, []float64{8}},
		{"CLEAR", []float64{1, 2}, []float64{}},
		{"DIV", []float64{8, 2}, []float64{4}},
		{"DROP", []float64{1, 2, 3}, []float64{1, 2}},
		{"DUP", []float64{1, 2, 3}, []float64{1, 2, 3, 3}},
		{"FACT", []float64{5}, []float64{120}},
		{"INV", []float64{2}, []float64{0.5}},
		{"LOG", []float64{10}, []float64{1}},
		{"MOD", []float64{5, 3}, []float64{2}},
		{"MUL", []float64{3, 5}, []float64{15}},
		{"NEG", []float64{3}, []float64{-3}},
		{"PI", []float64{}, []float64{3.1415926536}},
		{"POW", []float64{2, 3}, []float64{8}},
		{"SQRT", []float64{9}, []float64{3}},
		{"SUB", []float64{5, 3}, []float64{2}},
		{"SWAP", []float64{1, 2}, []float64{2, 1}},
	}

	c := NewCalculator()
	for _, tc := range tests {
		t.Run(tc.cmd, func(t *testing.T) {
			c.Clear()
			c.PushFloat64(tc.inputs...)
			testRun(c, tc.cmd)

			outputs := MapV(c.GetStack(), func(x Num) float64 { return x.InexactFloat64() })
			assert.Equal(t, tc.outputs, outputs)
		})
	}
}

func TestCommandUndo(t *testing.T) {
	c := NewCalculator()
	c.PushInt(1)
	c.snapshotForUndo()
	c.PushInt(2)
	testRun(c, "UNDO")
	assert.Equal(t, 1, c.Len())
	assert.Equal(t, 1, c.PopInt())
}

func TestCommandMaps(t *testing.T) {
	assert.Equal(t, "ADD", CommandsByKey["+"].Name)
	assert.Equal(t, "ADD", CommandsByName["ADD"].Name)
}

func TestCommandsValid(t *testing.T) {
	var (
		c = NewCalculator()
	)

	c.PushInt(-1)
	assert.Error(t, validFact(c))
	assert.Error(t, validGt0(c))
	assert.Error(t, validGte0(c))
	assert.NoError(t, validNot0(c))
	c.PushInt(0)
	assert.NoError(t, validFact(c))
	assert.Error(t, validGt0(c))
	assert.NoError(t, validGte0(c))
	assert.Error(t, validNot0(c))
	c.PushInt(1)
	assert.NoError(t, validFact(c))
	assert.NoError(t, validGt0(c))
	assert.NoError(t, validGte0(c))
	assert.NoError(t, validNot0(c))

	// few more errors cases for validFact
	for _, x := range []float64{0.5, 9999} {
		c.PushFloat64(x)
		assert.Error(t, validFact(c))
	}
}

func testRun(c *Calculator, name string) {
	switch fn := CommandsByName[name].fn.(type) {
	case func(*Calculator):
		fn(c)
	case func(*Calculator, Num):
		fn(c, c.Pop())
	case func(*Calculator) Num:
		c.Push(fn(c))
	case func(*Calculator, Num) Num:
		c.Push(fn(c, c.Pop()))
	case func(*Calculator, Num, Num):
		b, a := c.Pop(), c.Pop()
		fn(c, a, b)
	case func(*Calculator, Num, Num) Num:
		b, a := c.Pop(), c.Pop()
		c.Push(fn(c, a, b))
	}
}
