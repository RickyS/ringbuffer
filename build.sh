# Source me.  Bash function.

alias gt='go test -bench=. -benchmem -cover'
function cover () {
    go test -coverprofile=c.out && go tool cover -html=c.out
}

function rr () {
    echo this no longer compiles.
    go install -v ringbuffer  && go run runringbuffer.go
}
declare -f rr
declare -f cover
