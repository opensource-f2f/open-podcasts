package controllers

import (
	"github.com/SlyMarbo/rss"
	"reflect"
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

func Test_removeDuplicateStr(t *testing.T) {
	type args struct {
		strSlice []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{{
		name: "contains duplicated items",
		args: args{
			strSlice: []string{"a", "b", "a"},
		},
		want: []string{"a", "b"},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDuplicateStr(tt.args.strSlice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeDuplicateStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
