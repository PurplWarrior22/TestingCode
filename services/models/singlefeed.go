package models

//SingleFeed represents the service view of a single rss feed object
type SingleFeed struct {

	// the title of the feed
	Title string `json:"title"`

	// the unique id of the feed used to get by id
	Id string `json:"id,omitempty"`

	// the elasticsearch query used to populate the feed
	Query map[string]interface{} `json:"query"`
}

type Feeds struct {
	Feeds map[string]*SingleFeed `json:"feeds"`
}
