package main

import (
	"testing"
)

// Example text for benchmarking
var studentAnswer = "Paris is the capital city of France."
var correctAnswer = "Paris is the capital of France."

func BenchmarkNormalizeText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NormalizeText(studentAnswer)
	}
}

func BenchmarkTokenize(b *testing.B) {
	normalizedText := NormalizeText(studentAnswer)
	b.ResetTimer() // Reset timer to exclude normalization time
	for i := 0; i < b.N; i++ {
		Tokenize(normalizedText)
	}
}

//func BenchmarkWordFrequency(b *testing.B) {
//	tokens := Tokenize(NormalizeText(studentAnswer))
//	b.ResetTimer() // Reset timer to exclude tokenization time
//	for i := 0; i < b.N; i++ {
//		WordFrequency(tokens)
//	}
//}

//func BenchmarkCosineSimilarity(b *testing.B) {
//	freqMap1 := WordFrequency(Tokenize(NormalizeText(studentAnswer)))
//	freqMap2 := WordFrequency(Tokenize(NormalizeText(correctAnswer)))
//	b.ResetTimer() // Reset timer to exclude preparation time
//	for i := 0; i < b.N; i++ {
//		CosineSimilarity(freqMap1, freqMap2)
//	}
//}

func BenchmarkJaccardSimilarity(b *testing.B) {
	set1 := Tokenize(NormalizeText(studentAnswer))
	set2 := Tokenize(NormalizeText(correctAnswer))
	b.ResetTimer() // Reset timer to exclude preparation time
	for i := 0; i < b.N; i++ {
		JaccardSimilarity(set1, set2)
	}
}
