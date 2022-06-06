package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func init() {
	SchemeBuilder.Register(
		&Episode{}, &EpisodeList{},
		&RSS{}, &RSSList{},
		&Profile{}, &ProfileList{},
		&Notifier{}, &NotifierList{},
		&Subscription{}, &SubscriptionList{},
		&Category{}, &CategoryList{},
		&Show{}, &ShowList{},
		&ShowItem{}, &ShowItemList{},
		&Storage{}, &StorageList{},
		&Author{}, &AuthorList{},
	)
}

// SchemeGroupVersion is group version used to register these objects.
var SchemeGroupVersion = GroupVersion

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}
