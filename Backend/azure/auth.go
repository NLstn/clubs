package azure

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

var token azcore.AccessToken

func GetACSToken() string {

	if token.Token != "" && token.ExpiresOn.After(time.Now()) {
		return token.Token
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
