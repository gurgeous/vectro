Output "/tmp/demo.gif"

# see https://github.com/charmbracelet/vhs
Set BorderRadius 8
Set FontSize 15
Set FontFamily "Roboto"
Set Margin 10
Set Padding 0
Set Width 1000
Set Height 600
Set WindowBar Colorful
Set Theme "Catppuccin Frappe"

# boot app in VHS mode
Hide
  Env VHS "1"
  Type "just run" Enter Sleep 1
Show

#
# Vectro's VHS mode has a few special commands:
#
# [some text here] - append line to banner
# ctrl-e           - "erase" banner
#

# hi
Sleep 2
Type "[Welcome to Vectro, the rpn calculator." Sleep 1
Type " Let's calculate the golden ratio.][]" Sleep 1
Ctrl+E Sleep 1

# golden ratio math
Type "[The golden ratio is][(1 + sqrt(5)) / 2][]" Sleep 1
Type "[1 enter]" Sleep 1
  Type "1" Sleep 0.5
  Enter Sleep 1
Type "[5 sqrt (@ key is sqrt)]" Sleep 1
  Type "5" Sleep 0.5
  Type "@" Sleep 1
Type "[+]" Sleep 1
  Type "+" Sleep 1
Type "[2 /]" Sleep 1
  Type "2" Sleep 0.5
  Type "/" Sleep 1
Ctrl+E Sleep 1

# bye
Type "[That's the golden ratio.]" Sleep 1
Type "[Have fun!]" Sleep 3
