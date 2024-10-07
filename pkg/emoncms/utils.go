package emoncms

import (
	"strings"
)

func splitFeedDataString(targetSplit int, data string) (string, string) {
	// data string format: [[f,f],[f,f],[f,f]]
	if len(data) == 0 {
		return data, ""
	}

	targetSplit = min(targetSplit, 2+len(data)/2)
	halfdata := data[:targetSplit]

	splitIndex := strings.LastIndex(halfdata, "],[")
	if splitIndex == -1 {
		return data, ""
	}
	left := data[:splitIndex+1] + "]"
	right := "[" + data[splitIndex+2:]
	return left, right
}

func FeedNames(feeds []Feed) []string {
	names := make([]string, len(feeds))
	for i, feed := range feeds {
		names[i] = feed.Name
	}
	return names
}
