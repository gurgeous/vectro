package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/gurgeous/vectro/internal"
)

//
// Helper for rendering a `pane`, which has a border with title. The top border
// is called the `header` and we draw it manually. Lipgloss does the rest of it.
// Remember that lipgloss width/height does not include margins or borders.
//
// space               space
//   ↓                   ↓
// ╭─ the quick brown fox ────────────────────╮
// ↑↑       ↑                     ↑           ↑
// tl   	↑                     ↑           tr
//          ↑                     ↑
//        title                 filler
//

func RenderPane(style lipgloss.Style, titleInit string, text string) string {
	if style.GetWidth() == 0 || style.GetHeight() == 0 {
		return ""
	}

	var (
		ch = internal.BorderStyle.Render(style.GetBorderStyle().Top)
		tl = internal.BorderStyle.Render(style.GetBorderStyle().TopLeft) + ch
		tr = internal.BorderStyle.Render(style.GetBorderStyle().TopRight)

		paneWidth = style.GetWidth() + style.GetHorizontalBorderSize()
		spaces    = 2
	)

	// calculate truncate title and filler
	headerInner := paneWidth - lipgloss.Width(tl+tr)
	title := fmt.Sprintf(" %s ", ansi.Truncate(titleInit, headerInner-spaces, ""))
	filler := strings.Repeat(ch, headerInner-lipgloss.Width(title))

	// color title if necessary
	if len(title) == lipgloss.Width(title) {
		title = internal.BorderTitleStyle.Render(title)
	}

	header := lipgloss.NewStyle().
		Margin(style.GetMargin()).
		MarginBottom(0).
		Render(fmt.Sprintf("%s%s%s%s", tl, title, filler, tr))
	body := style.BorderTop(false).UnsetMarginTop().Render(text)

	return lipgloss.JoinVertical(lipgloss.Top, header, body)
}
