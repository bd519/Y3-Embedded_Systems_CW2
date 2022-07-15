#!/usr/bin/python3
# -*- mode: python; coding: utf-8 -*-

# Copyright (C) 2014, Oscar Acena <oscaracena@gmail.com>
# This software is under the terms of Apache License v2 or later.

from gattlib import BeaconService
import time
import random 

class Beacon(object):

    def __init__(self, data, address):
        self._uuid = data[0]
        self._major = data[1]
        self._minor = data[2]
        self._power = data[3]
        self._rssi = data[4]
        self._address = address

    def __str__(self):
        ret = "Beacon: address:{ADDR} uuid:{UUID} major:{MAJOR}"\
                " minor:{MINOR} txpower:{POWER} rssi:{RSSI}"\
                .format(ADDR=self._address, UUID=self._uuid, MAJOR=self._major,
                        MINOR=self._minor, POWER=self._power, RSSI=self._rssi)
        return ret

id_var = random.randint(1000,9999)

def broadcast():
    print("Broadcasting.")

    service = BeaconService("hci0")

    service.start_advertising("11111111-2222-3333-" + str(id_var) + "-555555555555",
                1, 1, 1, 10000000)
    time.sleep(5)
    service.stop_advertising()


def listen():
    print("Listening.")

    service = BeaconService("hci0")
    devices = service.scan(random.randint(2,15))

    for address, data in list(devices.items()):
        b = Beacon(data, address)
        print(b)

    print("Done.")

if __name__ == "__main__":
    while True:
        broadcast()
        listen()

