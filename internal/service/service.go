package service

import (
	"errors"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

const (
	// ValidSymbols contains all allowed characters for text input
	ValidSymbols = "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ0123456789.,:?'-/()\" "
	// ValidMorseChars contains all allowed characters for Morse code input
	ValidMorseChars = ".- "
)

// Error messages
const (
	ErrEmptyInput      = "empty input"
	ErrMixedInput      = "mixed input: Morse code contains invalid characters"
	ErrInvalidTextChar = "the text contains an invalid character"
	ErrInvalidMorse    = "incorrect Morse code: contains unrecognized characters"
	ErrMorseConversion = "error converting to Morse code"
	ErrTextConversion  = "error converting to text"
)

// TextHandler automatically detects and converts between text and Morse code.
// Returns error for invalid input or conversion failures.
func TextHandler(data string) (string, error) {
	if len(data) == 0 {
		return "", errors.New(ErrEmptyInput)
	}

	data = strings.TrimSpace(data)
	if len(data) == 0 {
		return "", errors.New(ErrEmptyInput)
	}

	isMorse, err := isMorseCode(data)
	if err != nil {
		return "", err
	}

	if isMorse {
		result := morse.ToText(data)
		if result == "" {
			return "", errors.New(ErrInvalidMorse)
		}
		return result, nil
	}

	if err := validateText(data); err != nil {
		return "", err
	}

	result := morse.ToMorse(data)
	if result == "" {
		return "", errors.New(ErrMorseConversion)
	}
	return result, nil
}

// isMorseCode checks if the input is valid Morse code
func isMorseCode(data string) (bool, error) {
	hasMorse := strings.ContainsAny(data, ".-")
	if !hasMorse {
		return false, nil
	}

	for _, r := range data {
		if !strings.ContainsRune(ValidMorseChars, r) {
			return false, errors.New(ErrMixedInput)
		}
	}
	return true, nil
}

// validateText checks if all characters are valid
func validateText(data string) error {
	for _, r := range data {
		if !strings.ContainsRune(ValidSymbols, r) {
			return errors.New(ErrInvalidTextChar + ": " + string(r))
		}
	}
	return nil
}
