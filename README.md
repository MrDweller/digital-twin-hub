# digital-twin-hub
 A system that allows the user to create digital twins by entering a json format of options for the twin. 
 The digital twin will connect to a physical twin and also set up a rest api that will act as a proxy for the physical twin.

## Setup
Create a `.env` file and add,

```
ADDRESS=<address>
PORT=<port>
SYSTEM_NAME=<system name>

SERVICE_REGISTRY_ADDRESS=<service registry address>
SERVICE_REGISTRY_PORT=<service registry port>
SERVICE_REGISTRY_IMPLEMENTATION=<service registry implementation>

DIGITAL_TWIN_REGISTRY_ADDRESS=<digital twin registry address>
DIGITAL_TWIN_REGISTRY_PORT=<digital twin registry port>
DIGITAL_TWIN_REGISTRY_IMPLEMENTATION=<digital twin registry implementation>

CERT_FILE_PATH=<path to cert .pem file>
KEY_FILE_PATH=<path to key .pem file>
TRUSTSTORE_FILE_PATH=<path to truststore .pem file>
```