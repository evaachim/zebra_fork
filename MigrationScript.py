## Preliminary script to get data from the racktables database into the Zebra tool.
##
## This could be extended for any type of resource for which there is data in the racktables database.

## Are we updating the zebra tool to have this data statically and then access it via the api calls? 
## Or we using the api calls to fetch from the db only that data which we need for the specific call?

## imports needed for access to mysql and the racktables database.
import os
import re
import mysql.connector as database

import requests
import mysql
import sys 
import pyodbc

## The local user to connect to the database.
username = os.environ.get("username")
password = os.environ.get("password")

## API routes.
types = "/api/v1/types"
labels = "/api/v1/labels"
query = "/api/v1/resources"
posts = "/api/v1/resources"
delete = "/api/v1/resources"

'''
## To be used if second implementation option is selected 

## get the response for the query.
queryResponse = requests.get(query)

## Make the response into json format.
criteria = queryResponse.json()

print(criteria)
'''

def getQuery():
    queryResponse = requests.get(query)
    ## This contains the id to look for when querring.
    ## Querry can be of different types - by id, by type, by label. These are in api.go.
    criteriaID = queryResponse.json()["id"]
    criteriaType = queryResponse.json()["type"]
    criteriaLabel = queryResponse.json()["label"]

    key = " "

    if len(criteriaID) != 0:
        key = criteriaID
    elif len(criteriaType) !=0:
        key = criteriaType
    elif len(criteriaLabel) != 0:
        key = criteriaLabel
    else: key = ""

    ## Return and use.
    return key


## Make the connection to the data base.
connection = database.connect(
    user = username,
    password = password,
    host = "localhost",
    database = "racktables"
)

## This is where the database cursor goes.
cursor = connection.cursor()

## Get data from database.

## Get all data for datacenter resources.
## Key must be the one the user passes in when the query is made. 
## The key comes form getQuery currently, but the format can change if necesary.
## Call this in the datacener area.
def get_db_resource_data(key):
    dcName, row, location = ""

    try:
        statement = "SELECT name, row_name, row_id, location_name FROM rackobject WHERE id=%s"
        ## selection is now by the object's id, it could be something else, as needed.
        data = (key,)
        cursor.execute(statement, data)
        for (name, row_name, row_id, location_name) in cursor:
            print("Retreived the data")
            
            dcName = name
            row = row_name
            rowID = row_id
            location = location_name

    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")
    return (dcName, row, rowID, location)

## Get data for network resources.
## get_net_resource_data() should go here.
## Should call this in the network area.


## function goes here and returns a tuple.

## Get data for compute resources.
## get_compute_resource_data() should go here.
## Should call this in the compute area.
## We don't have all of this data in the db, but it should be something like this:
def get_compute_resource_data(key):
    srvName, row, location, resType = ""
    try:
        ## rackobject is currently the only table containing such info, but it is not complete for our server needs.
        ## the db does not have complete data for various server types.
        statement = "SELECT name, serial, model, ip, location_name FROM <<some table name>> WHERE id=%s"
        ## selection is now by the object's id, it could be something else, as needed.
        data = (key,)
        cursor.execute(statement, data)
        for (name, serial, model, ip, location_name) in cursor:
            print("Retreived the data")
            srvName = name
            ## We don't have these in a server table.
          ## row = row_name
          ##  rowID = row_id
            location = location_name
            resType = determineType(name)
    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")
    return (srvName, serial, model, ip, location, resType)

    ## Get data by ID from dictionary.

def get_resource_data(key):
    srvName, row, location, resType = ""
    try:
        statement = "SELECT name, serial, model, ip, location_name FROM rackobject WHERE id=%s"

        data = (key,)
        cursor.execute(statement, data)

        for (name, serial, model, ip, location_name) in cursor:
            print("Retreived the data")
            srvName = name
            ## We don't have these in a server table.
          ## row = row_name
          ##  rowID = row_id
            location = location_name
            resType = determineType(name)
    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")
    return (srvName, serial, model, ip, location, resType)

## Get the type by id.
## pdu and ups can be considered as vm solutions so they're both added to vm.
## some are randomly dropped into 1504 but this objtype_id has no reference. Thus, name will be used too.
def determineIDMeaning(id, name):
    means, final, this = ""
    if id == "2" or id == "12":
        means = "VM"
    elif id == "3":
        means = "Shelf"
    elif id == "4":
        means = "Server"
    elif id == "8":
        means = "Switch"
    elif id == "1504":
        means = "Compute"
    elif id == "1503":
        means = "Other"
    else: means = "/"

    final = determineType(means, name)

    if final == "/":
        final = "unclassified"

    this = final

    return this

def determineType(means, name):
    name = name.lower()
    type = ""

    ### still need something for vapic* and vpod*, as well as FRODO*,  APPLIANCE-HOME1, CAPIC*, 
        # aci-github.cisco.com*, DMASHAL-VINTELLA*.
        # RESOLVED AS EXPLAINED BELOW:
    ### vpod uses VMware ESXi hosts, VMware vCenter, storage, networking and a Windows Console VM. 
            # => vcenter.
    ### vAPIC virtual machines use VMware vCenter ==> vcenter
        # Cisco ACI vCenter plug-in.
        # BUT also uses Cisco ACI Virtual Edge VM.
    ### Cisco Cloud APIC on Microsoft Azure is deployed and runs as an 
        # Microsoft Azure Virtual Machine => capic => VM.
    ### Frodo is enabled by default on VMs powered on after AOS 5.5.X => frodo => VM.
        # About frodo - VMware Technology Network VMTN
    ### Vintela -> VAS is Vintela's flagship product in a line that includes Vintela Management 
        # eXtensions (VMX), which extends Microsoft Systems Management Server => server.
    ### apic uses controllers and so does cisco aci but it is similar to switches => switch.
    
    if means == "Shelf":
        type = "Rack"
    elif means == "Compute":
        if "esx" in name:
            type = "ESX"
        elif "vm" in name:
            type = "vm"
        elif "jenkins" in name or "server" in name or "srv" in name or "vintella" in name:
            type = "Server"
        elif "bld" in name or "datacenter" in name:
            type = "Datacenter"
        elif "dmz" in name or "vlan" in name or "asa" in name or "bridge" in name:
            type = "VLAN"
        elif "vleaf" in name or "switch" in name or "sw" in name or "aci" in name:
            type = "Switch"
        elif "vm" in name or "capic" in name or "frodo" in name:
            type = "vm"
        elif "vapic" in name or "vpod" in name:
            type = "VCenter"
        elif "IPC" in name:
            type = "IPAddressPool"
    elif means == "Other":
        if "chasis" in name or "ixia" in name:
            type = "Rack"
        elif "nexus" in name or "switch" in name or "sw" in name or "n3k" in name:
            type = "Switch"
    else: type = means

    return type


    '''
    Reference model:
    * means applicable to zebra's current system.

    1 => array ('chapter_id' => 1, 'dict_value' => 'BlackBox'),
	2 => array ('chapter_id' => 1, 'dict_value' => 'PDU'),   *
	3 => array ('chapter_id' => 1, 'dict_value' => 'Shelf'), *
	4 => array ('chapter_id' => 1, 'dict_value' => 'Server'), *
	5 => array ('chapter_id' => 1, 'dict_value' => 'DiskArray'),
	6 => array ('chapter_id' => 1, 'dict_value' => 'TapeLibrary'),
	7 => array ('chapter_id' => 1, 'dict_value' => 'Router'),
	8 => array ('chapter_id' => 1, 'dict_value' => 'Network switch'), *
	9 => array ('chapter_id' => 1, 'dict_value' => 'PatchPanel'),
	10 => array ('chapter_id' => 1, 'dict_value' => 'CableOrganizer'),
	11 => array ('chapter_id' => 1, 'dict_value' => 'spacer'),
	12 => array ('chapter_id' => 1, 'dict_value' => 'UPS'), *
    '''

def getData():
    specificID, resName, theLabel, assets, problems, notes, resType, rackID = ""
    rowID, rowName, rackLocation = ""
    ip = "N/A"
    ## portName = "N/A"
    portID = 0
    
    try:
        statement = "SELECT id, name, label, objtype_id, asset_no, has_problems, comment FROM rackobject%s"

        ## data = (key,)
        cursor.execute(statement)

        for (id, name, label, objtype_id, asset_no, has_problems, comment) in cursor:
            print("Retreived the data")

            # resource specific data.
            specificID  = id
            resName = name

            # currently null for all, will be usefoul for containment model.
            theLabel = label

            # additional information.
            assets = asset_no
            problems = has_problems
            notes = comment

            # the resource type.
            resType = determineIDMeaning(specificID, resName)

            rackID = getRackDetails(specificID)

            rowID, rowName, rackLocation = getRowDetails(rackID)
            

            if resType == "Server" or resType == "ESX" or resType == "vm" or resType == "VCenter" or resType == "Switch":
                # get the IP data for all of the above.
                ip = getIPDetaiLs(specificID)

                # switches also have ports, get that data for them.
                if resType == "Switch":
                    ## portName = getPortDetails(specificID)
                    portID = getPortDetails(specificID)

    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")

    return (specificID, resName, theLabel, assets, problems, notes, resType, ip, portID, rowID, rowName, rackLocation)

## Some resources have IP details, get those from the right table, given an object_id.
def getIPDetaiLs(object_id):
    ipData = ""

    try:
        statement = "SELECT ip FROM IPv4Allocation WHERE object_id=%s"

        data = (object_id,)
        cursor.execute(statement, data)

        for ip in cursor:
            print("Retreived the data")
            ipData = ip

    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")
    return (ipData)

## Some resources have port details, get those from the right table, given an object_id.
## What do we do with resources that have multiple ports? We only support one port resources and
    # they need to be int. Do we change the structure of the existing code?  
    # currently using port_id because it has more "connections" throughout the data base.
def getPortDetails(object_id):
    ## portData = ""
    portID = 0

    try:
        ## statement = "SELECT id, name FROM Port WHERE object_id=%s"
        statement = "SELECT id FROM Port WHERE object_id=%s"

        data = (object_id,)
        cursor.execute(statement, data)

        for (id) in cursor:
            print("Retreived the data")
            portID = id            
            ## portData = name

    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")
    return (portID)

def getRackDetails(object_id):
    rackID = ""
    try:
        statement = "SELECT rack_id  FROM rackspace WHERE object_id=%s"

        data = (object_id,)
        cursor.execute(statement, data)

        for rack_id in cursor:
            print("Retreived the data")
            rackID = rack_id

    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")
    return (rackID)

## get details such as 
    # rowID, rowName, location for each rack and row.
def getRowDetails(id):
    rowID, rowName, rackLocation = ""
    try:
        statement = "SELECT row_id, row_name, location_name FROM Rack WHERE id=%s"

        data = (id,)
        cursor.execute(statement, data)

        for (row_id, row_name, location_name) in cursor:
            print("Retreived the data")
            rowID = row_id
            rowName = row_name
            rackLocation = location_name

    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")
    return (rowID, rowName, rackLocation)