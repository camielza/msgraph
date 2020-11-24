package resources

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"net/url"
	"westpac.co.nz/msgraph/pkg/helpers"
	"westpac.co.nz/msgraph/pkg/msgraph"
)

// GraphAPIV1UserListResponse User List Response
type GraphAPIV1UserListResponse struct {
	Users []GraphAPIV1UserResponse `json:"value"`
}

// GraphAPIV1UserResponse Graph API users resource response
type GraphAPIV1UserResponse struct {
	ID                string   `json:"id"`
	DisplayName       string   `json:"displayName"`
	Mail              string   `json:"mail"`
	BusinessPhones    []string `json:"businessPhones"`
	GivenName         string   `json:"givenName"`
	JobTitle          string   `json:"jobTitle"`
	MobilePhone       string   `json:"mobilePhone"`
	OfficeLocation    string   `json:"officeLocation"`
	PreferredLanguage string   `json:"preferredLanguage"`
	Surname           string   `json:"surname"`
	UserPrincipalName string   `json:"userPrincipalName"`
}

func (g GraphAPIV1UserResponse) ToString() string {
	return g.DisplayName
}

// UsersResource UsersResource
type UsersResource struct{}

func (g UsersResource) ConvertToResourceSlice(body []byte) []msgraph.Resource {
	var userList GraphAPIV1UserListResponse
	err := json.Unmarshal(body, &userList)
	helpers.ErrorHandlerFatal("JSON unmarshalling of response body failed:", err)

	log.Tracef("UNMASHALLED OBJECT: %+v", userList)

	return g.toResourcArr(userList)
}

func (g UsersResource) toResourcArr(userList GraphAPIV1UserListResponse) []msgraph.Resource {
	var resources = make([]msgraph.Resource, len(userList.Users))
	for index, value := range userList.Users {
		resources[index] = value
	}
	return resources
}

func (g UsersResource) CreateRequestPath(context cli.Context, args cli.Args) string {
	return "/v1.0/users"
}

func (g UsersResource) CreateQueryParams(context cli.Context, args cli.Args) url.Values {

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
