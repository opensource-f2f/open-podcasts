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
	"github.com/linuxsuren/open-podcasts/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// EventReconciler reconciles a Profile object
type EventReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups="",resources=events,verbs=get;list;watch
//+kubebuilder:rbac:groups=osf2f.my.domain,resources=profiles,verbs=list;get
//+kubebuilder:rbac:groups=osf2f.my.domain,resources=subscriptions,verbs=list;get
//+kubebuilder:rbac:groups=osf2f.my.domain,resources=notifiers,verbs=list;get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Profile object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *EventReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	_ = log.FromContext(ctx)

	receiveEvent := &v1.Event{}
	if err = r.Client.Get(ctx, req.NamespacedName, receiveEvent); err != nil {
		err = client.IgnoreNotFound(err)
		return
	}

	var notifiers []*v1alpha1.Notifier
	if notifiers, err = r.getNotifierList(ctx, receiveEvent.InvolvedObject); err == nil {
		for i := range notifiers {
			notifier := notifiers[i]

			if notifier.Spec.Slack != nil {
				err = notifier.Spec.Slack.Send(receiveEvent.Message)
			}

			if notifier.Spec.Feishu != nil {
				err = notifier.Spec.Feishu.Send(receiveEvent.Message)
			}
		}
	} else {
		// ignore if no notifier found
		err = nil
	}
	return
}

func (r *EventReconciler) getNotifierList(ctx context.Context, objectRef v1.ObjectReference) (notifiers []*v1alpha1.Notifier, err error) {
	profileList := &v1alpha1.ProfileList{}
	if err = r.List(ctx, profileList); err != nil {
		return
	}

	var notifierNames []v1.LocalObjectReference
	for i := range profileList.Items {
		profile := profileList.Items[i]

		sub := profile.Spec.Subscription
		if sub.Name == objectRef.Name {
			notifierNames = append(notifierNames, sub)
		}
	}

	// find all the notifiers, then filter the expected
	notifierList := &v1alpha1.NotifierList{}
	if err = r.Client.List(ctx, notifierList); err != nil {
		return
	}

	for i := range notifierList.Items {
		notifier := notifierList.Items[i]

		for j := range notifierNames {
			if notifierNames[j].Name == notifier.Name {
				notifiers = append(notifiers, &notifier)
				break
			}
		}
	}
	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *EventReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Event{}).
		WithEventFilter(&eventSourceFilter{component: "rss"}).
		Complete(r)
}

type eventSourceFilter struct {
	component string
}

func (f *eventSourceFilter) Create(e event.CreateEvent) bool {
	if receiveEvent, ok := e.Object.(*v1.Event); ok {
		if receiveEvent.Source.Component == f.component && receiveEvent.Type == v1.EventTypeNormal {
			return true
		}
	}
	return false
}

func (f *eventSourceFilter) Delete(e event.DeleteEvent) bool {
	return false
}

func (f *eventSourceFilter) Update(e event.UpdateEvent) bool {
	return false
}

func (f *eventSourceFilter) Generic(e event.GenericEvent) bool {
	return false
}
