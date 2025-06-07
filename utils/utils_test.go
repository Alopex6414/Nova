package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomAlphabetAndNumber(t *testing.T) {
	s := RandomAlphabetAndNumber(5)
	assert.Equal(t, 5, len(s))
}

func BenchmarkRandomAlphabetAndNumber(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := RandomAlphabetAndNumber(5)
		assert.Equal(b, 5, len(s))
	}
}

func BenchmarkRandomAlphabetAndNumberParallel(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s := RandomAlphabetAndNumber(5)
			assert.Equal(b, 5, len(s))
		}
	})
}

func TestRandomAlphabet(t *testing.T) {
	s := RandomAlphabet(5)
	assert.Equal(t, 5, len(s))
}

func BenchmarkRandomAlphabet(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := RandomAlphabet(5)
		assert.Equal(b, 5, len(s))
	}
}

func BenchmarkRandomAlphabetParallel(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s := RandomAlphabet(5)
			assert.Equal(b, 5, len(s))
		}
	})
}

func TestRandomNumber(t *testing.T) {
	s := RandomNumber(5)
	assert.Equal(t, 5, len(s))
}

func BenchmarkRandomNumber(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := RandomNumber(5)
		assert.Equal(b, 5, len(s))
	}
}

func BenchmarkRandomNumberParallel(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s := RandomNumber(5)
			assert.Equal(b, 5, len(s))
		}
	})
}
