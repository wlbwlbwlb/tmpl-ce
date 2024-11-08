package util

import (
	"os"
	"testing"
)

func TestPath(t *testing.T) {
	t.Log(os.Getenv("TMPL")) //env设置完后，shell或ide需要重启
}
