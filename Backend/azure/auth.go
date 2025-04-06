package azure

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

var token azcore.AccessToken

func GetACSToken() string {

	if token.Token != "" && token.ExpiresOn.After(time.Now()) {
		return token.Token
	}

	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	if tenantID == "" || clientID == "" || clientSecret == "" {
		return ""
	}

	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		log.Default().Println("Failed to create credential: " + err.Error())
		return ""
	}
	newToken, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{"https://communication.azure.com/.default"},
	})
	if err != nil {
		log.Default().Println("Failed to get token: " + err.Error())
		return ""
	}

	token = newToken

	return token.Token
}
