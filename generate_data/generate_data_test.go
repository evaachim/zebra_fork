/*
create 100 instances of each resource for some users
program to execute tests
*/

package dazzle_test

import (
	"net"
	"testing"

	dazzle "github.com/project-safari/zebra/generate_data"
	"github.com/stretchr/testify/assert"
)

func TestSetTypes(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	types := dazzle.AllResourceTypes()
	assert.NotNil(types)
}

func TestSetIP(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	samples := dazzle.IPsamples()
	assert.NotNil(samples)
}

func TestUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	user := dazzle.User()
	assert.NotNil(user)
}

func TestPass(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	pwd := dazzle.Password()
	assert.NotNil(pwd)
}

func TestName(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	name := dazzle.Name()

	assert.NotNil(name)
}

func TestRange(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	num := dazzle.Range()

	assert.NotNil(num)
}

func TestPorts(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	prt := dazzle.Ports()

	assert.NotNil(prt)
}

func TestModel(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	model := dazzle.Models()

	assert.NotNil(model)

	assert.NotEqual(model, " ")
}

func TestSerials(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ser := dazzle.Serials()

	assert.NotNil(ser)
	assert.NotEqual(ser, " ")
}

func TestRows(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	row := dazzle.Rows()

	assert.NotNil(row)
	assert.NotEqual(row, " ")
}

func TestAddresses(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	adr := dazzle.Addresses()

	assert.NotNil(adr)
	assert.NotEqual(adr, " ")
}

func TestOrder(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var a uint16 = 20

	var b uint16 = 5

	one, two := dazzle.Order(a, b)
	assert.True(one < two)
}

func TestCreateIPArr(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	IParr := dazzle.CreateIPArr(2)
	assert.NotEmpty(IParr)
}

func TestCreateVlanPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	VlanPool := dazzle.NewVlanPool("VlanPool")

	assert.NotNil(VlanPool)
	assert.NotEmpty(VlanPool)
}

func TestCreateVcenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	Vcenter := dazzle.NewVCenter("VlanPool", net.IP("192.222.004"))

	assert.NotNil(Vcenter)
}

func TestCreateSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	Switch := dazzle.NewSwitch("Switch", net.IP("192.222.004"))

	assert.NotNil(Switch)
}

func TestCreateIPAddressPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	IPs := dazzle.NewIPAddressPool("IPAddressPool", dazzle.CreateIPArr(2))

	assert.NotNil(IPs)

	assert.NotEmpty(IPs)
}

func TestCreateDatacenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	DataCenter := dazzle.NewDatacenter("Datacenter")

	assert.NotNil(DataCenter)
}

func TestCreateLabels(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := dazzle.CreateLabels()

	assert.NotNil(labels)
}

func TestCreateLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	Lab := dazzle.NewLab("Lab")

	assert.NotNil(Lab)
}

func TestCreateRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	Rack := dazzle.NewRack("Rack")

	assert.NotNil(Rack)

	assert.NotNil(Rack.BaseResource)
}

func TestIsGood(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	result := dazzle.IsGood(100)

	assert.NotNil(result)
	assert.False(result)

	errRes := dazzle.IsGood(0)
	assert.True(errRes)
}

func TestGeneration(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	num := 100
	result := dazzle.IsGood(num)

	assert.NotNil(dazzle.GenerateData(result, num))
	user, arr := dazzle.GenerateData(result, num)

	assert.NotNil(user)

	assert.NotEmpty(arr)
}
