// Package participant defines all the types and interfaces need to for
// handling operations realted to a participant at a hackathon
package participant

import (
	"net/url"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrUnsupportedScheme is a error return by a function when the url contains a
	// protocol scheme that is unsupported by the system right now
	ErrUnsupportedScheme = errors.New("participant: this database scheme not supported")
)

// SyncFromExternalDB should run as a goroutine it will run in the
// background thread and will uniformly sync all participants to the
// database for future use.
func SyncFromExternalDB(dbURI string, d time.Duration, quit <-chan struct{}) error {
	log.Debugf("will sync from external db: %s", dbURI)
	if err := validScheme(dbURI); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-time.Tick(d):
				log.Debug("starting sync process...")
			case <-quit:
				log.Debug("quitting sync process...")
				return
			}
		}
	}()
	return nil
}

func validScheme(dbURI string) error {
	url, err := url.Parse(dbURI)
	if err != nil {
		return errors.Wrap(err, "participant: error while parsing url for sheme validation")
	}
	if len(url.Scheme) > 0 && url.Scheme != "mongodb" {
		return ErrUnsupportedScheme
	}
	return nil
}
