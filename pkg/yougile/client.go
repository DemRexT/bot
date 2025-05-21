package yougile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lotBot/pkg/embedlog"
	"net/http"
)

type Config struct {
	Login    string `toml:"Login"`
	Password string `toml:"Password"`
	Token    string `toml:"Token"`
}

type YougileClient struct {
	embedlog.Logger
	cfg Config
}

func NewYougileClient(logger embedlog.Logger, cfg Config) *YougileClient {
	return &YougileClient{Logger: logger, cfg: cfg}
}

const baseURL = "https://ru.yougile.com/api-v2"

type TaskPayload struct {
	Title       string `json:"title"`
	ColumnID    string `json:"columnId"`
	Description string `json:"description"`
	Archived    bool   `json:"archived"`
	Completed   bool   `json:"completed"`
}
type TaskResponse struct {
	ID string `json:"id"`
}

func (c *YougileClient) CreateTask(task TaskPayload) (string, error) {
	url := fmt.Sprintf("%s/tasks", baseURL)

	data, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("ошибка сериализации данных: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("ошибка от API: статус %s", resp.Status)
	}

	var result TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("ошибка декодирования ответа: %w", err)
	}

	c.Printf("Задача создана. ID: %s", result.ID)
	return result.ID, nil
}

func (c *YougileClient) GetUserByID(userID string) ([]byte, error) {
	url := fmt.Sprintf("%s/users/%s", baseURL, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка от API: статус %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}

	return body, nil
}

func (c *YougileClient) GetTaskByID(taskID string) ([]byte, error) {
	url := fmt.Sprintf("%s/tasks/%s", baseURL, taskID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка от API: статус %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения тела ответа: %w", err)
	}

	return body, nil
}
