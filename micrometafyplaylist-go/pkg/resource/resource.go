package resource

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/mongodb/mongo-go-driver/bson/bsoncodec"

	"github.com/gin-gonic/gin"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"

	"github.com/mongodb/mongo-go-driver/mongo"
)

var dbPort int
var dbHost string

func init() {
	flag.IntVar(&dbPort, "dbport", 27017, "database connection port")
	flag.StringVar(&dbHost, "dbhost", "db", "database connection host")
}

type Resource struct {
	client     *mongo.Client
	collection *mongo.Collection
}

var resourceInstance *Resource
var resourceOnce sync.Once

func GetResourceInstance() *Resource {
	resourceOnce.Do(func() {
		client, err := mongo.NewClient(fmt.Sprintf("mongodb://%s:%d", dbHost, dbPort))
		if err != nil {
			log.Fatal(err)
		}
		err = client.Connect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}
		resourceInstance = &Resource{client, client.Database("micrometafyplaylist").Collection("playlists")}
	})
	return resourceInstance
}

func (r *Resource) Close() {
	r.client.Disconnect(context.TODO())
	log.Print("Close happened")
}

type dbPlaylist struct {
	ID     objectid.ObjectID `bson:"_id"`
	Name   string            `bson:"name"`
	Tracks []track           `bson:"tracks"`
}

func (p dbPlaylist) jsonPlaylist() jsonPlaylist {
	return jsonPlaylist{p.ID.Hex(), p.Name, p.Tracks}
}

type jsonPlaylist struct {
	ID     string  `json:"id"`
	Name   string  `json:"name" binding:"required"`
	Tracks []track `json:"tracks" binding:"required"`
}

type track struct {
	Name     string `bson:"name" json:"name" binding:"required"`
	Author   string `bson:"author" json:"author" binding:"required"`
	URL      string `bson:"url" json:"url" binding:"required"`
	Duration int64  `bson:"duration" json:"duration" binding:"required"`
	Origin   string `bson:"origin" json:"origin" binding:"required"`
}

func (r *Resource) NewPlaylistPOST(c *gin.Context) {
	var json jsonPlaylist
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := r.collection.InsertOne(context.Background(),
		dbPlaylist{ID: objectid.New(), Name: json.Name, Tracks: json.Tracks})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatal(err)
	}
	id := res.InsertedID.(*bson.Element)
	json.ID = id.Value().ObjectID().Hex()
	c.JSON(http.StatusOK, json)
}

func (r *Resource) PlaylistByIDGET(c *gin.Context) {
	id, err := objectid.FromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var result dbPlaylist
	filter := bson.NewDocument(bson.EC.ObjectID("_id", id))
	err = r.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result.jsonPlaylist())
}

func (r *Resource) AllPlaylistsGET(c *gin.Context) {
	filter := bson.NewDocument()
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var result []jsonPlaylist
	for cursor.Next(context.Background()) {
		var tmp dbPlaylist
		err = cursor.Decode(&tmp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		result = append(result, tmp.jsonPlaylist())
	}
	c.JSON(http.StatusOK, result)
}

func (r *Resource) RemovePlaylistDELETE(c *gin.Context) {
	id, err := objectid.FromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var result dbPlaylist
	filter := bson.NewDocument(bson.EC.ObjectID("_id", id))
	err = r.collection.FindOneAndDelete(context.Background(), filter).Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result.jsonPlaylist())
}

func (r *Resource) AddTrackToPlaylistPUT(c *gin.Context) {
	var json track
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := objectid.FromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tmp, err := bsoncodec.Marshal(json)
	jsonDocument, err := bson.ReadDocument(tmp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var result dbPlaylist
	filter := bson.NewDocument(bson.EC.ObjectID("_id", id))
	err = r.collection.FindOneAndUpdate(context.Background(), filter,
		bson.NewDocument(bson.EC.SubDocumentFromElements("$addToSet",
			bson.EC.SubDocument("tracks", jsonDocument)))).Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result.jsonPlaylist())
}
