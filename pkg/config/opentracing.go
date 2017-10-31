package config

import (
	"context"

	"github.com/hellofresh/gcloud-opentracing"
	"github.com/opentracing/basictracer-go"
)

// GoogleCloudTracing holds the Google Application Default Credentials
type GoogleCloudTracing struct {
	ProjectID    string `default:"" envconfig:"TRACING_GC_PROJECT_ID"`
	Email        string `default:"" envconfig:"TRACING_GC_EMAIL"`
	PrivateKey   string `default:"" envconfig:"TRACING_GC_PRIVATE_KEY"`
	PrivateKeyID string `default:"" envconfig:"TRACING_GC_PRIVATE_ID"`
}

// IsValid indicates if the google cloud configuration is valid
func (c *GoogleCloudTracing) IsValid() bool {
	return len(c.Email) > 0 &&
		len(c.PrivateKey) > 0 &&
		len(c.PrivateKeyID) > 0 &&
		len(c.ProjectID) > 0
}

// NewRecorder Create a new google cloud tracing recorder based of the config
func (c *GoogleCloudTracing) NewRecorder(logger gcloudtracer.Logger) (*gcloudtracer.Recorder, error) {
	return gcloudtracer.NewRecorder(
		context.Background(),
		gcloudtracer.WithLogger(logger),
		gcloudtracer.WithProject(c.ProjectID),
		gcloudtracer.WithJWTCredentials(gcloudtracer.JWTCredentials{
			Email:        c.Email,
			PrivateKey:   []byte(c.PrivateKey),
			PrivateKeyID: c.PrivateKeyID,
		}),
	)
}

// NoopRecorder implements the basictracer.Recorder interface.
type NoopRecorder struct{}

// RecordSpan complies with the basictracer.Recorder interface.
func (t *NoopRecorder) RecordSpan(span basictracer.RawSpan) {}
