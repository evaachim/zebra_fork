// Migration Script for data - it can be used to fetch,
// add, and use the data inside zebra.package migration
package migration

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	//nolint:gci

	"github.com/project-safari/zebra"

	// this is needed for mysql access.
	_ "github.com/go-sql-driver/mysql"
)

type Racktables struct {
	ID        string     `json:"object_id"` //nolint:tagliatelle
	Name      string     `json:"name"`
	Label     string     `json:"label"`
	ObjtypeID string     `json:"objtypeId"`
	AssetNo   string     `json:"assetNo"`
	Problems  string     `json:"hasProblems"`
	Comments  string     `json:"comment"`
	IP        string     `json:"ip"`
	Type      zebra.Type `json:"type"`
	Port      int        `json:"port"`
	RackID    string     `json:"rackId"`
	RowName   string     `json:"rowName"`
	Owner     string     `json:"owner"`
	RowID     string     `json:"rowId"`
	Location  string     `json:"locationName"`
	Serial    string     `json:"serialNumber"`
}

/* Beginning of Data filtering for migration. */

func lowerLetters(thisOne string) string {
	return strings.ToLower(thisOne)
}

func ComputeCase(name string) string {
	typ := ""

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
	} else {
		typ = "N/A"
	}

	return typ
}

func OtherCase(name string) string {
	this := ""

	if strings.Contains(name, "chasis") || strings.Contains(name, "ixia") || strings.Contains(name, "rack") {
		this = "dc.rack"
	} else if strings.Contains(name, "nexus") || strings.Contains(name, "sw") || strings.Contains(name, "switch") || strings.Contains(name, "n3k") {
		this = "network.switch"
	}

	return this
}

// Determine the specific type of a resource.
// nolint
func determineType(means string, resName string) string {
	name := lowerLetters(resName)
	typ := "N/A"

	if means == "Shelf" {
		typ = "dc.rack"
	} else if means == "Compute" {
		typ = ComputeCase(name)
	} else if means == "Other" {
		typ = OtherCase(name)
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

func SimpleOwnerFilter(resType string, ID string, db *sql.DB) (string, string) {
	var IPdata string
	var owner string

	if strings.Contains(resType, "compute") || resType == "network.switch" {
		IPdata = getIPDetaiLs(ID, db)

		owner = getUserDetails(IPdata, db)

		if IPdata == "" || net.ParseIP(IPdata) == nil {
			IPdata = "127.0.0.1"
		}
	} else {
		IPdata = "127.0.0.1"

		owner = "null"
	}

	return IPdata, owner
}

func checkEmpty(rt *Racktables) {
	if reflect.ValueOf(rt.AssetNo).IsZero() {
		rt.AssetNo = "NA"
	}

	if reflect.ValueOf(rt.Comments).IsZero() {
		rt.Comments = "NA"
	}

	if reflect.ValueOf(rt.ID).IsZero() {
		rt.ID = "NA"
	}

	if reflect.ValueOf(rt.Location).IsZero() {
		rt.Location = "NA"
	}

	if reflect.ValueOf(rt.Name).IsZero() {
		rt.Name = "N/A"
	}

	if reflect.ValueOf(rt.ObjtypeID).IsZero() {
		rt.ObjtypeID = "NA"
	}

	if reflect.ValueOf(rt.Owner).IsZero() {
		rt.Owner = "NA"
	}

	if reflect.ValueOf(rt.Port).IsZero() {
		rt.Port = 0
	}

	if reflect.ValueOf(rt.Problems).IsZero() {
		rt.Problems = "NA"
	}

	if reflect.ValueOf(rt.RackID).IsZero() {
		rt.RackID = "NA"
	}

	if reflect.ValueOf(rt.RowID).IsZero() {
		rt.RowID = "NA"
	}

	if reflect.ValueOf(rt.RowName).IsZero() {
		rt.RowName = "NA"
	}

	if reflect.ValueOf(rt.Type).IsZero() {
		rt.Type.Name = "NA"
		rt.Type.Description = "Not there ..."
	}

	if reflect.ValueOf(rt.IP).IsZero() {
		rt.IP = "1.1.1.1"
	}

	if reflect.ValueOf(rt.Label).IsZero() {
		rt.Label = "system.group"
	}

	if reflect.ValueOf(rt.Serial).IsZero() {
		rt.Serial = "NA"
	}
}

/* End of Data filtering for migration */

/* Beginning of data assignment for migration */

func Does() {
	var rt Racktables

	RackArr := []Racktables{}

	// Statement to query the db .
	statement := "SELECT rack_id, object_id  FROM rackspace"
	// to be filled in with appropriate user, password, and db name.
	db, err := sql.Open("mysql", "username:1234@/racktables")
	// if there is an error opening the connection, handle it
	if err != nil {
		log.Print(err.Error())
	}

	defer db.Close()

	// Execute the query
	results, err := db.Query(statement)
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&rt.RackID, &rt.ID)

		rt.ID, rt.Name, rt.Label, rt.ObjtypeID, rt.AssetNo, rt.Problems, rt.Comments = getMoreDetails(rt.ID, db)

		typeID := rt.ObjtypeID
		resType := determineIDMeaning(typeID, rt.Name)
		rt.Type.Name = resType
		rt.Type.Description = "migrated " + resType

		if reflect.ValueOf(rt.Type).IsZero() {
			rt.Type.Name = "NA"
			rt.Type.Description = "Not there ..."
		}

		if reflect.ValueOf(rt.Serial).IsZero() {
			rt.Serial = "NA"
		}

		rt.IP, rt.Owner = SimpleOwnerFilter(resType, rt.ID, db)

		if resType == "network.switch" {
			portInfo := getPortDetails(rt.ID, db)
			rt.Port = portInfo
		} else {
			portID := -1
			rt.Port = portID
		}

		rackID := getRackDetails(rt.ID, db)
		rowName, rowID, rowLocation := getRowDetails(rackID, db)

		rt.RowName = rowName
		rt.RowID = rowID
		rt.Location = rowLocation

		assetNumber := rt.AssetNo
		rt.AssetNo = assetNumber

		probs := rt.Problems
		rt.Problems = probs

		notes := rt.Comments
		rt.Comments = notes

		checkEmpty(&rt)

		RackArr = append(RackArr, rt)

		if err != nil {
			panic(err.Error())
		}
	}

	allData(RackArr)
}

func getToken() *http.Cookie {
	loginData := getAuth()

	fmt.Println("Got this from getAuth(): ", loginData)
	tokenData := &struct {
		JWT string `json:"jwt"`
	}{}

	err := json.Unmarshal(loginData, tokenData)
	if err != nil {
		fmt.Println("Got an error when unmarshaling login data")
	}

	cookie := &http.Cookie{
		Name:  "jwt",
		Value: tokenData.JWT,
	}

	return cookie
}

func getAuth() []byte {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	serverLoginUrl := "https://zebra.insieme.local:8000/login"

	migrationUser := []byte(`{
		"name":"admin",
		"password":"Riddikulus",
		"email":"admin@zebra.project-safari.io"
	}`)

	reader := bytes.NewReader(migrationUser)
	request, err := http.NewRequest("POST", serverLoginUrl, reader)
	if err != nil {
		fmt.Println("\nToken request got an error. This is it : ", err)
		os.Exit(1)
	}

	request.Close = true

	request.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(request)
	if err != nil {
		log.Print("\nToken resp got an error. This is it : ", err, "\n")
		os.Exit(1)
	}
	defer resp.Body.Close()

	loginBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		os.Exit(1)
	}

	fmt.Println("The response from the login is: ", resp.StatusCode)

	return loginBody
}

func PostIt() {
	Does()
	fmt.Println("\nDone with the migration!")
}

// Function that helps enforce minimum length for res. info.
func addChars(kind string, toAdd string) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	many := 0

	if strings.ToLower(kind) == "id" {
		many = 7 - len(toAdd) + 1
	}

	if strings.ToLower(kind) == "password" {
		many = 15 - len(toAdd) + 1
	}

	chars := make([]byte, many)
	for i := range chars {
		chars[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	extendedStr := toAdd + string(chars)

	return extendedStr
}

func allData(rackArr []Racktables) {
	serverUrl := "https://zebra.insieme.local:8000/api/v1/resources"
	token := getToken()
	var postIt []byte
	var res Racktables

	for i := 0; i < (len(rackArr)); i++ {
		res = rackArr[i]

		postIt = CreateResFromData(res)

		postable := bytes.NewReader(postIt)

		time.Sleep(20 * time.Millisecond)
		createRequests("POST", serverUrl, postable, token)
		time.Sleep(50 * time.Millisecond)
	}
}

func createRequests(method string, url string, body *bytes.Reader, token *http.Cookie) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("An error occured in the request. That is: %s", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Authorization", token.Value)

	req.AddCookie(token)

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("An error occurred in client. That's it: %s", err)
	} else if err == nil {

		_, err := ioutil.ReadAll(res.Body)
		fmt.Println("\n\nThis is the response itself: ", res)
		fmt.Println(res.StatusCode)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Get IPs from db based on type id.
func getIPDetaiLs(objectID string, db *sql.DB) string {
	var rt Racktables

	var IPnum string

	statement := "SELECT ip FROM IPv4Allocation WHERE object_id = ?"

	// Execute the query
	results, err := db.Query(statement, objectID)
	if err != nil {
		panic(err.Error())
	}

	defer results.Close()

	for results.Next() {
		// for each row, scan the result into our tag composite object
		err = results.Scan(&IPnum)

		rt.IP = IPnum

		if err != nil {
			panic(err.Error())
		}
	}

	return rt.IP
}

// Get port IDs from db based on type ID.
func getPortDetails(objectID string, db *sql.DB) int {
	var rt Racktables

	numPort := 0

	statement := "SELECT id FROM Port WHERE object_id = ?"

	// Execute the query
	results, err := db.Query(statement, objectID)
	if err != nil {
		panic(err.Error())
	}

	// defer results.Close()

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
func getMoreDetails(objectID string, db *sql.DB) (string, string, string, string, string, string, string) {
	var label sql.NullString

	var rt Racktables

	var assetNo sql.NullString

	var comment sql.NullString

	statement := "SELECT id, name, label, objtype_id, asset_no, has_problems, comment FROM rackobject WHERE id = ?"

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

func getRackDetails(objID string, db *sql.DB) string {
	var rt Racktables

	statement := "SELECT rack_id FROM RackSpace WHERE object_id = ?"

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
func getRowDetails(id string, db *sql.DB) (string, string, string) {
	statement := "SELECT row_id, row_name, location_name FROM Rack WHERE id = ?"

	var rt Racktables

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
func getUserDetails(resIP string, db *sql.DB) string {
	statement := "SELECT user FROM IPv4Log WHERE ip = ?"

	var rt Racktables

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

/* End of migration  */

/* Beginning of resource creation */

func CreateResFromData(res Racktables) []byte {
	resType := res.Type

	switch resType.Name {
	case "dc.dataceneter":
		theData := dcFiller(res)
		body, _ := json.Marshal(theData)

		return body

	case "dc.lab":
		theData := labFiller(res)
		body, _ := json.Marshal(theData)

		return body

	case "dc.rack", "dc.shelf":
		theData := rackFiller(res)
		body, _ := json.Marshal(theData)

		return body

	case "compute.server":
		theData := serverFiller(res)
		body, _ := json.Marshal(theData)

		return body

	case "compute.esx":
		theData := esxFiller(res)
		body, _ := json.Marshal(theData)

		return body

	case "compute.vm":
		theData := vmFiller(res)
		body, _ := json.Marshal(theData)

		return body

	case "compute.vceneter":
		theData := vcenterFiller(res)
		body, _ := json.Marshal(theData)

		return body

	case "network.switch":
		theData := switchFiller(res)
		body, _ := json.Marshal(theData)

		return body

	case "network.ipaddresspool":
		theData := addressPoolFiller(res)
		body, _ := json.Marshal(theData)

		return body

	case "network.vlanpool":
		theData := vlanFiller(res)
		body, _ := json.Marshal(theData)

		return body
	}

	return nil
}

/* End of resource creation */
