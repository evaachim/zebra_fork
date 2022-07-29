package display_test

import (
	"testing"

	"github.com/project-safari/zebra/cli/display"
	"github.com/stretchr/testify/assert"
)

var testCmd = display.NewZebra() //gochecknoglobals

/*
&cobra.Command{ //nolint:exhaustruct,exhaustivestruct,gochecknoglobals
	Use:     "show resources",
	Short:   "command to show zebra resources",
	Version: "test-version",
}
*/

func TestNewZebraCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cmd := display.NewZebra()
	assert.NotNil(cmd)
}

func TestNewNetCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	netC := display.NewNetCmd()
	assert.NotNil(netC)
}

func TestNewSrvCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	srvC := display.NewSrvCmd()
	assert.NotNil(srvC)
}

func TestNewDcCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dcC := display.NewDCCmd()
	assert.NotNil(dcC)
}

func TestShowSw(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"switch", "test-case"}

	sw := display.ShowSw(testCmd, args)

	assert.Nil(sw)
}

func TestShowServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"server", "test-case"}

	serv := display.ShowServ(testCmd, args)

	assert.Nil(serv)
}

func TestShowVC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"vcenters", "test-case"}

	vc := display.ShowVC(testCmd, args)

	assert.Nil(vc)
}

func TestShowVlan(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"vlan", "test-case"}

	vlan := display.ShowVlan(testCmd, args)

	assert.Nil(vlan)
}

func TestShowRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"rack", "test-case"}

	rack := display.ShowRack(testCmd, args)

	assert.Nil(rack)
}

func TestShowLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"lab", "test-case"}

	lab := display.ShowLab(testCmd, args)

	assert.Nil(lab)
}

func TestShowESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"esv", "test-case"}

	esv := display.ShowESX(testCmd, args)

	assert.Nil(esv)
}

func TestShowDC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"datacenter", "test-case"}

	dc := display.ShowDC(testCmd, args)

	assert.Nil(dc)
}

/*
func TestShowUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"user", "test-case"}

	user := display.ShowUsr(testCmd, args)

	assert.Nil(user)
}
*/
