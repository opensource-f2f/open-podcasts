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
	"github.com/go-logr/logr"
	"github.com/google/go-github/v44/github"
	"github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// StorageReconciler reconciles a Storage object
type StorageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	event  record.EventRecorder
	logger logr.Logger
}

//+kubebuilder:rbac:groups=osf2f.my.domain,resources=storages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=osf2f.my.domain,resources=storages/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=osf2f.my.domain,resources=storages/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Storage object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *StorageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	r.logger = log.FromContext(ctx)
	r.logger.Info("start to reconcile storage")

	storage := &v1alpha1.Storage{}
	if err = r.Get(ctx, req.NamespacedName, storage); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	storage.Status.State = "ok"
	errMap := r.verifyGitProviders(storage.Spec.GitProviderReleases)
	for key, val := range errMap {
		var msg string
		if val != nil {
			msg = val.Error()
			storage.Status.State = "error"
		}
		storage.Status.Conditions[key] = msg
	}

	r.storeRSSToGitProviders(storage.Annotations[v1alpha1.AnnotationKeyRSS], storage.Spec.GitProviderReleases)

	err = r.Status().Update(ctx, storage)
	return
}

func (r *StorageReconciler) verifyGitProviders(releases []v1alpha1.GitProviderRelease) (result map[string]error) {
	for i := range releases {
		item := releases[i]

		if item.Name == "" {
			continue
		}

		if item.Provider == "" {
			result[item.Name] = fmt.Errorf("provider is empty")
			continue
		}

		secret := &v1.Secret{}
		var err error
		if err = r.Get(context.Background(), types.NamespacedName{
			Namespace: item.Secret.Namespace,
			Name:      item.Secret.Name,
		}, secret); err != nil {
			result[item.Name] = err
			continue
		}

		gitClient, err := GetGitProvider(item.Provider, item.Server, string(secret.Data["token"]))
		if err != nil {
			result[item.Name] = err
		} else {
			var release *github.RepositoryRelease
			release, _, err = gitClient.Repositories.CreateRelease(context.Background(), item.Owner, item.Repo, &github.RepositoryRelease{
				Draft:   github.Bool(true),
				Name:    github.String("fake"),
				TagName: github.String("fake"),
			})
			if err != nil {
				result[item.Name] = err
			} else {
				_, _ = gitClient.Repositories.DeleteRelease(context.Background(), item.Owner, item.Repo, *release.ID)
			}
		}
	}
	return
}

func (r *StorageReconciler) storeRSSToGitProviders(rss string, releases []v1alpha1.GitProviderRelease) (result map[string]error) {
	if rss == "" {
		return
	}

	for i := range releases {
		item := releases[i]

		if item.Name == "" {
			continue
		}

		if item.Provider == "" {
			result[item.Name] = fmt.Errorf("provider is empty")
			continue
		}

		secret := &v1.Secret{}
		var err error
		if err = r.Get(context.Background(), types.NamespacedName{
			Namespace: item.Secret.Namespace,
			Name:      item.Secret.Name,
		}, secret); err != nil {
			result[item.Name] = err
			continue
		}

		gitClient, err := GetGitProvider(item.Provider, item.Server, string(secret.Data["token"]))
		if err != nil {
			result[item.Name] = err
		} else {
			branch := github.String("master")
			_, _, err = gitClient.Repositories.EnablePages(context.Background(), item.Owner, item.Repo, &github.Pages{
				Source: &github.PagesSource{
					Branch: branch,
					Path:   github.String("/docs"),
				},
			})
			if err != nil {
				r.logger.Error(err, "failed to enable page for %s/%s", item.Owner, item.Repo)
			}

			content, _, _, err := gitClient.Repositories.GetContents(context.Background(), item.Owner, item.Repo,
				"docs/rss.xml", &github.RepositoryContentGetOptions{})
			var sha *string
			if err == nil {
				sha = content.SHA
			}

			_, _, err = gitClient.Repositories.CreateFile(context.Background(), item.Owner, item.Repo,
				"docs/rss.xml",
				&github.RepositoryContentFileOptions{
					Branch:  branch,
					Message: github.String("init"),
					Author: &github.CommitAuthor{
						Name:  github.String("rick"),
						Email: github.String("rick@jenkins-zh.cn"),
					},
					SHA:     sha,
					Content: []byte(rss),
				})
			fmt.Println(err)
		}
	}
	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *StorageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.event = mgr.GetEventRecorderFor("storage")
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Storage{}).
		//WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}
