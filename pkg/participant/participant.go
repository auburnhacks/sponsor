// Package participant defines all the types and interfaces need to for
// handling operations realted to a participant at a hackathon
package participant

import (
	"context"
	"net/url"
	"strconv"
	"time"

	"github.com/auburnhacks/sponsor/pkg/db"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrUnsupportedScheme is a error return by a function when the url contains a
	// protocol scheme that is unsupported by the system right now
	ErrUnsupportedScheme = errors.New("participant: this database scheme not supported")
)

// Participant also implements the participant interface
type Participant struct {
	ID         string    `db:"id"`
	Name       string    `db:"name"`
	Email      string    `db:"email"`
	University string    `db:"university"`
	Major      string    `db:"major"`
	GradYear   int       `db:"grad_year"`
	Github     string    `db:"github"`
	Linkedin   string    `db:"linkedin"`
	Resume     string    `db:"resume"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type hacker struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Email   string             `json:"email" bson:"email"`
	Profile struct {
		Name       string `json:"name" bson:"name"`
		University string `json:"school" bson:"school"`
		GradYear   string `json:"graduationYear" bson:"graduationYear"`
	} `json:"profile" bson:"profile"`
	Confirmation struct {
		Github   string `json:"github" bson:"github"`
		Linkedin string `json:"twitter" bson:"twitter"`
		Major    string `json:"major" bson:"major"`
	} `json:"confirmation" bson:"confirmation"`
	Resume string `json:"url"`
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
			}
			pSlice, err = fetchResumes(pSlice, resumesURI)
			if err != nil {
				log.Errorf("error while fetching resumes: %v", err)
			}
			if err := saveToDB(pSlice); err != nil {
				log.Errorf("error while saving participants to the database: %v", err)
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
	cur, err := users.Find(ctx, bson.D{})
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
		if err := bson.Unmarshal(br, &h); err != nil {
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
		var test map[string]interface{}
		err = resumes.FindOne(ctx, bson.D{
			{"userid", p.ID.Hex()},
		}).Decode(&test)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				continue
			}
			return nil, err
		}
		url, ok := test["url"]
		if !ok {
			continue
		}
		p.Resume = url.(string)
	}
	return pSlice, nil
}

func saveToDB(pSlice []*hacker) error {
	log.Debug("saving participants to database")
	// delete all rows before inserting new ones
	// NOTE: DELETE FROM participants RETURNING *
	_, err := db.Conn.Exec("DELETE FROM participants")
	if err != nil {
		return errors.Wrap(err, "participant: error while deleteing old participants")
	}
	for _, h := range pSlice {
		query := `
		INSERT INTO participants
		(name, email, university, major, grad_year, github, linkedin, resume_url)
		VALUES(:name, :email, :university, :major, :grad_year, :github, :linkedin, :resume)`
		stmt, err := db.Conn.PrepareNamed(query)
		if err != nil {
			return errors.Wrap(err, "participant: error while inserting participant")
		}
		var p interface{}
		gradYear := 0
		if h.Profile.GradYear != "" {
			gradYear, err = strconv.Atoi(h.Profile.GradYear)
			if err != nil {
				return errors.Wrap(err, "participant: error while converting from string to integer")
			}
		}
		stmt.QueryRow(map[string]interface{}{
			"name":       h.Profile.Name,
			"email":      h.Email,
			"university": h.Profile.University,
			"major":      h.Confirmation.Major,
			"grad_year":  gradYear,
			"github":     h.Confirmation.Github,
			"linkedin":   h.Confirmation.Linkedin,
			"resume":     h.Resume,
		}).Scan(&p)
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

// List is a function that returns a slice of all participants
func List() ([]Participant, error) {
	query := `SELECT * FROM participants`
	rows, err := db.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	log.Debugf("%+v", rows)
	var pSlice []Participant
	for rows.Next() {
		var p Participant
		err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.University, &p.Major, &p.GradYear, &p.Github,
			&p.Linkedin, &p.Resume, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		pSlice = append(pSlice, p)
	}
	return pSlice, nil
}
