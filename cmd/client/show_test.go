package main //nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCmd = NewZebra() //nolint:gochecknoglobals

func TestNewZebraCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cmd := NewZebra()
	assert.NotNil(cmd)
}

func TestNewNetCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	netC := NewNetCmd(testCmd)
	assert.NotNil(netC)
}

func TestNewSrvCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	srvC := NewSrvCmd(testCmd)
	assert.NotNil(srvC)
}

func TestNewDcCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dcC := NewDCCmd(testCmd)
	assert.NotNil(dcC)
}

func TestShowServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"servers", "test-case"}

	srvCmd := NewSrvCmd(testCmd)
	rootCmd := New()
	rootCmd.AddCommand(srvCmd)
	serv := ShowServ(rootCmd, args)

	assert.NotNil(serv)
}

func TestShowVC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"vcenters", "test-case"}

	vc := NewSrvCmd(testCmd)
	rootCmd := New()

	rootCmd.AddCommand(vc)

	vcShow := ShowVC(rootCmd, args)
	assert.NotNil(vcShow)
}

func TestShowVlan(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"vlan", "test-case"}

	v := NewNetCmd(testCmd)
	rootCmd := New()

	rootCmd.AddCommand(v)

	vlan := ShowVlan(rootCmd, args)

	assert.NotNil(vlan)
}

func TestShowSw(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"switches", "test-case"}

	netCmd := NewNetCmd(testCmd)
	rootCmd := New()

	rootCmd.AddCommand(netCmd)
	sw := ShowSw(rootCmd, args)

	assert.NotNil(sw)
}

func TestShowRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"rack", "test-case"}

	rackCmd := NewDCCmd(testCmd)
	rootCmd := New()

	rootCmd.AddCommand(rackCmd)
	rack := ShowRack(rootCmd, args)

	assert.NotNil(rack)
}

func TestShowLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"lab", "test-case"}

	labCmd := NewDCCmd(testCmd)
	rootCmd := New()

	rootCmd.AddCommand(labCmd)

	lab := ShowLab(rootCmd, args)

	assert.NotNil(lab)
}

func TestShowESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"esv", "test-case"}

	esxCmd := NewDCCmd(testCmd)
	rootCmd := New()

	rootCmd.AddCommand(esxCmd)

	esv := ShowESX(rootCmd, args)

	assert.NotNil(esv)
}

func TestShowDC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"datacenter", "test-case"}

	dcCmd := NewDCCmd(testCmd)
	rootCmd := New()

	rootCmd.AddCommand(dcCmd)

	dc := ShowDC(rootCmd, args)

	assert.NotNil(dc)
}

func TestShowUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"user", "test-case"}

	rootCmd := New()
	rootCmd.AddCommand(testCmd)

	user := ShowUsr(rootCmd, args)

	assert.NotNil(user)
}
