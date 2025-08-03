package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Storage[T any] struct {
	filename string
}

func NewStorage[T any](filename string) *Storage[T] {
	return &Storage[T]{filename: filename}
}

func (s *Storage[T]) Save(data T) error {
	fileData, err := json.MarshalIndent(data, "", "  ")

	if err != nil {
		return fmt.Errorf("error marshaling data to JSON: %w", err)
	}

	return os.WriteFile(s.filename, fileData, 0644)
}

func (s *Storage[T]) Load() (T, error) {
	var data T

	// Check if the file exists, if not create an empty one
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		emptyFile, err := os.Create(s.filename)
		if err != nil {
			return data, fmt.Errorf("error creating file: %w", err)
		}

		emptyFile.Close()
	}

	// Read the file content
	fileData, err := os.ReadFile(s.filename)
	if err != nil {
		return data, fmt.Errorf("error reading file: %w", err)

	}

	// Check if the file is empty
	if len(fileData) == 0 {
		return data, nil
	}

	if err := json.Unmarshal(fileData, &data); err != nil {
		return data, fmt.Errorf("error unmarshaling JSON data: %w", err)
	}

	return data, nil
}
