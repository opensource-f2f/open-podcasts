package rss

import (
	"github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	"testing"
)

func TestGenerateRSSFromNonCRDFiles(t *testing.T) {
	type args struct {
		externalServer string
		show           string
		showItems      []string
	}
	tests := []struct {
		name        string
		args        args
		wantContent string
		wantErr     bool
	}{{
		name: "no items",
		args: args{
			externalServer: "https://linuxsuren-bot.github.io/devops-talk/",
			show:           "data/show.yaml",
			showItems:      nil,
		},
		wantContent: `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
  <channel>
	<title>DevOps Talk</title>
	<link></link>
	<description>This is a DevOps talk</description>
	<generator>Open Podcast (https://github.com/opensource-f2f/open-podcasts)</generator>
	<lastBuildDate>Sat, 28 May 2022 10:41:17 +0000</lastBuildDate>
	<pubDate>Sat, 28 May 2022 10:41:17 +0000</pubDate>
  </channel>
</rss>`,
		wantErr: false,
	}, {
		name: "with one item",
		args: args{
			externalServer: "https://linuxsuren-bot.github.io/devops-talk/",
			show:           "data/show.yaml",
			showItems:      []string{"data/item-1.yaml"},
		},
		wantContent: "",
		wantErr:     false,
	}, {
		name: "with a URL episode source",
		args: args{
			show:      "data/show.yaml",
			showItems: []string{"data/item-2-url.yaml"},
		},
		wantContent: "",
		wantErr:     false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GenerateRSSFromNonCRDFiles(tt.args.externalServer, tt.args.show, tt.args.showItems)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRSSFromNonCRDFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if gotContent != tt.wantContent {
			//	t.Errorf("GenerateRSSFromNonCRDFiles() gotContent = %v, want %v", gotContent, tt.wantContent)
			//}
		})
	}
}

func TestGetAudioFileURL(t *testing.T) {
	type args struct {
		externalServer string
		item           v1alpha1.ShowItem
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "github release case",
		args: args{
			externalServer: "https://linuxsuren-bot.github.io/devops-talk",
			item: v1alpha1.ShowItem{
				Spec: v1alpha1.ShowItemSpec{
					LocalStorage: "github-release",
					Filename:     "demo.mp3",
					Index:        1,
				},
			},
		},
		want: "https://linuxsuren-bot.github.io/devops-talk/releases/download/1/demo.mp3",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAudioFileURL(tt.args.externalServer, tt.args.item); got != tt.want {
				t.Errorf("GetAudioFileURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
