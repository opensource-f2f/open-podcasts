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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/opensource-f2f/open-podcasts/api/osf2f.my.domain/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// EpisodeLister helps list Episodes.
// All objects returned here must be treated as read-only.
type EpisodeLister interface {
	// List lists all Episodes in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Episode, err error)
	// Episodes returns an object that can list and get Episodes.
	Episodes(namespace string) EpisodeNamespaceLister
	EpisodeListerExpansion
}

// episodeLister implements the EpisodeLister interface.
type episodeLister struct {
	indexer cache.Indexer
}

// NewEpisodeLister returns a new EpisodeLister.
func NewEpisodeLister(indexer cache.Indexer) EpisodeLister {
	return &episodeLister{indexer: indexer}
}

// List lists all Episodes in the indexer.
func (s *episodeLister) List(selector labels.Selector) (ret []*v1alpha1.Episode, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Episode))
	})
	return ret, err
}

// Episodes returns an object that can list and get Episodes.
func (s *episodeLister) Episodes(namespace string) EpisodeNamespaceLister {
	return episodeNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// EpisodeNamespaceLister helps list and get Episodes.
// All objects returned here must be treated as read-only.
type EpisodeNamespaceLister interface {
	// List lists all Episodes in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Episode, err error)
	// Get retrieves the Episode from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Episode, error)
	EpisodeNamespaceListerExpansion
}

// episodeNamespaceLister implements the EpisodeNamespaceLister
// interface.
type episodeNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Episodes in the indexer for a given namespace.
func (s episodeNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Episode, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Episode))
	})
	return ret, err
}

// Get retrieves the Episode from the indexer for a given namespace and name.
func (s episodeNamespaceLister) Get(name string) (*v1alpha1.Episode, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("episode"), name)
	}
	return obj.(*v1alpha1.Episode), nil
}
