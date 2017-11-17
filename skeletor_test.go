package main

import (
	"skeletor/utils"
	"testing"
)

func BenchmarkHashPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		utils.HashPassword("password")
	}

}
