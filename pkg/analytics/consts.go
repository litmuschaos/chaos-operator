package analytics

const (
	// GAclientID contains TrackingID of the application
	GAclientID string = "UA-92076314-21"

	// supported event categories

	// CategoryLI category notifies installation of a component of Litmus Infrastructure
	CategoryLI string = "Litmus-Infra"
	// CategoryCE category notifies installation of a Litmus Experiment
	CategoryCE string = "Chaos-Experiment"

	// supported event actions

	// ActionI is sent when the installation is triggered
	ActionI string = "Installation"

	// supported event labels

	// LabelO denotes event is associated to which Litmus component
	LabelO string = "Chaos-Operator"
)
