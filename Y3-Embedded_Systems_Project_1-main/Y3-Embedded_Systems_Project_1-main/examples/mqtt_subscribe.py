import paho.mqtt.client as mqtt

def on_message(client, userdata, message):
    print("Received message:{} on topic {}".format(message.payload, message.topic))

client = mqtt.Client()
# client.tls_set(ca_certs="ca.crt", certfile = "server.crt", keyfile="server.key")
# ret_code = client.connect("146.169.200.162", port=8883)
ret_code = client.connect("146.169.172.254", port=8880)
# ret_code = client.connect("146.169.200.162", port=8880)
print(f"Connect ret code {ret_code}")
mqtt.error_string(ret_code)

client.loop_start()
client.on_message = on_message
client.subscribe("IC.es/JBMNsystems/#")

while True:
    pass