
from gattlib import BeaconService
import bluetooth_tracking as bt
import database as db
import random
import beacon

def test_operate_tracker(mocker):
    def mock_one_random_beacon(self, scan_time):
        mock_uuid = "1"
        mock_rssi = random.randint(-100,-20)
        mock_address="11:22:33:44:55:66"
        return {mock_address:[mock_uuid,0,0,0,mock_rssi]}
    mocker.patch(
        'gattlib.BeaconService.scan',
        mock_one_random_beacon
    )
    server_ip = "146.169.174.8"
    id_var = random.randint(1000,9999)
    host_uuid = "11111111-2222-3333-" + str(id_var) + "-555555555555"
    database = db.Database(server_ip, host_uuid)
    bt.operate_tracker(database, host_uuid)