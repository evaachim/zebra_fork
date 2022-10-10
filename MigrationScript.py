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

## Determine the subtype, that is, the final type of the resource; 
    # for example, a type can be a server and the final type can be esx.
def detSubType(name):
        if "ESX" in name or "VSHIELD" in name:
            ## Consider VSHIELD as esx sice it regulates traffic into the esx.
            type = "ESXServer"

        ## consider TRUNK as a vlan type
        elif "TRUNK" in name or "VLAN" in name:
            if "SPINE" in name or "LEAF" in name:
                type = "Switch"
            elif "IPC" in name:
                type = "IPAddressPool"

        elif "N3K" in name or "FEX" in name:
            # N3K specialises in IP Address Management, Active Directory & Cloud Management and Privileged Access Management & Security, 
                # allowing it to provide comprehensive support to IT organisations in these areas
            # FEX is a companion to a Nexus 5000 or Nexus 7000 switch
            type = "Switch"
        
        elif "ASA" in name or "BRIDGE" in name:
            ## Cisco ASA is a security device that combines firewall, antivirus, intrusion prevention, 
                # and virtual private network (VPN) capabilities

            ## Could we possibly consider BRIDGE as VLAN since it connects network segments?
            type = "VLAN"
        
        elif "LAB" in name:
            type = "Lab"
        
        elif "VCENTER" in name or "VC" in name:
            type = "VCenter"
        
        elif "JENKINS" in name or "Server" in name:
            type = "Server"
    
def determineType(name):
    type = ""
    prefix = name[0:4]

    ### Deal with prefixes first
    if "SW" in prefix:
        type = "switch"
        type = detSubType(name)
        
    ## elif "SWVM" in name:
      ##  type = "vm"
    ### so on.

    elif "VM" in prefix or "-VM" in prefix:
        type = "VM"
        type = detSubType(name)

    elif "APPL" in prefix:
        ## consider applience to be a server.
        type = "Server"
        type = detSubType(name)
    
    elif "BRIDG" in prefix:
        type = "VLAN"
        

    else: type =  detSubType(name)
    
    ## need some for sdk, cake, dev, vpod, 
    ## also for aci and noiro --> Cisco uses them for data centers.
    ## apic goes with sdn and is classified as a controller, is used for dcs, 
        #  and is in relation with switches.

    ## k8b is for managing a cluster of machines - could we place it under the vm category?

    ## What do we classify HV as ? - It's a network but we don't have a network struct, 
        # we do, however, have a network category.
    
    ## What category do we place WINDRIVER  in.

    ## What category should we place BRIDGE in ? Should it be vlan?

    ## Also, where do we classify centos?   

    ## Return the type and use it in the get functions.
    return type