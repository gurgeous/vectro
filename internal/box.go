package internal

import (
	"github.com/charmbracelet/lipgloss"
)

//
// Immutable geometric box. Remember that lipgloss width/height does not include
// margins or borders.
//

type Box struct {
	width, height int
}

func NewBox(width, height int) Box {
	return Box{width: width, height: height}
}

//
// return a copy of the box with new width/height

func (box Box) Width(v int) Box {
	return Box{width: v, height: box.height}
}

func (box Box) Height(v int) Box {
	return Box{width: box.width, height: v}
}

//
// accessors
//

func (box Box) GetWidth() int {
	return box.width
}
func (box Box) GetHeight() int {
	return box.height
}
func (box Box) Size() (int, int) {
	return box.width, box.height
}

//
// divide into 2 rows/cols
//

func (box Box) Rows() (Box, Box) {
	return box.CutTop(box.height / 2)
}

func (box Box) Cols() (Box, Box) {
	return box.CutLeft(box.width / 2)
}

//
// cut pieces, return two boxes
//

func (box Box) CutLeft(v int) (Box, Box) {
	return box.Width(v), box.Width(box.width - v)
}

func (box Box) CutRight(v int) (Box, Box) {
	return box.CutLeft(box.width - v)
}

func (box Box) CutTop(v int) (Box, Box) {
	return box.Height(v), box.Height(box.height - v)
}

func (box Box) CutBottom(v int) (Box, Box) {
	return box.CutTop(box.height - v)
}

//
// apply Box size to a lipgloss.Style
//

func (box Box) Apply(style lipgloss.Style) lipgloss.Style {
	width, height := box.width, box.height
	width -= style.GetHorizontalBorderSize() + style.GetHorizontalMargins()
	height -= style.GetVerticalBorderSize() + style.GetVerticalMargins()
	width, height = max(0, width), max(0, height)
	return style.Width(width).Height(height)
}

func (box Box) Render(style lipgloss.Style, str string) string {
	return box.Apply(style).Render(str)
}
