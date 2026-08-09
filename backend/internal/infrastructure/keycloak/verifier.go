package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"my-web-app.com/smart-logistic-hub/internal/infrastructure/config"
)

type JWTVerifier struct {
	cfg       *config.Config
	jwksCache map[string]interface{}
	mu        sync.RWMutex
	client    *http.Client
}

func NewJWTVerifier(cfg *config.Config) *JWTVerifier {
	return &JWTVerifier{
		cfg:    cfg,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (v *JWTVerifier) fetchJWKS() (map[string]interface{}, error) {
	v.mu.RLock()
	if v.jwksCache != nil {
		defer v.mu.RUnlock()
		return v.jwksCache, nil
	}
	v.mu.RUnlock()

	v.mu.Lock()
	defer v.mu.Unlock()

	if v.jwksCache != nil {
		return v.jwksCache, nil
	}

	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", v.cfg.KeycloakServerURL, v.cfg.KeycloakRealm)
	resp, err := v.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch jwks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch jwks: status %d", resp.StatusCode)
	}

	var jwks map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("decode jwks: %w", err)
	}

	v.jwksCache = jwks
	return v.jwksCache, nil
}

func (v *JWTVerifier) keyFunc(token *jwt.Token) (interface{}, error) {
	jwks, err := v.fetchJWKS()
	if err != nil {
		return nil, err
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("kid not found in token header")
	}

	keys, ok := jwks["keys"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("keys not found in jwks")
	}

	for _, k := range keys {
		keyMap, ok := k.(map[string]interface{})
		if !ok {
			continue
		}
		if keyMap["kid"] == kid {
			certStr, ok := keyMap["x5c"].([]interface{})
			if !ok || len(certStr) == 0 {
				return nil, fmt.Errorf("x5c not found for kid %s", kid)
			}
			certPEM := "-----BEGIN CERTIFICATE-----\n" + certStr[0].(string) + "\n-----END CERTIFICATE-----"
			return jwt.ParseRSAPublicKeyFromPEM([]byte(certPEM))
		}
	}

	return nil, fmt.Errorf("key not found for kid %s", kid)
}

func (v *JWTVerifier) VerifyToken(ctx context.Context, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, v.keyFunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(fmt.Sprintf("%s/realms/%s", v.cfg.KeycloakServerURL, v.cfg.KeycloakRealm)),
	)
	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
