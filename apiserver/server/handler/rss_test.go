package handler

import (
	"fmt"
	"github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_opmlRSSList_asXML(t *testing.T) {
	type fields struct {
		Title      string
		CreateDate time.Time
		Items      []v1alpha1.RSS
	}
	tests := []struct {
		name    string
		fields  fields
		wantXml []byte
		wantErr assert.ErrorAssertionFunc
	}{{
		name: "normal",
		fields: fields{
			Title:      "title",
			CreateDate: time.Time{},
			Items: []v1alpha1.RSS{{
				Spec: v1alpha1.RSSSpec{
					Title:       "title",
					Address:     "address",
					Link:        "link",
					Description: "description",
				},
			}},
		},
		wantXml: []byte(`<?xml version='1.0' encoding='UTF-8' standalone='no' ?>
<opml version="2.0">
  <head>
    <title>title</title>
    <dateCreated>0001-01-01 00:00:00 +0000 UTC</dateCreated>
  </head>
  <body>
    <outline text="description" title="title" type="rss" xmlUrl="address" htmlUrl="link"/>
  </body>
</opml>
`),
		wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
			assert.Nil(t, err)
			return false
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &opmlRSSList{
				Title:      tt.fields.Title,
				CreateDate: tt.fields.CreateDate,
				Items:      tt.fields.Items,
			}
			gotXml, err := o.asXML()
			if tt.wantErr(t, err, fmt.Sprintf("asXML()")) {
				return
			}
			assert.Equalf(t, tt.wantXml, gotXml, "asXML()")
		})
	}
}
