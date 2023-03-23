package migration //nolint:testpackage

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/project-safari/zebra/cmd/script"
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
	testEmpty1, _, testEmpty2 := CreateResFromData(rt)

	assert.Nil(testEmpty1)

	assert.Equal(testEmpty2, "")

	rt.AssetNo = "1"
	rt.ID = "123"
	rt.IP = "1.1.1.1"
	rt.Name = "test-switch"
	rt.ObjtypeID = "8"

	// test for creating a switch
	rt.Type = "network.switch"
	testCreateSwitch, _, _ := CreateResFromData(rt)
	assert.NotNil(testCreateSwitch)

	// test for creating a dc
	rt.Type = "dc.dataceneter"
	testCreateDC, _, _ := CreateResFromData(rt)
	assert.NotNil(testCreateDC)

	// test for creating a lab
	rt.Type = "dc.lab"
	testCreateLab, _, _ := CreateResFromData(rt)
	assert.NotNil(testCreateLab)

	// test for creating a rack with shelf type
	rt.Type = "dc.shelf"
	testCreateShelf, _, _ := CreateResFromData(rt)
	assert.NotNil(testCreateShelf)

	// test for creating a vm
	rt.Type = "compute.vm"
	testCreateVM, _, _ := CreateResFromData(rt)
	assert.NotNil(testCreateVM)

	// test for creating a vc
	rt.Type = "compute.vceneter"
	testCreateVC, _, _ := CreateResFromData(rt)
	assert.NotNil(testCreateVC)

	// test for creating a server
	rt.Type = "compute.server"
	testCreateSrv, _, _ := CreateResFromData(rt)
	assert.NotNil(testCreateSrv)

	// test for creating an esx server
	rt.Type = "compute.esx"
	testCreateESX, _, _ := CreateResFromData(rt)
	assert.NotNil(testCreateESX)

	// test for creating a IPAddressPool
	rt.Type = "network.ipaddresspool"
	testCreateIP, _, _ := CreateResFromData(rt)
	assert.NotNil(testCreateIP)

	// test for creating a vlanPool
	rt.Type = "network.vlanpool"
	testCreateVP, _, _ := CreateResFromData(rt)
	assert.NotNil(testCreateVP)
}

func TestFiller(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	var rt Racktables

	rt.AssetNo = "1"
	rt.ID = "123"
	rt.IP = "1.1.1.1"
	rt.Name = "test-switch"
	rt.ObjtypeID = "8"

	rt.Type = "network.switch"
	testSwitchFiller := switchFiller(rt)
	assert.NotNil(testSwitchFiller)

	rt.Name = "test-server"
	rt.ObjtypeID = "8"

	rt.Type = "compute.server"

	testServerFiller := serverFiller(rt)
	assert.NotNil(testServerFiller)

	rt.Name = "test-esx"
	rt.ObjtypeID = "9"

	rt.Type = "compute.esx"

	testESXfiller := esxFiller(rt)
	assert.NotNil(testESXfiller)

	rt.Name = "test-vc"
	rt.ObjtypeID = "9"

	rt.Type = "compute.vcenter"

	testVCfiller := vcenterFiller(rt)
	assert.NotNil(testVCfiller)

	rt.Name = "test-vm"
	rt.ObjtypeID = "9"

	rt.Type = "compute.vm"

	testVMfiller := vmFiller(rt)
	assert.NotNil(testVMfiller)

	rt.Name = "test-vlan"
	rt.ObjtypeID = "10"

	rt.Type = "network.vlan"

	testVLANfiller := vlanFiller(rt)
	assert.NotNil(testVLANfiller)
}

var errFake = errors.New("fake error")

type fakeReader struct {
	err bool
}

func (f fakeReader) Read(b []byte) (int, error) {
	if f.err {
		return 0, errFake
	}

	return 0, io.EOF
}

func TestReadJSON(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	req := script.MakeLabelRequest(assert, nil, "a", "b", "c")

	labelReq := &struct {
		Labels []string `json:"labels"`
	}{Labels: []string{}}

	assert.Nil(script.ReadJSON(context.Background(), req, labelReq))

	// Bad IO reader
	req.Body = ioutil.NopCloser(fakeReader{err: true})
	assert.NotNil(script.ReadJSON(context.Background(), req, nil))

	// Empty Body
	req.Body = ioutil.NopCloser(fakeReader{err: false})
	assert.NotNil(script.ReadJSON(context.Background(), req, nil))
}
