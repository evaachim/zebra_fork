package main //nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineType(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// compute category
	means := "Compute"
	resName := "esxServer"

	result := determineType(means, resName)
	assert.Equal(result, "compute.esx")

	resName = "JENKINS"
	result = determineType(means, resName)
	assert.Equal(result, "compute.server")

	resName = "BLD123"
	result = determineType(means, resName)
	assert.Equal(result, "dc.datacenter")

	resName = "VLAN"
	result = determineType(means, resName)
	assert.Equal(result, "network.vlanPool")

	resName = "switchA"
	result = determineType(means, resName)
	assert.Equal(result, "network.switch")

	resName = "capic-1"
	result = determineType(means, resName)
	assert.Equal(result, "compute.vm")

	resName = "xYvapic/122"
	result = determineType(means, resName)
	assert.Equal(result, "compute.vcenter")

	resName = "Ipc"
	result = determineType(means, resName)
	assert.Equal(result, "network.ipAddressPool")

	// larger other category
	means = "Other"
	resName = "ixia"

	result = determineType(means, resName)
	assert.Equal(result, "dc.rack")

	resName = "nexus"

	result = determineType(means, resName)
	assert.Equal(result, "network.switch")

	// no category
	means = ""
	resName = ""

	result = determineType(means, resName)
	assert.Equal(result, "")
}

//nolint:funlen
func TestDetermineIDMeaning(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// test for vm.
	id := "2"
	name := "VM"

	result := determineIDMeaning(id, name)
	assert.Equal(result, "compute.vm")

	// test for rack with name shelf.
	id = "30"
	name = "Shelf"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "dc.rack")

	// test for rack with name rack.
	name = "Rack"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "dc.rack")

	// test for vc.
	id = "38"
	name = "VC"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "compute.vcenter")

	// test for server.
	id = "4"
	name = "server"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "compute.server")

	// test for sw.
	id = "8"
	name = "sw"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "network.switch")

	// tests for compute's id.
	id = "1504"
	name = "sw"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "network.switch")

	id = "1504"
	name = "/"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "")

	// test for other's id.
	id = "1503"
	name = "chasis"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "dc.rack")

	// test for wrong id.
	id = "0"
	name = "chasis"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "unclassified")
}

/*
// Test for reading json.

	func TestRead(t *testing.T) {
		t.Parallel()

		assert := assert.New(t)
		req := makeLabelRequest(assert, nil, "a", "b", "c")

		labelReq := &struct {
			Labels []string `json:"labels"`
		}{Labels: []string{}}

		assert.Nil(readJSONdata(context.Background(), req, labelReq))

		// Bad IO reader
		req.Body = ioutil.NopCloser(mockReader{err: true})
		assert.NotNil(readJSONdata(context.Background(), req, nil))

		// Empty Body
		req.Body = ioutil.NopCloser(mockReader{err: false})
		assert.NotNil(readJSONdata(context.Background(), req, nil))
	}


type mockReader struct {
	err bool
}

func (f mockReader) Read(b []byte) (int, error) {
	if f.err {
		return 0, errors.New("mock error") //nolint:goerr113
	}

	return 0, io.EOF
}

func makeLabelRequest(assert *assert.Assertions, resources *ResourceAPI, labels ...string) *http.Request {
	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)
	ctx = context.WithValue(ctx, AuthCtxKey, "test-label")

	req, err := http.NewRequestWithContext(ctx, "GET", "/api/v1/labels", nil)
	assert.Nil(err)
	assert.NotNil(req)

	v := map[string][]string{"labels": labels}
	b, e := json.Marshal(v)
	assert.Nil(e)

	req.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	return req
}
*/

func TestAllData(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var rack Racktables

	rackArr := []Racktables{}

	rack.ID = "123"
	rack.Name = "test-rack"

	rackArr = append(rackArr, rack)

	assert.NotNil((rackArr))
}

//nolint:funlen
func TestCreateRes(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var rt Racktables

	// test for creating an empty resource
	testEmpty1, testEmpty2 := createResFromData(rt)

	assert.Nil(testEmpty1)

	assert.Equal(testEmpty2, "")

	rt.AssetNo = "1"
	rt.ID = "123"
	rt.IP = "1.1.1.1"
	rt.Name = "test-switch"
	rt.ObjtypeID = "8"

	// test for creating a switch
	rt.Type = "network.switch"
	testCreateSwitch, _ := createResFromData(rt)
	assert.NotNil(testCreateSwitch)

	// test for creating a dc
	rt.Type = "dc.datacenetr"
	testCreateDC, _ := createResFromData(rt)
	assert.NotNil(testCreateDC)

	// test for creating a lab
	rt.Type = "dc.lab"
	testCreateLab, _ := createResFromData(rt)
	assert.NotNil(testCreateLab)

	// test for creating a rack with shelf type
	rt.Type = "dc.shelf"
	testCreateShelf, _ := createResFromData(rt)
	assert.NotNil(testCreateShelf)

	// test for creating a vm
	rt.Type = "compute.vm"
	testCreateVM, _ := createResFromData(rt)
	assert.NotNil(testCreateVM)

	// test for creating a vc
	rt.Type = "compute.vcenetr"
	testCreateVC, _ := createResFromData(rt)
	assert.NotNil(testCreateVC)

	// test for creating a server
	rt.Type = "compute.server"
	testCreateSrv, _ := createResFromData(rt)
	assert.NotNil(testCreateSrv)

	// test for creating an esx server
	rt.Type = "compute.esx"
	testCreateESX, _ := createResFromData(rt)
	assert.NotNil(testCreateESX)

	// test for creating a IPAddressPool
	rt.Type = "network.ipaddresspool"
	testCreateIP, _ := createResFromData(rt)
	assert.NotNil(testCreateIP)

	// test for creating a vlanPool
	rt.Type = "network.vlanpool"
	testCreateVP, _ := createResFromData(rt)
	assert.NotNil(testCreateVP)
}
