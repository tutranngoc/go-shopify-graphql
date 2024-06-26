package graphqlclient

import (
	"fmt"
	"net/http"

	"github.com/gempages/go-shopify-graphql/graphql"
)

const (
	shopifyAccessTokenHeader           = "X-Shopify-Access-Token"
	shopifyStoreFrontAccessTokenHeader = "X-Shopify-Storefront-Access-Token"
)

var (
	apiProtocol          = "https"
	defaultAPIPathPrefix = "admin/api"
	defaultAPIVersion    = "2024-01"
	apiEndpoint          = "graphql.json"
)

// Option is used to configure options
type Option func(t *transport)

// WithVersion optionally sets the API version if the passed string is valid
func WithVersion(graphqlApiVersion string) Option {
	return func(t *transport) {
		if graphqlApiVersion != "" && graphqlApiVersion != "latest" {
			t.apiVersion = graphqlApiVersion
		}
	}
}

func WithStoreFrontVersion(apiVersion string) Option {
	return func(t *transport) {
		if apiVersion != "" && apiVersion != "latest" {
			t.apiVersion = apiVersion
		}
	}
}

// WithToken optionally sets oauth token
func WithToken(token string) Option {
	return func(t *transport) {
		t.accessToken = token
	}
}

func WithStoreFrontToken(token string) Option {
	return func(t *transport) {
		t.storeFrontAccessToken = token
		t.apiPath = "api"
	}
}

// WithPrivateAppAuth optionally sets private app credentials
func WithPrivateAppAuth(apiKey string, password string) Option {
	return func(t *transport) {
		t.apiKey = apiKey
		t.password = password
	}
}

type transport struct {
	accessToken           string
	storeFrontAccessToken string
	apiKey                string
	password              string
	apiVersion            string
	apiPath               string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.accessToken != "" {
		req.Header.Set(shopifyAccessTokenHeader, t.accessToken)
	} else if t.apiKey != "" && t.password != "" {
		req.SetBasicAuth(t.apiKey, t.password)
	} else if t.storeFrontAccessToken != "" {
		req.Header.Set(shopifyStoreFrontAccessTokenHeader, t.storeFrontAccessToken)
	}

	return http.DefaultTransport.RoundTrip(req)
}

// NewClient creates a new client (in fact, just a simple wrapper for a graphql.Client)
func NewClient(shopifyDomain string, opts ...Option) *graphql.Client {
	trans := &transport{
		apiPath:    defaultAPIPathPrefix,
		apiVersion: defaultAPIVersion,
	}

	for _, opt := range opts {
		opt(trans)
	}

	httpClient := &http.Client{Transport: trans}
	url := buildAPIEndpoint(shopifyDomain, trans.apiPath, trans.apiVersion)
	graphClient := graphql.NewClient(url, httpClient)
	return graphClient
}

func buildAPIEndpoint(domain string, path string, version string) string {
	return fmt.Sprintf("%s://%s/%s/%s/%s", apiProtocol, domain, path, version, apiEndpoint)
}
