from read_temperature import ReadTemperature
from datetime import datetime

def main():
	start_time = datetime.now()
	temperature = ReadTemperature()
	measurement_every_x_seconds = 1
	while True:
		if (datetime.now() - start_time).total_seconds() > measurement_every_x_seconds:
			temp = temperature.read()
			start_time = datetime.now()
			print(f"Temperature is {temp}")






if __name__ == '__main__':
	main()