package database

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type DatabaseDriver struct {
	Driver neo4j.DriverWithContext
}

// Constructor function to create a new database driver instance
func Connect() *DatabaseDriver {
	uri := "bolt://localhost:7687" // Bolt protocol
	username := "neo4j"
	password := "password"

	// Create a driver instance
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		log.Fatalf("Failed to create driver: %v", err)
	}

	// Return the driver instance
	return &DatabaseDriver{Driver: driver}
}

// Close function to close the database driver
func (d *DatabaseDriver) Close(c context.Context) {
	if err := d.Driver.Close(c); err != nil {
		log.Printf("Failed to close Neo4j driver: %v", err)
	}
}
