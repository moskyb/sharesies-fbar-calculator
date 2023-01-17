package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	MFAToken string `json:"mfa_token"`
	Remember bool   `json:"remember"`
}

type LoginResponse struct {
	RakaiaToken string `json:"rakaia_token"`
	User        struct {
		PortfolioID string `json:"portfolio_id"`
	} `json:"user"`
}

type LoginError struct {
	InvalidMFAToken bool   `json:"invalid_mfa_token"`
	Type            string `json:"type"`
}

func Login(i LoginInput) (LoginResponse, error) {
	j, err := json.Marshal(i)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to marshal login input: %w", err)
	}

	req, err := http.NewRequest("POST", "https://app.sharesies.nz/api/identity/login", bytes.NewBuffer(j))
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "https://github.com/moskyb/shareies-fbar-caluclator (please send me an email at ben@mosk.nz if i'm causing trouble!)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to post login: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to read login response: %w", err)
	}

	if resp.StatusCode != 200 {
		return LoginResponse{}, fmt.Errorf("login failed with code: %s: %s", resp.Status, string(body))
	}

	var loginResponse LoginResponse
	err = json.Unmarshal(body, &loginResponse)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("failed to unmarshal login response: %w", err)
	}

	if loginResponse == (LoginResponse{}) {
		// We got a 200, but the response wasn't in the format we expected. This is probably because we got a 200, but the
		// login failed. Let's check if we got a login error.
		var loginError LoginError
		err = json.Unmarshal(body, &loginError)
		if err != nil {
			return LoginResponse{}, fmt.Errorf("failed to unmarshal login error: %w", err)
		}

		if loginError == (LoginError{}) {
			// By golly, i don't know what happened. Let's just return the body as an error.
			return LoginResponse{}, fmt.Errorf("the login returned a successful status code, but the response wasn't something that I could parse: %s ", string(body))
		}

		switch loginError.Type {
		case "identity_mfa_required":
			if loginError.InvalidMFAToken {
				return LoginResponse{}, fmt.Errorf("invalid MFA token")
			}
			return LoginResponse{}, fmt.Errorf("MFA required")
		case "identity_anonymous":
			return LoginResponse{}, fmt.Errorf("invalid email or password")
		default:
			return LoginResponse{}, fmt.Errorf("unknown login error: %s", loginError.Type)
		}
	}

	return loginResponse, nil
}
