/*
Copyright 2022 The open-podcasts Authors. All rights reserved.

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
	"github.com/eduncan911/podcast"
	"github.com/google/go-github/v44/github"
	"github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"os"
	"path"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ShowItemReconciler reconciles a ShowItem object
type ShowItemReconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	ExternalServer string
}

//+kubebuilder:rbac:groups=osf2f.my.domain,resources=showitems,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=osf2f.my.domain,resources=showitems/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=osf2f.my.domain,resources=showitems/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ShowItem object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *ShowItemReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	_ = log.FromContext(ctx)

	showItem := &v1alpha1.ShowItem{}
	if err = r.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      req.Name,
	}, showItem); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	show := &v1alpha1.Show{}
	if err = r.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      showItem.Spec.ShowRef,
	}, show); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	showItemList := &v1alpha1.ShowItemList{}
	if err = r.List(ctx, showItemList, client.MatchingLabels(map[string]string{
		v1alpha1.LabelKeyShowRef: show.Name,
	})); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	rssXML := generateRSS(r.ExternalServer, show, showItemList)
	err = r.saveRSSXML(rssXML, show.Namespace, show.Spec.Storage)

	err = r.storeDataFile(show.Namespace, show.Spec.Storage, showItem)
	return
}

func generateRSS(externalServer string, show *v1alpha1.Show, showItems *v1alpha1.ShowItemList) string {
	ti, l, d := show.Spec.Title, show.Spec.Link, show.Spec.Description
	pubDate, updatedDate := show.CreationTimestamp.Time, show.CreationTimestamp.Time

	// instantiate a new Podcast
	p := podcast.New(ti, l, d, &pubDate, &updatedDate)
	p.Language = show.Spec.Language
	p.Generator = "Open Podcast (https://github.com/opensource-f2f/open-podcasts)"

	for i := range showItems.Items {
		item := showItems.Items[i]

		_, _ = p.AddItem(podcast.Item{
			Title:       item.Spec.Title,
			Description: item.Spec.Description,
			Comments:    "notes",
			PubDate:     &item.CreationTimestamp.Time,
			Enclosure: &podcast.Enclosure{
				URL: fmt.Sprintf("%s/showitems/%s/download", externalServer, item.Name),
			},
		})
	}
	return p.String()
}

func (r *ShowItemReconciler) saveRSSXML(rssXML, namespace string, storageRef *v1.LocalObjectReference) (err error) {
	if storageRef == nil {
		return
	}

	ctx := context.Background()
	storage := &v1alpha1.Storage{}
	if err = r.Get(ctx, types.NamespacedName{
		Namespace: namespace,
		Name:      storageRef.Name,
	}, storage); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	storage.Annotations[v1alpha1.AnnotationKeyRSS] = rssXML
	err = r.Update(ctx, storage)
	return
}

func (r *ShowItemReconciler) storeDataFile(ns string, storageRef *v1.LocalObjectReference, showItem *v1alpha1.ShowItem) (
	err error) {
	downloadURL := showItem.Annotations[v1alpha1.AnnotationKeyDownloadURL]
	if storageRef == nil || downloadURL == "" {
		return
	}

	ctx := context.Background()
	storage := &v1alpha1.Storage{}
	if err = r.Get(ctx, types.NamespacedName{
		Namespace: ns,
		Name:      storageRef.Name,
	}, storage); err != nil {
		return
	}

	releases := storage.Spec.GitProviderReleases
	for i := range releases {
		item := releases[i]

		if item.Name == "" || item.Provider == "" {
			continue
		}

		secret := &v1.Secret{}
		var err error
		if err = r.Get(context.Background(), types.NamespacedName{
			Namespace: item.Secret.Namespace,
			Name:      item.Secret.Name,
		}, secret); err != nil {
			continue
		}

		gitClient, err := GetGitProvider(item.Provider, item.Server, string(secret.Data["token"]))
		if err != nil {
			continue
		} else {
			name := github.String(fmt.Sprintf("%d", showItem.Spec.Index))
			var release *github.RepositoryRelease
			release, _, err = gitClient.Repositories.CreateRelease(context.Background(), item.Owner, item.Repo, &github.RepositoryRelease{
				Name:    name,
				TagName: name,
			})
			if err == nil {
				var f *os.File
				f, err = getLocalFile(downloadURL, path.Join(os.TempDir(), showItem.Spec.Filename))
				if err == nil {
					_, _, err = gitClient.Repositories.UploadReleaseAsset(context.Background(), item.Owner, item.Repo, *release.ID,
						&github.UploadOptions{}, f)
				}
			}

			fmt.Println(err)
		}
	}
	return
}

func getLocalFile(downloadURL, localFile string) (f *os.File, err error) {
	if _, err = os.Stat(localFile); err != nil {
		var resp *http.Response
		if resp, err = http.Get(downloadURL); err != nil {
			return
		}

		var data []byte
		if data, err = ioutil.ReadAll(resp.Body); err == nil {
			err = ioutil.WriteFile(localFile, data, 0644)
		}
	}

	if err == nil {
		// TODO do the hash check
		f, err = os.OpenFile(localFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	}
	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *ShowItemReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.ShowItem{}).
		Complete(r)
}
