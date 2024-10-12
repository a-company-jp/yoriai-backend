package firestore

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/a-company/yoriai-backend/pkg/config"
	"google.golang.org/api/option"
	"log"
	"os"
)

var (
	client *firestore.Client
)

func New() *firestore.Client {
	if client == nil {
		ctx := context.Background()
		if config.Config.Firestore.JsonCredentialFile == "" {
			// this means production environment
			client, _ = firestore.NewClient(ctx, config.Config.Firestore.ProjectID)
			return client
		}
		data, err := os.ReadFile(config.Config.Firestore.JsonCredentialFile)
		options := option.WithCredentialsJSON(data)
		dbID := "(default)"
		if config.Config.Firestore.DatabaseID != "" {
			dbID = config.Config.Firestore.DatabaseID
		}
		client, err = firestore.NewClientWithDatabase(ctx, config.Config.Firestore.ProjectID, dbID, options)
		if err != nil {
			log.Fatalf("firebase.NewClient err: %v", err)
		}
	}
	return client
}
