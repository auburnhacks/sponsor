// Package participant defines all the types and interfaces need to for
// handling operations realted to a participant at a hackathon
package participant

import (
	"context"
	"net/url"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/bsoncodec"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
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
	ID      objectid.ObjectID `json:"id" bson:"_id"`
	Profile struct {
		Name string `json:"name"`
	} `json:"profile"`
	Confirmation struct {
		Github   string `json:"github"`
		Linkedin string `json:"twitter"`
	} `json:"confirmation"`
	Resume string `json:"url"`
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
			pSlice, err := fetchParticipants(quillURI, resumesURI)
			if err != nil {
				log.Errorf("error while fetching participants: %v", err)
				continue
			}
			pSlice, err = fetchResumes(pSlice, resumesURI)
			if err != nil {
				log.Errorf("error while fetching resumes: %v", err)
				continue
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

func fetchParticipants(quillURI, resumesURI string) ([]*hacker, error) {
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

	pSlice := []*hacker{}
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
	return pSlice, nil
}

func fetchResumes(pSlice []*hacker, resumesDBURI string) ([]*hacker, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := mongo.NewClient(resumesDBURI)
	if err != nil {
		return nil, err
	}
	log.Debug("connecting to resumes mongodb")
	if err := client.Connect(ctx); err != nil {
		return nil, err
	}
	resumes := client.Database("resumes").Collection("resumes_19")
	for _, p := range pSlice {
		if len(p.Profile.Name) == 0 {
			continue
		}
		err = resumes.FindOne(ctx,
			bson.NewDocument(
				bson.EC.String("userid", p.ID.Hex()),
			),
		).Decode(&p)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue
			}
			return nil, err
		}
	}
	return pSlice, nil
}

func saveToDB(pSlice []Participant) error {
	log.Debug("saving participants to database")
	// delete all rows before inserting new ones
	// NOTE: DELETE FROM participants RETURNING *
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
