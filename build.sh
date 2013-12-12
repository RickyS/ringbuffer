# Source me.  Bash function.

alias test='go test -bench=. -benchmem -cover'
function rr () {
    echo this no longer compiles.
    go install -v ringbuffer  && go run runringbuffer.go
}
declare -f rr
