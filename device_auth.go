package githubcopilotapi

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func DeviceLogin(ctx context.Context, clientID string) (string, error) {
	config := &oauth2.Config{
		ClientID: clientID,
		Endpoint: github.Endpoint,
		Scopes:   []string{"read:user"},
	}

	resp, err := config.DeviceAuth(ctx)
	if err != nil {
		return "", err
	}

	if resp.VerificationURIComplete != "" {
		fmt.Printf("Please visit %s to authenticate.\n", resp.VerificationURIComplete)
	} else {
		fmt.Printf("Please take this code %q to authenticate at %s.\n", resp.UserCode, resp.VerificationURI)
	}
	fmt.Println("Press 'y' to continue, or any other key to abort.")
	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	if input != "y" && input != "Y" {
		return "", fmt.Errorf("aborted")
	}
	fmt.Println("Authenticating, please wait...")

	token, err := config.DeviceAccessToken(ctx, resp)
	if err != nil {
		return "", err
	}
	if token.AccessToken != "" {
		return token.AccessToken, nil
	}

	return "", fmt.Errorf("failed to get token")
}
