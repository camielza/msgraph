package msgraph

import (
	"github.com/urfave/cli/v2"
	"net/url"
)

// Resource Resource
type Resource interface {
	ToString() string
}

// BaseResourceAPI BaseResourceAPI
type ResourceAPI interface {
	ConvertToResourceSlice(body []byte) []Resource
	CreateQueryParams(context cli.Context, args cli.Args) url.Values
	CreateRequestPath(context cli.Context, args cli.Args) string
}
