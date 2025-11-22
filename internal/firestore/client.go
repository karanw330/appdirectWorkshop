package firestore

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type Client struct {
	*firestore.Client
	subcollectionID string
}

// NewClient creates a new Firestore client
// For local development: uses service account file
// For Cloud Run: uses Application Default Credentials (ADC)
func NewClient(ctx context.Context, projectID, serviceAccountPath string) (*Client, error) {
	var app *firebase.App
	var err error

	// Check if we're running in Cloud Run (or using ADC)
	// If serviceAccountPath is empty or "ADC", use Application Default Credentials
	if serviceAccountPath == "" || serviceAccountPath == "ADC" {
		log.Println("Using Application Default Credentials (ADC)")
		app, err = firebase.NewApp(ctx, &firebase.Config{
			ProjectID: projectID,
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Use service account file for local development
		log.Printf("Using service account file: %s", serviceAccountPath)
		opt := option.WithCredentialsFile(serviceAccountPath)
		app, err = firebase.NewApp(ctx, &firebase.Config{
			ProjectID: projectID,
		}, opt)
		if err != nil {
			return nil, err
		}
	}

	client, err := app.Firestore(ctx)
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
