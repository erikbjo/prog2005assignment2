package db

import (
	"assignment-2/server/shared"
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

// Collection names in Firestore
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
AddDocument Structures data by the provided struct and sends it to Firestore to be registered as a
document.
*/
func AddDocument[T any](
	data interface{}, collection string,
) error {

	// Assert type to target struct
	target, ok := data.(T)
	if !ok {
		return fmt.Errorf("data does not match target struct")
	}

	// Add document to Firestore
	_, _, err := client.Collection(collection).Add(ctx, target)
	if err != nil {
		return err
	}

	log.Println("Document was successfully added to DB")

	return nil
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
		doc, err2 := getDocumentByID(id, collection)
		if err2 != nil {
			log.Println("Error extracting body of returned document" + id)
			return nil, err2
		}

		if err4 := doc.DataTo(&data); err4 != nil {
			log.Println("Error unmarshalling document mapOfContent:", err4)
			return nil, err4
		}
		// A document map with string keys
	} else {
		log.Println("No valid ID was provided")
		return nil, fmt.Errorf("no valid ID was provided")
	}
	return data, nil
}

/*
GetAllDocuments Returns all documents in collection.
*/
func GetAllDocuments(collection string) (
	[]interface{},
	error,
) {
	// interface of document content
	var allData []interface{}

	// Collective retrieval of documents
	iter := client.Collection(collection).Documents(
		ctx,
	)
	// Loop through all entries in provided collection
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

		// Append the document to the slice
		allData = append(allData, mapOfContent)
	}
	return allData, nil
}

/*
UpdateDashboardConfigDocument Updates a document with the provided ID, if found.
*/
func UpdateDashboardConfigDocument(w http.ResponseWriter, r *http.Request, collection string) error {
	// TODO: make this use struct instead of map (because e.g. id key is saved as "id" in Firestore)
	var updates *shared.DashboardConfig

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

	if ok, err := documentExists(ctx, collection, documentID); ok && err == nil {

		// Adds id and lastChange field
		updates.ID = documentID
		updates.LastChange = time.Now()

		// Find document with matching ID
		foundDocument, err3 := getDocumentByID(documentID, collection)
		if err3 != nil {
			log.Println("Error trying to find document with ID: " + documentID)
			return err3
		}

		// Get the firebase ID of the document
		firebaseID := foundDocument.Ref.ID

		// TODO: maybe go back to update function, use a loop to update each key and value
		// Add element in embedded structure.
		// Note: this structure is defined by the client, not the server!; it exemplifies the use of a complex structure
		// and illustrates how you can use Firestore features such as Firestore timestamps.
		_, err2 := client.Collection(collection).Doc(firebaseID).Set(ctx, updates)
		if err2 != nil {
			// Error handling
			log.Printf("Error when updating document. Error: %s", err2.Error())
			return err2
		}
	} else if !ok && err == nil {
		log.Printf(
			"A document with the provided ID: %s, was not found in the collection: %s.\n",
			documentID, collection,
		)
	} else {
		log.Println("Error while trying to find document: ", err.Error())
		return err
	}
	return nil
}

/*
DeleteDocument Deletes a document with the provided ID, if found.
*/
func DeleteDocument(w http.ResponseWriter, r *http.Request, collection string) error {
	documentID := r.PathValue("id")

	// Checks if a document with the provided ID exists in the collection
	if ok, err := documentExists(ctx, collection, documentID); ok && err == nil {
		// Find document with matching ID
		foundDocument, err3 := getDocumentByID(documentID, collection)
		if err3 != nil {
			log.Println("Error trying to find document with ID: " + documentID)
			return err3
		}

		// Get the firebase ID of the document
		firebaseID := foundDocument.Ref.ID

		// Delete specified document
		_, err2 := client.Collection(collection).Doc(firebaseID).Delete(ctx)
		if err2 != nil {
			log.Println("Error while deleting document:" + documentID)
			return err2
		}
	} else if !ok && err == nil {
		log.Printf(
			"A document with the provided ID: %s, was not found in the collection: %s.\n",
			documentID, collection,
		)
	} else {
		log.Println("Error while trying to find document: ", err.Error())
		return err
	}
	log.Printf("The document: %s, was successfully deleted.", documentID)
	return nil
}

/*
documentExists Checks if a document with the provided ID exists in the collection.
*/
func documentExists(ctx context.Context, collection, documentID string) (bool, error) {
	// Query documents based on the "id" field
	iter := client.Collection(collection).Where("ID", "==", documentID).Documents(ctx)

	// Get the first document from the query iterator
	_, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return false, fmt.Errorf("document with ID %s not found in collection %s", documentID, collection)
		}
		return false, err
	}

	return true, nil
}

/*
getDocumentByID Retrieves a document with the provided ID from the collection.
*/
func getDocumentByID(id string, collection string) (*firestore.DocumentSnapshot, error) {
	// Query documents based on the "id" field
	iter := client.Collection(collection).Where("ID", "==", id).Documents(ctx)

	// Get the first document from the query iterator
	docSnap, err := iter.Next()
	if err != nil {
		if errors.Is(err, iterator.Done) {
			return nil, fmt.Errorf("document with ID %s not found in collection %s", id, collection)
		}
		return nil, err
	}

	return docSnap, nil
}

/*
Initialize Initializes the Firestore client.
*/
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

/*
Close Closes the Firestore client.
*/
func Close() {
	errClose := client.Close()
	if errClose != nil {
		log.Fatal("Closing of the Firebase client failed. Error:", errClose)
	}
}
