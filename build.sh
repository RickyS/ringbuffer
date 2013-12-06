# Source me.  Bash function.

function rr () {
    go install -v ringbuffer  && go run runringbuffer.go
}
