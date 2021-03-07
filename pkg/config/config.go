package config

// Configuration represents Global Config
type Configuration struct {
	Kubeconfig *string
	Namespace  *string
	TTL        *int32
	Shell      *string
	Memory     *string
	CPU        *string
	Pod        *string
}
