package common

import "github.com/opensource-f2f/open-podcasts/generated/clientset/versioned"

type CommonOption struct {
	Client           *versioned.Clientset
	DefaultNamespace string
}
