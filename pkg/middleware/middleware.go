package middleware

type key int

const (
	// OpenTracingSpanContext Key used to store the opentracing span.
	OpenTracingSpanContext key = iota
	// DatabaseConnection is the key that stores the current Postgres connection.
	DatabaseConnection
)

var (
	// Debug indicates that middleware should run in debug mode
	Debug bool
)
