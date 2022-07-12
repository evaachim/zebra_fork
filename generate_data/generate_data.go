
/*

create 100 of each resource for some users
display results for each respective type

*/

package generate_data

import (
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/network"
)

var resourceTypes = []string{"VLANPool", "Switch", "IPAddressPool", "Datacenter", "Lab", "Rack", "Server", "ESX", "VM", " "}
var SampleIpAddr = []string{"192.332.11.05", "192.232.11.37", "192.232.22.05", "192.225.11.05", "192.0.0.0", "192.192.192.192", "225.225.225.225", "192.192.64.08"}

// create user roles

func User() string {

	// some usernames
	name_list := []string{
		"user1",
		"user2",
		"user3",
		"manager",
		"admin",
		"guest",
		"another-user",
		" ",
		"engineer1",
		"designer",
		"private_user",
		"ceo",
		"staff",
	}

	username := RandData(name_list)

	return username

}

//create some passwords

func Password() string {
	p_list := []string{
		"pass123",
		"pass222",
		"pass321",
		"m3g2r1",
		"passthrough",
		" ",
		"another",
		"no_pass",
		"somepass",
		"v)entry2",
		"p%3eqr3",
		"random_pass",
		"validate000",
	}
	pwd := RandData(p_list)

	return pwd
}

//random selection from lists

func RandData(res []string) string {
	rand.Seed(time.Now().UnixNano())
	var length = len(res)
	var ind int = rand.Intn(length - 1)
	typ := res[ind]
	return typ
}

// creating random starting and end points

func Range() uint16 {

	nums := []uint16{0, 11, 121, 44, 32, 234, 9, 64, 33, 2, 5, 8, 16, 200, 14, 77}
	rand.Seed(time.Now().UnixNano())
	var length = len(nums)
	var ind int = rand.Intn(length - 1)
	num := nums[ind]
	return num
}

// creating random ports

func Ports() uint32 {
	nums := []uint32{1, 2, 3, 4, 6, 8, 9, 16, 27, 32, 36, 54, 72, 64, 81, 128, 162, 216, 256, 512}
	rand.Seed(time.Now().UnixNano())
	var length = len(nums)
	var ind int = rand.Intn(length - 1)
	port := nums[ind]
	return port

}

// create random sample model names

func Models() string {
	models := []string{"model A", "modelB", "modelC", "modelD", "modelA.1", "modelA.1.1", "modelE", "modelD.010", "modelF", "modelG", "modelG.1", "modelG.1.1"}
	model := RandData(models)
	return model
}

// create random serial codes

func Serials() string {
	nums := []string{"00000", "00001", "00002", "00003", "00004", "00005", "00006", "00007", "00008", "00009", "00010", "00020", "00030", "00040", "00050", "00060", "00070", "00080", "00090", "00100", "00200", "00300", "00400", "00500", "01000", "02000", "03000"}
	return RandData(nums)
}

//creating various resource types

func CreateVlanPool(theType string) (*network.VLANPool, zebra.NamedResource) {
	var start uint16 = Range()
	theLabels := CreateLabels()

	theRes := zebra.NewBaseResource(theType, theLabels)

	ret := &network.VLANPool{

		BaseResource: theRes,
		RangeStart:   start,
		RangeEnd:     (start + 100),
	}
	return ret, theRes
}

func CreateSwitch(theType string, ip net.IP) (*network.Switch, zebra.NamedResource) {
	var serial string = Serials()
	var model string = Models()
	var ports uint32 = Ports()
	theLabels := CreateLabels()
	theRes := zebra.NewBaseResource(theType, theLabels)

	ret := &network.Switch{

		BaseResource: theRes,
		ManagementIP: ip,
		SerialNumber: serial,
		Model:        model,
		NumPorts:     ports,
	}
	return ret, theRes
}

func CreateIpAddressPool(theType string, ip net.IP) (*network.IPAddressPool, zebra.NamedResource) {

	theLabels := CreateLabels()
	theRes := zebra.NewBaseResource(theType, theLabels)

	ret := &network.IPAddressPool{
		BaseResource: theRes,
		//Subnets:      []net.IPNet,
	}
	return ret, theRes
}

func CreateDatacenter() (*dc.Datacenter, zebra.NamedResource) {
	namedRes := new(zebra.NamedResource)

	ret := &dc.Datacenter{
		NamedResource: *namedRes,
		Address:       "sample address",
	}

	return ret, *namedRes
}

func CreateLab() (*dc.Lab, zebra.NamedResource) {
	namedRes := new(zebra.NamedResource)

	ret := &dc.Lab{
		NamedResource: *namedRes,
	}
	return ret, *namedRes
}

func CreateRack() (*dc.Rack, zebra.NamedResource) {

	namedRes := new(zebra.NamedResource)

	ret := &dc.Rack{
		NamedResource: *namedRes,
		Row:           "sample row",
	}
	return ret, *namedRes
}

//sample labels

func CreateLabels() map[string]string {
	codes := make(map[string]string)
	col := " "
	let := " "
	colors := []string{"red", "yellow", "green", "blue", "white", "magenta", "black", "purple", "brown", "orange", "pink", "grey"}
	letters := []string{"alpha", "beta", "gamma", "delta", "epsilon", "eta", "theta", "Iota", "Kappa", "Lambda", "Mu", "Nu"}

	col = RandData(colors)
	let = RandData(letters)
	codes[let] = col

	// get the key

	return codes
}

//put it all together

func Generate_Data() {

	// go through each resource type

	for i := 0; i < len(resourceTypes); i++ {

		theType := resourceTypes[i]

		// 100 resources of each type

		for each := 0; each < 100; each++ {

			creds := new(auth.User)

			sampleIP := RandData(SampleIpAddr)

			creds.PasswordHash = Password()

			creds.Role.Name = User()

			if theType == "VLANPool" {

				Res, base := CreateVlanPool(theType)
				creds.NamedResource = base
				fmt.Println("Information: ", each, "\nThe data with user info and named resource ", creds, "\nThe data with complete resource info: ", Res)

			} else if theType == "Switch" {

				Res, base := CreateSwitch(theType, net.IP(sampleIP))
				creds.NamedResource = base
				fmt.Println("Information: ", each, "\nThe data with user info and named resource ", creds, "\nThe data with complete resource info: ", Res)

			} else if theType == "IPAddressPool" {

				Res, base := CreateIpAddressPool(theType, net.IP(sampleIP))
				creds.NamedResource = base
				fmt.Println("Information: ", each, "\nThe data with user info and named resource ", creds, "\nThe data with complete resource info: ", Res)

			} else if theType == "Datacenter" {

				Res, base := CreateDatacenter()
				creds.NamedResource = base
				fmt.Println("Information: ", each, "\nThe data with user info and named resource ", creds, "\nThe data with complete resource info: ", Res)

			} else if theType == "Lab" {

				Res, base := CreateLab()
				creds.NamedResource = base
				fmt.Println("Information: ", each, "\nThe data with user info and named resource ", creds, "\nThe data with complete resource info: ", Res)

			} else if theType == "Rack" {

				Res, base := CreateRack()
				creds.NamedResource = base
				fmt.Println("Information: ", each, "\nThe data with user info and named resource ", creds, "\nThe data with complete resource info: ", Res)

			}

		}

	}

}
