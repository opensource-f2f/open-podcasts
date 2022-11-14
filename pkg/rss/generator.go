package rss

import (
	"fmt"
	"github.com/eduncan911/podcast"
	"github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

func GenerateRSSFromNonCRDFiles(externalServer, show string, showItems []string) (content string, err error) {
	// load and parse the show data
	var data []byte
	if data, err = ioutil.ReadFile(show); err != nil {
		return
	}

	showSpec := &v1alpha1.ShowSpec{}
	if err = yaml.Unmarshal(data, showSpec); err != nil {
		return
	}

	var showItemSpecList []*v1alpha1.ShowItemSpec
	// load and parse the show items data
	for i := range showItems {
		item := showItems[i]

		if data, err = ioutil.ReadFile(item); err != nil {
			return
		}

		showItem := &v1alpha1.ShowItemSpec{}
		if err = yaml.Unmarshal(data, showItem); err != nil {
			return
		}
		showItemSpecList = append(showItemSpecList, showItem)
	}

	content = GenerateRSSFromNonCRD(externalServer, showSpec, showItemSpecList)
	return
}

func GenerateRSSFromNonCRD(externalServer string, showSpec *v1alpha1.ShowSpec, showItemSpecList []*v1alpha1.ShowItemSpec) string {
	show := &v1alpha1.Show{
		Spec: *showSpec,
	}
	showItems := &v1alpha1.ShowItemList{
		Items: make([]v1alpha1.ShowItem, len(showItemSpecList)),
	}
	for i := range showItemSpecList {
		item := showItemSpecList[i]
		showItems.Items[i] = v1alpha1.ShowItem{
			Spec: *item,
		}
	}
	return GenerateRSS(externalServer, show, showItems)
}

// GenerateRSS generate the RSS string content from custom resource Show and ShowItemList
func GenerateRSS(externalServer string, show *v1alpha1.Show, showItems *v1alpha1.ShowItemList) string {
	ti, l, d := show.Spec.Title, show.Spec.Link, show.Spec.Description
	pubDate, updatedDate := show.CreationTimestamp.Time, show.CreationTimestamp.Time

	// instantiate a new Podcast
	p := podcast.New(ti, l, d, &pubDate, &updatedDate)
	p.Language = show.Spec.Language
	p.Generator = "Open Podcast (https://github.com/opensource-f2f/open-podcasts)"
	p.Link = show.Spec.Link
	if len(show.Spec.Categories) > 0 {
		p.Category = show.Spec.Categories[0]
		for i := range show.Spec.Categories {
			category := show.Spec.Categories[i]
			p.ICategories = append(p.ICategories, &podcast.ICategory{Text: category})
		}
	}
	if show.Spec.Image != "" {
		p.Image = &podcast.Image{
			URL:   show.Spec.Image,
			Title: "cover",
			Link:  show.Spec.Image,
		}
	}

	for i := range showItems.Items {
		item := showItems.Items[i]

		_, _ = p.AddItem(podcast.Item{
			Title:       item.Spec.Title,
			Description: item.Spec.Description,
			Comments:    "notes",
			PubDate:     &item.CreationTimestamp.Time,
			Enclosure: &podcast.Enclosure{
				URL: GetAudioFileURL(externalServer, item),
			},
		})
	}
	return p.String()
}

func GetAudioFileURL(externalServer string, item v1alpha1.ShowItem) string {
	itemSpec := item.Spec
	externalServer = strings.TrimSuffix(externalServer, "/")

	switch itemSpec.LocalStorage {
	case "github-release":
		return fmt.Sprintf("%s/releases/download/%d/%s", externalServer, itemSpec.Index, itemSpec.Filename)
	case "url":
		return itemSpec.Filename
	default:
		return fmt.Sprintf("%s/showitems/%s/download", externalServer, item.Name)
	}
}
