package server

import (
	"api-save-fact/model"
	"context"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type Server struct {
	server *http.Server
	// todo Buffer
}

// NewServer запускает HTTP-сервер с использованием указанной конфигурации и буфера
func NewServer(cfg model.Config, buffer chan<- model.Data) *Server {
	// Устанавливаем обработчик для указанного в конфигурации маршрута
	http.HandleFunc("/_api/facts/save_fact", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		Save(w, r, buffer)
	})

	return &Server{server: &http.Server{Addr: cfg.Host + ":" + cfg.Port}}
}

func (s *Server) Start() error {
	log.Println("Starting server...")
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	// todo buffer close
	return s.server.Shutdown(ctx)
}

func (s *Server) Close() error {
	log.Println("Closing server...")
	// todo buffer close
	return s.server.Close()
}

// Save сохраняет данные в буфер
func Save(w http.ResponseWriter, r *http.Request, buffer chan<- model.Data) {
	var data model.Data

	// Получаем данные из запроса
	data.PeriodStart = r.FormValue("period_start")
	data.PeriodEnd = r.FormValue("period_end")
	data.PeriodKey = r.FormValue("period_key")
	data.IndicatorToMoId = r.FormValue("indicator_to_mo_id")
	data.IndicatorToMoFactId = r.FormValue("indicator_to_mo_fact_id")
	data.Value = r.FormValue("value")
	data.FactTime = r.FormValue("fact_time")
	data.IsPlan = r.FormValue("is_plan")
	data.AuthUserID = r.FormValue("auth_user_id")
	data.Comment = r.FormValue("comment")

	// Проверяем наличие данных
	if err := validator.New().Struct(data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	// Добавляем данные в буфер
	buffer <- data

	// Отправляем ответ
	w.WriteHeader(http.StatusCreated)
}
