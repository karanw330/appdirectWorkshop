package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

type Client struct {
	*firestore.Client
	subcollectionID string
}

// NewClient creates a new Firestore client
// For local development: uses service account file
// For Cloud Run: uses Application Default Credentials (ADC)
func NewClient(ctx context.Context, projectID, databaseID, serviceAccountPath string) (*Client, error) {
	var client *firestore.Client
	var err error

	// If using default database, we can keep it simple or unify the logic.
	// However, firebase.NewApp doesn't easily expose WithDatabase. 
	// So we switch to using the firestore client directly which supports it more transparently,
	// or we use the specific constructor if available. 
	// Given the potential dependency versions, let's try to construct it with options.

	opts := []option.ClientOption{}
	
	if serviceAccountPath != "" && serviceAccountPath != "ADC" {
		log.Printf("Using service account file: %s", serviceAccountPath)
		opts = append(opts, option.WithCredentialsFile(serviceAccountPath))
	} else {
		log.Println("Using Application Default Credentials (ADC)")
	}

	if databaseID != "" && databaseID != "(default)" {
		log.Printf("Connecting to specific database: %s", databaseID)
		client, err = firestore.NewClientWithDatabase(ctx, projectID, databaseID, opts...)
	} else {
		client, err = firestore.NewClient(ctx, projectID, opts...)
	}

	if err != nil {
		return nil, err
	}

	return &Client{
		Client:          client,
		subcollectionID: "workshop_attendees", // Default, can be overridden
	}, nil
}

func (c *Client) SetSubcollectionID(id string) {
	c.subcollectionID = id
}

func (c *Client) GetSubcollectionID() string {
	return c.subcollectionID
}

// Helper function to get collection reference
func (c *Client) GetCollection(ctx context.Context, name string) *firestore.CollectionRef {
	return c.Collection(name)
}

// Helper function to get subcollection reference
func (c *Client) GetSubcollection(ctx context.Context, parentDocID, subcollectionName string) *firestore.CollectionRef {
	return c.Collection("workshops").Doc(parentDocID).Collection(subcollectionName)
}
