package secretsengine

type playerDataClient struct {
}

// newClient creates a new client to access HashiCups
// and exposes it for any secrets or roles to use.
func newClient(config *playerDataConfig) (*playerDataClient, error) {
	return &playerDataClient{}, nil
}
