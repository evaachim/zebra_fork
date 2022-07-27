package show_test

import (
	"testing"

	"github.com/project-safari/zebra/cli/show"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var testCmd = &cobra.Command{ //nolint:exhaustruct,exhaustivestruct
	Use:     "show resources",
	Short:   "command to show zebra resources",
	Version: "test-version",
}

func TestNewZebraCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cmd := show.NewZebra()
	assert.NotNil(cmd)
}

func TestNew(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cmd := show.New()
	assert.NotNil(cmd)
}

func TestNewNetCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	netC := show.NewNetCmd()
	assert.NotNil(netC)
}

func TestNewSrvCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	srvC := show.NewSrvCmd()
	assert.NotNil(srvC)
}

func TestNewDcCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dcC := show.NewDCCmd()
	assert.NotNil(dcC)
}

func TestShowSw(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"switch", "test-case"}

	sw := show.ShowSw(testCmd, args)

	assert.Nil(sw)
}

func TestShowServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"server", "test-case"}

	serv := show.ShowServ(testCmd, args)

	assert.Nil(serv)
}

func TestShowVC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"vcenters", "test-case"}

	vc := show.ShowVC(testCmd, args)

	assert.Nil(vc)
}

func TestShowVlan(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"vlan", "test-case"}

	vlan := show.ShowVlan(testCmd, args)

	assert.Nil(vlan)
}

func TestShowRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"rack", "test-case"}

	rack := show.ShowRack(testCmd, args)

	assert.Nil(rack)
}

func TestShowLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"lab", "test-case"}

	lab := show.ShowLab(testCmd, args)

	assert.Nil(lab)
}

func TestShowESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"esv", "test-case"}

	esv := show.ShowESX(testCmd, args)

	assert.Nil(esv)
}

func TestShowDC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"datacenter", "test-case"}

	dc := show.ShowDC(testCmd, args)

	assert.Nil(dc)
}
