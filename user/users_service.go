package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"maclt/graphdb/neo4j/database"
)

type UserService struct {
	DatabaseDriver *database.DatabaseDriver
}

func (us *UserService) RegisterUser(c *gin.Context) {
	type request struct {
		Name string `json:"name" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	name := req.Name

	ctx := c.Request.Context()
	session := us.DatabaseDriver.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := "CREATE (u:User {name: $name}) RETURN u"
		_, err := tx.Run(ctx, query, map[string]interface{}{"name": name})
		return nil, err
	})

	if err != nil {
		log.Printf("Failed to create user: %v", err)
	}
}

func (us *UserService) GetUser(c *gin.Context) {
	ctx := c.Request.Context()
	session := us.DatabaseDriver.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	name := c.Param("name")

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := "MATCH (u:User {name: $name}) RETURN u"
		records, err := tx.Run(ctx, query, map[string]interface{}{"name": name})
		if err != nil {
			return nil, err
		}

		// Fetch the first record
		if records.Next(ctx) {
			node := records.Record().Values[0].(neo4j.Node)
			return node.Props, nil
		}
		return nil, records.Err() // Explicitly return nil if no record is found
	})

	if err != nil {
		log.Printf("Failed to get user: %v", err)
	}

	user := result.(map[string]interface{})

	c.JSON(http.StatusOK, user)
}

func (us *UserService) MarryUser(c *gin.Context) {
	type request struct {
		Wife    string `json:"wife" binding:"required"`
		Husband string `json:"husband" binding:"required"`
	}

	var req request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wife := req.Wife
	husband := req.Husband

	ctx := c.Request.Context()
	session := us.DatabaseDriver.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := `
		  MATCH (a:User {name: $name1}), (b:User {name: $name2})
		  CREATE (a)-[:MARRY]->(b), (b)-[:MARRY]->(a)
		`
		_, err := tx.Run(ctx, query, map[string]interface{}{"name1": wife, "name2": husband})
		return nil, err
	})

	if err != nil {
		log.Printf("Failed to marry users each other: %v", err)
	}
}

func (us *UserService) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	session := us.DatabaseDriver.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	name := c.Param("name")

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		query := "MATCH (n:User {name: $name}) DETACH DELETE n"
		_, err := tx.Run(ctx, query, map[string]interface{}{"name": name})
		return nil, err
	})

	if err != nil {
		log.Printf("Failed to delete user: %v", err)
	}
}
