package api

import (
	"fmt"
	"tiny-go/db"
	"tiny-go/util"

	"github.com/aidarkhanov/nanoid/v2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Engine *gin.Engine = gin.Default()
)

func SetupService(config *util.Configuration) {
	Engine.SetTrustedProxies(nil)

	// Get all links
	Engine.GET("/api/list", apiList)

	// Short links
	Engine.POST("/api/shorten", apiShorten)

	// Resolve short links
	Engine.GET("/api/:shortId", api)

	Engine.Run(fmt.Sprintf(":%d", config.Service.Port))
}

func apiList(c *gin.Context) {
	shortLinks, err := db.GetShortLinks()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, shortLinks)
}

func apiShorten(c *gin.Context) {
	var shortLink db.ShortLink
	if err := c.ShouldBindJSON(&shortLink); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if shortId, err := nanoid.GenerateString(nanoid.DefaultAlphabet, 6); err == nil {
		shortLink.Id = primitive.NewObjectID()
		shortLink.ShortId = shortId
		shortLink.Views = 0
	} else {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if _, err := db.GetLinkCollection().InsertOne(c, &shortLink); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"shortId": &shortLink.ShortId})
}

func api(c *gin.Context) {
	shortId, isSet := c.Params.Get("shortId")
	if !isSet {
		c.JSON(400, gin.H{"error": "shortId is required"})
		return
	}

	shortLink, err := db.ResolveShortLink(shortId)
	if err != nil {
		c.JSON(404, gin.H{"error": "Not found."})
		return
	}

	c.Redirect(302, shortLink.Url)

	if err := shortLink.IncrementViews(); err != nil {
		panic(err)
	}
}
