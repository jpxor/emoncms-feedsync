package main

import (
	"flag"
	"fmt"
	"jpxor/emoncms/feedsync/pkg/emoncms"
	"os"
	"time"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "./config.yaml", "path to the configuration file")
	flag.Parse()

	config, err := readConfig(configPath)
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	filters := NewFilterMap()

	local := emoncms.NewClient(config.Local.Host, config.Local.APIKey)
	remote := emoncms.NewClient(config.Remote.Host, config.Remote.APIKey)

	localFeeds, err := local.Feed.List()
	if err != nil {
		fmt.Println("Error getting local feeds:", err)
		os.Exit(1)
	}

	localFeedsMap := make(map[string]emoncms.Feed, len(localFeeds))
	for _, feed := range localFeeds {
		localFeedsMap[feed.Name] = feed
	}

	remoteFeeds, err := remote.Feed.List()
	if err != nil {
		fmt.Println("Error getting remote feeds:", err)
		os.Exit(1)
	}

	remoteFeeds = filterByNames(remoteFeeds, config.FeedsFilter)
	remoteFeeds = filterByNames(remoteFeeds, emoncms.FeedNames(localFeeds))

	if config.UrlLimit > 0 {
		remote.SetUrlLimit(int(config.UrlLimit))
	}

	if config.Start != 0 {
		for i := range remoteFeeds {
			remotefeed := &remoteFeeds[i]
			remotefeed.LastUpdate = config.Start
		}
	} else {
		for i := range remoteFeeds {
			remotefeed := &remoteFeeds[i]
			err := remote.Feed.TimeValue(remotefeed)
			if err != nil {
				fmt.Println("Error failed to fetch LastUpdated time and value:", err)
				os.Exit(1)
			}
			if isUnixMilli(remotefeed.LastUpdate) {
				fmt.Println("WARN: time in milliseconds!?", remotefeed.Name, remotefeed.LastUpdate)
				remotefeed.LastUpdate /= 1000
			}
		}
	}

	for {
		now := time.Now().Unix()
		for i := range remoteFeeds {
			remotefeed := &remoteFeeds[i]

			localfeed, ok := localFeedsMap[remotefeed.Name]
			if !ok {
				fmt.Printf("no matching local feed found for %s\n", remotefeed.Name)
				os.Exit(1)
			}

			fmt.Printf("\nSync: %s (interval: %ds)\n", remotefeed.Name, localfeed.Interval)
			end := int64(0)

			for end < now {
				var data string

				// may only return a portion of the data, returns actual end
				data, end, err = local.Feed.Data(localfeed, remotefeed.LastUpdate, now)
				if err != nil {
					fmt.Println("ERROR failed to get data:", err)
					break
				}

				// filter data (ie remove outliers)
				data, err = filters.Apply(remotefeed.Name, data)
				if err != nil {
					fmt.Println("ERROR failed to filter data:", err)
					break
				}

				err = remote.Feed.Insert(*remotefeed, data)
				if err != nil {
					fmt.Println("ERROR failed to send data:", err)
					break
				}

				remotefeed.LastUpdate = end
				fmt.Println("Sent", len(data), "bytes:", trimString(data, 80))
			}
		}

		fmt.Println("Sleeping for", config.Interval, "seconds.")
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}

func isUnixMilli(unixtime int64) bool {
	// assumes we're not dealing with times from the future
	return unixtime > time.Now().Unix() &&
		unixtime <= time.Now().UnixMilli()
}

func filterByNames(feedList []emoncms.Feed, feedNames []string) []emoncms.Feed {
	if len(feedNames) == 0 {
		return feedList
	}
	var result []emoncms.Feed
	nameSet := make(map[string]bool)

	for _, name := range feedNames {
		nameSet[name] = true
	}
	for _, feed := range feedList {
		if nameSet[feed.Name] {
			result = append(result, feed)
		}
	}
	return result
}

func trimString(str string, size int) string {
	if len(str) <= size {
		return str
	}
	return str[:size] + "..."
}
