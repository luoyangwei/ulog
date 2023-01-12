package ulog

import (
	"fmt"
	"testing"
)

func TestColorFiexdhead(t *testing.T) {

	head, _ := colorFiexdhead("info")
	fmt.Println(head)
}
