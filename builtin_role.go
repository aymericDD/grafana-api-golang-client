package gapi

import (
	"bytes"
	"encoding/json"
)

type BuiltRole struct {
	BuiltinRole string `json:"builtInRole"`
	RoleUID     string `json:"roleUid"`
}


// GetBuiltInRoles gets all built-in role grants. Available only in Grafana Enterprise.
func (c *Client) GetBuiltInRoles() (map[string][]*Role, error) {
	builtInRoles := make(map[string][]*Role, 0)
	err := c.request("GET", "/api/access-control/builtin-roles", nil, nil, &builtInRoles)
	if err != nil {
		return nil, err
	}
	return builtInRoles, nil
}

// NewBuiltInRole creates a new grant for a built-in role. Available only in Grafana Enterprise.
func (c *Client) NewBuiltInRole(builtInRole BuiltRole) (*BuiltRole, error) {
	data, err := json.Marshal(builtInRole)
	if err != nil {
		return nil, err
	}

	created := &BuiltRole{}

	err = c.request("POST", "/api/access-control/builtin-roles", nil, bytes.NewBuffer(data), &created)
	if err != nil {
		return nil, err
	}

	return created, err
}

// DeleteBuiltInRole delete the grant from built-in role. Available only in Grafana Enterprise.
func (c *Client) DeleteBuiltInRole(builtInRole BuiltRole) error {
	data, err := json.Marshal(builtInRole)
	if err != nil {
		return err
	}

	err = c.request("DELETE", "/api/access-control/builtin-roles", nil, bytes.NewBuffer(data), nil)

	return err
}
