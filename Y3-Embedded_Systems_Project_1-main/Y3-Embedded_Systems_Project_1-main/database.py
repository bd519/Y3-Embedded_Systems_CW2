import time
from influxdb_client import InfluxDBClient, Point, WritePrecision
from influxdb_client.client.write_api import SYNCHRONOUS
import beacon

class Database:
    def __init__(self, server_ip : str, host_uuid: str):
        self.server_ip = server_ip
        self.found_beacons = dict()
        self.host_uuid = host_uuid
    
    def get_database_details(self):
        token = "vdRXm4rGle9XaeOhuah2WHWOyXi81-AMQuuFZTq42HnI-kU5uchjnm6duk6GvO9sdES6YFT7CpV75qxei_0fIw=="
        org = "jbmn"
        bucket = "iot"
        return (token, org, bucket)

    #Returns true if beacon has been found before
    def report_beacon(self, new_beacon : beacon.Beacon):
        if new_beacon.address not in self.found_beacons:
            self.found_beacons[new_beacon.address] = new_beacon
            return True
        else:
            self.found_beacons[new_beacon.address].found_again(new_beacon.rssi_strengths[0], new_beacon.scan_times[0])
            return False
    
    def upload_beacons(self):
        token, org, bucket = self.get_database_details()

        with InfluxDBClient(url="http://" + str(self.server_ip) + ":8086", token=token, org=org) as client:
            write_api = client.write_api(write_options=SYNCHRONOUS)
            for address,beacon in list(self.found_beacons.items()):
                for time, rssi in zip(beacon.scan_times, beacon.rssi_strengths):
                    point = Point("mem") \
                    .tag("host", self.host_uuid) \
                    .tag("neighbour", beacon.uuid) \
                    .field("RSSI", rssi) \
                    .time(time, WritePrecision.MS)
                    write_api.write(bucket, org, point)
            self.found_beacons = dict()
            client.close()
            return True
        return False
    
    def report_temperature(self, temperature:float):
        token, org, bucket = self.get_database_details()

        with InfluxDBClient(url="http://" + str(self.server_ip) + ":8086", token=token, org=org) as client:
            write_api = client.write_api(write_options=SYNCHRONOUS)
           
            point = Point("temperature") \
            .tag("host", self.host_uuid) \
            .field("temperature", temperature) \
            .time(int(time.time()*1000), WritePrecision.MS)
            write_api.write(bucket, org, point)
            client.close()
            return True
        return False