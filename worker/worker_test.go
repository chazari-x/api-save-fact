package worker

import (
	"api-save-fact/mocks"
	"api-save-fact/model"
	"github.com/tjarratt/babble"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestS_send(t *testing.T) {
	// Создание отправителя
	sender := mocks.NewSender(t)

	// Создание воркера
	w := &Worker{
		Buffer: make(chan model.Data, 1000),
		Sender: sender,
	}

	// Запуск воркера
	w.Start()

	// Генерация данных
	babbler := babble.NewBabbler()
	babbler.Separator = " "

	for i := range 3000 {
		babbler.Count = rand.Intn(3) + 3
		start := time.Now().AddDate(rand.Intn(5), rand.Intn(12), rand.Intn(31))
		end := start.AddDate(rand.Intn(5), rand.Intn(12), rand.Intn(31))
		data := model.Data{
			PeriodStart:         start.Format("2006-01-02"),
			PeriodEnd:           end.Format("2006-01-02"),
			PeriodKey:           "month",
			IndicatorToMoId:     strconv.Itoa(i + 1),
			IndicatorToMoFactId: strconv.Itoa(i + 1),
			Value:               strconv.Itoa(rand.Intn(9)),
			FactTime:            end.Format("2006-01-02"),
			IsPlan:              strconv.Itoa(rand.Intn(2)),
			AuthUserID:          strconv.Itoa(rand.Intn(30) + 10),
			Comment:             babbler.Babble(),
		}

		sender.On("Send", data).Return(nil)

		w.Buffer <- data
	}

	// Завершение работы воркера, с ожиданием завершения отправки данных
	w.Shutdown()
}
