## Preliminary script to get data from the racktables database into the Zebra tool.
##
## This could be extended for any type of resource for which there is data in the racktables database.

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

## get the response for the query.
queryResponse = requests.get(query)

## Make the response into json format.
criteria = queryResponse.json()

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
def get_db_resource_data(criteria):
    dbName, row, location = ""
    try:
        statement = "SELECT name, row_name, row_id, location_name FROM rackobject WHERE id=%s"
        ## selection is now by the object's id, it could be something else, as needed.
        data = (criteria,)
        cursor.execute(statement, data)
        for (name, row_name, row_id, location_name) in cursor:
            print("Retreived the data")
            dbName = name
            row = row_name
            rowID = row_id
            location = location_name
    except database.Error as e:
        print(f"Error retreiving entry from database: {e}")
    return (dbName, row, rowID, location)


## Get data for network resources.


## get_net_resource_data() should go here.



## Get data for compute resources.

## get_compute_resource_data() should go here.