package store

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/filestore"
	"github.com/project-safari/zebra/idstore"
	"github.com/project-safari/zebra/labelstore"
	"github.com/project-safari/zebra/typestore"
)

// ErrNilResource prevents nil resources from being created.
var ErrNilResource = errors.New("nil resource not allowed")

// ResourceStore is a struct used for creation of new resource stores.
// a resource store contains a lock, a root path for storage,
// a  zebra.ResourceFactory,
// and pointers to filestore.FileStore, idstore.IDStore, typestore.TypeStore.
type ResourceStore struct {
	lock        sync.RWMutex
	StorageRoot string
	Factory     zebra.ResourceFactory
	fs          *filestore.FileStore
	ids         *idstore.IDStore
	ls          *labelstore.LabelStore
	ts          *typestore.TypeStore
}

// Function uses the ResourceStore struct to generate a new resource store.
// It takes in a string to the root path and a  zebra.ResourceFactory.
// Returns a pointer to ResourceStore.
func NewResourceStore(root string, factory zebra.ResourceFactory) *ResourceStore {
	return &ResourceStore{
		lock:        sync.RWMutex{},
		StorageRoot: root,
		Factory:     factory,
		fs:          nil,
		ids:         nil,
		ls:          nil,
		ts:          nil,
	}
}

// Operation function on ResourceStore: initialization.
func (rs *ResourceStore) Initialize() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	rs.fs = filestore.NewFileStore(rs.StorageRoot, rs.Factory)
	if err := rs.fs.Initialize(); err != nil {
		return err
	}

	resources, err := rs.fs.Load()
	if err != nil {
		return err
	}

	rs.ids = idstore.NewIDStore(resources)
	rs.ls = labelstore.NewLabelStore(resources)
	rs.ts = typestore.NewTypeStore(resources)

	return nil
}

// Operation function on ResourceStore: wiping.
func (rs *ResourceStore) Wipe() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	rs.fs = nil
	rs.ids = nil
	rs.ls = nil
	rs.ts = nil

	return nil
}

// Operation function on ResourceStore: clearing.
func (rs *ResourceStore) Clear() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	if err := rs.fs.Clear(); err != nil {
		return err
	}

	if err := rs.ids.Clear(); err != nil {
		return err
	}

	if err := rs.ls.Clear(); err != nil {
		return err
	}

	if err := rs.ts.Clear(); err != nil {
		return err
	}

	return nil
}

// Return ResourceMap with resource type as key and list of resources as val.
func (rs *ResourceStore) Load() (*zebra.ResourceMap, error) {
	rs.lock.RLock()
	defer rs.lock.RUnlock()

	return rs.ts.Load()
}

// Operation function on ResourceStore: create.
func (rs *ResourceStore) Create(res zebra.Resource) error {
	if res == nil {
		return ErrNilResource
	}

	if err := res.Validate(context.Background()); err != nil {
		return err
	}

	rs.lock.Lock()
	defer rs.lock.Unlock()

	err := rs.fs.Create(res)
	if err != nil {
		return err
	}

	err = rs.ids.Create(res)
	if err != nil {
		return err
	}

	err = rs.ls.Create(res)
	if err != nil {
		return err
	}

	err = rs.ts.Create(res)
	if err != nil {
		return err
	}

	return nil
}

// Operation function on ResourceStore: delete.
func (rs *ResourceStore) Delete(resource zebra.Resource) error {
	if resource == nil || resource.Validate(context.Background()) != nil {
		return zebra.ErrInvalidResource
	}

	rs.lock.Lock()
	defer rs.lock.Unlock()

	err := rs.fs.Delete(resource)
	if err != nil {
		return err
	}

	err = rs.ids.Delete(resource)
	if err != nil {
		return err
	}

	err = rs.ls.Delete(resource)
	if err != nil {
		return err
	}

	err = rs.ts.Delete(resource)
	if err != nil {
		return err
	}

	return nil
}

// Return all resources in a ResourceMap.
func (rs *ResourceStore) Query() *zebra.ResourceMap {
	rs.lock.RLock()
	defer rs.lock.RUnlock()

	resMap, err := rs.ts.Load()
	if err != nil {
		return nil
	}

	retMap := zebra.NewResourceMap(resMap.GetFactory())

	zebra.CopyResourceMap(retMap, resMap)

	return retMap
}

// Return resources with matching UUIDs.
func (rs *ResourceStore) QueryUUID(uuids []string) *zebra.ResourceMap {
	rs.lock.RLock()
	defer rs.lock.RUnlock()

	resMap := rs.ids.Query(uuids)
	retMap := zebra.NewResourceMap(resMap.GetFactory())

	zebra.CopyResourceMap(retMap, resMap)

	return retMap
}

// Return resources with matching types.
func (rs *ResourceStore) QueryType(types []string) *zebra.ResourceMap {
	rs.lock.RLock()
	defer rs.lock.RUnlock()

	resMap := rs.ts.Query(types)
	retMap := zebra.NewResourceMap(resMap.GetFactory())

	zebra.CopyResourceMap(retMap, resMap)

	return retMap
}

// Return resources with matching label.
func (rs *ResourceStore) QueryLabel(query zebra.Query) (*zebra.ResourceMap, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}

	rs.lock.RLock()
	defer rs.lock.RUnlock()

	resMap := rs.ls.Query(query)
	retMap := zebra.NewResourceMap(resMap.GetFactory())

	zebra.CopyResourceMap(retMap, resMap)

	return retMap, nil
}

// Return resources which match given property/value(s).
// Naive search implementation, >= O(n) for n resources.
func (rs *ResourceStore) QueryProperty(query zebra.Query) (*zebra.ResourceMap, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}

	rs.lock.RLock()
	defer rs.lock.RUnlock()

	if query.Op == zebra.MatchEqual || query.Op == zebra.MatchIn {
		return rs.propertyMatch(query, true)
	}

	return rs.propertyMatch(query, false)
}

// Function to check for matches in properties.
// Returns a pointer to zebra.ResourceMap and an error where appropriate (or nil otherwise).
func (rs *ResourceStore) propertyMatch(query zebra.Query, inVals bool) (*zebra.ResourceMap, error) {
	resMap, err := rs.ts.Load()
	if err != nil {
		return nil, err
	}

	retMap := zebra.NewResourceMap(rs.Factory)

	for t, l := range resMap.Resources {
		for _, res := range l.Resources {
			val := FieldByName(reflect.ValueOf(res).Elem(), query.Key).String()
			inList := zebra.IsIn(val, query.Values)

			if inVals && inList {
				retMap.Add(res, t)
			} else if !inVals && !inList {
				retMap.Add(res, t)
			}
		}
	}

	return retMap, nil
}

// Filter given map by uuids.
func FilterUUID(uuids []string, resMap *zebra.ResourceMap) (*zebra.ResourceMap, error) {
	retMap := zebra.NewResourceMap(resMap.GetFactory())

	for t, l := range resMap.Resources {
		for _, res := range l.Resources {
			if zebra.IsIn(res.GetID(), uuids) {
				retMap.Add(res, t)
			}
		}
	}

	return retMap, nil
}

// Filter given map by types.
func FilterType(types []string, resMap *zebra.ResourceMap) (*zebra.ResourceMap, error) {
	f := resMap.GetFactory()
	retMap := zebra.NewResourceMap(f)

	for _, t := range types {
		l, ok := resMap.Resources[t]
		if !ok {
			continue
		}

		copyL := zebra.NewResourceList(f)

		zebra.CopyResourceList(copyL, l)
		retMap.Resources[t] = copyL
	}

	return retMap, nil
}

// Filter given map by label name and val.
func FilterLabel(query zebra.Query, resMap *zebra.ResourceMap) (*zebra.ResourceMap, error) {
	if err := query.Validate(); err != nil {
		return resMap, err
	}

	retMap := zebra.NewResourceMap(resMap.GetFactory())

	inVals := false

	if query.Op == zebra.MatchEqual || query.Op == zebra.MatchIn {
		inVals = true
	}

	for t, l := range resMap.Resources {
		for _, res := range l.Resources {
			labels := res.GetLabels()
			matchIn := labels.MatchIn(query.Key, query.Values...)

			if (inVals && matchIn) || (!inVals && !matchIn) {
				retMap.Add(res, t)
			}
		}
	}

	return retMap, nil
}

// Filter given map by property name (case insensitive) and val.
func FilterProperty(query zebra.Query, resMap *zebra.ResourceMap) (*zebra.ResourceMap, error) {
	if err := query.Validate(); err != nil {
		return resMap, err
	}

	retMap := zebra.NewResourceMap(resMap.GetFactory())

	inVals := false

	if query.Op == zebra.MatchEqual || query.Op == zebra.MatchIn {
		inVals = true
	}

	for t, l := range resMap.Resources {
		for _, res := range l.Resources {
			val := FieldByName(reflect.ValueOf(res).Elem(), query.Key).String()
			matchIn := zebra.IsIn(val, query.Values)

			if (inVals && matchIn) || (!inVals && !matchIn) {
				retMap.Add(res, t)
			}
		}
	}

	return retMap, nil
}

// Ignore case in returning value of given field.
func FieldByName(v reflect.Value, field string) reflect.Value {
	field = strings.ToLower(field)

	return v.FieldByNameFunc(
		func(found string) bool {
			return strings.ToLower(found) == field
		})
}
