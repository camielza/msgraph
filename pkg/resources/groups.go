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

// GraphAPIV1GroupListResponse Group List Response
type GraphAPIV1GroupListResponse struct {
	Groups []GraphAPIV1GroupResponse `json:"value"`
}

// GraphAPIV1GroupResponse Graph API groups resource response
type GraphAPIV1GroupResponse struct {
	ID                    string    `json:"id"`
	DeletedDate           time.Time `json:"deletedDateTime"`
	Classification        string    `json:"classification"`
	CreatedDateTime       time.Time `json:"createdDateTime"`
	Description           string    `json:"description"`
	DisplayName           string    `json:"displayName"`
	Mail                  string    `json:"mail"`
	MailEnabled           bool      `json:"mailEnabled"`
	MailNickName          string    `json:"mailNickname"`
	LastSyncDateTime      time.Time `json:"onPremisesLastSyncDateTime"`
	SecurityIdentifier    string    `json:"onPremisesSecurityIdentifier"`
	SyncEnabled           bool      `json:"onPremisesSyncEnabled"`
	PreferredDataLocation string    `json:"preferredDataLocation"`
	RenewedDateTime       time.Time `json:"renewedDateTime"`
	SecurityEnabled       bool      `json:"securityEnabled"`
	Visibility            string    `json:"visibility"`
}

func (g GraphAPIV1GroupResponse) ToString() string {
	return g.DisplayName
}

// GroupsResource GroupsResource
type GroupsResource struct{}

func (g GroupsResource) ConvertToResourceSlice(body []byte) []msgraph.Resource {
	var groupList GraphAPIV1GroupListResponse
	err := json.Unmarshal(body, &groupList)
	helpers.ErrorHandlerFatal("JSON unmarshalling of response body failed:", err)

	log.Tracef("UNMASHALLED OBJECT: %+v", groupList)

	return g.toResourceArr(groupList)
}

func (g GroupsResource) toResourceArr(groupList GraphAPIV1GroupListResponse) []msgraph.Resource {
	var resources = make([]msgraph.Resource, len(groupList.Groups))
	for index, value := range groupList.Groups {
		resources[index] = value
	}
	return resources
}

func (g GroupsResource) CreateRequestPath(context cli.Context, args cli.Args) string {
	return "/v1.0/groups"
}

func (g GroupsResource) CreateQueryParams(context cli.Context, args cli.Args) url.Values {

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
