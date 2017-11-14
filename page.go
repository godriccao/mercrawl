package mercrawl

import (
	"sync"
)

// PageState stores crawled pages
type PageState struct {
	crawled map[string]bool
	Mux     *sync.Mutex
}

// Set marks a page as stored
func (ps *PageState) Set(index string) {
	ps.Mux.Lock()
	ps.crawled[index] = true
	ps.Mux.Unlock()
}

// Get provide the crawled state of a certain page
func (ps *PageState) Get(index string) bool {
	return ps.crawled[index]
}
