package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kubernetes-sigs/aws-iam-authenticator/pkg/token"
)

var gen token.Generator
var clusterID string

func handler(w http.ResponseWriter, r *http.Request) {
	var tok token.Token
	var err error
	tok, err = gen.Get(clusterID)
	if err != nil {
		fmt.Fprintf(w, "Failed to retrieve token: %v", err)
	}
	log.Printf("Got token %v", gen.FormatJSON(tok))
	fmt.Fprintf(w, "%v\n", gen.FormatJSON(tok))
}

func init() {
	var err error
	gen, err = token.NewGenerator(false)
	if err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}

	clusterID = os.Getenv("EKS_CLUSTER_ID")
	if clusterID == "" {
		log.Fatal("EKS_CLUSTER_ID must be set")
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Info("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
