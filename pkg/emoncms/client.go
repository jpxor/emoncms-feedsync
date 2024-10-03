package emoncms

type client struct {
	host   string
	apikey string
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
	c := client{host, apikey}
	return Client{
		Feed: &c,
	}
}
