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
	firebaseAuth        = "./serviceAccountKey.json"
	dashboardCollection = "dashboards"
)

/*
Reads a string from the body in plain-text and sends it to Firestore to be registered as a document.
*/
func AddDashboardConfigDocument(w http.ResponseWriter, r *http.Request) (string, error) {
	// very generic way of reading body; should be customized to specific use case
	// e.g. decode the body into dashboard config

	/*
		contentBuffer, err := io.Copy(&contentBuffer)
		if err != nil {
			log.Println("Reading payload from body failed.")
			http.Error(w, "Reading payload failed.", http.StatusInternalServerError)
			return
		}
	*/

	config := map[string]interface{}{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&config); err != nil {
		http.Error(w, "Error while decoding json.", http.StatusInternalServerError)
		log.Println("Error while decoding json: ", err.Error())
		return "", err
	}

	log.Println("Received request to add document for content ", fmt.Sprint(config))
	if len(fmt.Sprint(config)) == 0 {
		http.Error(
			w,
			"Your payload (to be stored as document) appears to be empty. Ensure to terminate URI with /.",
			http.StatusBadRequest,
		)
		return "", fmt.Errorf("content appears to be empty")
	} else {
		randomDocID := utils.GenerateRandomID()

		// Add element in embedded structure.
		// Note: this structure is defined by the client, not the server!; it exemplifies the use of a complex structure
		// and illustrates how you can use Firestore features such as Firestore timestamps.
		_, err2 := client.Collection(dashboardCollection).Doc(randomDocID).Set(
			ctx,
			config,
		)
		if err2 != nil {
			// Error handling
			log.Println("Error when adding document " + fmt.Sprint(config) + ", Error: " + err2.Error())
			http.Error(
				w, "Error when adding document "+fmt.Sprint(config)+", Error: "+err2.Error(),
				http.StatusBadRequest,
			)
			return "", err2
		} else {
			// Returns document ID in body
			http.Error(w, randomDocID, http.StatusCreated)
			return randomDocID, fmt.Errorf(
				"Document added to dashboardCollection. " +
					"Identifier of returned document: " + randomDocID,
			)

		}
	}
}

/*
Lists all the documents in the dashboardCollection (see constant above) to the user.
*/
func DisplayDocument(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("entered displayDocument")
	// Gets dashboard ID from given URL
	dashboardId := r.PathValue("id")

	// map to store found dashboard configs
	var m map[string]interface{}

	if len(dashboardId) != 0 {
		// Extract individual dashboard

		// Retrieve specific dashboard based on id (Firestore-generated hash)
		res := client.Collection(dashboardCollection).Doc(dashboardId)

		// Retrieve reference to document
		doc, err2 := res.Get(ctx)
		if err2 != nil {
			log.Println("Error extracting body of returned document of dashboard " + dashboardId)
			http.Error(
				w,
				"Error extracting body of returned document of dashboard "+dashboardId,
				http.StatusInternalServerError,
			)
			return err2
		}

		// A dashboard map with string keys. Each key is one field, like "content" or "timestamp"
		m = doc.Data()
		_, err3 := fmt.Fprintln(w, m["content"]) // here we retrieve the field containing the originally stored payload
		if err3 != nil {
			log.Println("Error while writing response body of dashboard " + dashboardId)
			http.Error(
				w, "Error while writing response body of dashboard "+dashboardId,
				http.StatusInternalServerError,
			)
			return err3
		}

		// Encodes the found document
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(m); err != nil {
			http.Error(w, "Error while encoding to json.", http.StatusInternalServerError)
			log.Println("Error while encoding to json: ", err.Error())
		}
	} else {
		// Collective retrieval of dashboards
		iter := client.Collection(dashboardCollection).Documents(ctx) // Loop through all entries in dashboardCollection "dashboards"

		for {
			doc, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				log.Printf("Failed to iterate: %v", err)
				return err
			}
			// Note: You can access the document ID using "doc.Ref.ID"

			// A dashboard map with string keys. Each key is one field, like "content" or "timestamp"
			m = doc.Data()
			_, err = fmt.Fprintln(w, m["content"])
			if err != nil {
				log.Println("Error while writing response body (Error: " + err.Error() + ")")
				http.Error(
					w,
					"Error while writing response body (Error: "+err.Error()+")",
					http.StatusInternalServerError,
				)
				return err
			}
		}
	}
	return nil
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
