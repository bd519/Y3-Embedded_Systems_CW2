import RPi.GPIO as GPIO

def setup_alert():
	GPIO.setmode(GPIO.BCM)
	GPIO.setup(18, GPIO.OUT)

def mqtt_callback(client, userdata, message):
	if message.payload == b'True' or message.payload == 'True' or message.payload == True:
		alert()
	else:
		turnoff()

def alert():
	GPIO.output(18,True)

def turnoff():
	GPIO.output(18,False)
