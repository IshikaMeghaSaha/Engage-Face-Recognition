package main

import (
	"context"
	"fmt"
	"net/http"

	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v1.0/face"
	"github.com/Azure/go-autorest/autorest"
	uuid "github.com/gofrs/uuid"
)

func compare(w http.ResponseWriter, r *http.Request) {

	// A global context for use in all samples
	faceContext := context.Background()

	/*
		Authenticate
	*/
	subscriptionKey := "9e3b9f3dbd8345a5a1f2fc68412e7bfb"
	endpoint := "https://face-recognition-app.cognitiveservices.azure.com/"

	// Client used for Detect Faces, Find Similar, and Verify examples.
	client := face.NewClient(endpoint)
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(subscriptionKey)

	fmt.Println()
	fmt.Println("-----------------------------")
	fmt.Println("VERIFY")
	// <snippet_ver_images>
	// Create a slice list to hold the target photos of the same person

	targetImageFileName := "./img2.jpg"

	// The source photos contain this person, maybe
	sourceImageFileName1 := "./img1.jpg"
	urlSource1, err := os.Open(sourceImageFileName1)
	if err != nil {
		fmt.Println("Cannot open file")
	}
	returnFaceIDVerify := true
	returnFaceLandmarksVerify := false
	returnRecognitionModelVerify := false

	detectedVerifyFaces1, dErrV1 := client.DetectWithStream(faceContext, urlSource1, &returnFaceIDVerify, &returnFaceLandmarksVerify, nil, face.Recognition03, &returnRecognitionModelVerify, face.Detection02)
	if dErrV1 != nil {
		log.Fatal(dErrV1)
	}
	// Dereference the result, before getting the ID
	dVFaceIds1 := *detectedVerifyFaces1.Value
	// Get ID of the detected face
	imageSource1Id := dVFaceIds1[0].FaceID
	fmt.Println(fmt.Printf("%v face(s) detected from image: %v", len(dVFaceIds1), sourceImageFileName1))

	var detectedVerifyFacesIds [1]uuid.UUID
	//for i, imageFileName := range targetImageFileNames {
	urlSource, er := os.Open(targetImageFileName)
	if er != nil {
		fmt.Println("Cannot open file")
	}

	detectedVerifyFaces, dErrV := client.DetectWithStream(faceContext, urlSource, &returnFaceIDVerify, &returnFaceLandmarksVerify, nil, face.Recognition03, &returnRecognitionModelVerify, face.Detection02)
	if dErrV != nil {
		log.Fatal(dErrV)
	}
	// Dereference *[]DetectedFace from Value in order to loop through it.
	dVFaces := *detectedVerifyFaces.Value
	// Add the returned face's face ID
	detectedVerifyFacesIds[0] = *dVFaces[0].FaceID
	fmt.Println(fmt.Printf("%v face(s) detected from image: %v", len(dVFaces), targetImageFileName))

	//var verifyRequestBody2
	verifyRequestBody2 := face.VerifyFaceToFaceRequest{FaceID1: imageSource1Id, FaceID2: &detectedVerifyFacesIds[0]}
	var verifyResultDiff, vErrDiff = client.VerifyFaceToFace(faceContext, verifyRequestBody2)
	if vErrDiff != nil {
		log.Fatal(vErrDiff)
	}
	// Check if the faces are from the same person.
	if *verifyResultDiff.IsIdentical {
		fmt.Fprintf(w, "Authorized!")

	} else {
		// Low confidence means they are more differant than same.
		fmt.Fprintf(w, "Not Authorized!")

	}

}

func main() {
	http.HandleFunc("/Result", compare)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}
