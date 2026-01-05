package googleauth

import (
	"context"

	"golang.org/x/oauth2/google"
)

type Checker interface {
	Check(ctx context.Context) error
}

type ADCChecker struct {
	Scopes []string
}

func (c ADCChecker) Check(ctx context.Context) error {
	creds, err := google.FindDefaultCredentials(ctx, c.Scopes...)
	if err != nil {
		return err
	}

	_, tokenErr := creds.TokenSource.Token()
	if tokenErr != nil {
		return tokenErr
	}

	return nil
}
