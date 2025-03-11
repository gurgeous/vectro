package internal

import (
	"math"
	"os"
	"regexp"
	"slices"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/shopspring/decimal"
)

//
// Num and friends
//

type Num = decimal.Decimal

var (
	// decimal constants
	Half    = decimal.NewFromFloat(0.5)
	Ln10    = decimal.NewFromFloat(math.Log(10))
	One     = decimal.NewFromFloat(1)
	Pi      = decimal.NewFromFloat(math.Pi)
	Epsilon = decimal.NewFromFloat(1e-6)
)

// is this Num an int?
func IsInt(value Num) bool {
	return value.Sub(value.Round(0)).Abs().LessThan(Epsilon)
}

// if x seems to be an Int, round it
func Normalize(x Num) Num {
	x = x.Round(int32(Precision)) //nolint:gosec
	if IsInt(x) {
		x = x.Round(0)
	}
	return x
}

// x!
func Factorial(x Num) Num {
	if x.IsNegative() {
		panic("factorial of negative number")
	}
	if !x.IsInteger() {
		panic("factorial of non-integer")
	}
	var acc = One
	for ii := One; ii.Cmp(x) <= 0; ii = ii.Add(One) {
		acc = acc.Mul(ii)
	}
	return acc
}

// ln(x)
func Ln(x Num) Num {
	y, _ := x.Ln(10)
	return y
}

func Pow(x, y Num) Num {
	z, _ := x.PowWithPrecision(y, int32(Precision)) //nolint:gosec
	return z
}

//
// files
//

// does file exist?
func FileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

//
// array generics
//

func Dup[E any](s []E) []E {
	return append([]E(nil), s...)
}

func Filter[E any](s []E, fn func(E) bool) []E {
	result := []E{}
	for _, e := range s {
		if fn(e) {
			result = append(result, e)
		}
	}
	return result
}

// map from one array to another
func Map[E, F any](s []E, fn func(E) F) []F {
	return MapWithIndex(s, func(_ int, e E) F {
		return fn(e)
	})
}

func Last[E any](s []E) E {
	return s[len(s)-1]
}

// map from one array to another
func MapWithIndex[E, F any](s []E, fn func(int, E) F) []F {
	result := make([]F, len(s))
	for ii, e := range s {
		result[ii] = fn(ii, e)
	}
	return result
}

func Pop[E any](s []E) (E, []E) {
	return s[len(s)-1], s[:len(s)-1]
}

func Push[E any](s []E, values ...E) []E {
	return append(s, values...)
}

func Repeat[E any](v E, len int) []E {
	return Map(Sequence(len), func(_ int) E { return v })
}

// return reversed copy of array
func Reverse[E any](s []E) []E {
	result := make([]E, len(s))
	copy(result, s)
	slices.Reverse(result)
	return result
}

func Sequence(len int) []int {
	result := make([]int, 0, len)
	for ii := range len {
		result = append(result, ii)
	}
	return result
}

func Shift[E any](s []E) (E, []E) {
	return s[0], s[1:]
}

// Truncate an array, but remove stuff from the start
func TruncateStart[E any](s []E, maxLen int) []E {
	if len(s) > maxLen {
		s = s[len(s)-maxLen:]
	}
	return s
}

func Truncate[E any](s []E, maxLen int) []E {
	if len(s) > maxLen {
		s = s[:maxLen]
	}
	return s
}

//
// styling
//

// look for **xxx**, apply a style
func StyleBetweenStars(str string, style lipgloss.Style) string {
	var re = regexp.MustCompile(`(?s)\*\*(.*?)\*\*`)
	return re.ReplaceAllStringFunc(str, func(s string) string {
		return style.Render(s[2 : len(s)-2])
	})
}

// fit lines into w/h of style. Truncates both horizontally and vertically.
func ClipLines(lines []string, style lipgloss.Style) []string {
	w := style.GetWidth() - style.GetHorizontalPadding()
	h := style.GetHeight() - style.GetVerticalPadding()
	if w <= 0 || h <= 0 {
		return nil
	}
	return Map(Truncate(lines, h), func(s string) string {
		return ansi.Truncate(s, w, "...")
	})
}
