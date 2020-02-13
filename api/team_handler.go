package api

import (
	"errors"
	"net/http"

	"github.com/InVisionApp/rye"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// @Title Get Team
// @Summary Gets information about a team.
// @Description Uses the user's session to get a Team object, if the user has access to the team. The session ID is in a JWT within an HTTP cookie.
// @Tags login
// @Accept json
// @Produce json
// @Param Calling-Service header string true "The caller for the request, with valid requests set to value: example-bff"
// @Param TeamID path string true "ID of team to get info."
// @Success 200 {object} team.Team "The team information"
// @Failure 400 {object} rye.JSONStatus "Invalid request"
// @Failure 401 {object} rye.JSONStatus "Invalid session"
// @Failure 403 {object} rye.JSONStatus "No access"
// @Failure 500 {object} rye.JSONStatus "Unexpected server error, refer to the error message"
// @Router /api/v1/teams/{TeamID} [get]
func (a *API) getTeamHandler(rw http.ResponseWriter, r *http.Request) *rye.Response {
	teamID := mux.Vars(r)["TeamID"]
	if len(teamID) != 25 {
		return &rye.Response{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("Invalid request"),
		}
	}

	// Find session ID from HTTP cookie.
	sessionID, err := getSessionID(r.Cookies())
	if err != nil {
		return &rye.Response{
			StatusCode: http.StatusUnauthorized,
			Err:        errors.New("Invalid session"),
		}
	}

	// Get teams in session.
	teams := a.apiDAL.SessionTeams(sessionID)
	if len(teams) == 0 {
		return &rye.Response{
			StatusCode: http.StatusForbidden,
			Err:        errors.New("No access"),
		}
	}

	// Get info about the team.
	teamInfo := a.apiDAL.GetTeamByID(teamID)
	if teamInfo == nil {
		return &rye.Response{
			StatusCode: http.StatusForbidden,
			Err:        errors.New("No access"),
		}
	}

	return respondAsJSON(rw, http.StatusOK, teamInfo)
}

func getSessionID(cookies []*http.Cookie) (string, error) {
	for i := range cookies {
		if cookies[i].Name == jwtCookieName {
			token, err := jwtgo.ParseWithClaims(cookies[i].Value, &jwtgo.MapClaims{},
				func(token *jwtgo.Token) (interface{}, error) {
					return []byte(jwtSecret), nil
				})
			if err != nil {
				return "", err
			}

			claims, ok := token.Claims.(*jwtgo.MapClaims)
			if !ok {
				return "", errors.New("Invalid JWT claim")
			}
			if err := claims.Valid(); err != nil {
				return "", err
			}

			sessionID, ok := (*claims)[jwtSessionName].(string)
			if !ok {
				return "", errors.New("No session ID found in claim")
			}

			return sessionID, nil
		}
	}
	return "", errors.New("No cookies")
}
