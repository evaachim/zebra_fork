// Migration Script for data - it can be used to fetch,
// add, and use the data inside zebra.
package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	//nolint:gci

	"github.com/go-logr/logr"

	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model"
	"github.com/project-safari/zebra/model/compute"
	"github.com/project-safari/zebra/model/dc"
	"github.com/project-safari/zebra/model/network"

	// this is needed for mysql access.
	_ "github.com/go-sql-driver/mysql"
)

var ErrEmptyBody = errors.New("empty request body, cannot proceed")

// Racktables struct is the struct that contains info from the racktables table in the mysql db.
//
// It contains the id, ip, name, and type of the item in the racktable.
type Racktables struct {
	ID        string `json:"object_id"` //nolint:tagliatelle
	Name      string `json:"name"`
	Label     string `json:"label"`
	ObjtypeID string `json:"objtypeId"`
	AssetNo   string `json:"assetNo"`
	Problems  string `json:"hasProblems"`
	Comments  string `json:"comment"`
	IP        string `json:"ip"`
	Type      string `json:"type"`
	Port      int    `json:"port"`
	RackID    string `json:"rackId"`
	RowName   string `json:"rowName"`
	Owner     string `json:"owner"`
	RowID     string `json:"rowId"`
	Location  string `json:"locationName"`
}

type ResourceAPI struct {
	factory zebra.ResourceFactory
	Store   zebra.Store
}

type CtxKey string

const (
	ResourcesCtxKey = CtxKey("resources")
)

func NewResourceAPI(factory zebra.ResourceFactory) *ResourceAPI {
	return &ResourceAPI{
		factory: factory,
		Store:   nil,
	}
}

// Determine the specific type of a resource.
// nolint
func determineType(means string, resName string) string {
	name := strings.ToLower(resName)
	typ := ""

	if means == "Shelf" {
		typ = "dc.rack"
	} else if means == "Compute" {
		if strings.Contains(name, "esx") {
			typ = "compute.esx"
		} else if strings.Contains(name, "jenkins") || strings.Contains(name, "server") || strings.Contains(name, "srv") || strings.Contains(name, "vintella") {
			typ = "compute.server"
		} else if strings.Contains(name, "datacenter") || strings.Contains(name, "dc") || strings.Contains(name, "bld") {
			typ = "dc.datacenter"
		} else if strings.Contains(name, "dmz") || strings.Contains(name, "vlan") || strings.Contains(name, "asa") || strings.Contains(name, "bridge") {
			typ = "network.vlanPool"
		} else if strings.Contains(name, "vleaf") || strings.Contains(name, "switch") || strings.Contains(name, "sw") || strings.Contains(name, "aci") {
			typ = "network.switch"
		} else if strings.Contains(name, "vm") || strings.Contains(name, "capic") || strings.Contains(name, "frodo") {
			typ = "compute.vm"
		} else if strings.Contains(name, "vapic") || strings.Contains(name, "vpod") {
			typ = "compute.vcenter"
		} else if strings.Contains(name, "ipc") {
			typ = "network.ipAddressPool"
		}
	} else if means == "Other" {
		if strings.Contains(name, "chasis") || strings.Contains(name, "ixia") || strings.Contains(name, "rack") {
			typ = "dc.rack"
		} else if strings.Contains(name, "nexus") || strings.Contains(name, "sw") || strings.Contains(name, "switch") || strings.Contains(name, "n3k") {
			typ = "network.switch"
		}
	} else {
		typ = means
	}

	return typ
}

// Get resource type by id.
// nolint
func determineIDMeaning(id string, name string) string {
	means := ""
	final := ""
	this := ""

	if id == "2" || id == "27" {
		means = "compute.vm"
	} else if id == "30" || id == "31" || id == "34" || id == "3" {
		means = "dc.rack"
	} else if id == "38" {
		means = "compute.vcenter"
	} else if id == "4" || id == "13" || id == "36" {
		means = "compute.server"
	} else if id == "8" || id == "12" || id == "14" || id == "21" || id == "26" || id == "32" || id == "33" {
		means = "network.switch"
	} else if id == "1504" {
		means = "Compute"
	} else if id == "1503" {
		means = "Other"
	} else {
		means = "/"
	}

	final = determineType(means, name)

	if final == "/" {
		final = "unclassified"
	}

	this = final

	return this
}

//nolint:funlen
func main() {
	var rt Racktables

	RackArr := []Racktables{}

	// Statement to query the db - currently only one rack, 76.
	statement := "SELECT rack_id, object_id  FROM rackspace WHERE rack_id = 76"
	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:1234@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&rt.RackID, &rt.ID)

		rt.ID, rt.Name, rt.Label, rt.ObjtypeID, rt.AssetNo, rt.Problems, rt.Comments = getMoreDetails(rt.ID)

		typeID := rt.ObjtypeID
		resType := determineIDMeaning(typeID, rt.Name)
		rt.Type = resType

		if strings.Contains(resType, "compute") || resType == "network.switch" {
			rt.IP = getIPDetaiLs(rt.ID)

			ownedBy := getUserDetails(rt.IP)
			rt.Owner = ownedBy
		} else {
			rt.IP = "null"

			ownedBy := "null"
			rt.Owner = ownedBy
		}

		if resType == "network.switch" {
			portInfo := getPortDetails(rt.ID)
			rt.Port = portInfo
		} else {
			portID := -1
			rt.Port = portID
		}

		rackID := getRackDetails(rt.ID)
		rowName, rowID, rowLocation := getRowDetails(rackID)

		rt.RowName = rowName
		rt.RowID = rowID
		rt.Location = rowLocation

		assetNumber := rt.AssetNo
		rt.AssetNo = assetNumber

		probs := rt.Problems
		rt.Problems = probs

		notes := rt.Comments
		rt.Comments = notes

		RackArr = append(RackArr, rt)

		if err != nil {
			panic(err.Error())
		}
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	allData(RackArr)
}

func allData(rackArr []Racktables) {
	factory := zebra.Factory()

	myAPI := NewResourceAPI(factory)

	h := processPost()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	for i := 0; i < (len(rackArr)); i++ {
		res := rackArr[i]

		_, eachRes := createResFromData(res)

		// Create new resource on zebra with post request.
		req := createRequest("POST", "/resources", eachRes, myAPI)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
	}
}

func createRequest(method string, url string,
	body string, api *ResourceAPI,
) *http.Request {
	ctx := context.WithValue(context.Background(), ResourcesCtxKey, api)
	req, _ := http.NewRequestWithContext(ctx, method, url, nil)

	if body != "" {
		req.Body = ioutil.NopCloser(bytes.NewBufferString(body))
		print("Added   ", body, "  successfully!\n")
	}

	return req
}

func readJSONdata(ctx context.Context, req *http.Request, data interface{}) error {
	log := logr.FromContextOrDiscard(ctx)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	log.Info("request", "body", string(body))

	if len(body) > 0 {
		err = json.Unmarshal(body, data)
	} else {
		err = ErrEmptyBody
	}

	return err
}

func addAndPost(resMap *zebra.ResourceMap, f func(zebra.Resource) error) error {
	for _, l := range resMap.Resources {
		for _, r := range l.Resources {
			if err := f(r); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateRes(ctx context.Context, resMap *zebra.ResourceMap) error {
	// Check all resources to make sure they are valid
	for _, l := range resMap.Resources {
		for _, r := range l.Resources {
			if err := r.Validate(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

func processPost() httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := req.Context()
		log := logr.FromContextOrDiscard(ctx)
		api, ok := ctx.Value(ResourcesCtxKey).(*ResourceAPI)

		if !ok {
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		resMap := zebra.NewResourceMap(model.Factory())

		// Read request, return error if applicable
		if err := readJSONdata(ctx, req, resMap); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("resources could not be created, could not read request")

			return
		}

		if validateRes(ctx, resMap) != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("resources could not be created, found invalid resource(s)")

			return
		}

		// Add all resources to store
		if addAndPost(resMap, api.Store.Create) != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Info("internal server error while creating resources")

			return
		}

		log.Info("successfully created resources")

		res.WriteHeader(http.StatusOK)
	}
}

// Get IPs from db based on type id.
func getIPDetaiLs(objectID string) string {
	var rt Racktables

	statement := "SELECT ip FROM IPv4Allocation WHERE object_id = ?"

	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:1234@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement, objectID)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer results.Close()

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&rt.IP)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}

	return rt.IP
}

// Get port IDs from db based on type ID.
func getPortDetails(objectID string) int {
	var rt Racktables

	numPort := 0

	statement := "SELECT id FROM Port WHERE object_id = ?"

	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:1234@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement, objectID)
	if err != nil {
		panic(err.Error())
	}

	defer results.Close()

	for results.Next() {
		err = results.Scan(&rt.Port)

		if err != nil {
			panic(err.Error())
		}

		numPort++
	}

	return numPort
}

// Get rack details using the resource's specific ID.
func getMoreDetails(objectID string) (string, string, string, string, string, string, string) {
	var rt Racktables

	var label sql.NullString

	var assetNo sql.NullString

	var comment sql.NullString

	statement := "SELECT id, name, label, objtype_id, asset_no, has_problems, comment FROM rackobject WHERE id = ?"
	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:1234@/racktables")

	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement, objectID)

	if err != nil {
		panic(err.Error())
	}

	defer results.Close()

	for results.Next() {
		err = results.Scan(&rt.ID, &rt.Name, &label, &rt.ObjtypeID, &assetNo, &rt.Problems, &comment)

		rt.Label = label.String

		rt.AssetNo = assetNo.String

		rt.Comments = comment.String

		if err != nil {
			panic(err.Error())
		}

		log.Print(rt.RackID)
	}

	return rt.ID, rt.Name, rt.Label, rt.ObjtypeID, rt.AssetNo, rt.Problems, rt.Comments
}

func getRackDetails(objID string) string {
	var rt Racktables

	statement := "SELECT rack_id FROM RackSpace WHERE object_id = ?"
	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:1234@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement, objID)

	if err != nil {
		panic(err.Error())
	}

	defer results.Close()

	for results.Next() {
		err = results.Scan(&rt.RackID)

		if err != nil {
			panic(err.Error())
		}
	}

	return rt.RackID
}

// Get row and location details based on rack info (rack ID).
func getRowDetails(id string) (string, string, string) {
	var rt Racktables

	statement := "SELECT row_id, row_name, location_name FROM Rack WHERE id = ?"

	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:1234@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement, id)
	if err != nil {
		panic(err.Error())
	}
	defer results.Close()

	for results.Next() {
		err = results.Scan(&rt.RowID, &rt.RowName, &rt.Location)

		if err != nil {
			panic(err.Error())
		}
	}

	return rt.RowID, rt.RowName, rt.Location
}

// Get owner / user details based on the resource's IP.
func getUserDetails(resIP string) string {
	var rt Racktables

	statement := "SELECT user FROM IPv4Log WHERE ip = ?"

	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "eachim:1234@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement, resIP)
	if err != nil {
		panic(err.Error())
	}

	defer results.Close()

	for results.Next() {
		err = results.Scan(&rt.Owner)

		if err != nil {
			panic(err.Error())
		}
	}

	return rt.Owner
}

// Function to create a resource given data obtained from db, guven a certain type.
//
// Returns a zebra.Resource and a string version of the resource struct to be used with APIs.
//
//nolint:cyclop, funlen, lll
func createResFromData(res Racktables) (zebra.Resource, string) {
	// dbResources, err := GetData()
	// var resType = ""
	resType := res.Type

	switch resType {
	case "dc.datacenetr":
		addR := dc.NewDatacenter(res.Location, res.Name, res.Owner, "system.group-datacenter")
		print("Added dc " + res.Name + "\n")
		this := `{"datacenter":[{"id":` + res.ID + `,"type:"` + resType + `,"name:"` + res.Name + `,owner:` + res.Owner + "}]}"

		return addR, this

	case "dc.lab":
		addR := dc.NewLab(res.Name, res.Owner, "system.group-datacenter-lab")
		print("Added lab " + res.Name + "\n")
		this := `{"lab":[{"id":` + res.ID + `,"type":` + resType + `,"name":` + res.Name + `,"owner":` + res.Owner + "}]}"

		return addR, this

	case "dc.rack", "dc.shelf":
		addR := dc.NewRack(res.RowName, res.RowID, res.Name, res.Location, res.Owner, "system.group-datacenter-lab-rack")
		print("Added rack " + res.Name + "\n")
		this := `{"rack":[{"id":` + res.ID + `,"type":` + resType + `,"name":` + res.Name + `,"owner":` + res.Owner + `,"Row":` + res.RowName + `,"RowID":` + res.RowID + `,"Asset":` + res.AssetNo + `,"RowID":` + res.RowID + `,"Problems":` + res.Problems + `,"Location":` + res.Location + "}]}"

		return addR, this

	case "compute.server":
		addR := compute.NewServer("serial", "model", res.Name, res.Owner, "system.group-server")
		print("Added server " + res.Name + "\n")
		this := `{"server":[{"id":` + res.ID + `,"type":` + resType + `,"name":` + res.Name + `,"owner":` + res.Owner + `,"boardIP":` + res.IP + "}]}"

		return addR, this

	case "compute.esx":
		addR := compute.NewESX(res.ID, res.Name, res.Owner, "system.group-server-esx")
		print("Added esx " + res.Name + "\n")
		this := `{"esx":[{"id":` + res.ID + `,"type":` + resType + `,"name":` + res.Name + `,"owner":` + res.Owner + `,"ip":` + res.IP + "}]}"

		return addR, this

	case "compute.vm":
		addR := compute.NewVM("esx??", res.Name, res.Owner, "system.group-server-vcenter-vm")
		print("Added esx" + res.Name + "\n")
		this := `{"vm":[{"id":` + res.ID + `,"type":` + resType + `,"name":` + res.Name + `,"owner":` + res.Owner + `,"ip":` + res.IP + "}]}"

		return addR, this

	case "compute.vcenetr":
		addR := compute.NewVCenter(res.Name, res.Owner, "system.group-server-vcenter")
		print("Added vc " + res.Name + "\n")
		this := `{"vcenter":[{"id":` + res.ID + `,"type":` + resType + `,"name":` + res.Name + `,"owner":` + res.Owner + `,"ip":` + res.IP + "}]}"

		return addR, this

	case "network.switch":
		addR := network.NewSwitch(res.Name, res.Owner, "system.group-vlan-switch")
		print("Added sw " + res.Name + "\n")
		this := `{"switch":[{"id":` + res.ID + `,"type":` + resType + `,"name":` + res.Name + `,"owner":` + res.Owner + `,"managementIp":` + res.IP + `,"numPorts":` + strconv.Itoa(res.Port) + "}]}"

		return addR, this

	case "network.ipaddresspool":
		addR := network.NewIPAddressPool(res.Name, res.Owner, "system.group-vlan-ipaddrpool")
		print("Added IPpool" + res.Name + "\n")
		this := `{"IPAddressPool":[{"id":` + res.ID + `,"type":` + resType + `,"name":` + res.Name + `,"owner":` + res.Owner + "}]}"

		return addR, this

	case "network.vlanpool":
		addR := network.NewVLANPool(res.Name, res.ObjtypeID, "system.group-vlan")
		print("Added vlan" + res.Name + "\n")
		this := `{"VLANPool":[{"id":` + res.ID + `,"type":` + resType + `,"name":` + res.Name + `,"owner":` + res.Owner + "}]}"

		return addR, this
	}

	return nil, ""
}
