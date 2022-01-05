package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShortLink struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`

	ShortId string `bson:"shortId"`

	Url string `json:"url" bson:"url" binding:"required"`

	Views uint64 `bson:"views"`
}

func (s *ShortLink) Delete() error {
	_, err := GetLinkCollection().DeleteOne(Context, bson.M{"_id": s.Id})
	return err
}

func (s *ShortLink) IncrementViews() error {
	_, err := GetLinkCollection().UpdateOne(Context, bson.M{"_id": s.Id}, bson.M{"$inc": bson.M{"views": 1}})
	return err
}

func ResolveShortLink(shortId string) (*ShortLink, error) {
	var shortLink ShortLink
	err := GetLinkCollection().FindOne(Context, bson.M{"shortId": shortId}).Decode(&shortLink)
	return &shortLink, err
}

func GetShortLinks() ([]ShortLink, error) {
	var shortLinks []ShortLink
	cursor, err := GetLinkCollection().Find(Context, bson.M{})
	if err != nil {
		return shortLinks, err
	}

	err = cursor.All(Context, &shortLinks)
	return shortLinks, err
}
