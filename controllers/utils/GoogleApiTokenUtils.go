package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/themane/MMOServer/constants"
	"github.com/themane/MMOServer/models"
	"github.com/themane/MMOServer/mongoRepository/exceptions"
	"google.golang.org/api/idtoken"
	"io"
	"net/http"
)

func ValidateIdToken(c *gin.Context) (*models.UserSocialDetails, error) {
	idToken := extractToken(c)
	return ParseIdToken(idToken)
}

func ParseIdToken(idToken string) (*models.UserSocialDetails, error) {
	userDetails, err := ParseGoogleIdToken(idToken)
	if err != nil {
		userDetails, err = ParseFacebookIdToken(idToken)
		if err != nil {
			return nil, err
		}
		return userDetails, nil
	}
	return userDetails, nil
}

func ParseGoogleIdToken(idToken string) (*models.UserSocialDetails, error) {
	payload, err := idtoken.Validate(context.Background(), idToken, "")
	if err != nil {
		return nil, &exceptions.NoSuchCombinationError{Message: err.Error()}
	}
	userDetails := models.UserSocialDetails{
		Id:            fmt.Sprintf("%v", payload.Claims["sub"]),
		Name:          fmt.Sprintf("%v", payload.Claims["name"]),
		Email:         fmt.Sprintf("%v", payload.Claims["email"]),
		PictureUrl:    fmt.Sprintf("%v", payload.Claims["picture"]),
		Authenticator: constants.GoogleAuthenticator,
	}
	return &userDetails, nil
}

func ParseFacebookIdToken(idToken string) (*models.UserSocialDetails, error) {
	fbUserDetailsUrl := "https://graph.facebook.com/me?fields=id,name,email,picture&access_token=" + idToken
	response, err := http.Get(fbUserDetailsUrl)
	if err != nil {
		return nil, err
	}
	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		return nil, &exceptions.NoSuchCombinationError{Message: "invalid login"}
	}
	var fbUserDetails models.FbUserDetails
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&fbUserDetails)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	userDetails := models.UserSocialDetails{
		Id:            fbUserDetails.Id,
		Name:          fbUserDetails.Name,
		Email:         fbUserDetails.Email,
		PictureUrl:    fbUserDetails.Picture.Data.Url,
		Authenticator: constants.FacebookAuthenticator,
	}
	return &userDetails, nil
}
