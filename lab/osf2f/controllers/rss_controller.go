/*
Copyright 2022 The osf2f Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"github.com/SlyMarbo/rss"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1alpha1 "github.com/linuxsuren/goplay/lab/osf2f/api/v1alpha1"
)

// RSSReconciler reconciles a RSS object
type RSSReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=osf2f.my.domain,resources=rsses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=osf2f.my.domain,resources=rsses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=osf2f.my.domain,resources=rsses/finalizers,verbs=update
//+kubebuilder:rbac:groups=osf2f.my.domain,resources=episodes,verbs=get;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *RSSReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	_ = log.FromContext(ctx)

	rssObj := &v1alpha1.RSS{}
	if err = r.Client.Get(ctx, req.NamespacedName, rssObj); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	address := rssObj.Spec.Address
	if address == "" {
		result = ctrl.Result{RequeueAfter: time.Minute}
		err = r.errorAndRecord(rssObj, v1.EventTypeWarning, "Failed to fetch RSS",
			fmt.Sprintf("the address of the RSS: %s is empty", req.NamespacedName.String()))
		return
	}

	err = r.fetchByRSS(address, rssObj)
	return
}

func (r *RSSReconciler) fetchByRSS(address string, rssObject *v1alpha1.RSS) (err error) {
	var feed *rss.Feed
	if feed, err = rss.Fetch(address); err != nil {
		err = r.errorAndRecord(rssObject, v1.EventTypeWarning, "Failed to fetch RSS",
			fmt.Sprintf("failed to fetch RSS by address: %s, error is %v", address, err))
		return
	}

	rssObject.Spec.Title = feed.Title
	rssObject.Spec.Description = feed.Description
	rssObject.Spec.Link = feed.Link
	if feed.Image != nil {
		if feed.Image.URL != "" {
			rssObject.Spec.Image = feed.Image.URL
		} else if feed.Image.Href != "" {
			rssObject.Spec.Image = feed.Image.Href
		}
	}
	if err = r.Client.Update(context.Background(), rssObject); err != nil {
		return
	}

	if err = r.storeEpisodes(feed.Items, rssObject.ObjectMeta); err == nil {
		err = r.setLastUpdateTime(rssObject.Namespace, rssObject.Name)
	}
	return
}

func (r *RSSReconciler) setLastUpdateTime(ns, name string) (err error) {
	rssObj := &v1alpha1.RSS{}
	if err = r.Client.Get(context.Background(), types.NamespacedName{
		Namespace: ns,
		Name:      name,
	}, rssObj); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}
	rssObj.Status.LastUpdateTime = metav1.NewTime(time.Now())
	err = r.Client.Status().Update(context.Background(), rssObj)
	return
}

func (r *RSSReconciler) storeEpisodes(items []*rss.Item, meta metav1.ObjectMeta) (err error) {
	for i, _ := range items {
		rssMeta := meta.DeepCopy()
		episodeName := fmt.Sprintf("%s-%d", meta.Name, i)

		if err = r.storeEpisode(items[i], rssMeta, episodeName); err != nil {
			return
		}
	}
	return
}

func (r *RSSReconciler) storeEpisode(item *rss.Item, meta *metav1.ObjectMeta, episodeName string) (err error) {
	var audioSource string
	if len(item.Enclosures) > 0 {
		audioSource = item.Enclosures[0].URL
	}

	episode := &v1alpha1.Episode{}
	if err = r.Client.Get(context.Background(), types.NamespacedName{
		Namespace: meta.Namespace,
		Name:      episodeName,
	}, episode); err != nil {
		episode := &v1alpha1.Episode{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: meta.Namespace,
				Name:      episodeName,
				Labels: map[string]string{
					"rss": meta.Name,
				},
				OwnerReferences: []metav1.OwnerReference{{
					Name:       meta.Name,
					UID:        meta.UID,
					Kind:       "RSS",
					APIVersion: "osf2f.my.domain/v1alpha1",
				}},
			},
			Spec: v1alpha1.EpisodeSpec{
				Title:       item.Title,
				Summary:     item.Summary,
				Content:     item.Content,
				AudioSource: audioSource,
				Link:        item.Link,
			},
		}
		if item.Image != nil {
			if item.Image.URL != "" {
				episode.Spec.CoverImage = item.Image.URL
			} else if item.Image.Href != "" {
				episode.Spec.CoverImage = item.Image.Href
			}
		}
		err = r.Client.Create(context.Background(), episode)
	} else {
		episode.OwnerReferences = []metav1.OwnerReference{{
			Name:       meta.Name,
			UID:        meta.UID,
			Kind:       "RSS",
			APIVersion: "osf2f.my.domain/v1alpha1",
		}}
		episode.Labels = map[string]string{
			"rss": meta.Name,
		}
		if item.Image != nil {
			if item.Image.URL != "" {
				episode.Spec.CoverImage = item.Image.URL
			} else if item.Image.Href != "" {
				episode.Spec.CoverImage = item.Image.Href
			}
		}
		err = r.Client.Update(context.Background(), episode)
	}
	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *RSSReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.recorder = mgr.GetEventRecorderFor("rss")
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.RSS{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}

func (r *RSSReconciler) errorAndRecord(object runtime.Object, eventType, reason, msg string) error {
	r.recorder.Eventf(object, eventType, reason, msg)
	return fmt.Errorf(msg)
}