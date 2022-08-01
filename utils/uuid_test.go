package utils

import "testing"

func TestUUID(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(GenerateRandomID())
	}
}

func TestGenerateSixFigure(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(GenerateSixFigure())
	}
}

func TestRandomHalfRate(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(RandomHalfRate())
	}
}
