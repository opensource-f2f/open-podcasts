package handler

import (
	"github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
	"time"
)

func Test_sortWithDesc(t *testing.T) {
	type args struct {
		items []v1alpha1.Episode
	}
	tests := []struct {
		name   string
		args   args
		verify func(t *testing.T, items []v1alpha1.Episode)
	}{{
		name: "normal case",
		args: args{
			items: []v1alpha1.Episode{{
				Spec: v1alpha1.EpisodeSpec{
					Title: "1",
					Date:  metav1.Time{Time: time.Now()},
				},
			}, {
				Spec: v1alpha1.EpisodeSpec{
					Title: "2",
					Date:  metav1.Time{Time: time.Now().Add(time.Minute)},
				},
			}},
		},
		verify: func(t *testing.T, items []v1alpha1.Episode) {
			assert.Equal(t, "2", items[0].Spec.Title)
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortWithDesc(tt.args.items)
			tt.verify(t, tt.args.items)
		})
	}
}
