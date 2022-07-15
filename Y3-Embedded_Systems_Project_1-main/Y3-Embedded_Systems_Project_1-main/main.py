import random
import database
import bluetooth_tracking
from read_temperature import ReadTemperature
import argparse 
import alert
import yaml
import os
import paho.mqtt.client as mqtt

demo_mode = False

def extract_uuid_from_yaml(username : str):
    
    with open(os.path.abspath("./server/lookup_table.yml")) as file:
        documents = yaml.full_load(file)
        print(documents)
        user_ids = documents["users"]
        user_ids = {v: k for k, v in user_ids.items()}
        return user_ids[username]

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("server_ip", type=str,
                    help="serer_ip")

    parser.add_argument("user_name", type=str,
                    help="user_name")
    parser.add_argument("--notemp", help="Flag to disable temp", action='store_true')
    args = parser.parse_args()
    
    host_uuid = extract_uuid_from_yaml(args.user_name)
    database = database.Database(args.server_ip, host_uuid)
    temperature = ReadTemperature()
    alert.setup_alert()

    client = mqtt.Client()
    ret_code = client.connect(args.server_ip, port=1883)
    print(f"Connect ret code {ret_code}")
    mqtt.error_string(ret_code)

    client.loop_start()
    client.on_message = alert.mqtt_callback
    client.subscribe("IC.es/JBMNsystems/"+str(host_uuid))

    while True:
        bluetooth_tracking.operate_tracker(database, host_uuid)
        if not args.notemp:
            database.report_temperature(temperature.read())

           