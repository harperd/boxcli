package boxclient

// Runtime options for Box Client
// Method: GET, PUT, POST, DELETE, etc.
// Database: FHIR or DOC
// Resource: A valid FHIR resource
// Color: If true, JSON output is syntax highlighted
// Unformatted: If true, JSON output is not formatted
// Count: If true, only the count of the results are returned
// Index: If true, only the resource at the specified index is returned
// Query: The JSON query to apply
type Options struct {
	Address string
	Method string
	Database string
	Resource string
	Color bool
	Unformatted bool
	OmitNulls bool
	Count bool
	Index string
	Query string
}