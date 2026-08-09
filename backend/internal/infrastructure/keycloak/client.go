package keycloak

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"my-web-app.com/smart-logistic-hub/internal/infrastructure/config"
)

var (
	jwksCache map[string]interface{}
	jwksMutex sync.RWMutex
)

func FetchJWKS(cfg *config.Config) (map[string]interface{}, error) {
	jwksMutex.RLock()
	if jwksCache != nil {
		defer jwksMutex.RUnlock()
		return jwksCache, nil
	}
	jwksMutex.RUnlock()

	jwksMutex.Lock()
	defer jwksMutex.Unlock()

	if jwksCache != nil {
		return jwksCache, nil
	}

	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", cfg.KeycloakServerURL, cfg.KeycloakRealm)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch jwks: status %d", resp.StatusCode)
	}

	var jwks map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	jwksCache = jwks
	return jwksCache, nil
}
