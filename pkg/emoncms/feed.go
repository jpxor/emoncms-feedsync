package emoncms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Feed struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Interval   int64   `json:"interval"`
	Value      float64 `json:"-"`
	LastUpdate int64   `json:"-"`
}

// end-point called 'list', but return a map for convenience
func (c client) List() ([]Feed, error) {
	url := fmt.Sprintf("http://%s/feed/list.json?meta=1&apikey=%s", c.host, c.apikey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var feeds []Feed
	err = json.NewDecoder(resp.Body).Decode(&feeds)
	return feeds, err
}

// update feed with last updated time and value
func (c client) TimeValue(feed *Feed) error {
	url := fmt.Sprintf("https://%s/feed/timevalue.json?id=%s&apikey=%s", c.host, feed.ID, c.apikey)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var timeValue struct {
		Time  int64   `json:"time"`
		Value float64 `json:"value"`
	}
	err = json.Unmarshal(body, &timeValue)
	if err != nil {
		return err
	}

	feed.Value = timeValue.Value
	feed.LastUpdate = timeValue.Time
	return nil
}

// returns the data and the actual end time (incase there where too many data points)
func (c client) Data(feed Feed, start, end int64) (string, int64, error) {
	if feed.Interval == 0 {
		return "", start, fmt.Errorf("feed interval is 0: %s", feed.Name)
	}

	// there is a limit on number of data points returned per transaction
	const limit int64 = 8928
	count := (end - start) / feed.Interval
	if count > limit {
		end = start + (limit-1)*feed.Interval
	}

	format := "http://%s/feed/data.json?id=%s&start=%d&end=%d&interval=%d&average=0&timeformat=unix&skipmissing=1&apikey=%s"
	url := fmt.Sprintf(format, c.host, feed.ID, start, end, feed.Interval, c.apikey)

	resp, err := http.Get(url)
	if err != nil {
		return "", start, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", start, err
	}
	return string(body), end, nil
}

func (c client) Insert(feed Feed, data string) error {
	const format = "https://%s/feed/insert.json?id=%s&apikey=%s&data=%s"

	if data == "[]" || data == "" { // no data to send
		return nil
	}

	// ensure the url stays under an acceptable limit (eg, 2000 bytes long)
	if len(data)+len(format) > 2000 {
		// split the data string and send two requests
		left, right := splitFeedDataString(data)

		// oh ya, we're recursing
		err := c.Insert(feed, left)
		if err != nil {
			return err
		}
		// self imposed rate limiting
		time.Sleep(33 * time.Millisecond)

		// and recurse again to get the second half
		return c.Insert(feed, right)
	}

	url := fmt.Sprintf(format, c.host, feed.ID, c.apikey, data)
	resp, err := http.Post(url, "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error status: %s", resp.Status)
	}
	return nil
}
