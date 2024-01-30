package token

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
)

// EndpointOpts specifies a "token" Ironic Inspector Endpoint.
type EndpointOpts struct {
	IronicInspectorEndpoint string
	IronicInspectorToken    string
}

func initClientOpts(client *gophercloud.ProviderClient, eo EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc := new(gophercloud.ServiceClient)
	if eo.IronicInspectorEndpoint == "" {
		return nil, fmt.Errorf("IronicInspectorEndpoint is required")
	}
	if eo.IronicInspectorToken == "" {
		return nil, fmt.Errorf("Token required")
	}

	sc.MoreHeaders = map[string]string{"X-Auth-Token": eo.IronicInspectorToken}
	sc.Endpoint = gophercloud.NormalizeURL(eo.IronicInspectorEndpoint)
	sc.ProviderClient = client
	return sc, nil
}

// NewBareMetalIntrospectionHTTPToken creates a ServiceClient that may be used to access a
// "token" bare metal introspection service.
func NewBareMetalIntrospectionHTTPToken(eo EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(&gophercloud.ProviderClient{}, eo)
	if err != nil {
		return nil, err
	}

	sc.Type = "baremetal-inspector"

	return sc, nil
}
