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
                        "description": "DigitalTwinDTO JSON",
                        "name": "DigitalTwin",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/manufacturer.DigitalTwinDTO"
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
        },
        "/upload-certificates": {
            "post": {
                "description": "Upload certificate files as zip, to be used by a digital twin. Takes cert.pem and key.pem files and gives ` + "`" + `certId` + "`" + `.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Management"
                ],
                "summary": "Upload certificate files as zip",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Cert file",
                        "name": "cert",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "file",
                        "description": "Key file",
                        "name": "key",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/manufacturer.CertificateDTO"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "manufacturer.AnomalyDTO": {
            "type": "object",
            "properties": {
                "anomalyType": {
                    "type": "string",
                    "default": "stuck"
                }
            }
        },
        "manufacturer.CertificateDTO": {
            "type": "object",
            "properties": {
                "certificateId": {
                    "type": "string"
                }
            }
        },
        "manufacturer.ConnectionDTO": {
            "type": "object",
            "properties": {
                "connectionModel": {
                    "type": "object",
                    "additionalProperties": {}
                },
                "connectionType": {
                    "type": "string",
                    "default": "simple-CoAP"
                }
            }
        },
        "manufacturer.ControllCommandDTO": {
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
        "manufacturer.DigitalTwinDTO": {
            "type": "object",
            "properties": {
                "certificateId": {
                    "type": "string"
                },
                "controlCommands": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/manufacturer.ControllCommandDTO"
                    }
                },
                "handleableAnomalies": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/manufacturer.AnomalyDTO"
                    }
                },
                "physicalTwinConnection": {
                    "$ref": "#/definitions/manufacturer.ConnectionDTO"
                },
                "sensedProperties": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/manufacturer.SensedPropertyDTO"
                    }
                },
                "systemName": {
                    "type": "string"
                }
            }
        },
        "manufacturer.SensedPropertyDTO": {
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
