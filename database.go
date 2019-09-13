package main

import (
	"gopkg.in/couchbase/gocb.v1"
)

// player holds information about a player.
type player struct {
	ID     string      `json:"id"`
	Team   *team       `json:"team"`
	Stats  *statistics `json:"stats"`
	Age    int         `json:"age"`
	Number string      `json:"number"`
	Email  string      `json:"email"`
}

// team holds all team related information.
type team struct {
	Name  string `json:"name"`
	Wins  int    `json:"wins"`
	Loses int    `json:"loses"`
}

// game contains the location of where it was played and the
// winning team.
type game struct {
	Location string `json:"location"`
	Winner   *team  `json:"winner"`
}

// season contains all games for that given period.
type season struct {
	Games []game `json:"games"`
}

// statistics represents player statistics.
type statistics struct {
	Goals   int `json:"goals"`
	Assists int `json:"assists"`
}

// database maintains database state.
type database struct {
	db *gocb.Cluster
}

// newDatabase creates a new value of type database pointer
// with a containing value to access the underlying datastore.
func newDatabase(c *config) (*database, error) {
	cluster, err := gocb.Connect("couchbase://" + c.db.host)
	if err != nil {
		return nil, err
	}
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: c.db.user,
		Password: c.db.pass,
	})
	return &database{
		db: cluster,
	}, nil
}

// AddPlayer adds the given player to the database.
func (d *database) AddPlayer(p *player) error {
	bucket, err := d.db.OpenBucket("players", "")
	if err != nil {
		return err
	}
	bucket.Manager("", "").CreatePrimaryIndex("", true, false)
	if _, err := bucket.Upsert("u:"+p.ID, p, 0); err != nil {
		return err
	}
	return nil
}

// Players retrieves all players from the database.
func (d *database) Players() ([]player, error) {
	bucket, err := d.db.OpenBucket("players", "")
	if err != nil {
		return nil, err
	}
	query := gocb.NewN1qlQuery("SELECT * FROM players")
	rows, err := bucket.ExecuteN1qlQuery(query, []interface{}{})
	if err != nil {
		return nil, err
	}
	var players []player
	var row interface{}
	for rows.Next(&row) {
		players = append(players, row)
	}
	return players, nil
}

// PlayerByID retrieves a player with the given ID.
func (d *database) PlayerByID(id string) (*player, error) {
	bucket, err := d.db.OpenBucket("players", "")
	if err != nil {
		return nil, err
	}
	var player player
	if _, err := bucket.Get("u:"+id, &player); err != nil {
		return nil, err
	}
	return &player, nil
}
