package boxclient

// Runtime options for Box Client
// Address: The HTTP address to the AidBox REST-base services
// Method: GET, PUT, POST, DELETE, etc.
// Resource: A valid FHIR resource
// Color: If true, JSON output is syntax highlighted
// Unformatted: If true, JSON output is not formatted
// Query: The JSON query to apply
type Options struct {
	Address string
	Method string
	Resource string
	Color bool
	Unformatted bool
	Query string
}