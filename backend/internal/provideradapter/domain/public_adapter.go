package domain

// PublicAdapter extends the base Adapter for providers that don't require authentication
type PublicAdapter interface {
	Adapter
	// No additional methods needed as these providers don't require authentication
}
