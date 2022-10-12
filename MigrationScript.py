## Preliminary script to get data from the racktables database into the Zebra tool.
##
## This could be extended for any type of resource for which there is data in the racktables database.

## Are we updating the zebra tool to have this data statically and then access it via the api calls? => V.
## Or we using the api calls to fetch from the db only that data which we need for the specific call?

## imports needed for access to mysql and the racktables database.
import os
from pickletools import uint8
import re
import mysql.connector as database

import requests
import mysql
import sys 
import pyodbc

from collections import namedtuple

## Struct that will be used to get/access data.
zebraData = namedtuple('Resource', ['resType','specificID', 'resName', 'theLabel', 'assets', 'problems', 'notes', 'ip', 'owner', 'portID', 'rowID', 'rowName', 'rackID', 'rackLocation'])

## List with all of the data structs i.e all resources.
zebraList = []

## The local user to connect to the database.
## Must set these up as local variables to use as below.
username = os.environ.get("username")
password = os.environ.get("password")

## Option to only import by api query (not using this method currently)=> too costly
## API routes.
types = "/api/v1/types"
labels = "/api/v1/labels"
query = "/api/v1/resources"
posts = "/api/v1/resources"
delete = "/api/v1/resources"

## Make the connection to the correct data base.
connection = database.connect(
    user = username,
    password = password,
    host = "localhost",
    database = "racktables"
)

## This is where the database cursor goes.
cursor = connection.cursor()

'''
## Option to only import by api query => too costly
def getQType():
    key, kind = getQuery()
    query = ""

    if kind == "id":
        query = "ID"
    elif kind == "type":
        query = "type"
    elif kind == "label":
        query = "label"
    else: query = ""

    return query
'''

'''
## Option to only import by api query pert 2 => too costly
def getQuery():
    queryResponse = requests.get(query)
    ## This contains the id to look for when querring.
    ## Querry can be of different types - by id, by type, by label. These are in api.go.
    criteriaID = queryResponse.json()["id"]
    criteriaType = queryResponse.json()["type"]
    criteriaLabel = queryResponse.json()["label"]

    key, kind = " "

    if len(criteriaID) != 0:
        key = criteriaID
        kind = "id"
    elif len(criteriaType) !=0:
        key = criteriaType
        kind = "type"
    elif len(criteriaLabel) != 0:
        key = criteriaLabel
        kind = "label"
    else: 
        key = ""
        kind = ""
    ## Return and use.
    return key, kind
'''

## Get the type by id.
## pdu and ups can be considered as vm solutions so they're both added to vm.
## some are randomly dropped into 1504 but this objtype_id has no reference. Thus, name will be used too.
def determineIDMeaning(id, name):
    means, final, this = ""
    id = str(id)

    if id == "2" or id == "27":
        means = "VM"
    elif id == "30" or id == "31" or id == "34":
        means = "Rack"
    elif id == "3":
        means = "Shelf"
    elif id == "38":
        means = "VCenter"
    elif id == "4" or id == "13" or id == "36":
        means = "Server"
    elif id == "8" or id == "12" or id == "14" or id == "21" or id == "26" or id == "32" or id == "33":
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

## determine the specific type of a resource.
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

    ## vcenter is a management interface type => management interface type = vcenter
    
    if means == "Shelf":
        type = "Rack"
    elif means == "Compute":
        if "esx" in name:
            type = "ESX"
        elif "jenkins" in name or "server" in name or "srv" in name or "vintella" in name:
            type = "Server"
        elif "bld" in name or "datacenter" in name or "dc" in name:
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
        if "chasis" in name or "ixia" in name or "rack" in name:
            type = "Rack"
        elif "nexus" in name or "switch" in name or "sw" in name or "n3k" in name:
            type = "Switch"
    else: type = means

    return type

'''
## Option to only import by api query pert 3 => too costly
## This would not work because there are several id's that mean the same thing and also several types that go into the same typeID category 
    # It'd be too costly to calculate the distinction.

def convertTypeToTypeID(kind):
    theTypeID = -1

    if kind == "Shelf":
        theTypeID = 30
    elif kind ==  "Rack":
        theTypeID = 3
    elif kind == "jenkins" or kind == "server" or  kind == "srv" or kind == "vintella":
        theTypeID = 4
    elif kind == "esx":
        theTypeID = 1504

    return theTypeID
'''

## Get data from database.
def getData():
    ## variables to store info temporarily.
    type, rack_ID, port_ID, owned_by, row_ID, row_Name, rack_Location = ""
    IP = "N/A"
    port_ID = 0
    
    queryWOkey = "SELECT id, name, label, objtype_id, asset_no, has_problems, comment FROM rackobject%s"

    try:
        statement = queryWOkey

        ## data = (key,)
        cursor.execute(statement)

        for (id, name, label, objtype_id, asset_no, has_problems, comment) in cursor:
            print("Retreived the data")

            # get the resource's type based on its data.
            type = determineIDMeaning(objtype_id, name) 
            zebraData.resType = type

            # resource specific data.
            zebraData.specificID  = id
            zebraData.resName = name
            
            # currently null for all, will be usefoul for containment model.
            zebraData.theLabel = label

            # additional information.
            zebraData.assets = asset_no
            zebraData.problems = has_problems
            zebraData.notes = comment

            rack_ID = getRackDetails(id)
            # the resource's rack information.
            zebraData.rackID = rack_ID

            # further details about the rack and row.
            row_ID, row_Name, rack_Location = getRowDetails(rack_ID)
            zebraData.rowID = row_ID 
            zebraData.rowName = row_Name 
            zebraData.rackLocation = rack_Location
            
            if zebraData.resType == "Server" or zebraData.resType == "ESX" or zebraData.resType == "vm" or zebraData.resType == "VCenter" or zebraData.resType == "Switch":
                # get the IP data for all of the above.
                IP = getIPDetaiLs(id)
                zebraData.ip = IP

                # all resources with ip need to have their owner's info extracted.
                owned_by = getUserDetails(IP)
                zebraData.owner = owned_by

                # switches also have ports, get that data for them.
                if zebraData.resType == "Switch":
                    port_ID = getPortDetails(id)
                    zebraData.portID = port_ID

            ## add this struct to the list, there might be many.
            zebraList.append(zebraData)

    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")

    # return (specificID, resName, theLabel, assets, problems, notes, resType, ip, owner, portID, rowID, rowName, rackLocation)
    return zebraList

## Some resources have IP details, get those from the right table, given an object_id.
def getIPDetaiLs(object_id):
    ipData = ""

    try:
        statement = "SELECT ip FROM IPv4Allocation WHERE object_id =%s"

        data = (object_id,)
        cursor.execute(statement, data)

        for ip in cursor:
            print("Retreived the data")
            ipData = ip
            ## yield (ipData)

    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")

    return ipData 

## Some resources have port details, get those from the right table, given an object_id.
## What do we do with resources that have multiple ports? We only support one port resources and
    # they need to be int. Do we change the structure of the existing code?  
    # currently using port_id because it has more "connections" throughout the data base.
def getPortDetails(object_id):
    ## portData = ""
    portID = 0
    portList = []
    numPorts = 0

    try:
        ## statement = "SELECT id, name FROM Port WHERE object_id=%s"
        statement = "SELECT id FROM Port WHERE object_id=%s"

        data = (object_id,)
        cursor.execute(statement, data)

        for (id) in cursor:
            print("Retreived the data")
            portID = id            
            ## portData = name
            ## yield (portID)
            portList.append(portID)

        ## Instead of returning the actual port numbers, return how many ports this switch has.
            # Needed for our current implementation of zebra.
        numPorts = len(portList)

    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")
    
    return numPorts

## details about each rack, depending on the object's id, will be used in further queries.
def getRackDetails(object_id):
    RackID = ""
    # rackList = []

    try:
        statement = "SELECT rack_id  FROM rackspace WHERE object_id =%s"

        data = (object_id,)
        cursor.execute(statement, data)

        for rack_id in cursor:
            print("Retreived the data")
            rackID = rack_id
            # rackList.append(rackID)

        '''
        ## Removed to simplify 
            # Not necesary since although object_ids may repeat, they will be on the same rack
                # since representing the same resource.
        if(len(set(rackList)) !=1):
            rackID = rackID
            print("Error")
        '''

    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")

    return rackID

## get details such as 
    # rowID, rowName, location for each rack and row.
def getRowDetails(id):
    rowID, rowName, rackLocation = ""
    try:
        statement = "SELECT row_id, row_name, location_name FROM Rack WHERE id =%s"

        data = (id,)
        cursor.execute(statement, data)

        for (row_id, row_name, location_name) in cursor:
            print("Retreived the data")
            rowID = row_id
            rowName = row_name
            rackLocation = location_name
            ## yield (rowID, rowName, rackLocation)

            if rowID == "":
                rowID = "N/A"
            elif rowName == "":
                rowName = "N/A"
            elif rackLocation == "":
                rackLocation = "N/A"

    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")
    
    return (rowID, rowName, rackLocation)

## for user or owner: IPv4Log .
## get user owner / user data.
def getUserDetails(resIP):
    ownedBy = ""
    try:
        statement = "SELECT user FROM IPv4Log WHERE ip=%s"

        data = (resIP,)
        cursor.execute(statement, data)

        for (user) in cursor:
            print("Retreived the data")
            ownedBy = user
            ## yield (ownedBy)

    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")
    
    return (ownedBy)