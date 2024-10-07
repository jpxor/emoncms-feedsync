package emoncms

type client struct {
	host     string
	apikey   string
	urlLimit int
}

type FeedAPI interface {
	List() ([]Feed, error)
	TimeValue(feed *Feed) error
	Data(feed Feed, start, end int64) (string, int64, error)
	Insert(feed Feed, data string) error
}

type Client struct {
	Feed FeedAPI
}

func NewClient(host, apikey string) Client {
	c := client{host, apikey, 2000}
	return Client{
		Feed: &c,
	}
}

func (c *Client) SetUrlLimit(limit int) {
	internalClient := c.Feed.(*client)
	internalClient.urlLimit = limit
}
