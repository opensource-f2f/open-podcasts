package controllers

import (
	"github.com/SlyMarbo/rss"
	"testing"
)

func Test_getFixedLink(t *testing.T) {
	type args struct {
		source string
		feed   *rss.Feed
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "ximalaya",
		args: args{
			source: "http://www.ximalaya.com/album/53320813.xml",
			feed: &rss.Feed{
				Link: "http://www.ximalaya.com",
			},
		},
		want: "http://www.ximalaya.com/album/53320813",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFixedLink(tt.args.source, tt.args.feed); got != tt.want {
				t.Errorf("getFixedLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
