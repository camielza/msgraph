package main

import (
	"encoding/json"
	"fmt"
	"github.com/rodaine/table"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"reflect"
	"strings"
	"westpac.co.nz/msgraph/pkg/helpers"
	"westpac.co.nz/msgraph/pkg/msauth"
	"westpac.co.nz/msgraph/pkg/msgraph"
	"westpac.co.nz/msgraph/pkg/resources"
	"westpac.co.nz/msgraph/pkg/slices"
)

var resourceMap map[string]msgraph.ResourceAPI

func init() {

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.ErrorLevel)

	resourceMap = make(map[string]msgraph.ResourceAPI)

	var groupsAPI msgraph.ResourceAPI = resources.GroupsResource{}
	var usersAPI msgraph.ResourceAPI = resources.UsersResource{}
	var applicationsAPI msgraph.ResourceAPI = resources.ApplicationsResource{}

	resourceMap["groups"] = groupsAPI
	resourceMap["users"] = usersAPI
	resourceMap["applications"] = applicationsAPI

}

func stringCompare(item1 interface{}, item2 interface{}) bool {
	return strings.ToLower(string(item1.(string))) == strings.ToLower(string(item2.(string)))
}

func list(resourceName string, tenantID string, clientID string, clientSecret string, context cli.Context, args cli.Args) {

	log.Debug("Retrieving token...")

	var baseResource = msgraph.BaseResource{
		HTTPClient: msauth.GetOAuth2Client(tenantID, clientID, clientSecret),
		Version:    "v1.0",
	}

	resourceAPI := resourceMap[resourceName]

	resources := baseResource.List(resourceAPI, context, args)

	log.Debug("Fetched resource:", len(resources))

	var fields = []interface{}{}
	if context.IsSet("fields") {
		untrimmedFields := strings.Split(context.String("fields"), ",")
		for index := range untrimmedFields {
			fields = append(fields, strings.TrimSpace(untrimmedFields[index]))
		}
		log.Debug("FIELDS", fields)
	}

	var resourceStr string = ""
	switch context.String("output") {
	case "json":
		resourceBytes, err := json.Marshal(resources)
		helpers.ErrorHandlerFatal("Unable to Marshall resource to JSON", err)
		resourceStr = string(resourceBytes)
	case "text":
		for _, resource := range resources {
			resourceStr += fmt.Sprintf("%+v\n", resource)
		}
	case "table":
		var tbl table.Table
		headerDisplayed := false
		for _, resource := range resources {

			displayFields := slices.Union(getResourceFields(resource), fields, stringCompare)
			log.Debug("FIELDS", displayFields)

			if !headerDisplayed {
				tbl = table.New(displayFields...)
				headerDisplayed = true
			}

			values := getResourceValues(resource, displayFields)
			tbl.AddRow(values...)
		}
		tbl.Print()
	case "string":
		for _, resource := range resources {
			resourceStr += fmt.Sprintf("%s\n", resource.ToString())
		}
	}
	fmt.Println(resourceStr)
}

func getResourceValues(resource msgraph.Resource, headers []interface{}) []interface{} {
	var values []interface{}
	reflected := reflect.ValueOf(resource)
	for _, fieldName := range headers {
		values = append(values, reflected.FieldByName(fieldName.(string)))
	}
	return values
}

func getResourceFields(resource msgraph.Resource) []interface{} {
	var headers = make([]interface{}, 0)
	reflected := reflect.ValueOf(resource)
	for i := 0; i < reflected.NumField(); i++ {
		fieldName := reflected.Type().Field(i).Name
		headers = append(headers, fieldName)
	}
	return headers
}

func main() {

	var app = &cli.App{
		Usage:   "Azure MSGraph API",
		Version: "1.00",
		Authors: []*cli.Author{
			&cli.Author{
				Name: "Camiel de Vleeschauwer",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "tenant",
				Aliases:  []string{"t"},
				Usage:    "The Azure tenant ID to use for the query",
				EnvVars:  []string{"AZ_TENANT"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "clientID",
				Aliases:  []string{"c"},
				Usage:    "The Azure client ID for the SPN",
				EnvVars:  []string{"AZ_CLIENTID"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "clientSecret",
				Aliases:  []string{"s"},
				Usage:    "The Azure client Secret for the SPN",
				EnvVars:  []string{"AZ_CLIENTSECRET"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "verbose",
				Aliases:  []string{"V"},
				Usage:    fmt.Sprintf("How loud should I be: (%s)", []string{"trace", "debug", "info", "warning", "error"}),
				Required: false,
				Value:    log.InfoLevel.String(),
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    fmt.Sprintf("output format: (%s)", []string{"json", "text", "table", "string"}),
				Required: false,
				Value:    "string",
			},
			&cli.StringFlag{
				Name:     "fields",
				Aliases:  []string{"f"},
				Usage:    "command separated list of field to output",
				Required: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "groups",
				Aliases:     []string{"g"},
				Usage:       "The Azure Active Directory 'groups' resource",
				Description: "Actions for the groups resource",
				Subcommands: []*cli.Command{
					{
						Name:      "list",
						Aliases:   []string{"l"},
						Usage:     "list groups, optionally filtering by start string",
						ArgsUsage: "[filter - group name start]",
						Action: func(c *cli.Context) error {

							if c.IsSet("verbose") {
								level, err := log.ParseLevel(c.String("verbose"))
								helpers.ErrorHandlerFatal("Could not parse verbosity ", err)
								log.SetLevel(level)
							}
							list(
								"groups",
								c.String("tenant"),
								c.String("clientID"),
								c.String("clientSecret"),
								*c,
								c.Args(),
							)
							return nil
						},
					},
				},
			},
			{
				Name:        "users",
				Aliases:     []string{"u"},
				Usage:       "The Azure Active Directory 'users' resource",
				Description: "Actions for the users resource",
				Subcommands: []*cli.Command{
					{
						Name:      "list",
						Aliases:   []string{"l"},
						Usage:     "list users, optionally filtering by start string",
						ArgsUsage: "[filter - users name start]",
						Action: func(c *cli.Context) error {

							if c.IsSet("verbose") {
								level, err := log.ParseLevel(c.String("verbose"))
								helpers.ErrorHandlerFatal("Could not parse verbosity ", err)
								log.SetLevel(level)
							}
							list(
								"users",
								c.String("tenant"),
								c.String("clientID"),
								c.String("clientSecret"),
								*c,
								c.Args(),
							)
							return nil
						},
					},
				},
			},
			{
				Name:        "applications",
				Aliases:     []string{"a"},
				Usage:       "The Azure Active Directory 'applications' resource",
				Description: "Actions for the applications resource",
				Subcommands: []*cli.Command{
					{
						Name:      "list",
						Aliases:   []string{"l"},
						Usage:     "list applications, optionally filtering by start string",
						ArgsUsage: "[filter - applications name start]",
						Action: func(c *cli.Context) error {

							if c.IsSet("verbose") {
								level, err := log.ParseLevel(c.String("verbose"))
								helpers.ErrorHandlerFatal("Could not parse verbosity ", err)
								log.SetLevel(level)
							}
							list(
								"applications",
								c.String("tenant"),
								c.String("clientID"),
								c.String("clientSecret"),
								*c,
								c.Args(),
							)
							return nil
						},
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
