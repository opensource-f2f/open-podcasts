package v1alpha1

func init() {
	SchemeBuilder.Register(
		&Episode{}, &EpisodeList{},
		&RSS{}, &RSSList{},
		&Profile{}, &ProfileList{},
		&Notifier{}, &NotifierList{},
	)
}
