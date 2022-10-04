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
	left   *dc.Datacenter
	right  *compute.Server
	mid    *network.Switch
	// then use each of these as subtrees.
}

// subtrees for the left child tree of BaseResource (i.e. Datacenter) example of how to (re)structure
type Datacenter struct {
	zebra.BaseResource
	Address string `json:"address"`
	root   *BaseResource
	left   *dc.Lab
}

type Lab struct{
	zebra.BaseResource
	root   *Datacenter
	left   *dc.Rack
}

// subtrees for the mid child tree of BaseResource (i.e. Switch) example of how to (re)structure
type Switch struct {
	zebra.BaseResource
	Credentials  zebra.Credentials `json:"credentials"`
	ManagementIP net.IP            `json:"managementIp"`
	SerialNumber string            `json:"serialNumber"`
	Model        string            `json:"model"`
	NumPorts     uint32            `json:"numPorts"`
	root   *BaseResource
	left   *compute.VLANPool
	right *compute.IPAddressPool
}

type VLANPool struct {
	zebra.BaseResource
	RangeStart uint16 `json:"rangeStart"`
	RangeEnd   uint16 `json:"rangeEnd"`
	root   *Switch
}

type IPAddressPool struct {
	zebra.BaseResource
	Subnets []net.IPNet `json:"subnets"`
	root   *Switch
}

// subtrees for the mid child tree of BaseResource (i.e. Server) example of how to (re)structure
//
type Server struct {
	zebra.BaseResource
	Credentials  zebra.Credentials `json:"credentials"`
	SerialNumber string            `json:"serialNumber"`
	BoardIP      net.IP            `json:"boardIp"`
	Model        string            `json:"model"`
	root   *ResourceContents
	left   *compute.ESX
	right  *compute.VCenter
	mid    *network.VM
}

type ESX struct {
	zebra.BaseResource
	Credentials zebra.Credentials `json:"credentials"`
	ServerID    string            `json:"serverId"`
	IP          net.IP            `json:"ip"`
	root   *Server
}

type VCenter struct {
	zebra.BaseResource
	Credentials zebra.Credentials `json:"credentials"`
	IP          net.IP            `json:"ip"`
	root   *Server
}

type VM struct {
	zebra.BaseResource
	Credentials  zebra.Credentials `json:"credentials"`
	ESXID        string            `json:"esxId"`
	ManagementIP net.IP            `json:"managementIp"`
	VCenterID    string            `json:"vCenterId"`
	root   *Server
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
