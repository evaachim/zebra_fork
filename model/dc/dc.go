// Package dc provides structs and functions pertaining to datacenter resources.
package dc

import (
	"context"
	"errors"

	"github.com/project-safari/zebra"
)

var ErrAddressEmpty = errors.New("address is empty")

var ErrRowEmpty = errors.New("row is empty")

// Function that returns a zabra type of name datacenter and dc category.
func DatacenterType() zebra.Type {
	return zebra.Type{
		Name:        "dc.datacenter",
		Description: "data center",
	}
}

func EmptyDatacenter() zebra.Resource {
	d := new(Datacenter)
	d.Meta.Type = DatacenterType()

	return d
}

// A Datacenter represents the physical building. It is a named resource also
// with a building address.
type Datacenter struct {
	zebra.BaseResource
	Address string `json:"address"`
}

// Create new dc resources.
// Function that creates a new resource of type datacenter.
//
// It takes in a name, an owner, and a group,
// and returns a pointer to Datacenter.
// Create new dc resources.
func NewDatacenter(address, name, owner, group string) *Datacenter {
	return &Datacenter{
		BaseResource: *zebra.NewBaseResource(DatacenterType(), name, owner, group),
		Address:      address,
	}
}

// Validate returns an error if the given Datacenter object has incorrect values.
// Else, it returns nil.
func (dc *Datacenter) Validate(ctx context.Context) error {
	if dc.Address == "" {
		return ErrAddressEmpty
	}

	if dc.Meta.Type.Name != "dc.datacenter" {
		return zebra.ErrWrongType
	}

	return dc.BaseResource.Validate(ctx)
}

// Function that returns a zabra type of name lab and dc category.
func LabType() zebra.Type {
	return zebra.Type{
		Name:        "dc.lab",
		Description: "data center lab",
	}
}

func EmptyLab() zebra.Resource {
	l := new(Lab)
	l.Meta.Type = LabType()

	return l
}

// A Lab represents the lab consisting of a name and an ID.
type Lab struct{ zebra.BaseResource }

// Create new dc resources.
// Function that creates a new resource of type lab.
//
// It takes in a name, an owner, and a group,
// and returns a pointer to Lab.
// Create new dc resources.
func NewLab(name, owner, group string) *Lab {
	return &Lab{
		BaseResource: *zebra.NewBaseResource(LabType(), name, owner, group),
	}
}

// Function to validate an lab resource, given a pointer to the lab struct.
func (l *Lab) Validate(ctx context.Context) error {
	if l.Meta.Type.Name != "dc.lab" {
		return zebra.ErrWrongType
	}

	return l.BaseResource.Validate(ctx)
}

// Function that returns a zabra type of name rack and dc category.
func RackType() zebra.Type {
	return zebra.Type{
		Name:        "dc.rack",
		Description: "server rack",
	}
}

func EmptyRack() zebra.Resource {
	r := new(Rack)
	r.Meta.Type = RackType()

	return r
}

// A Rack represents a datacenter rack. It consists of a name, ID, and associated
// row.
type Rack struct {
	zebra.BaseResource
	// Might want to change the json to row_name to distinguish and to match the db data.
	Row string `json:"row"`

	// adding data from data base.
	RowID string `json:"rowId"`
	// Asset must be of different type - to look into.
	Asset    int    `json:"assetNo"`
	Problems string `json:"hasProblems"`
	Comment  string `json:"comment"`
	// This should be different from address and is of format address/labNo.
	Location string `json:"locationId"`
}

// Validate returns an error if the given Rack object has incorrect values.
// Else, it returns nil.
func (r *Rack) Validate(ctx context.Context) error {
	if r.Row == "" {
		return ErrRowEmpty
	}

	if r.Meta.Type.Name != "dc.rack" {
		return zebra.ErrWrongType
	}

	return r.BaseResource.Validate(ctx)
}

// Function that creates a new resource of type rack.
//
// It takes in a name, an owner, and a group,
// and returns a pointer to Rack.
func NewRack(row, rowID, name, locate, owner, group string) *Rack {
	return &Rack{
		BaseResource: *zebra.NewBaseResource(RackType(), name, owner, group),
		Row:          row,
		// added info to correspond to data in the db.
		RowID: rowID,
		// Asset: assets,
		// Problems: prob,
		// Comment:  comments,
		Location: locate,
	}
}
