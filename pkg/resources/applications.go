package resources

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net/url"
	"time"
	"westpac.co.nz/msgraph/pkg/helpers"
	"westpac.co.nz/msgraph/pkg/msgraph"
)

// GraphAPIV1ApplicationListResponse Application List Response
type GraphAPIV1ApplicationListResponse struct {
	Applications []GraphAPIV1ApplicationResponse `json:"value"`
}

// GraphAPIV1ApplicationResponse Graph API applications resource response
type GraphAPIV1ApplicationResponse struct {
	ID                        string    `json:"id"`
	DisplayName               string    `json:"displayName"`
	DeletedDataTime           time.Time `json:"deletedDateTime"`
	AppID                     string    `json:"appId"`
	ApplicationTemplateID     string    `json:"applicationTemplateId"`
	CreatedDateTime           time.Time `json:"createdDateTime"`
	Description               string    `json:"description"`
	IdentifierURIs            []string  `json:"identifierUris"`
	IsDeviceOnlyAuthSupported bool      `json:"isDeviceOnlyAuthSupported"`
	IsFallbackPublicClient    bool      `json:"isFallbackPublicClient"`
	Notes                     string    `json:"notes"`
	PublisherDomain           string    `json:"publisherDomain"`
	SignInAudience            string    `json:"signInAudience"`
	Tags                      []string  `json:"tags"`
	TokenEncryptionKeyID      string    `json:"tokenEncryptionKeyId"`
}

func (g GraphAPIV1ApplicationResponse) ToString() string {
	return g.DisplayName
}

// ApplicationsResource ApplicationsResource
type ApplicationsResource struct{}

func (g ApplicationsResource) ConvertToResourceSlice(body []byte) []msgraph.Resource {
	var applicationList GraphAPIV1ApplicationListResponse
	err := json.Unmarshal(body, &applicationList)
	helpers.ErrorHandlerFatal("JSON unmarshalling of response body failed:", err)

	log.Tracef("UNMASHALLED OBJECT: %+v", applicationList)

	return g.toResourceArr(applicationList)
}

func (g ApplicationsResource) toResourceArr(applicationList GraphAPIV1ApplicationListResponse) []msgraph.Resource {
	var resources = make([]msgraph.Resource, len(applicationList.Applications))
	for index, value := range applicationList.Applications {
		resources[index] = value
	}
	return resources
}

func (g ApplicationsResource) CreateRequestPath(context cli.Context, args cli.Args) string {
	return "/v1.0/applications"
}

func (g ApplicationsResource) CreateQueryParams(context cli.Context, args cli.Args) url.Values {

	startWith := ""
	if args.Len() > 0 {
		startWith = args.Get(0)
	}

	if startWith != "" {
		filter := new(msgraph.FilterCriteria)
		criteria := filter.StartWith("displayName", startWith)
		return msgraph.CreateURLFilterParams(criteria)
	}
	return nil
}
