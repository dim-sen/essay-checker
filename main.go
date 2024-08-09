package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"unicode"

	_ "github.com/lib/pq"
)

// NormalizeText normalizes the text by converting it to lowercase and removing non-alphanumeric characters
func NormalizeText(text string) string {
	var normalized strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			normalized.WriteRune(unicode.ToLower(r))
		}
	}
	return normalized.String()
}

// Tokenize splits the text into a slice of words
func Tokenize(text string) []string {
	return strings.Fields(text)
}

// WordFrequency calculates the frequency of each word in the slice
func WordFrequency(words []string) map[string]int {
	freq := make(map[string]int)
	for _, word := range words {
		freq[word]++
	}
	return freq
}

// CosineSimilarity calculates the cosine similarity between two frequency maps
func CosineSimilarity(freqMap1, freqMap2 map[string]int) float64 {
	dotProduct := 0
	magnitude1 := 0
	magnitude2 := 0

	for word, freq1 := range freqMap1 {
		freq2 := freqMap2[word]
		dotProduct += freq1 * freq2
		magnitude1 += freq1 * freq1
	}

	for _, freq2 := range freqMap2 {
		magnitude2 += freq2 * freq2
	}

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0.0
	}

	return float64(dotProduct) / (math.Sqrt(float64(magnitude1)) * math.Sqrt(float64(magnitude2)))
}

// ProcessAnswer processes a single answer and calculates its similarity
func ProcessAnswer(id int, uncheckAnswer, correctAnswer string, wg *sync.WaitGroup) {
	defer wg.Done()

	answerNormalized := NormalizeText(uncheckAnswer)
	correctNormalized := NormalizeText(correctAnswer)

	tokenizeUncheck := Tokenize(answerNormalized)
	tokenizeCorrect := Tokenize(correctNormalized)

	uncheckFrequency := WordFrequency(tokenizeUncheck)
	correctFrequency := WordFrequency(tokenizeCorrect)

	similarity := CosineSimilarity(uncheckFrequency, correctFrequency)
	fmt.Printf("ID: %d, Cosine Similarity: %f\n", id, similarity)
}

func main() {
	connStr := "user=admin password=admin dbname=db_essay sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to database")
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, student_answer, correct_answer FROM answers")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Rows:", rows)
	}
	defer rows.Close()

	var wg sync.WaitGroup
	for rows.Next() {
		var id int
		var uncheckAnswer, correctAnswer string
		err := rows.Scan(&id, &uncheckAnswer, &correctAnswer)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		go ProcessAnswer(id, uncheckAnswer, correctAnswer, &wg)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
