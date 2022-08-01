package main //nolint:testpackage

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func test() *cobra.Command {
	testCmd := NewZebra()

	return testCmd
}

func TestNewZebraCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cmd := NewZebra()
	assert.NotNil(cmd)
}

func TestNewNetCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	netC := NewNetCmd(test())
	assert.NotNil(netC)
}

func TestNewSrvCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	srvC := NewSrvCmd(test())
	assert.NotNil(srvC)
}

func TestNewDcCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dcC := NewDCCmd(test())
	assert.NotNil(dcC)
}

func TestShowServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"servers", "test-case"}

	srvCmd := NewSrvCmd(test())
	rootCmd := New()
	rootCmd.AddCommand(srvCmd)
	serv := ShowServ(rootCmd, args)

	assert.NotNil(serv)
}

func TestShowVC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"vcenters", "test-case"}

	vc := NewSrvCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(vc)

	vcShow := ShowVC(rootCmd, args)
	assert.NotNil(vcShow)
}

func TestShowVlan(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"vlans", "test-case"}

	v := NewNetCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(v)

	vlan := ShowVlan(rootCmd, args)

	assert.NotNil(vlan)
}

//

func TestShowSw(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"switch", "test-case"}

	netCmd := NewNetCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(netCmd)
	sw := ShowSw(rootCmd, args)

	assert.NotNil(sw)
}

func TestShowRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"racks", "test-case"}

	rackCmd := NewDCCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(rackCmd)
	rack := ShowRack(rootCmd, args)

	assert.NotNil(rack)
}

func TestShowLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"labs", "test-case"}

	labCmd := NewDCCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(labCmd)

	lab := ShowLab(rootCmd, args)

	assert.NotNil(lab)
}

//

func TestShowESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"esx", "test-case"}

	esxCmd := NewSrvCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(esxCmd)

	esv := ShowESX(rootCmd, args)

	assert.NotNil(esv)
}

func TestShowDC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"dc", "test-case"}

	dcCmd := NewDCCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(dcCmd)

	dc := ShowDC(rootCmd, args)

	assert.NotNil(dc)
}

func TestShowUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"users", "test-case"}

	rootCmd := New()
	rootCmd.AddCommand(test())

	user := ShowUsr(rootCmd, args)

	assert.NotNil(user)
}

func TestShowReg(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"registrations", "test-case"}

	rootCmd := New()
	rootCmd.AddCommand(test())

	reg := ShowReg(rootCmd, args)

	assert.NotNil(reg)
}

func TestShowIP(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"ip", "test-case"}

	rootCmd := New()
	rootCmd.AddCommand(test())

	addr := ShowIP(rootCmd, args)

	assert.NotNil(addr)
}

//

func TestShowVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"vms", "test-case"}

	rootCmd := New()
	rootCmd.AddCommand(test())

	machine := ShowVM(rootCmd, args)

	assert.NotNil(machine)
}
