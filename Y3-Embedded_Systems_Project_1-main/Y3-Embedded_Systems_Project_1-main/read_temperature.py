import smbus2
import time

class ReadTemperature:
	"""Class that checks if the user has fever (i.e. average temp over 38 C). """
	def __init__(self, N_measurements=10):
		self.data = []
		self.N_measurements = N_measurements
		self.average_temperature = None

	def read(self):
		"""Reads the temperature. Returns the temperature in Celsius.
	 	Should not be done too often."""
		bus = smbus2.SMBus(1)
		bytes_temp = bus.read_i2c_block_data(0x40, 0xE3, 2)
		return (bytes_temp[0]*256 + bytes_temp[1])*175.72/65536 - 46.85

	def check_high_temperature(self):
		"""Checks if the wearer has a suspected temperature 
		based on the last self.N_measurements temp. measurements.
		Returns True if the average has been over 38 AND the last value is over 38, 
		False otherwise."""

		# Read and store the temperature
		temp = self.read()
		if len(self.data) >= self.N_measurements: # Then pop the oldest measurement and append
			self.data.pop(0)
			self.data.append(temp)
		else: # Then just append the temp measurement
			self.data.append(temp)

		# Compute the average
		self.average = sum(self.data) / len(self.data)

		# Return
		if len(self.data) < self.N_measurements: return False
		else: return (self.average > 38 and self.data[-1] > 38)

	def get_average(self):
		"""Returns the most recent estimate of the average temperature."""
		return self.average


	def _read_10_values(self):
		print("Measuring")
		for i in range(10):
			print(self.read())
			time.sleep(0.1)
		print("All done")

if __name__ == '__main__':
	temperatureObj = ReadTemperature()
	temperatureObj._read_10_values()
