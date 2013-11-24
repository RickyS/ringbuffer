#! /bin/bash
# Source me.

function rr () {
    go install -v ringbuffer  && go run ringbuffer.go
}
