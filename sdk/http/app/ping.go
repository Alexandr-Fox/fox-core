package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Alexandr-Fox/fox-core/internal/controllers"
	"net/http"
)

func Ping() (*controllers.App, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/app.ping", fqdn), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // Важно: закрыть тело ответа

	// Проверяем статус ответа
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(resp.Status)
	}

	// Декодируем JSON в структуру
	var result controllers.App
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
