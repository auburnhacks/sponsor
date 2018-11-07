// Package participant defines all the types and interfaces need to for
// handling operations realted to a participant at a hackathon
package participant

import (
	"context"
	"net/url"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/bsoncodec"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrUnsupportedScheme is a error return by a function when the url contains a
	// protocol scheme that is unsupported by the system right now
	ErrUnsupportedScheme = errors.New("participant: this database scheme not supported")
)

// New returns a struct that implements the Participant interface
func New(name, linkedinProfile, githubProfile string) (Participant, error) {
	return nil, errors.New("error not implemented")
}

// Participant is an interface that any participant must satisfy
type Participant interface {
	Info() string
	Links() []string
}

type hacker struct {
	Profile struct {
		Name string `json:"name"`
	} `json:"profile"`
	Confirmation struct {
		Github   string `json:"github"`
		Linkedin string `json:"twitter"`
	} `json:"confirmation"`
	Resume string
}

func (p *hacker) Info() string {
	return p.Profile.Name
}

func (p *hacker) Links() []string {
	return []string{p.Confirmation.Github, p.Confirmation.Linkedin}
}

// Hacker also implements the participant interface
type Hacker struct {
	ID        string
	Name      string
	Github    string
	Linkedin  string
	Resume    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Sync should run as a goroutine it will run in the
// background thread and will uniformly sync all participants to the
// database for future use.
func Sync(quillURI, resumesURI string, d time.Duration, quit <-chan struct{}) error {
	log.Debugf("will sync from external db: %s, %s", quillURI, resumesURI)
	if err := validScheme(quillURI); err != nil {
		return err
	}
	if err := validScheme(resumesURI); err != nil {
		return err
	}
	go func() {
		for {
			pSlice, err := fetchParticipantsFromMongo(quillURI, resumesURI)
			if err != nil {
				log.Errorf("error while fetching participants: %v", err)
				return
			}
			if err := saveToDB(pSlice); err != nil {
				log.Errorf("erro while saving to database: %v", err)
				return
			}
			select {
			case <-time.Tick(d):
				continue
			case <-quit:
				return
			}
		}
	}()
	return nil
}

func fetchParticipantsFromMongo(quillURI, resumesURI string) ([]Participant, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Debug("creating mongo client")
	client, err := mongo.NewClient(quillURI)
	if err != nil {
		return nil, errors.Wrap(err, "participant: error while creating mongo client")
	}
	log.Debug("connecting to mongodb")
	if err := client.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "participant: error while connecting to mongodb")
	}
	// instance of the users collection on the mongodb for quill
	// TODO: make this more abstract
	users := client.Database("quill").Collection("users")
	cur, err := users.Find(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "participant: error while find participants from collection")
	}
	defer cur.Close(ctx)

	pSlice := []Participant{}
	for cur.Next(ctx) {
		br, err := cur.DecodeBytes()
		if err != nil {
			return nil, errors.Wrap(err, "participant: error while decoding bytes from mongo")
		}
		var h *hacker
		if err := bsoncodec.Unmarshal(br, &h); err != nil {
			return nil, errors.Wrap(err, "participant: error while unmarshaling bson to struct")
		}
		pSlice = append(pSlice, h)
	}
	// at this point we have all the information from quill
	// query the resumes database for other information
	return pSlice, nil
}

func saveToDB(pSlice []Participant) error {
	log.Debug("saving participants to database")
	// delete all rows before inserting new ones
	// NOTE: DELETE FROM participants RETURNING *
	for _, p := range pSlice {
		log.Debugf("%s", p.Info())
	}
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
