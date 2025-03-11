_default:
  @just --list

check:
  @just _banner check: lint... ; just lint
  @just _banner check: test... ; just test
  @just _banner check: done


init:
  asdf plugin add golang || true
  asdf plugin add nodejs || true
  asdf install
  brew install gifsicle golangci goreleaser watchexec

lint *ARGS:
  @golangci-lint run ./... {{ARGS}}

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
  @if [ -z "${GITHUB_TOKEN:-}" ]; then just _fatal "GITHUB_TOKEN is required" ; fi
  @just _banner "test..."
  just test
  @just _banner "goreleaser release..."
  goreleaser release --clean
  @just _banner "done"

snapshot:
  @clear
  @just _banner "test..."
  just test
  @just _banner "goreleaser --snapshot..."
  goreleaser release --snapshot --clean
  @just _banner "done"

#
# internal
#

_banner *ARGS:
  @printf '\e[42;37;1m[%s] %-72s \e[m\n' "$(date +%H:%M:%S)" "{{ARGS}}"
_fatal *ARGS:
  @printf '\e[41;37;1m[%s] %-72s \e[m\n' "$(date +%H:%M:%S)" "{{ARGS}}"
  @exit 1
