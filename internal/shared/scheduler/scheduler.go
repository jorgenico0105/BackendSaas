package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type JobFunc func(ctx context.Context)

func StartCron(job JobFunc) {
	loc, _ := time.LoadLocation("America/Guayaquil")

	c := cron.New(
		cron.WithLocation(loc),
		cron.WithSeconds(),
	)

	_, err := c.AddFunc("0 0 0 * * *", func() {
		log.Println("Ejecutando job programado...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		job(ctx)
	})

	if err != nil {
		log.Fatal(err)
	}

	c.Start()
}
