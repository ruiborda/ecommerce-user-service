package database

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"os"
)

var Client *firestore.Client

func GetFirestoreClient() *firestore.Client {
	if Client != nil {
		return Client
	}
	ctx := context.Background()

	GcpCredentialJsonBase64 := os.Getenv("GCP_CREDENTIAL_JSON_BASE64")
	if GcpCredentialJsonBase64 == "" {
		log.Fatalln("GCP_CREDENTIAL_JSON_BASE64 not found")
	}
	credentialJson, err := base64.StdEncoding.DecodeString(GcpCredentialJsonBase64)

	sa := option.WithCredentialsJSON(credentialJson)

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	Client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// todo use Client.Close() in main.go
	//defer client.Close()
	return Client
}
