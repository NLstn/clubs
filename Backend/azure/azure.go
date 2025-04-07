package azure

import (
	"errors"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

var tenantID string
var clientID string
var clientSecret string

var cred *azidentity.ClientSecretCredential

func Init() error {
	tenantID = os.Getenv("AZURE_TENANT_ID")
	clientID = os.Getenv("AZURE_CLIENT_ID")
	clientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	if tenantID == "" || clientID == "" || clientSecret == "" {
		return errors.New("AZURE_TENANT_ID, AZURE_CLIENT_ID, and AZURE_CLIENT_SECRET must be set")
	}

	var err error
	cred, err = azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		return errors.New("Failed to create credential: " + err.Error())
	}

	log.Default().Println("Azure SDK initialized")
	return nil
}
