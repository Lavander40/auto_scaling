package event_consumer

import (
	"auto_scaling/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Pocessor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Pocessor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		newEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("error during consumer start: %s", err.Error())
			continue
		}

		if len(newEvents) == 0 {
			time.Sleep(2 * time.Second)
			continue
		}

		if err := c.handleEvents(newEvents); err != nil {
			log.Print(err)
			continue
		}

	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		//log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event %s", err.Error())
			continue
		}
	}

	return nil
}
