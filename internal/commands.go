package internal

import (
	"errors"

	"github.com/atotto/clipboard"
	"github.com/shopspring/decimal"
)

type Command struct {
	Name  string
	key   string
	fn    interface{}
	fmt   string
	valid func(*Calculator) error
}

//
// the commands
//

var Commands = []Command{
	{Name: "ADD", key: "+", fn: add, fmt: "%s + %s = %s"},
	{Name: "CLEAR", key: "esc", fn: clear},
	{Name: "DIV", key: "/", fn: div, valid: validNot0, fmt: "%s / %s = %s"},
	{Name: "DROP", fn: drop},
	{Name: "DUP", key: "xxx", fn: dup},
	{Name: "FACT", key: "!", fn: fact, valid: validFact, fmt: "%s! = %s"},
	{Name: "INV", key: "i", fn: inv, fmt: "1 / %s = %s"},
	{Name: "LN", fn: ln, valid: validGt0, fmt: "ln(%s) = %s"}, // bad key, don't do it
	{Name: "LOG", key: "l", fn: log, valid: validGt0, fmt: "log(%s) = %s"},
	{Name: "MOD", key: "%", fn: mod, fmt: "%s mod %s = %s"},
	{Name: "MUL", key: "*", fn: mul, fmt: "%s * %s = %s"},
	{Name: "NEG", key: "n", fn: neg},
	{Name: "PI", key: "p", fn: pi},
	{Name: "POW", key: "^", fn: pow, fmt: "%s ^ %s = %s"},
	{Name: "SQRT", key: "@", fn: sqrt, valid: validGte0, fmt: "sqrt(%s) = %s"},
	{Name: "SUB", key: "-", fn: sub, fmt: "%s - %s = %s"},
	{Name: "SWAP", key: "s", fn: swap},
	{Name: "YANK", key: "y", fn: yank},
	{Name: "UNDO", key: "z", fn: undo, valid: validUndo},
}

var CommandsByKey = map[string]Command{}
var CommandsByName = map[string]Command{}

// these are sometimes run directly
const (
	DROP = "DROP"
	DUP  = "DUP"
	NEG  = "NEG"
	UNDO = "UNDO"
	YANK = "YANK"
)

func init() {
	for _, cmd := range Commands {
		if cmd.key != "" {
			CommandsByKey[cmd.key] = cmd
		}
		CommandsByName[cmd.Name] = cmd
	}
}

//
// commands
//

func add(c *Calculator, a, b Num) Num { return a.Add(b) }
func clear(c *Calculator)             { c.Clear() }
func div(c *Calculator, a, b Num) Num { return a.Div(b) }
func drop(c *Calculator, a Num)       { /* nop */ }
func dup(c *Calculator, a Num)        { c.Push(a, a) }
func fact(c *Calculator, a Num) Num   { return Factorial(a) }
func swap(c *Calculator, a, b Num)    { c.Push(b, a) }
func inv(c *Calculator, a Num) Num    { return One.Div(a) }
func ln(c *Calculator, a Num) Num     { return Ln(a) }
func log(c *Calculator, a Num) Num    { return Ln(a).Div(Ln10) }
func mod(c *Calculator, a, b Num) Num { return a.Mod(b) }
func mul(c *Calculator, a, b Num) Num { return a.Mul(b) }
func neg(c *Calculator, a Num) Num    { return a.Neg() }
func pi(c *Calculator) Num            { return Pi }
func pow(c *Calculator, a, b Num) Num { return Pow(a, b) }
func sqrt(c *Calculator, a Num) Num   { return Pow(a, Half) }
func sub(c *Calculator, a, b Num) Num { return a.Sub(b) }

// more complicated
func yank(c *Calculator, a Num) {
	c.Push(a)
	clipboard.WriteAll(a.String())
}

func undo(c *Calculator) {
	c.Undo()
}

//
// helpers
//

func validFact(c *Calculator) error {
	a := c.Peek()
	if a.IsNegative() || !IsInt(a) {
		return errors.New("not a positive int")
	}
	if a.GreaterThan(decimal.NewFromFloat(100)) {
		return errors.New("too large")
	}
	return nil
}
func validGt0(c *Calculator) error {
	if !c.Peek().IsPositive() {
		return errors.New("not positive")
	}
	return nil
}
func validGte0(c *Calculator) error {
	if c.Peek().IsNegative() {
		return errors.New("not positive")
	}
	return nil
}
func validNot0(c *Calculator) error {
	if c.Peek().IsZero() {
		return errors.New("divide by zero")
	}
	return nil
}

func validUndo(c *Calculator) error {
	if len(c.undo) == 0 {
		return errors.New("nothing to undo")
	}
	return nil
}
