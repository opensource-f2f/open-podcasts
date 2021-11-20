package main

import (
	"fmt"
	"github.com/SlyMarbo/rss"
)

func main() {
	feed, err := rss.Fetch("http://www.ximalaya.com/album/53320813.xml")
	if err != nil {
		// handle error.
	}

	for _, item := range feed.Items {
		fmt.Println(item)
	}

	// ... Some time later ...

	err = feed.Update()
	if err != nil {
		// handle error.
	}
}
