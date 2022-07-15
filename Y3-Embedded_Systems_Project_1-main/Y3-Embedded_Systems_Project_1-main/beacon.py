class Beacon(object):
    def __init__(self, data, address, time_mili):
        self.uuid = data[0]
        self.address = address

        self.rssi_strengths = [data[4]]
        self.scan_times = [time_mili]

    def __str__(self):
        ret = "Beacon: address:{ADDR} uuid:{UUID} major:{MAJOR}"\
                " minor:{MINOR} txpower:{POWER} rssi:{RSSI}"\
                .format(ADDR=self.address, UUID=self.uuid, MAJOR=self._major,
                        MINOR=self._minor, POWER=self._power, RSSI=self.rssi_strengths[0])
        return ret
    
    def found_again(self, rssi, time_mili):
        self.scan_times.append(time_mili)
        self.rssi_strengths.append(rssi)