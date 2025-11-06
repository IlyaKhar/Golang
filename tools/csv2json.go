package tools

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"os"
)

func Csv2Json(csvPath string, jsonPath string) error {
	csvFile, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)

	headers, err := csvReader.Read() // читает первую строку CSV
	if err != nil {
		return err
	}

	jsonFile, err := os.Create(jsonPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var result []map[string]string // массив объектов для JSON

	for {
		row, err := csvReader.Read() // читаем строку
		if err == io.EOF {           // если файл кончился
			break
		}
		if err != nil {
			return err
		}

		// Создаём объект: сопоставляем заголовок → значение
		obj := make(map[string]string)
		for i, header := range headers {
			if i < len(row) {
				obj[header] = row[i] // "name" → "Alice"
			} else {
				obj[header] = "" // если значений меньше заголовков
			}
		}

		result = append(result, obj) // добавляем в массив
	}

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ") // красивое форматирование
	if err := encoder.Encode(result); err != nil {
		return err
	}
	return nil
}
