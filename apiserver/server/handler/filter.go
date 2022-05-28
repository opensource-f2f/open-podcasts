package handler

import (
	"github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	"strings"
)

type rssFilter interface {
	filter([]v1alpha1.RSS) []v1alpha1.RSS
}

type rssNonFilter struct {
}

func (f *rssNonFilter) filter(items []v1alpha1.RSS) []v1alpha1.RSS {
	return items
}

type rssCategoryFilter struct {
	keyword string
}

func (f *rssCategoryFilter) filter(items []v1alpha1.RSS) (result []v1alpha1.RSS) {
	for i := range items {
		item := items[i]
		for j := range item.Spec.Categories {
			category := item.Spec.Categories[j]
			if strings.ToLower(category) == f.keyword {
				result = append(result, item)
				break
			}
		}
	}
	return items
}
