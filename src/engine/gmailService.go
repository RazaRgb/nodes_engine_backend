package engine

import (
	"backend/src/db"
	"backend/src/models"
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Endpoint:     google.Endpoint,
}

func getGoogleAuthToken(ownerID uuid.UUID) (models.AccessToken, error) {
	dbToken, err := db.GetServiceAccessToken(ownerID, "google", "")
	if err != nil {
		return models.AccessToken{}, err // AUTH_REQUIRED error
	}

	token := &oauth2.Token{
		AccessToken:  string(dbToken.Token),
		RefreshToken: string(dbToken.RefreshToken),
		Expiry:       dbToken.Expiry,
	}

	tokenSource := GoogleOauthConfig.TokenSource(context.Background(), token)

	validToken, err := tokenSource.Token()
	if err != nil {
		return models.AccessToken{}, fmt.Errorf("failed to validate or refresh google token: %v", err)
	}

	finalToken := models.AccessToken{
		ID:                uuid.NewString(),
		ProviderAccountID: "",
		Owner:             ownerID,
		Token:             []byte(validToken.AccessToken),
		RefreshToken:      []byte(validToken.RefreshToken),
		Provider:          "google",
		Expiry:            validToken.Expiry,
	}

	if validToken.AccessToken != string(dbToken.Token) {
		db.UpdateServiceAccessToken(finalToken)
	}

	return finalToken, nil
}
