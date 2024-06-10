package worker

import (
	"api-save-fact/model"
	"api-save-fact/sender"
	"log"
	"sync"
	"time"
)

type Node struct {
	Data model.Data
	Sent bool
	Next *Node
}

type LinkedList struct {
	First *Node
	Last  *Node
}

type Worker struct {
	Buffer         chan model.Data
	Sender         sender.Sender
	linkedList     *LinkedList
	readerCancel   bool
	readerShutdown bool
	wg             sync.WaitGroup
}

func NewWorker(buffer chan model.Data, cfg model.Config) *Worker {
	return &Worker{
		Buffer: buffer,
		Sender: sender.NewSender(cfg, buffer),
	}
}

func (w *Worker) Start() {
	go w.startReader()
	go w.startBalancer()
}

// startReader запускает чтение данных из буфера
func (w *Worker) startReader() {
	log.Println("Starting reader...")

	w.wg.Add(1)

	// Отложенное завершение работы
	defer func() {
		w.wg.Done()
		w.readerCancel = true
		log.Println("Reader stopped")
	}()

	for {
		select {
		case data, ok := <-w.Buffer:
			if !ok {
				return
			}

			// Чтение данных из буфера
			if w.linkedList == nil {
				w.linkedList = &LinkedList{Last: &Node{
					Data: data,
					Next: nil,
				}}
				w.linkedList.First = &Node{
					Next: w.linkedList.Last,
				}
			} else {
				w.linkedList.Last.Next = &Node{
					Data: data,
				}
				w.linkedList.Last = w.linkedList.Last.Next
			}
		case <-time.After(time.Second):
			if w.readerShutdown {
				return
			}
		}
	}
}

// startBalancer балансирует данные между отправкой и чтением
func (w *Worker) startBalancer() {
	log.Println("Starting balancer...")

	w.wg.Add(1)
	// Отложенное завершение работы
	defer func() {
		w.wg.Done()
		log.Println("Balancer stopped")
	}()

	for {
		// Проверка на наличие данных
		if w.linkedList == nil || w.linkedList.First == nil {
			if w.readerCancel {
				return
			}

			continue
		}

		if w.linkedList.First.Data.PeriodStart != "" && !w.linkedList.First.Sent {
			log.Println("Sending data...", w.linkedList.First.Data)
			if err := w.Sender.Send(w.linkedList.First.Data); err != nil {
				log.Println("Error sending data:", err)
			}
		}

		w.linkedList.First.Sent = true

		// Удаление отправленных данных
		if w.linkedList.First.Next != nil {
			w.linkedList.First = w.linkedList.First.Next
		} else if w.readerCancel {
			return
		}
	}
}

// Shutdown завершает работу воркера, дожидаясь завершения всех задач
func (w *Worker) Shutdown() {
	log.Println("Shutting down worker...")

	w.readerShutdown = true

	w.wg.Wait()
}
