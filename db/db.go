package db

import (
	"assignment-2/server/shared"
	"cloud.google.com/go/firestore" // Firestore-specific support
	"context"                       // State handling across API boundaries; part of native GoLang API
	"errors"
	firebase "firebase.google.com/go" // Generic firebase support
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
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
	firebaseAuth = "./serviceAccountKey.json"
	collection   = "dashboards"
)

/*
Handler for all database operations
*/
func HandleDB(w http.ResponseWriter, r *http.Request) {
	implementedMethods := []string{http.MethodPost, http.MethodGet}

	switch r.Method {
	case http.MethodPost:
		addDashboardConfigDocument(w, r)
	case http.MethodGet:
		displayDocument(w, r)
	default:
		http.Error(
			w, fmt.Sprintf(
				"REST Method '%s' not supported. Currently only '%v' are supported.", r.Method,
				implementedMethods,
			), http.StatusNotImplemented,
		)
		return
	}
}

/*
Reads a string from the body in plain-text and sends it to Firestore to be registered as a document.
*/
func addDashboardConfigDocument(w http.ResponseWriter, r *http.Request) {
	// very generic way of reading body; should be customized to specific use case
	content, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Reading payload from body failed.")
		http.Error(w, "Reading payload failed.", http.StatusInternalServerError)
		return
	}
	log.Println("Received request to add document for content ", string(content))
	if len(string(content)) == 0 {
		log.Println("Content appears to be empty.")
		http.Error(
			w,
			"Your payload (to be stored as document) appears to be empty. Ensure to terminate URI with /.",
			http.StatusBadRequest,
		)
		return
	} else {
		// Add element in embedded structure.
		// Note: this structure is defined by the client, not the server!; it exemplifies the use of a complex structure
		// and illustrates how you can use Firestore features such as Firestore timestamps.
		id, _, err2 := client.Collection(collection).Add(
			ctx,
			shared.DashboardConfig{},
		)
		if err2 != nil {
			// Error handling
			log.Println("Error when adding document " + string(content) + ", Error: " + err2.Error())
			http.Error(w, "Error when adding document "+string(content)+", Error: "+err2.Error(), http.StatusBadRequest)
			return
		} else {
			// Returns document ID in body
			log.Println("Document added to collection. Identifier of returned document: " + id.ID)
			http.Error(w, id.ID, http.StatusCreated)
			return
		}
	}
}

/*
Lists all the documents in the collection (see constant above) to the user.
*/
func displayDocument(w http.ResponseWriter, r *http.Request) {

	// Gets dashboard ID from given URL
	dashboardId := r.PathValue("id")

	if len(dashboardId) != 0 {
		// Extract individual dashboard

		// Retrieve specific dashboard based on id (Firestore-generated hash)
		res := client.Collection(collection).Doc(dashboardId)

		// Retrieve reference to document
		doc, err2 := res.Get(ctx)
		if err2 != nil {
			log.Println("Error extracting body of returned document of dashboard " + dashboardId)
			http.Error(
				w,
				"Error extracting body of returned document of dashboard "+dashboardId,
				http.StatusInternalServerError,
			)
			return
		}

		// A dashboard map with string keys. Each key is one field, like "content" or "timestamp"
		m := doc.Data()
		_, err3 := fmt.Fprintln(w, m["content"]) // here we retrieve the field containing the originally stored payload
		if err3 != nil {
			log.Println("Error while writing response body of dashboard " + dashboardId)
			http.Error(
				w, "Error while writing response body of dashboard "+dashboardId,
				http.StatusInternalServerError,
			)
			return
		}
	} else {
		// Collective retrieval of dashboards
		iter := client.Collection(collection).Documents(ctx) // Loop through all entries in collection "dashboards"

		for {
			doc, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				log.Printf("Failed to iterate: %v", err)
				return
			}
			// Note: You can access the document ID using "doc.Ref.ID"

			// A dashboard map with string keys. Each key is one field, like "content" or "timestamp"
			m := doc.Data()
			_, err = fmt.Fprintln(w, m["content"])
			if err != nil {
				log.Println("Error while writing response body (Error: " + err.Error() + ")")
				http.Error(
					w,
					"Error while writing response body (Error: "+err.Error()+")",
					http.StatusInternalServerError,
				)
				return
			}
		}
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

	// Close down client at the end of the function
	defer func() {
		errClose := client.Close()
		if errClose != nil {
			log.Fatal("Closing of the Firebase client failed. Error:", errClose)
		}
	}()
}
