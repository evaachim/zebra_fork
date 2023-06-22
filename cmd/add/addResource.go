package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"project-safari/zebra/model/compute"
	"project-safari/zebra/model/dc"

	"github.com/project-safari/zebra/model/network"
	"github.com/project-safari/zebra/model/user"
)

type addRequest struct {
	TheType   string `json:"addType"`
	TheName   string `json:"addName"`
	TheOwner  string `json:"addOwner`
	TheStatus string `json:"addStatus"`
}

func assignaAddress() string {
	in := []string{"SJC17/LAB120", "SJC15/LAB157", "SJC14/LAB140"}
	randomIndex := rand.Intn(len(in))
	pickAddress := in[randomIndex]

	return pickAddress
}

func assignSerial() string {
	const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const numericBytes = "1234567890"

	chars := make([]byte, 2)
	for i := range chars {
		chars[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	nums := make([]byte, 3)
	for i := range chars {
		chars[i] = numericBytes[rand.Int63()%int64(len(numericBytes))]
	}

	serialStr := string(chars) + string(nums)

	return serialStr
}

func assignServerID() string {
	const serverBytes = "AB1CD2EF3GH4IJ5KL6MN7OP8QR9ST0UVWXYZ"

	chars := make([]byte, 5)
	for i := range chars {
		chars[i] = serverBytes[rand.Int63()%int64(len(serverBytes))]
	}

	serverID := string(chars)

	return serverID
}

func addPOST(method string, url string, body *bytes.Reader, token *http.Cookie) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Printf("An error occurred in the request. That is: %s", err)
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

func setResName(desiredName string) string {
	var resName string
	/*
		fmt.Println("Input the name to assign to this resource:")
		_, err := fmt.Scanln(&resName)
		if err != nil {
			log.Fatal(err)
		}
	*/
	resName = desiredName

	return resName
}

func esxFiller(userinfo user.User, req addRequest) {
	theESX := compute.NewESX(assignServerID(), req.TheName, req.TheOwner, "compute.esx")

	this, err := json.Marshal(theESX)
	if err != nil {
		fmt.Println("Encountered issues when unmarshaling POST data.")
	}

	addPOST("POST", "https://zebra.insieme.local:8000", bytes.NewReader(this), userinfo.Key.PublicKey())
}

func vcenterFiller(usernfo user.User, req addRequest) {
	resName := setResName("")

	user := usernfo.Role.Name

	theVCenter := compute.NewVCenter(resName, user, "compute.vcenter")

	this, err := json.Marshal(theVCenter)
	if err != nil {
		fmt.Println("Encountered issues when unmarshaling POST data.")
	}

	addPOST("POST", "https://zebra.insieme.local:8000", bytes.NewReader(this), userinfo.Key.PublicKey())
}

func serverFiller(userinfo user.User, req addRequest) {
	serverSerial := assignSerial()

	serverModel := "db-compute-server" + req.TheName

	theServer := compute.NewServer(serverModel, serverSerial, req.TheName, req.TheOwner, "compute.server")

	this, err := json.Marshal(theServer)
	if err != nil {
		fmt.Println("Encountered issues when unmarshaling POST data.")
	}

	addPOST("POST", "https://zebra.insieme.local:8000", bytes.NewReader(this), userinfo.Key.PublicKey())
}

func vmFiller(userinfo user.User, req addRequest) {
	esx := ""

	theVM := compute.NewVM(esx, req.TheName, req.TheOwner, "compute.vm")

	this, err := json.Marshal(theVM)

	if err != nil {
		fmt.Println("Encountered issues when unmarshaling POST data.")
	}

	addPOST("POST", "https://zebra.insieme.local:8000", bytes.NewReader(this), userinfo.Key.PublicKey())
}

func labFiller(userinfo user.User, req addRequest) {
	theLab := dc.NewLab(req.TheName, req.TheOwner, "dc.lab")

	this, err := json.Marshal(theLab)

	if err != nil {
		fmt.Println("Encountered issues when unmarshaling POST data.")
	}

	addPOST("POST", "https://zebra.insieme.local:8000", bytes.NewReader(this), userinfo.Key.PublicKey())
}

func dcFiller(userinfo user.User, req addRequest) {
	address := assignaAddress()

	group := "dc.datacenter"

	theDC := dc.NewDatacenter(address, req.TheName, req.TheOwner, group)

	this, err := json.Marshal(theDC)

	if err != nil {
		fmt.Println("Encountered issues when unmarshaling POST data.")
	}

	addPOST("POST", "https://zebra.insieme.local:8000", bytes.NewReader(this), userinfo.Key.PublicKey())
}

func switchFiller(userinfo user.User, req addRequest) {
	theSwitch := network.NewSwitch(req.TheName, req.TheOwner, "network.switch")

	this, err := json.Marshal(theSwitch)

	if err != nil {
		fmt.Println("Encountered issues when unmarshaling POST data.")
	}

	addPOST("POST", "https://zebra.insieme.local:8000", bytes.NewReader(this), userinfo.Key.PublicKey())
}

func addressPoolFiller(userinfo user.User, req addRequest) {
	theIPpool := network.NewIPAddressPool(req.TheName, req.TheOwner, "network.ipAddressPool")

	this, err := json.Marshal(theIPpool)

	if err != nil {
		fmt.Println("Encountered issues when unmarshaling POST data.")
	}

	addPOST("POST", "https://zebra.insieme.local:8000", bytes.NewReader(this), userinfo.Key.PublicKey())
}

func vlanFiller(userinfo user.User, req addRequest) {
	theVLAN := network.NewVLANPool(req.TheName, req.TheOwner, "network.vlanPool")

	this, err := json.Marshal(theVLAN)

	if err != nil {
		fmt.Println("Encountered issues when unmarshaling POST data.")
	}

	addPOST("POST", "https://zebra.insieme.local:8000", bytes.NewReader(this), userinfo.Key.PublicKey())
}
