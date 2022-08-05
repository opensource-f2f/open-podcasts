package controllers

import (
	"github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func Test_generateRSS(t *testing.T) {
	date, err := time.Parse(time.RFC1123Z, "Tue, 03 May 2022 04:02:17 +0000")
	assert.Nil(t, err)

	type args struct {
		show      *v1alpha1.Show
		showItems *v1alpha1.ShowItemList
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "normal",
		args: args{
			show: &v1alpha1.Show{
				ObjectMeta: v1.ObjectMeta{
					CreationTimestamp: v1.NewTime(date),
				},
				Spec: v1alpha1.ShowSpec{
					Title:       "title",
					Description: "desc",
					Link:        "link",
				},
			},
			showItems: &v1alpha1.ShowItemList{
				Items: []v1alpha1.ShowItem{{
					ObjectMeta: v1.ObjectMeta{
						CreationTimestamp: v1.NewTime(date),
					},
					Spec: v1alpha1.ShowItemSpec{
						Title:       "title",
						Description: "desc",
						Image:       "image",
					},
				}},
			},
		},
		want: `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
  <channel>
    <title>title</title>
    <link>link</link>
    <description>desc</description>
    <generator>Open Podcast (https://github.com/opensource-f2f/open-podcasts)</generator>
    <lastBuildDate>Tue, 03 May 2022 04:02:17 +0000</lastBuildDate>
    <pubDate>Tue, 03 May 2022 04:02:17 +0000</pubDate>
    <item>
      <guid>audio</guid>
      <title>title</title>
      <link>audio</link>
      <description>desc</description>
      <comments>notes</comments>
      <pubDate>Tue, 03 May 2022 04:02:17 +0000</pubDate>
      <enclosure url="/showitems//download" length="0" type="audio/x-m4a"></enclosure>
    </item>
  </channel>
</rss>`,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateRSS("", tt.args.show, tt.args.showItems); got != tt.want {
				t.Errorf("generateRSS() = %v, want %v", got, tt.want)
			}
		})
	}
}
