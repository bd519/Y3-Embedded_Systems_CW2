#!/usr/bin/python3
# -*- mode: python; coding: utf-8 -*-

# Copyright (C) 2014, Oscar Acena <oscaracena@gmail.com>
# This software is under the terms of Apache License v2 or later.

from gattlib import BeaconService
import time
import random 
from datetime import datetime

import beacon
import database as db

def broadcast( beacon_uuid:str, broadcast_length_min: int, broadcast_length_max :int):
    print("Broadcasting.")

    service = BeaconService("hci0")

    service.start_advertising(beacon_uuid,
                1, 1, 1, 10000000)
    time.sleep(random.randint(broadcast_length_min, broadcast_length_max))
    service.stop_advertising()

def scan(database : db.Database, scan_length_min : int, scan_length_max : int):
    print("Scanning.")

    service = BeaconService("hci0")
    devices = service.scan(random.randint(scan_length_min, scan_length_max))

    found_new_beacon = False
    for address, data in list(devices.items()):
        print("Found" + str(address))
        time_mili = int(time.time()*1000)
        new_beacon = beacon.Beacon(data, address, time_mili)
        found_new_beacon = database.report_beacon(new_beacon)

    return found_new_beacon

scan_length_min = 2
scan_length_max = 10
number_of_scans = 6

broadcast_length_min = 2
broadcast_length_max = 10

def operate_tracker(database : db.Database, host_uuid : str):
    for i in range(number_of_scans):
        found_beacon = scan(database, scan_length_min, scan_length_max)
        if found_beacon: 
            broadcast(host_uuid,scan_length_max,scan_length_max)
        
    broadcast(host_uuid,broadcast_length_min,broadcast_length_max)
    database.upload_beacons()