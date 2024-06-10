package sender

import (
	"api-save-fact/model"
	"bytes"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
)

//go:generate mockery --name Sender
type Sender interface {
	Send(data model.Data) error
}

type S struct {
	cfg    model.Config
	buffer chan model.Data
}

func NewSender(cfg model.Config, buffer chan model.Data) Sender {
	return &S{cfg: cfg, buffer: buffer}
}

// Send отправляет данные на сервер
func (s *S) Send(data model.Data) error {
	// Создаем тело запроса
	var requestBody bytes.Buffer

	// Создаем multipart writer
	writer := multipart.NewWriter(&requestBody)

	// Заполняем тело запроса данными
	if err := writeFiled(writer, data); err != nil {
		return fmt.Errorf("error writing field: %v", err)
	}

	// Закрываем writer
	if err := writer.Close(); err != nil {
		return fmt.Errorf("error closing writer: %v", err)
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", s.cfg.Href, &requestBody)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if s.cfg.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.cfg.BearerToken)
	}

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		log.Println("data sent successfully")
		return nil
	case http.StatusBadRequest, http.StatusBadGateway:
		return fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

// writeFiled записывает данные в multipart writer
func writeFiled(writer *multipart.Writer, data model.Data) error {
	if err := writer.WriteField("period_start", data.PeriodStart); err != nil {
		return err
	}
	if err := writer.WriteField("period_end", data.PeriodEnd); err != nil {
		return err
	}
	if err := writer.WriteField("period_key", data.PeriodKey); err != nil {
		return err
	}
	if err := writer.WriteField("indicator_to_mo_id", data.IndicatorToMoId); err != nil {
		return err
	}
	if err := writer.WriteField("indicator_to_mo_fact_id", data.IndicatorToMoFactId); err != nil {
		return err
	}
	if err := writer.WriteField("value", data.Value); err != nil {
		return err
	}
	if err := writer.WriteField("fact_time", data.FactTime); err != nil {
		return err
	}
	if err := writer.WriteField("is_plan", data.IsPlan); err != nil {
		return err
	}
	if err := writer.WriteField("auth_user_id", data.AuthUserID); err != nil {
		return err
	}
	if err := writer.WriteField("comment", data.Comment); err != nil {
		return err
	}

	return nil
}
