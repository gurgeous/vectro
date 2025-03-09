//
// global constants and styles
//

package internal

import "github.com/charmbracelet/lipgloss"

var (
	// how many lines of the stack should we show?
	StackSize = 6
	// truncate history/stack after this to avoid memory issues
	MaxArraySize = 50
	// size of undo stack
	UndoSize = 50
	// how many digits of precision?
	Precision = 10

	LG = lipgloss.NewStyle() // just to make things easy

	// panes & borders
	PaneStyle        = LG.Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(Blue400)
	BorderStyle      = LG.Foreground(PaneStyle.GetBorderTopForeground())
	BorderTitleStyle = LG.Foreground(Gray400)

	// stack
	ErrorStyle  = LG.Foreground(White).Background(Red600)
	SayStyle    = LG.Foreground(White).Background(Green500)
	StackStyle  = PaneStyle.PaddingTop(2).PaddingBottom(1)
	IndexStyle  = LG.Foreground(Gray600)
	CursorStyle = LG.Foreground(lipgloss.AdaptiveColor{Light: string(Yellow500), Dark: string(Yellow300)})

	// help
	HelpStyle    = LG.Foreground(Gray700)
	HelpKeyStyle = LG.Foreground(Green500).Bold(true)

	// status
	StatusStyle = LG.
			Foreground(White).
			Background(Blue600).
			Bold(true).
			AlignHorizontal(lipgloss.Center).
			Padding(0, 1)

	// vhs banner
	BannerStyle = LG.
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Foreground(Yellow400).
			Bold(true).
			Padding(0, 10)

	// cramped
	CrampedStyle = LG.
			Foreground(Gray400).
			Background(Violet900).
			Bold(true).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			Padding(0, 3)

	// same len as StackSize, for simplicity
	GradientColors = []lipgloss.TerminalColor{
		lipgloss.AdaptiveColor{Light: string(Gray200), Dark: string(Gray600)},
		lipgloss.AdaptiveColor{Light: string(Gray300), Dark: string(Gray500)},
		lipgloss.AdaptiveColor{Light: string(Gray400), Dark: string(Gray400)},
		lipgloss.AdaptiveColor{Light: string(Gray500), Dark: string(Gray300)},
		lipgloss.AdaptiveColor{Light: string(Gray600), Dark: string(Gray200)},
		lipgloss.AdaptiveColor{Light: string(Black), Dark: string(White)},
	}
	GradientStyles = Foregrounds(GradientColors)

	TitleColors = []lipgloss.TerminalColor{
		Red600, Yellow600, Blue600, Green600, Orange600, Purple500,
	}
	TitleStyles = Foregrounds(TitleColors)
)
