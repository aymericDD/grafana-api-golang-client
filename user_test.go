package gapi

import (
	"testing"

	"github.com/gobs/pretty"
)

const (
	getUsersJSON       = `[{"id":1,"email":"users@localhost","isAdmin":true}]`
	getUserJSON        = `{"id":2,"email":"user@localhost","isGrafanaAdmin":false}`
	getUserByEmailJSON = `{"id":3,"email":"userByEmail@localhost","isGrafanaAdmin":true}`
	getUserUpdateJSON  = `{"id":4,"email":"userUpdate@localhost","isGrafanaAdmin":false}`
	createUserRoleJSON = `
{
  "message":"Role added to the user."
}
`
	deleteUserRoleJSON = `
{
  "message":"Role removed from the user."
}
`
	getUserRolesJSON = `
[
{
    "orgId": 1,
    "uid": "vc3SCSsGz",
    "name": "test:policy",
	"version": 1,
    "description": "Test policy description",
    "permissions": [
        {
            "id": 6,
            "permission": "test:self",
            "scope": "test:self",
            "updated": "2021-02-22T16:16:05.646913+01:00",
            "created": "2021-02-22T16:16:05.646912+01:00"
        }
    ],
    "updated": "2021-02-22T16:16:05.644216+01:00",
    "created": "2021-02-22T16:16:05.644216+01:00"
}
]
`
)

func TestUsers(t *testing.T) {
	server, client := gapiTestTools(t, 200, getUsersJSON)
	defer server.Close()

	resp, err := client.Users()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pretty.PrettyFormat(resp))

	if len(resp) != 1 {
		t.Fatal("No users were returned.")
	}

	user := resp[0]

	if user.Email != "users@localhost" ||
		user.ID != 1 ||
		user.IsAdmin != true {
		t.Error("Not correctly parsing returned users.")
	}
}

func TestUser(t *testing.T) {
	server, client := gapiTestTools(t, 200, getUserJSON)
	defer server.Close()

	user, err := client.User(1)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pretty.PrettyFormat(user))

	if user.Email != "user@localhost" ||
		user.ID != 2 ||
		user.IsAdmin != false {
		t.Error("Not correctly parsing returned user.")
	}
}

func TestUserByEmail(t *testing.T) {
	server, client := gapiTestTools(t, 200, getUserByEmailJSON)
	defer server.Close()

	user, err := client.UserByEmail("admin@localhost")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pretty.PrettyFormat(user))

	if user.Email != "userByEmail@localhost" ||
		user.ID != 3 ||
		user.IsAdmin != true {
		t.Error("Not correctly parsing returned user.")
	}
}

func TestUserUpdate(t *testing.T) {
	server, client := gapiTestTools(t, 200, getUserUpdateJSON)
	defer server.Close()

	user, err := client.User(4)
	if err != nil {
		t.Fatal(err)
	}
	user.IsAdmin = true
	err = client.UserUpdate(user)
	if err != nil {
		t.Error(err)
	}
}

func TestNewUserRole(t *testing.T) {
	server, client := gapiTestTools(t, 201, createUserRoleJSON)
	t.Cleanup(func() {
		server.Close()
	})

	id := int64(1)

	if err := client.NewUserRole(id, "vc3SCSsGz"); err != nil {
		t.Error(err)
	}
}

func TestGetUserRoles(t *testing.T) {
	server, client := gapiTestTools(t, 200, getUserRolesJSON)
	t.Cleanup(func() {
		server.Close()
	})

	id := int64(1)

	resp, err := client.GetUserRoles(id)

	if err != nil {
		t.Error(err)
	}

	expected := []*Role{
		{
			OrgID:       1,
			Version:     1,
			UID:         "vc3SCSsGz",
			Name:        "test:policy",
			Description: "Test policy description",
			Permissions: []Permission{
				{
					Action: "test:self",
					Scope:  "test:self",
				},
			},
		}}

	for i, expect := range expected {
		t.Run("check response data", func(t *testing.T) {
			if expect.UID != resp[i].UID || expect.Name != resp[i].Name {
				t.Error("Not correctly parsing returned user roles.")
			}
		})
	}
}

func TestDeleteUserRole(t *testing.T) {
	server, client := gapiTestTools(t, 200, deleteUserRoleJSON)
	t.Cleanup(func() {
		server.Close()
	})

	id := int64(1)

	if err := client.DeleteUserRole(id, "vc3SCSsGz"); err != nil {
		t.Error(err)
	}
}
