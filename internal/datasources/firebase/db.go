package firebase

import (
	"assignment-2/internal/utils"
	"cloud.google.com/go/firestore" // Firestore-specific support
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"context" // State handling across API boundaries; part of native GoLang API
	"errors"
	firebase "firebase.google.com/go" // Generic firebase support
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
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

type dummyStruct struct {
	Dummy string
	ID    string
}

func GetStatusCodeOfCollection(w http.ResponseWriter, collection string) int {
	// Check if the Firestore client is initialized
	if client == nil {
		// If client is nil, return 503 Service Unavailable status code
		// http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
		log.Println("Firestore client is not initialized")
		return http.StatusServiceUnavailable
	}

	// Send a dummy document to the collection to check if the database is available
	id := utils.GenerateRandomID()
	err := AddDocument[dummyStruct](
		dummyStruct{Dummy: "dummy", ID: id},
		collection,
	)
	if err != nil {
		log.Println("Error while trying to add dummy document to the collection: ", err.Error())
		return http.StatusServiceUnavailable
	}

	defer func() {
		err := DeleteDocument(id, collection)
		if err != nil {
			log.Println(
				"Error while trying to delete dummy document from the collection: ",
				err.Error(),
			)
		}
	}()

	// Check if the Firestore client is connected by performing a simple query
	iter := client.Collection(collection).Documents(ctx)
	defer iter.Stop()

	// Attempt to retrieve the first document
	_, err = iter.Next()
	if err != nil {
		// If there's an error connecting to the database, return 503 Service Unavailable status code
		// http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
		log.Println(
			"Error while trying to retrieve the dummy document from the collection: ",
			err.Error(),
		)
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

	return nil
}

/*
GetDocument Returns the document that matches with the provided ID from a collection
*/
func GetDocument[T any](
	id string,
	collection string,
) (T, error) {
	// interface of document content
	var data T

	if len(id) != 0 {

		// Extract individual document
		doc, err2 := getDocumentByID(id, collection)
		if err2 != nil {
			log.Println("Error extracting body of returned document" + id)
			return data, err2
		}

		if err4 := doc.DataTo(&data); err4 != nil {
			log.Println("Error unmarshalling document mapOfContent:", err4)
			return data, err4
		}
		// A document map with string keys
	} else {
		log.Println("No valid ID was provided")
		return data, fmt.Errorf("no valid ID was provided")
	}
	return data, nil
}

/*
GetAllDocuments Returns all documents in collection.
*/
func GetAllDocuments[T any](collection string) (
	[]T,
	error,
) {
	// interface of document content
	var allData []T

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

		var data T
		if err3 := doc.DataTo(&data); err3 != nil {
			log.Println("Error unmarshalling document data:", err3)
			return nil, err3
		}

		// Append the document to the slice
		allData = append(allData, data)
	}
	return allData, nil
}

/*
UpdateDocument Updates a document with the provided ID, if found.
*/
func UpdateDocument[T any](
	updatedDocument interface{},
	documentID string,
	collection string,
) error {
	if ok, err := documentExists(ctx, collection, documentID); ok && err == nil {

		// Find document with matching ID
		foundDocument, err3 := getDocumentByID(documentID, collection)
		if err3 != nil {
			log.Println("Error trying to find document with ID: " + documentID)
			return err3
		}

		// Get the firebase ID of the document
		firebaseID := foundDocument.Ref.ID

		data, ok := updatedDocument.(T)
		if !ok {
			return fmt.Errorf("data does not match target struct")
		}

		// Add element in embedded structure.
		_, err2 := client.Collection(collection).Doc(firebaseID).Set(ctx, data)
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
func DeleteDocument(id string, collection string) error {
	// Checks if a document with the provided ID exists in the collection
	if ok, err := documentExists(ctx, collection, id); ok && err == nil {
		// Find document with matching ID
		foundDocument, err3 := getDocumentByID(id, collection)
		if err3 != nil {
			log.Println("Error trying to find document with ID: " + id)
			return err3
		}

		// Get the firebase ID of the document
		firebaseID := foundDocument.Ref.ID

		// Delete specified document
		_, err2 := client.Collection(collection).Doc(firebaseID).Delete(ctx)
		if err2 != nil {
			log.Println("Error while deleting document:" + id)
			return err2
		}
	} else if !ok && err == nil {
		log.Printf(
			"A document with the provided ID: %s, was not found in the collection: %s.\n",
			id, collection,
		)
	} else {
		log.Println("Error while trying to find document: ", err.Error())
		return err
	}

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
			return false, fmt.Errorf(
				"document with ID %s not found in collection %s",
				documentID,
				collection,
			)
		}
		return false, err
	}

	return true, nil
}

/*
numOfDocumentsInCollection Returns the number of documents in the collection.
*/
func NumOfDocumentsInCollection(collection string) (int, error) {
	result, err := client.Collection(collection).NewAggregationQuery().WithCount("all").Get(ctx)
	if err != nil {
		log.Println("firestore: error while trying to get count of documents in collection")
		return -1, fmt.Errorf("firestore: error while trying to get count of documents in collection")
	}

	count, ok := result["all"]
	if !ok {
		log.Println("firestore: couldn't get alias for COUNT from results")
		return -1, fmt.Errorf("firestore: couldn't get alias for COUNT from results")
	}

	countValue := count.(*firestorepb.Value)

	return int(countValue.GetIntegerValue()), nil
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
			return nil, fmt.Errorf("document not found in collection")
		}
		return nil, err
	}

	return docSnap, nil
}

func InitializeForTesting() {
	// Set the Firestore emulator host
	err := os.Setenv(
		"FIRESTORE_EMULATOR_HOST",
		"localhost:8080",
	) // Default port for Firestore emulator
	if err != nil {
		log.Fatalf("Failed to set Firestore emulator host: %v", err)
	} else {
		log.Println("Firestore emulator host set successfully")
	}

	// Initialize Firestore client
	ctx = context.Background()
	app, err := firebase.NewApp(
		ctx, &firebase.Config{
			ProjectID: "prog2005-assignment-2-c2e5c",
		},
		option.WithoutAuthentication(),
	)
	if err != nil {
		log.Fatalf("Failed to create Firestore app: %v", err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	// Assign the client to your package-level variable for Firestore client
	client = firestoreClient

	log.Println("Firestore client initialized for testing")
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
		log.Println("Failed to create Firebase app: ", err)
		return
	}

	// Instantiate client
	client, err = app.Firestore(ctx)

	// Alternative setup, directly through Firestore (without initial reference to Firebase); but requires Project ID; useful if multiple projects are used
	// client, err := firestore.NewClient(ctx, projectID)

	// Check whether there is an error when connecting to Firestore
	if err != nil {
		log.Println("Failed to create Firestore client: ", err)
		return
	}
	log.Println("Firestore client initialized normally")
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
