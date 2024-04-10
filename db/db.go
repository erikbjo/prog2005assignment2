package db

import (
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
AddDocument Reads a string from the body in plain-text and sends it to Firestore to be registered as a
document.
*/
func AddDocument(w http.ResponseWriter, r *http.Request, collection string) (string, error) {
	content := map[string]interface{}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&content); err != nil {
		http.Error(w, "Error while decoding json.", http.StatusInternalServerError)
		log.Println("Error while decoding json: ", err.Error())
		return "", err
	}

	log.Println("Received request to add document for content ", fmt.Sprint(content))
	if content == nil {
		http.Error(
			w,
			"Your payload (to be stored as document) appears to be empty. Ensure to terminate URI with /.",
			http.StatusBadRequest,
		)
		return "", fmt.Errorf("content appears to be empty")
	} else {
		randomDocumentID := utils.GenerateRandomID()

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
			http.Error(
				w, "Error when adding document "+fmt.Sprint(content)+", Error: "+err2.Error(),
				http.StatusBadRequest,
			)
			return "", err2
		} else {
			// Returns document ID in body
			http.Error(w, randomDocumentID, http.StatusCreated)
			return randomDocumentID, fmt.Errorf(
				"Document added to dashboardCollection. " +
					"Identifier of returned document: " + randomDocumentID,
			)

		}
	}
}

/*
DisplayDocument Returns a document if specific ID is provided or all documents in collection.
*/
func DisplayDocument(
	w http.ResponseWriter,
	r *http.Request,
	collection string,
) (map[string]interface{}, error) {
	// Gets document ID from given URL
	documentId := r.PathValue("id")

	// map to store found document
	var m map[string]interface{}

	if len(documentId) != 0 {
		// Extract individual document

		// Retrieve specific document based on id
		res := client.Collection(collection).Doc(documentId)

		// Retrieve reference to document
		doc, err2 := res.Get(ctx)
		if err2 != nil {
			log.Println("Error extracting body of returned document" + documentId)
			http.Error(
				w,
				"Error extracting body of returned document"+documentId,
				http.StatusInternalServerError,
			)
			return nil, err2
		}

		// A document map with string keys. Each key is one field, like "content" or "timestamp"
		m = doc.Data()
		_, err3 := fmt.Fprintln(
			w,
			m["content"],
		) // here we retrieve the field containing the originally stored payload
		if err3 != nil {
			log.Println("Error while writing response body of document " + documentId)
			http.Error(
				w, "Error while writing response body of document "+documentId,
				http.StatusInternalServerError,
			)
			return nil, err3
		}
	} else {
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
			// Note: You can access the document ID using "doc.Ref.ID"

			// A document map with string keys. Each key is one field, like "content" or "timestamp"
			m = doc.Data()
			_, err = fmt.Fprintln(w, m["content"])
			if err != nil {
				log.Println("Error while writing response body (Error: " + err.Error() + ")")
				http.Error(
					w,
					"Error while writing response body (Error: "+err.Error()+")",
					http.StatusInternalServerError,
				)
				return nil, err
			}
		}
	}
	return m, nil
}

func UpdateDocument(w http.ResponseWriter, r *http.Request, collection string) error {
	var newContent []firestore.Update

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newContent); err != nil {
		http.Error(w, "Error while decoding json.", http.StatusInternalServerError)
		log.Println("Error while decoding json: ", err.Error())
		return err
	}
	if newContent == nil {
		http.Error(
			w,
			"Your payload (to be stored as document) appears to be empty. Ensure to terminate URI with /.",
			http.StatusBadRequest,
		)
		return fmt.Errorf("content appears to be empty")
	}

	// Get ID from the URL provided in the request
	documentId := r.PathValue("id")

	// Add element in embedded structure.
	// Note: this structure is defined by the client, not the server!; it exemplifies the use of a complex structure
	// and illustrates how you can use Firestore features such as Firestore timestamps.
	_, err2 := client.Collection(collection).Doc(documentId).Update(ctx, newContent)
	if err2 != nil {
		// Error handling
		log.Println("Error when adding document " + fmt.Sprint(newContent) + ", Error: " + err2.Error())
		http.Error(
			w, "Error when adding document "+fmt.Sprint(newContent)+", Error: "+err2.Error(),
			http.StatusBadRequest,
		)
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
			http.Error(
				w,
				"Error extracting body of returned document"+documentId,
				http.StatusInternalServerError,
			)
			return err2
		}
	} else if !ok && err == nil {
		http.Error(
			w, fmt.Sprintf(
				"A document with the provided ID: %s, was not found in the collection: %s.\n",
				documentId, collection,
			), http.StatusBadRequest,
		)
		log.Printf(
			"A document with the provided ID: %s, was not found in the collection: %s.\n",
			documentId, collection,
		)
	} else {
		http.Error(w, "Error while trying to find document.", http.StatusInternalServerError)
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
