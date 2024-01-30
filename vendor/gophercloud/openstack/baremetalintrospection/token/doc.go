/*
Package token provides support for http_basic bare metal endpoints.

Example of obtaining and using a client:

      client, err := token.NewBareMetalIntrospectionHTTPToken(token.Endpoints{
		IronicInspectorEndpoint:     "http://localhost:6385/v1/",
		IronicInspectorToken:        "token",
	})
	if err != nil {
		panic(err)
	}

	client.Microversion = "1.50"
	nodes.ListDetail(client, nodes.listOpts{})
*/
package token
