package event_consumer

import (
	"log"
	"time"

	"ZavtraVebinar/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int //Размер пачки, говорит, сколько событий будем обрабатывать за раз
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {

	for i := 0; i < 1; i++ {
		c.processor.ProcessExistingCHATIDList()
	}

	for {

		gotEvents, err := c.fetcher.Fetch(c.batchSize) //отлавливает новые события, добавляет новые chatID в список
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("Новое событие:%s", event.Text)

		if err := c.processor.ProcessNewCHATID(event); err != nil {
			log.Printf("Невозможно обработать событие:%s", err.Error())

			continue
		}
	}

	return nil
}
