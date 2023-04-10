package event_consumer

import (
	"bot/events"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start(url string, port int) error {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      nil,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		log.Println("caught signal", map[string]string{
			"signal": s.String(),
		})

		ctx, cancelc := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelc()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		log.Println("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})

		cancel()
		wg.Wait()
		shutdownError <- nil
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		log.Printf("Starting server on port %d", port)
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
		}
		err = <-shutdownError
		if err != nil {
			log.Println(err)
		}

		log.Println(err)

		log.Println("stopped server", map[string]string{
			"addr": srv.Addr,
		})
	}()
	if url != "" {
		gotEvents, err := c.fetcher.ListenWebhooks(url)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Println(err)
		}
	} else {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				gotEvents, err := c.fetcher.Fetch(c.batchSize)
				if err != nil {
					log.Printf("[ERR] consumer: %s", err.Error())
					continue
				}

				if len(gotEvents) == 0 {
					time.Sleep(1 * time.Second)
					continue
				}

				if err := c.handleEvents(gotEvents); err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}
	return nil
}

func (c *Consumer) handleEvents(events []events.Event[events.TelegramMeta]) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can not handle event: %s", err.Error())
			continue
		}
	}
	return nil
}
