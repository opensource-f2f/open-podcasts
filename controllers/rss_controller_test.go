package controllers

import (
	"github.com/SlyMarbo/rss"
	"github.com/stretchr/testify/assert"
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

func Test_getSingleCategory(t *testing.T) {
	type args struct {
		category string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{{
		name: "multiple categories",
		args: args{
			category: "society & culture",
		},
		wantResult: "society",
	}, {
		name: "single category",
		args: args{
			category: "tech",
		},
		wantResult: "tech",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantResult, getSingleCategory(tt.args.category), "getSingleCategory(%v)", tt.args.category)
		})
	}
}
