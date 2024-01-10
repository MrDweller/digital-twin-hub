# digital-twin-hub
 A system that allows the user to create digital twins by entering a json format of options for the twin. 
 The digital twin will connect to a physical twin and also set up a rest api that will act as a proxy for the physical twin.

## Setup
Create a `.env` file and add,

```
ADDRESS=<address>
PORT=<port>

SYSTEM_NAME=<system name>


MONGO_DB_CONNECTION_STRING=<uri to connect to a mongoDB instance>

SERVER_MODE=<what security the server uses>

SERVICE_REGISTRY_ADDRESS=<service registry address>
SERVICE_REGISTRY_PORT=<service registry port>
SERVICE_REGISTRY_IMPLEMENTATION=<service registry implementation>

DIGITAL_TWIN_REGISTRY_ADDRESS=<digital twin registry address>
DIGITAL_TWIN_REGISTRY_PORT=<digital twin registry port>
DIGITAL_TWIN_REGISTRY_IMPLEMENTATION=<digital twin registry implementation>

CERT_FILE_PATH=<path to cert .pem file>
KEY_FILE_PATH=<path to key .pem file>
TRUSTSTORE_FILE_PATH=<path to truststore .pem file>
AUTHENTICATION_INFO=<authentication info>
```

* `ADDRESS`: defines what address the system will start on, and also register in the service registry.
* `PORT`: defines what port the system will start on, and also register in the service registry.

* `SYSTEM_NAME`: defines what name the system will use to register in the service registry. 

* `MONGO_DB_CONNECTION_STRING`: connection uri of a mongoDB instance that will be connected to. Here data about running digital twins will be stored. 

* `SERVER_MODE`: what security the server uses. Currently these are implemented, 
    * `unsecure`: http is used for all communications.
    * `secure`: https is used for all communications. Also `CERT_FILE_PATH`, `KEY_FILE_PATH` and `TRUSTSTORE_FILE_PATH` is recuired if this is enabled.

* `SERVICE_REGISTRY_ADDRESS`: address of the service registry that will be used to register this system.
* `SERVICE_REGISTRY_PORT`: port of the service registry that will be used to register this system.
* `SERVICE_REGISTRY_IMPLEMENTATION`: the type of service registry that will be used. Currently these are implemented, 
    * `serviceregistry-arrowhead-4.6.1`: uses the service registry of arrowhead version 4.6.1

* `DIGITAL_TWIN_REGISTRY_ADDRESS`: address of the digital twin registry that will be used to register all digital twins.
* `DIGITAL_TWIN_REGISTRY_PORT`: port of the digital twin registry that will be used to register all digital twins.
* `DIGITAL_TWIN_REGISTRY_IMPLEMENTATION`: the type of digital twin registry that will be used. Currently these are implemented, 
    *  `digital-twin-registry-arrowhead-4.6.1`: uses the service registry of arrowhead version 4.6.1 as a digital twin registry.

* `CERT_FILE_PATH`: path to cert .pem file, that will be used in setting up ssl communication. This is necessary if the system uses secure mode.
* `KEY_FILE_PATH`: path to key .pem file, that will be used in setting up ssl communication. This is necessary if the system uses secure mode.
* `TRUSTSTORE_FILE_PATH`: path to truststore .pem file, this defines what systems to trust.This is necessary if the system uses secure mode.
* `AUTHENTICATION_INFO`: authentication info is the public key of the cert file. Used with arrowhead to give the public key to other systems.

## Documentation 
* [Digital Twin System of System description](https://github.com/MrDweller/digital-twin-hub/blob/master/docs/Digital-Twin_SoSD.pdf)
* [Digital Twin System description](https://github.com/MrDweller/digital-twin-hub/blob/master/docs/DigitalTwin/digital-twin-hub_sysd.pdf)
* [Digital Twin Hub System description](https://github.com/MrDweller/digital-twin-hub/blob/master/docs/DigitalTwin/digital-twin_sysd.pdf)
* [Create Digital Twin Service description](https://github.com/MrDweller/digital-twin-hub/blob/master/docs/DigitalTwin/create-digital-twin_sd.pdf)
* [Create Digital Twin Interface design description](https://github.com/MrDweller/digital-twin-hub/blob/master/docs/DigitalTwin/create-digital-twin_idd.pdf)
* [Remove Digital Twin Service description](https://github.com/MrDweller/digital-twin-hub/blob/master/docs/DigitalTwin/remove-digital-twin_sd.pdf)
* [Remove Digital Twin Interface design description](https://github.com/MrDweller/digital-twin-hub/blob/master/docs/DigitalTwin/remove-digital-twin_idd.pdf)
* [Controller System description](https://github.com/MrDweller/controller-system/blob/master/docs/controller_sysd.pdf)
* [Sensor retrieval System description](https://github.com/MrDweller/sensor-retrieval-system/blob/master/docs/sensor-retrieval_sysd.pdf)