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
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	osf2fv1alpha1 "github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
)

// EpisodeReconciler reconciles a Episode object
type EpisodeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const podcastTitle = "title.podcast"

//+kubebuilder:rbac:groups=osf2f.my.domain,resources=episodes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=osf2f.my.domain,resources=episodes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=osf2f.my.domain,resources=episodes/finalizers,verbs=update
//+kubebuilder:rbac:groups=osf2f.my.domain,resources=rsses,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Episode object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *EpisodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	_ = log.FromContext(ctx)
	episode := &osf2fv1alpha1.Episode{}
	if err = r.Get(ctx, req.NamespacedName, episode); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	var ownerRef v1.OwnerReference
	if len(episode.OwnerReferences) > 0 {
		ownerRef = episode.OwnerReferences[0]
	} else {
		return
	}

	rss := &osf2fv1alpha1.RSS{}
	if err = r.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      ownerRef.Name,
	}, rss); err != nil {
		// TODO should remove the RSS or create an event for it?
		return
	}

	title := episode.Annotations[podcastTitle]
	if rss.Spec.Title != "" && rss.Spec.Title != title {
		if episode.Annotations == nil {
			episode.Annotations = map[string]string{}
		}
		episode.Annotations[podcastTitle] = rss.Spec.Title
		err = r.Update(ctx, episode)
	}
	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *EpisodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&osf2fv1alpha1.Episode{}).
		Complete(r)
}
