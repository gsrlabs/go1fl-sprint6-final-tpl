package service

import (
	"errors"
	"strings"
	
	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

const (
	ValidSymbols    = "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ0123456789.,:?'-/()\""
	ValidMorseChars = ".- "
)

// TextHandler automatically detects whether a string is Morse code or plain text,
// and converts it to the opposite format.
func TextHandler(data string) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	isMorse, err := isMorseCode(data)
	if err != nil {
		return "", err
	}

	if isMorse {
		result := morse.ToText(data)
		if result == "" {
			return "", errors.New("некорректный код Морзе: содержит нераспознанные символы")
		}
		return result, nil
	}

	if err := validateText(data); err != nil {
		return "", err
	}

	result := morse.ToMorse(data)
	if result == "" {
		return "", errors.New("ошибка конвертации в код Морзе")
	}
	return result, nil
}


func isMorseCode(data string) (bool, error) {
	hasMorse := strings.ContainsAny(data, ".-")
	if !hasMorse {
		return false, nil
	}

	for _, r := range data {
		if !strings.ContainsRune(ValidMorseChars, r) {
			return false, errors.New("смешанный ввод: код Морзе содержит недопустимые символы")
		}
	}
	return true, nil
}


func validateText(data string) error {
	for _, r := range data {
		if r != ' ' && !strings.ContainsRune(ValidSymbols, r) {
			return errors.New("текст содержит недопустимый символ: " + string(r))
		}
	}
	return nil
}
