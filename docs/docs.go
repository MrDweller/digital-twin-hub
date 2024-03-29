// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "url": "https://github.com/MrDweller/digital-twin-hub"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/create-digital-twin": {
            "post": {
                "description": "Create a new digital twin based on the given JSON object. This will create a connection to the physical twin based on the connection info given, this will also include generating endpoints to controll and view sensed data.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Management"
                ],
                "summary": "Create a new digital twin",
                "parameters": [
                    {
                        "description": "DigitalTwinModel JSON",
                        "name": "DigitalTwinModel",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/manufacturer.DigitalTwinModelDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/manufacturer.SystemDefinitionDTO"
                        }
                    }
                }
            }
        },
        "/remove-digital-twin": {
            "delete": {
                "description": "Delete a digital twin based on the given address and port.",
                "tags": [
                    "Management"
                ],
                "summary": "Delete a digital twin",
                "parameters": [
                    {
                        "type": "string",
                        "description": "address",
                        "name": "address",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "port ",
                        "name": "port",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        }
    },
    "definitions": {
        "manufacturer.ConnectionDTO": {
            "type": "object",
            "properties": {
                "connectionModel": {
                    "$ref": "#/definitions/manufacturer.ConnectionModelDTO"
                },
                "connectionType": {
                    "type": "string",
                    "default": "simple-CoAP"
                }
            }
        },
        "manufacturer.ConnectionModelDTO": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string",
                    "default": "localhost"
                },
                "port": {
                    "type": "integer",
                    "default": 5000
                }
            }
        },
        "manufacturer.ControllPropertiesDTO": {
            "type": "object",
            "properties": {
                "serviceDefinition": {
                    "type": "string",
                    "default": "lamp"
                },
                "serviceUri": {
                    "type": "string",
                    "default": "/lamp"
                }
            }
        },
        "manufacturer.DigitalTwinModelDTO": {
            "type": "object",
            "properties": {
                "controlCommands": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/manufacturer.ControllPropertiesDTO"
                    }
                },
                "physicalTwinConnection": {
                    "$ref": "#/definitions/manufacturer.ConnectionDTO"
                },
                "sensedProperties": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/manufacturer.SensedPropertiesDTO"
                    }
                }
            }
        },
        "manufacturer.SensedPropertiesDTO": {
            "type": "object",
            "properties": {
                "intervalTime": {
                    "type": "integer",
                    "default": 10
                },
                "sensorEndpointMode": {
                    "type": "string",
                    "default": "INTERVAL_RETRIEVAL"
                },
                "serviceDefinition": {
                    "type": "string",
                    "default": "temperature"
                },
                "serviceUri": {
                    "type": "string",
                    "default": "/temperature"
                }
            }
        },
        "manufacturer.SystemDefinitionDTO": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string",
                    "default": "localhost"
                },
                "authenticationInfo": {
                    "type": "string"
                },
                "port": {
                    "type": "integer",
                    "default": 5000
                },
                "systemName": {
                    "type": "string",
                    "default": "my-digital-twin"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Digital Twin Hub",
	Description:      "This page shows the REST interfaces offered by the Digital Twin Hub.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
