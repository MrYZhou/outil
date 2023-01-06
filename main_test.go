package main

import (
	"testing"

	. "github.com/MrYZhou/outil/command"
)
func TestRun(t *testing.T) {
	Run(".","docker stats")
}