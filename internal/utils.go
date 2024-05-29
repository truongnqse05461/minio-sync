package internal

import (
	"encoding/csv"
	"log"
	"os"
)

// ReadBucketName read bucket from .csv file
func ReadBucketName(f string) []string {
	file, err := os.Open(f)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read CSV file: %s", err)
	}
	var ids []string
	// Iterate through the records
	for i, record := range records {
		// Skip the header row
		if i == 0 {
			continue
		}
		ids = append(ids, record[0])
	}
	return ids
}
