package db

import (
	"assignment-2/server/shared"
	"assignment-2/server/utils"
	"cloud.google.com/go/firestore" // Firestore-specific support
	"context"                       // State handling across API boundaries; part of native GoLang API
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go" // Generic firebase support
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"time"
)

/*
This server shows an example of how to interact with Firebase directly, including
storing and retrieval of content.
*/

// Firebase context and client used by Firestore functions throughout the program.
var ctx context.Context
var client *firestore.Client

// Collection name in Firestore
const (
	firebaseAuth           = "./serviceAccountKey.json"
	DashboardCollection    = "dashboards"
	NotificationCollection = "notifications"
)

func GetStatusCodeOfCollection(w http.ResponseWriter, collection string) int {
	// Check if the Firestore client is initialized
	if client == nil {
		// If client is nil, return 503 Service Unavailable status code
		http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
		return http.StatusServiceUnavailable
	}

	// Check if the Firestore client is connected by performing a simple query
	iter := client.Collection(collection).Documents(ctx)
	defer iter.Stop()

	// Attempt to retrieve the first document
	_, err := iter.Next()
	if err != nil {
		// If there's an error connecting to the database, return 503 Service Unavailable status code
		http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
		return http.StatusServiceUnavailable
	}

	// If code reaches this point, the database is available
	return http.StatusOK
}

/*
AddDashboardConfigDocument Reads a string from the body in plain-text and sends it to Firestore to be registered as a
document.
*/
func AddDashboardConfigDocument(w http.ResponseWriter, r *http.Request, collection string) (
	string,
	*shared.DashboardConfig, error,
) {
	var content *shared.DashboardConfig

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&content); err != nil {
		log.Println("Error while decoding json: ", err.Error())
		return "", nil, err
	}

	log.Println("Received request to add document for content: ", content)
	if content == nil {
		log.Println("content appears to be empty")
		return "", nil, fmt.Errorf("content appears to be empty")
	} else {
		randomDocumentID := utils.GenerateRandomID()

		// Sets the ID field to the new document ID
		content.ID = randomDocumentID

		// Sets the lastChange field to the current time stamp
		content.LastChange = time.Now()

		// Add element in embedded structure.
		// Note: this structure is defined by the client, not the server!; it exemplifies the use of a complex structure
		// and illustrates how you can use Firestore features such as Firestore timestamps.
		_, err2 := client.Collection(collection).Doc(randomDocumentID).Set(
			ctx,
			content,
		)
		if err2 != nil {
			// Error handling
			log.Println("Error when adding document " + fmt.Sprint(content) + ", Error: " + err2.Error())
			return "", nil, err2
		} else {
			// Returns document ID, and map of content
			return randomDocumentID, content, nil
		}
	}
}

/*
GetDocument Returns the document that matches with the provided ID from a collection
*/
func GetDocument(
	id string,
	collection string,
) (interface{}, error) {
	// interface of document content
	var data interface{}

	if len(id) != 0 {
		// Extract individual document

		// Retrieve specific document based on id
		res := client.Collection(collection).Doc(id)

		// Retrieve reference to document
		doc, err2 := res.Get(ctx)
		if err2 != nil {
			log.Println("Error extracting body of returned document" + id)
			return nil, err2
		}

		var mapOfContent map[string]interface{}
		if err4 := doc.DataTo(&mapOfContent); err4 != nil {
			log.Println("Error unmarshalling document mapOfContent:", err4)
			return nil, err4
		}
		// A document map with string keys
		data = mapOfContent
		fmt.Printf("content is: %v", data)
	} else {
		log.Println("No valid ID was provided")
		return nil, fmt.Errorf("no valid ID was provided")
	}
	return data, nil
}

/*
GetAllDocuments Returns all documents in collection.
*/
func GetAllDocuments(w http.ResponseWriter, r *http.Request, collection string) (
	[]interface{},
	error,
) {
	// interface of document content
	var data interface{}
	var allData []interface{}

	// Collective retrieval of documents
	iter := client.Collection(collection).Documents(
		ctx,
	) // Loop through all entries in provided collection

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			log.Printf("Failed to iterate: %v", err)
			return nil, err
		}

		var mapOfContent map[string]interface{}
		if err3 := doc.DataTo(&mapOfContent); err3 != nil {
			log.Println("Error unmarshalling document mapOfContent:", err3)
			return nil, err3
		}

		// A document map with string keys. Each key is one field, like "content" or "timestamp"
		data = mapOfContent
		allData = append(allData, data)
	}
	return allData, nil
}

func UpdateDocument(w http.ResponseWriter, r *http.Request, collection string) error {
	// TODO: Update lastChange field to new current time
	var updates map[string]interface{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updates); err != nil {
		log.Println("Error while decoding json: ", err.Error())
		return err
	}
	if updates == nil {
		return fmt.Errorf("content appears to be empty")
	}

	// Get ID from the URL provided in the request
	// TODO: maybe use utils.GetIDFromRequest(r) or take id as a parameter
	documentID := r.PathValue("id")

	// Adds id and lastChange field
	updates["id"] = documentID
	updates["lastChange"] = time.Now()

	// TODO: maybe go back to update function, use a loop to update each key and value
	// Add element in embedded structure.
	// Note: this structure is defined by the client, not the server!; it exemplifies the use of a complex structure
	// and illustrates how you can use Firestore features such as Firestore timestamps.
	_, err2 := client.Collection(collection).Doc(documentID).Set(ctx, updates)
	if err2 != nil {
		// Error handling
		log.Printf("Error when updating document. Error: %s", err2.Error())
		return err2
	}
	return nil
}

// DeleteDocument with the provided ID, if found.
func DeleteDocument(w http.ResponseWriter, r *http.Request, collection string) error {
	documentId := r.PathValue("id")

	// Checks if a document with the provided ID exists in the collection
	if ok, err := documentExists(ctx, collection, documentId); ok && err == nil {
		// Delete specified document
		_, err2 := client.Collection(collection).Doc(documentId).Delete(ctx)
		if err2 != nil {
			log.Println("Error extracting body of returned document" + documentId)
			return err2
		}
	} else if !ok && err == nil {
		log.Printf(
			"A document with the provided ID: %s, was not found in the collection: %s.\n",
			documentId, collection,
		)
	} else {
		log.Println("Error while trying to find document: ", err.Error())
		return err
	}
	return nil
}

// documentExists checks if a document exists in a Firestore collection.
func documentExists(ctx context.Context, collection, documentID string) (bool, error) {
	iter := client.Collection(collection).Where("id", "==", documentID).Documents(ctx)
	defer iter.Stop()

	// Iterate over the result to see if any document matches the ID.
	for {
		_, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			// Document not found.
			return false, nil
		}
		if err != nil {
			return false, err
		}
		// Document found.
		return true, nil
	}
}

func Initialize() {
	// Firebase initialization
	ctx = context.Background()

	// We use a service account, load credentials file that you downloaded from your project's settings menu.
	// It should reside in your project directory.
	// Make sure this file is git-ignored, since it is the access token to the database.
	sa := option.WithCredentialsFile(firebaseAuth)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Println(err)
		return
	}

	// Instantiate client
	client, err = app.Firestore(ctx)

	// Alternative setup, directly through Firestore (without initial reference to Firebase); but requires Project ID; useful if multiple projects are used
	// client, err := firestore.NewClient(ctx, projectID)

	// Check whether there is an error when connecting to Firestore
	if err != nil {
		log.Println(err)
		return
	}
}

// Close down client
func Close() {
	errClose := client.Close()
	if errClose != nil {
		log.Fatal("Closing of the Firebase client failed. Error:", errClose)
	}
}
