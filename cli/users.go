package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// GetUsersPaginated - Returns all registered users in HoneyBadger using pagination
func (hbc *HoneyBadgerClient) GetUsersPaginated(pagePath string, hbUserList []HoneyBadgerUser) ([]HoneyBadgerUser, error) {
	var hbUsers HoneyBadgerUsers
	urlPath := fmt.Sprintf("%s/%s", hbc.HostURL, pagePath)

	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return hbUsers.Users, err
	}

	body, err := hbc.doRequest(req)
	if err != nil {
		return hbUsers.Users, err
	}

	err = json.Unmarshal(body, &hbUsers)
	if err != nil {
		return hbUsers.Users, err
	}

	if hbUsers.Links.NextPage != "" {
		return hbc.GetUsersPaginated(hbUsers.Links.NextPage, hbUsers.Users)
	}

	hbUserList = append(hbUserList, hbUsers.Users...)

	return hbUserList, nil
}

// GetUsers - Returns all registered users in HoneyBadger
func (hbc *HoneyBadgerClient) GetUsers() ([]HoneyBadgerUser, error) {
	var hbUsers HoneyBadgerUsers

	return hbc.GetUsersPaginated("", hbUsers.Users)
}

// FindUserByID - Returns a user by ID
func (hbc *HoneyBadgerClient) FindUserByID(userID int) (HoneyBadgerUser, error) {
	hbUsers, err := hbc.GetUsers()
	if err != nil {
		return HoneyBadgerUser{}, err
	}

	for _, user := range hbUsers {
		if user.Id == userID {
			return user, nil
		}
	}
	return HoneyBadgerUser{}, errors.New("User not found")
}

// CreateUser - Crea a HoneyBadger User
func (hbc *HoneyBadgerClient) CreateUser(userEmail string) (int, error) {
	var hbUser HoneyBadgerUser
	var jsonPayload = []byte(`{"team_invitation":"` + userEmail + `"}`)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/teams/ID/team_invitations", "http://localhost:8080"), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return 0, err
	}

	body, err := hbc.doRequest(req)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(body, &hbUser)
	if err != nil {
		return 0, err
	}

	return hbUser.Id, nil
}

// UpdateUser - Update HoneyBadger User Information
func (hbc *HoneyBadgerClient) UpdateUser(userID int, isAdmin bool) error {
	var jsonPayload = []byte(`{"team_member":{"admin":` + strconv.FormatBool(isAdmin) + `}}`)

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v2/teams/ID/team_members/ID", "http://localhost:8080"), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	_, err = hbc.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser - Delete HoneyBadger User
func (hbc *HoneyBadgerClient) DeleteUser(userID int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v2/teams/ID/team_members/ID", "http://localhost:8080"), nil)
	if err != nil {
		return err
	}

	_, err = hbc.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
