package types

import "github.com/golang-jwt/jwt/v5"

type TypeDataToken struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
}

type TypeAuthResponseJson struct {
	TypeResponse
	Data *TypeDataToken `json:"data"`
}

type TypeAuthErrorResponseJson struct {
	TypeResponse
	Data map[string][]string `json:"data"`
}

type TypeAuthExpiredAccessToken struct {
	TypeResponse
	RedirectLogin bool `json:"redirectLogin"`
}

type ClaimsRefreshToken struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

type ClaimsAccessToken struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}
