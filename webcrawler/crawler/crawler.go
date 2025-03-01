package crawler

import (
	"log"
	"slices"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"

	"book-search/webcrawler/config"
	"book-search/webcrawler/services/redis"
)

var allowedDomains = []string{"www.naiin.com", "www.chulabook.com", "www.amazon.com"}

func Crawl(seedURLs []string) error {
	// get and check env
	env, err := config.GetEnv()
	if err != nil {
		return err
	}

	storage := redis.GetStorage()

	c := colly.NewCollector(
		colly.AllowedDomains(allowedDomains...),
		colly.UserAgent(env.CrawlerUserAgent),
		colly.Async(true),
	)

	c.SetRequestTimeout(10 * time.Second)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 4, // Concurrent request limit is not the same as the number of crawler consumer threads
		RandomDelay: 5 * time.Second,
		Delay:       1 * time.Second,
	})

	err = c.SetStorage(storage)
	if err != nil {
		return err
	}

	q, err := queue.New(env.CrawlerThreads, storage)
	if err != nil {
		return err
	}

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if slices.Contains(allowedDomains, e.Request.URL.Host) {
			err = q.AddURL(e.Request.AbsoluteURL(e.Attr("href")))
			if err != nil {
				log.Println("Error adding URL to queue:", err)
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error visiting", r.Request.URL.String(), "with status", r.StatusCode)
	})

	for _, url := range seedURLs {
		if err = q.AddURL(url); err != nil {
			return err
		}
	}

	for {
		// Handle the queue with async requests,
		// wait for all requests to complete and check if the queue is empty
		if err = q.Run(c); err != nil {
			return err
		}
		c.Wait()
		if q.IsEmpty() {
			log.Println("Crawl completed, queue is empty")
			break
		}
	}

	return nil
}
