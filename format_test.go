package ulog

import (
	"fmt"
	"runtime"
	"testing"
)

func TestColorFiexdhead(t *testing.T) {
	c := callerInfoBuilder(runtime.Caller(0))
	c.level = InfoLevel
	head, _ := colorFiexdhead(c)
	fmt.Println(head)
}
