package gapi

import (
	"testing"

	"github.com/gobs/pretty"
)

const (
	searchTeamJSON = `
{
  "totalCount": 1,
  "teams": [
    {
      "id": 1,
      "orgId": 1,
      "name": "MyTestTeam",
      "email": "",
      "avatarUrl": "/avatar/3f49c15916554246daa714b9bd0ee398",
      "memberCount": 1
    }
  ],
  "page": 1,
  "perPage": 1000
}
`
	getTeamJSON = `
{
  "id": 1,
  "orgId": 1,
  "name": "MyTestTeam",
  "email": "",
  "avatarUrl": "avatar/abcdef",
  "memberCount": 1,
  "permission": 0
}
`
	addTeamsJSON = `
{
  "message":"Team created",
  "teamId":2
}
`
	updateTeamJSON     = `{"message":"Team updated"}`
	deleteTeamJSON     = `{"message":"Team deleted"}`
	getTeamMembersJSON = `
[
  {
    "orgId": 1,
    "teamId": 1,
    "userID": 3,
    "auth_module": "oauth_github",
    "email": "user1@email.com",
    "login": "user1",
    "avatarUrl": "/avatar/1b3c32f6386b0185c40d359cdc733a79",
    "labels": [],
    "permission": 0
  },
  {
    "orgId": 1,
    "teamId": 1,
    "userID": 2,
    "auth_module": "oauth_github",
    "email": "user2@email.com",
    "login": "user2",
    "avatarUrl": "/avatar/cad3c68da76e45d10269e8ef02f8e73e",
    "labels": [],
    "permission": 0
  }
]
`
	addTeamMemberJSON = `
{
  "userID": 2
}
`
	removeMemberFromTeamJSON = `{"message":"Team Member removed"}`
	getTeamPreferencesJSON   = `
{
  "theme": "",
  "homeDashboardID": 0,
  "timezone": ""
}
`
	updateTeamPreferencesJSON = `
{
  "message":"Preferences updated"
}
`

	createTeamRoleJSON = `
{
  "message":"Role added to the team."
}
`
	deleteTeamRoleJSON = `
{
  "message":"Role removed from the team."
}
`
	getTeamRolesJSON = `
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

func TestSearchTeam(t *testing.T) {
	server, client := gapiTestTools(t, 200, searchTeamJSON)
	defer server.Close()

	query := "myteam"
	resp, err := client.SearchTeam(query)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pretty.PrettyFormat(resp))

	expect := &SearchTeam{
		TotalCount: 1,
		Teams: []*Team{
			{
				ID:          1,
				OrgID:       1,
				Name:        "MyTestTeam",
				Email:       "",
				AvatarURL:   "avatar/3f49c15916554246daa714b9bd0ee398",
				MemberCount: 1,
				Permission:  0,
			},
		},
		Page:    1,
		PerPage: 1000,
	}
	t.Run("check data", func(t *testing.T) {
		if expect.TotalCount != resp.TotalCount || expect.Teams[0].Name != resp.Teams[0].Name {
			t.Error("Not correctly parsing returned team search.")
		}
	})
}

func TestTeam(t *testing.T) {
	server, client := gapiTestTools(t, 200, getTeamJSON)
	defer server.Close()

	id := int64(1)
	resp, err := client.Team(id)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pretty.PrettyFormat(resp))

	expect := &Team{
		ID:          1,
		OrgID:       1,
		Name:        "MyTestTeam",
		Email:       "",
		AvatarURL:   "avatar/abcdef",
		MemberCount: 1,
		Permission:  0,
	}
	t.Run("check data", func(t *testing.T) {
		if resp.ID != expect.ID || resp.Name != expect.Name {
			t.Error("Not correctly parsing returned team.")
		}
	})
}

func TestAddTeam(t *testing.T) {
	server, client := gapiTestTools(t, 200, addTeamsJSON)
	defer server.Close()

	name := "TestTeam"
	email := ""

	id, err := client.AddTeam(name, email)
	if err != nil {
		t.Error(err)
	}
	if id == 0 {
		t.Error("AddTeam returned an invalid ID")
	}
}

func TestUpdateTeam(t *testing.T) {
	server, client := gapiTestTools(t, 200, updateTeamJSON)
	defer server.Close()

	id := int64(1)
	name := "TestTeam"
	email := ""

	err := client.UpdateTeam(id, name, email)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteTeam(t *testing.T) {
	server, client := gapiTestTools(t, 200, deleteTeamJSON)
	defer server.Close()

	id := int64(1)

	err := client.DeleteTeam(id)
	if err != nil {
		t.Error(err)
	}
}

func TestTeamMembers(t *testing.T) {
	server, client := gapiTestTools(t, 200, getTeamMembersJSON)
	defer server.Close()

	id := int64(1)

	resp, err := client.TeamMembers(id)
	if err != nil {
		t.Fatal(err)
	}
	expects := []*TeamMember{
		{
			OrgID:      1,
			TeamID:     1,
			UserID:     3,
			Email:      "user1@email.com",
			Login:      "user1",
			AvatarURL:  "/avatar/1b3c32f6386b0185c40d359cdc733a79",
			Permission: 0,
		},
		{
			OrgID:      1,
			TeamID:     1,
			UserID:     2,
			Email:      "user2@email.com",
			Login:      "user2",
			AvatarURL:  "/avatar/cad3c68da76e45d10269e8ef02f8e73e",
			Permission: 0,
		},
	}

	for i, expect := range expects {
		t.Run("check data", func(t *testing.T) {
			if expect.Email != resp[i].Email || expect.AvatarURL != resp[i].AvatarURL {
				t.Error("Not correctly parsing returned team members.")
			}
		})
	}
}

func TestAddTeamMember(t *testing.T) {
	server, client := gapiTestTools(t, 200, addTeamMemberJSON)
	defer server.Close()

	id := int64(1)
	userID := int64(2)

	if err := client.AddTeamMember(id, userID); err != nil {
		t.Error(err)
	}
}

func TestRemoveMemberFromTeam(t *testing.T) {
	server, client := gapiTestTools(t, 200, removeMemberFromTeamJSON)
	defer server.Close()

	id := int64(1)
	userID := int64(2)

	if err := client.RemoveMemberFromTeam(id, userID); err != nil {
		t.Error(err)
	}
}

func TestTeamPreferences(t *testing.T) {
	server, client := gapiTestTools(t, 200, getTeamPreferencesJSON)
	defer server.Close()

	id := int64(1)

	resp, err := client.TeamPreferences(id)
	if err != nil {
		t.Fatal(err)
	}
	expect := &Preferences{
		Theme:           "",
		HomeDashboardID: 0,
		Timezone:        "",
	}

	t.Run("check data", func(t *testing.T) {
		if expect.Theme != resp.Theme || expect.HomeDashboardID != resp.HomeDashboardID {
			t.Error("Not correctly parsing returned team preferences.")
		}
	})
}

func TestUpdateTeamPreferences(t *testing.T) {
	server, client := gapiTestTools(t, 200, updateTeamPreferencesJSON)
	defer server.Close()

	id := int64(1)
	preferences := Preferences{
		Theme:           "",
		HomeDashboardID: int64(0),
		Timezone:        "",
	}

	if err := client.UpdateTeamPreferences(id, preferences); err != nil {
		t.Error(err)
	}
}

func TestNewTeamRole(t *testing.T) {
	server, client := gapiTestTools(t, 201, createTeamRoleJSON)
	t.Cleanup(func() {
		server.Close()
	})

	id := int64(1)

	if err := client.NewTeamRole(id, "vc3SCSsGz"); err != nil {
		t.Error(err)
	}
}

func TestGetTeamRoles(t *testing.T) {
	server, client := gapiTestTools(t, 200, getTeamRolesJSON)
	t.Cleanup(func() {
		server.Close()
	})

	id := int64(1)

	resp, err := client.GetTeamRoles(id)

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
				t.Error("Not correctly parsing returned team roles.")
			}
		})
	}
}

func TestDeleteTeamRole(t *testing.T) {
	server, client := gapiTestTools(t, 200, deleteTeamRoleJSON)
	t.Cleanup(func() {
		server.Close()
	})

	id := int64(1)

	if err := client.DeleteTeamRole(id, "vc3SCSsGz"); err != nil {
		t.Error(err)
	}
}
