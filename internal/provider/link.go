package provider

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	sharedmodels "github.com/hashicorp/cloud-sdk-go/clients/cloud-shared/v1/models"
)

// newLink constructs a new Link from the passed arguments. ID should be the
// user specified resource ID.
//
// Adapted from https://github.com/hashicorp/cloud-api-internal/blob/master/helper/hashicorp/cloud/location/link.go#L10-L23
func newLink(loc *sharedmodels.HashicorpCloudLocationLocation, svcType string, id string) *sharedmodels.HashicorpCloudLocationLink {
	return &sharedmodels.HashicorpCloudLocationLink{
		Type:     svcType,
		ID:       id,
		Location: loc,
	}
}

// linkURL generates a URL from the passed link. If the link is invalid, an
// error is returned. The Link URL is a globally unique, human readable string
// identifying a resource.
// This version of the function includes all location data (org, project, provider, and region).
//
// Adapted from https://github.com/hashicorp/cloud-api-internal/blob/master/helper/hashicorp/cloud/location/link.go#L25-L60
func linkURL(l *sharedmodels.HashicorpCloudLocationLink) (string, error) {
	if l == nil {
		return "", errors.New("nil link")
	}

	if l.Location == nil {
		return "", errors.New("link missing Location")
	}

	// Validate that the link contains the necessary information
	if l.Location.ProjectID == "" {
		return "", errors.New("link missing project ID")
	} else if l.Location.OrganizationID == "" {
		return "", errors.New("link missing organization ID")
	} else if l.Location.Region.Provider == "" {
		return "", errors.New("link missing provider")
	} else if l.Location.Region.Region == "" {
		return "", errors.New("link missing region")
	} else if l.Type == "" {
		return "", errors.New("link missing resource type")
	}

	// Determine the ID of the resource
	id := l.ID
	if id == "" {
		return "", errors.New("link missing resource ID")
	}

	// Generate the URL
	urn := fmt.Sprintf("/organization/%s/project/%s/provider/%s/region/%s/%s/%s",
		l.Location.OrganizationID,
		l.Location.ProjectID,
		l.Location.Region.Provider,
		l.Location.Region.Region,
		l.Type,
		id)

	return urn, nil
}

// parseLinkURL parses a link URL into a link. If the URL is malformed, an
// error is returned.
func parseLinkURL(urn string) (*sharedmodels.HashicorpCloudLocationLink, error) {
	match, _ := regexp.MatchString("^/organization/.+/project/.+/provider/.+/region/.+/.+/.+$", urn)
	if !match {
		return nil, errors.New("url is not in the correct format: /organization/{org_id}/project/{project_id}/provider/{provider}/region/{region}/{type}/{id}")
	}

	components := strings.Split(urn, "/")

	return &sharedmodels.HashicorpCloudLocationLink{
		Type: components[9],
		ID:   components[10],
		Location: &sharedmodels.HashicorpCloudLocationLocation{
			OrganizationID: components[2],
			ProjectID:      components[4],
			Region: &sharedmodels.HashicorpCloudLocationRegion{
				Provider: components[6],
				Region:   components[8],
			},
		},
	}, nil
}