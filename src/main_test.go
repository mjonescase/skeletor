package main

import (
	"testing"
)

func TestSucceed(t *testing.T) {
}

func TestFail(t *testing.T) {
	t.Errorf("wrong, Wrong, WRONG!!!")
}
