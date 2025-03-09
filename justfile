_default:
  @just --list

init:
  asdf plugin add golang || true
  asdf plugin add nodejs || true
  asdf install
  brew install gifsicle goreleaser watchexec

# go run
run *ARGS:
  @go run ./... {{ARGS}}

# go test
test *ARGS:
  @go test ./... {{ARGS}}

# go test, and watch for changes
test-watch *ARGS:
  @watchexec --watch . --clear=reset just test "{{ARGS}}"

#
# publish
#

demo:
  @clear
  @just _banner "running vhs..."
  @vhs demo.tape
  @cd /tmp && gifsicle -O3 --lossy -o demo-small.gif demo.gif ; ls -lh demo.gif demo-small.gif
  @mv /tmp/demo-small.gif demo.gif
  @just _banner "done"
  @qlmanage -p demo.gif

release:
  @clear
  @just _banner "test..."
  just test
  @just _banner "goreleaser release..."
  goreleaser release --clean
  @just _banner "done"

snapshot:
  @just _banner "goreleaser --snapshot..."
  goreleaser release --snapshot --clean
  @just _banner "done"

#
# internal
#

_banner *ARGS:
  @printf '\e[42;37;1m[%s] %-72s \e[m\n' "$(date +%H:%M:%S)" "{{ARGS}}"
