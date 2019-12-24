package analytics

const (
	// GAclientID contains TrackingID of the application
	GAclientID string = "UA-127388617-2"

	// supported event categories

	// Install event is sent when install occurs
	Install string = "install"

	// ChaosOperator event is sent when Reconcile function ends its job
	ChaosOperator string = "chaos-operator"

	// Ping event is sent periodically
	Ping string = "ping"

	// RunnerCreation event is sent when a Runner pod is created
	RunnerCreation string = "runner-creation"

	// AppName event
	AppName string = "litmus-installations"
)
