package zebra

import (
	"context"
)

/*

Alt containment method

Something like this:

type ResourceContents struct {
	left  *ResourceContents
	right *ResourceContents
	// data  BaseResource
	data interface{}
}

type BaseResource struct {
	Meta   Meta   `json:"meta"`
	Status Status `json:"status,omitempty"`
	root   *ResourceContents
	left   *dc.DataCenterType
	right  *compute.ServerType
	mid    *network.VLANPoolType
	// then use each of these as subtrees.
}


func getRes(n *BaseResource, key string) *node {
	var (
		p *node
		i int
	)
	for i = 0; i < len(key) && n != nil; {
		c := key[i]
		p = n
		if c < n.c {
			n = n.left
		} else if c > n.c {
			n = n.right
		} else {
			n = n.mid
			i++
		}
	}
	if i == len(key) {
		return p
	}
	return n
}
*/

// Resource interface is implemented by all resources and provides resource
// validation and label selection methods.
type Resource interface {
	Validate(ctx context.Context) error
	GetMeta() Meta
	GetStatus() Status
}

// A BaseResource struct represents a basic resource  with the appropriate meta data and a status.
type BaseResource struct {
	Meta   Meta   `json:"meta"`
	Status Status `json:"status,omitempty"`
}

// Function that calidates a basic resource (BaseResource).
// It returns an error or nil in the absence thereof.
func (r *BaseResource) Validate(ctx context.Context) error {
	if err := r.Meta.Validate(); err != nil {
		return err
	}

	if err := r.Status.Validate(); err != nil {
		return err
	}

	return nil
}

// Function on a pointer to BaseResource to get the meta data from the given resource.
func (r *BaseResource) GetMeta() Meta {
	return r.Meta
}

// Function on a pointer to BaseResource to get the status of the given resource.
func (r *BaseResource) GetStatus() Status {
	return r.Status
}

// Function that returns a pointer to a BaseResource struct
// with meta data (type, name, group, owner) for the resource and its status.
func NewBaseResource(rType Type, name, owner, group string) *BaseResource {
	return &BaseResource{
		Meta:   NewMeta(rType, name, group, owner),
		Status: DefaultStatus(),
	}
}
