package instances

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/db/v1/flavors"
	"github.com/rackspace/gophercloud/openstack/db/v1/users"
	"github.com/rackspace/gophercloud/pagination"
)

// Volume represents information about an attached volume for a database instance.
type Volume struct {
	// The size in GB of the volume
	Size int

	Used float64
}

// Instance represents a remote MySQL instance.
type Instance struct {
	// Indicates the datetime that the instance was created
	Created string //time.Time

	// Indicates the most recent datetime that the instance was updated.
	Updated string //time.Time

	// Indicates the hardware flavor the instance uses.
	Flavor flavors.Flavor

	// A DNS-resolvable hostname associated with the database instance (rather
	// than an IPv4 address). Since the hostname always resolves to the correct
	// IP address of the database instance, this relieves the user from the task
	// of maintaining the mapping. Note that although the IP address may likely
	// change on resizing, migrating, and so forth, the hostname always resolves
	// to the correct database instance.
	Hostname string

	// Indicates the unique identifier for the instance resource.
	ID string

	// Exposes various links that reference the instance resource.
	Links []gophercloud.Link

	// The human-readable name of the instance.
	Name string

	// The build status of the instance.
	Status string

	// Information about the attached volume of the instance.
	Volume Volume
}

type commonResult struct {
	gophercloud.Result
}

// CreateResult represents the result of a Create operation.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a Get operation.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a Delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Extract will extract an Instance from various result structs.
func (r commonResult) Extract() (*Instance, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		Instance Instance `mapstructure:"instance"`
	}

	err := mapstructure.Decode(r.Body, &response)

	return &response.Instance, err
}

// InstancePage represents a single page of a paginated instance collection.
type InstancePage struct {
	pagination.LinkedPageBase
}

// IsEmpty checks to see whether the collection is empty.
func (page InstancePage) IsEmpty() (bool, error) {
	instances, err := ExtractInstances(page)
	if err != nil {
		return true, err
	}
	return len(instances) == 0, nil
}

// NextPageURL will retrieve the next page URL.
func (page InstancePage) NextPageURL() (string, error) {
	type resp struct {
		Links []gophercloud.Link `mapstructure:"instances_links"`
	}

	var r resp
	err := mapstructure.Decode(page.Body, &r)
	if err != nil {
		return "", err
	}

	return gophercloud.ExtractNextURL(r.Links)
}

// ExtractInstances will convert a generic pagination struct into a more
// relevant slice of Instance structs.
func ExtractInstances(page pagination.Page) ([]Instance, error) {
	casted := page.(InstancePage).Body

	var response struct {
		Instances []Instance `mapstructure:"instances"`
	}

	err := mapstructure.Decode(casted, &response)

	return response.Instances, err
}

// UserRootResult represents the result of an operation to enable the root user.
type UserRootResult struct {
	gophercloud.Result
}

// Extract will extract root user information from a UserRootResult.
func (r UserRootResult) Extract() (*users.User, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		User users.User `mapstructure:"user"`
	}

	err := mapstructure.Decode(r.Body, &response)

	return &response.User, err
}

// ActionResult represents the result of action requests, such as: restarting
// an instance service, resizing its memory allocation, and resizing its
// attached volume size.
type ActionResult struct {
	gophercloud.ErrResult
}
