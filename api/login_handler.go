package api

import (
	"net/http"
	"time"

	"github.com/InVisionApp/rye"
)

const (
	jwtCookieName            = "invision_jwt"
	jwtCookieExpiresDuration = 30 * time.Minute
	jwtSecret                = "b981f158e5b7313508c356af0b7e34f92a4f044d" // For signing the JWT
	jwtIssuer                = "InVision"
	jwtExpiresAtDuration     = 24 * time.Hour // Expiration 1 day
	jwtSessionName           = "session_id"
)

type LoginRequest struct {
	Email    string `json:"email" format:"email"` // Required: Email address
	Password string `json:"password"`             // Required: Plain text password
}

type LoginResponse struct {
	Status  string `json:"status"`        // Values: "ok" on success, "error" on failure.
	Message string `json:"message"`       // Description of the status
	JWT     string `json:"jwt,omitempty"` // Set to JWT if query ?jwtbody or ?jwtboth is set.
}

// @Title Login
// @Summary Login a user using credentials and get a JWT with a session ID.
// @Description The login endpoint is how you get a JWT with a session ID.
// @Tags login
// @Accept json
// @Produce json
// @Param Calling-Service header string true "The caller for the request, with valid requests set to value: example-bff"
// @Param input body api.LoginRequest true "User credentials"
// @Param jwtbody query bool false "Upon a successful login, it will return the JWT in the response body and NOT set a Set-Cookie header for the JWT."
// @Param jwtboth query bool false "Upon a successful login, it will return the JWT in the response body and SET a Set-Cookie header for the JWT."
// @Success 200 {object} api.LoginResponse "Login successful and session created"
// @Failure 400 {object} rye.JSONStatus "Incorrect payload, refer to error message"
// @Failure 401 {object} rye.JSONStatus "Unauthorized: invalid credentials"
// @Failure 403 {object} rye.JSONStatus "Caller not whitelisted"
// @Failure 500 {object} rye.JSONStatus "Unexpected server error, refer to the error message"
// @Router /api/v1/login [post]
func (a *API) loginHandler(rw http.ResponseWriter, r *http.Request) *rye.Response {
	// Hint: Returning a *rye.Response will generate a rye.JSONStatus for the response.
	//       Just need to set the StatusCode and Err fields.

	// Hint: Use a.apiDAL.Authenticate() to verify credentials.

	return nil
}
