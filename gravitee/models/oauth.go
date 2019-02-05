package models

import "database/sql"

type OauthClient struct {
	BaseModel
	ClientKey    string         `sql:"type:varchar(254);unique;not null"`
	ClientSecret string         `sql:"type:varchar(60);not null"`
	RedirectURI  sql.NullString `sql:"type:varchar(200)"`
}

func (c *OauthClient) TableName() string {
	return "oauth_clients"
}

type OauthScope struct {
}

func (s *OauthScope) TableName() string {
	return "oauth_scopes"
}

type OauthRole struct {
}

func (r *OauthRole) TableName() string {
	return "oauth_roles"
}

type OauthUser struct {
	BaseModel
	RoleID sql.NullString
}

func (u *OauthUser) TableName() string {
	return "oauth_users"
}

type OauthRefreshToken struct {
}

func (rt *OauthRefreshToken) TableName() string {
	return "oauth_refresh_tokens"
}

type OauthAccessToken struct {
}

func (at *OauthAccessToken) TableName() string {
	return "oauth_access_tokens"
}

type OauthAuthorizationCode struct {
}

func (ac *OauthAuthorizationCode) TableName() string {
	return "oauth_authorization_codes"
}
