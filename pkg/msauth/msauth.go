package msauth

import (
	"context"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/oauth2/microsoft"
	"net/http"
)

//GetOAuth2Client Returns a OAuth2 http.Client instance that will inject Bearer token
func GetOAuth2Client(tenant string, clientID string, clientSecret string) *http.Client {

	microsoftEndpoints := microsoft.AzureADEndpoint(tenant)
	conf := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://graph.microsoft.com/.default"},
		TokenURL:     microsoftEndpoints.TokenURL,
	}

	return conf.Client(context.Background())

}
