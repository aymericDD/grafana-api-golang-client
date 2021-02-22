package gapi

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Role struct {
	OrgID       int64        `json:"orgId"`
	Version     int64        `json:"version"`
	UID         string       `json:"uid,omitempty"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	Action string `json:"permission"`
	Scope  string `json:"scope"`
}

// GetRole gets a role with permissions. Available only in Grafana Enterprise.
func (c *Client) GetRole(uid string) (*Role, error) {
	r := &Role{}
	err := c.request("GET", buildUrl(uid), nil, nil, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// NewRole creates a new role with permissions. Available only in Grafana Enterprise.
func (c *Client) NewRole(role Role) (*Role, error) {
	data, err := json.Marshal(role)
	if err != nil {
		return nil, err
	}

	created := &Role{}

	err = c.request("POST", "/api/access-control/roles", nil, bytes.NewBuffer(data), &created)
	if err != nil {
		return nil, err
	}

	return created, err
}

// UpdateRole updates the role and permissions. Available only in Grafana Enterprise.
func (c *Client) UpdateRole(role Role) error {
	data, err := json.Marshal(role)
	if err != nil {
		return err
	}

	err = c.request("PUT", buildUrl(role.UID), nil, bytes.NewBuffer(data), nil)

	return err
}

// DeleteRole deletes the role with it's permissions. Available only in Grafana Enterprise.
func (c *Client) DeleteRole(uid string) error {
	return c.request("DELETE", buildUrl(uid), nil, nil, nil)
}

func buildUrl(uid string) string {
	const rootUrl = "/api/access-control/roles"
	return fmt.Sprintf("%s/%s", rootUrl, uid)
}
