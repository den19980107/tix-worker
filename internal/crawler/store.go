package crawler

import "fmt"

type CrawlerStore struct {
	crawlerMap map[int]*Crawler
}

func InitStore() CrawlerStore {
	return CrawlerStore{
		crawlerMap: map[int]*Crawler{},
	}
}

func (c *CrawlerStore) Get(id int) (*Crawler, error) {
	cw, exist := c.crawlerMap[id]
	if !exist {
		return nil, fmt.Errorf("crawler %d not exist", id)
	}

	return cw, nil
}

func (c *CrawlerStore) Create(id int) error {
	_, exist := c.crawlerMap[id]
	if exist {
		return fmt.Errorf("crawler %d already exist", id)
	}

	cw := Create()
	c.crawlerMap[id] = &cw

	return nil
}

func (c *CrawlerStore) Delete(id int) {
	delete(c.crawlerMap, id)
}
