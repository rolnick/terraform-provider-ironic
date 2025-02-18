package token

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
)

// EndpointOpts specifies a "http_basic" Ironic Endpoint
type EndpointOpts struct {
	IronicEndpoint string
	IronicToken    string
}

func initClientOpts(client *gophercloud.ProviderClient, eo EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc := new(gophercloud.ServiceClient)
	if eo.IronicEndpoint == "" {
		return nil, fmt.Errorf("IronicEndpoint is required")
	}
	if eo.IronicToken == "" {
		return nil, fmt.Errorf("Token is required")
	}
	sc.MoreHeaders = map[string]string{"X-Auth-Token": eo.IronicToken}
	sc.Endpoint = gophercloud.NormalizeURL(eo.IronicEndpoint)
	sc.ProviderClient = client
	return sc, nil
}

// NewBareMetalHTTPBasic creates a ServiceClient that may be used to access a
// "http_basic" bare metal service.
func NewBareMetalHTTPToken(eo EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc, err := initClientOpts(&gophercloud.ProviderClient{}, eo)
	if err != nil {
		return nil, err
	}

	sc.Type = "baremetal"

	return sc, nil
}
