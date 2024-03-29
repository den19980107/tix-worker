package crawlerpool

import (
	"fmt"
	"log"
	"tix-worker/internal/crawler"
)

type CrawlerPool struct {
	m map[int]*crawler.Crawler
}

func New() CrawlerPool {
	return CrawlerPool{
		m: map[int]*crawler.Crawler{},
	}
}

func (pool *CrawlerPool) Set(id int, crawler crawler.Crawler) {
	pool.m[id] = &crawler
}

func (pool *CrawlerPool) Get(id int) (*crawler.Crawler, error) {
	c, exist := pool.m[id]
	if !exist {
		return nil, fmt.Errorf("crawler of id %d not exist", id)
	}

	return c, nil
}

func (pool *CrawlerPool) Remove(id int) {
	_, exist := pool.m[id]
	if !exist {
		log.Printf("crawler of id %d not exist", id)
		return
	}

	delete(pool.m, id)
}
