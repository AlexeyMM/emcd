package vault

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	vaultapi "github.com/hashicorp/vault/api"
)

type Provider struct {
	client *vaultapi.Client
}

func NewProvider() (*Provider, error) {
	config := vaultapi.DefaultConfig()
	if err := config.ReadEnvironment(); err != nil {
		return nil, fmt.Errorf("read config vault: %w", err)
	}
	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Provider{client: client}, err
}

func (v Provider) Value(u *url.URL) (string, error) {
	// maybe need optimization: "use local cache for re-access"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	elementsPath := strings.Split(u.Path, "/")
	if len(elementsPath) == 3 {
		r, err := v.client.KVv2(u.Host).Get(ctx, elementsPath[1])
		if err != nil {
			return "", fmt.Errorf("recive data from vault: %w", err)
		}
		s, ok := r.Data[elementsPath[2]].(string)
		if !ok {
			return "", fmt.Errorf("elemetn %s not string value", elementsPath[2])
		}
		return s, nil
	}
	return "", fmt.Errorf("incorrect path secret")
}

func (v Provider) Scheme() string {
	return "secret-vault"
}
