// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/contract/tokens": {
            "get": {
                "description": "Request token list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contract",
                    "tokens"
                ],
                "summary": "Get Tokens",
                "operationId": "getTokens",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.Token"
                            }
                        }
                    },
                    "507": {
                        "description": "Insufficient Storage",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    }
                }
            },
            "post": {
                "description": "Add list of tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Contract",
                    "tokens"
                ],
                "summary": "Add Tokens",
                "operationId": "addTokens",
                "parameters": [
                    {
                        "description": "Add tokens",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.Token"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.Token"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    }
                }
            }
        },
        "/provider/pairs": {
            "get": {
                "description": "Get list of pairs",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Provider",
                    "pairs"
                ],
                "summary": "Get Pairs",
                "operationId": "getPairs",
                "parameters": [
                    {
                        "description": "Get pairs",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/RequestPairs"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.TradePair"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    }
                }
            },
            "post": {
                "description": "Add list of pairs",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Provider",
                    "pairs"
                ],
                "summary": "Add Pairs",
                "operationId": "addPairs",
                "parameters": [
                    {
                        "description": "Add pairs",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.TradePair"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.TradePair"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    }
                }
            }
        },
        "/storage/pools": {
            "get": {
                "description": "Get list of pools",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Storage",
                    "pools"
                ],
                "summary": "Get Pools",
                "operationId": "getPools",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/ListPools"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    }
                }
            },
            "post": {
                "description": "Add list of pools",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Storage",
                    "pools"
                ],
                "summary": "Add Pools",
                "operationId": "addPools",
                "parameters": [
                    {
                        "description": "Add pools",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.TradePool"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.TradePool"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    },
                    "507": {
                        "description": "Insufficient Storage",
                        "schema": {
                            "$ref": "#/definitions/Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "Error": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "message"
                }
            }
        },
        "ListPools": {
            "description": "Request list of trading pools",
            "type": "object",
            "properties": {
                "pools": {
                    "description": "list of trade pools",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.TradePool"
                    }
                }
            }
        },
        "RequestPairs": {
            "description": "Request for searching trade pair",
            "type": "object",
            "properties": {
                "protocol": {
                    "description": "trade protocol",
                    "allOf": [
                        {
                            "$ref": "#/definitions/entities.SwapProtocol"
                        }
                    ]
                },
                "tokenPair": {
                    "description": "pair of tokens",
                    "allOf": [
                        {
                            "$ref": "#/definitions/entities.TokenPair"
                        }
                    ]
                }
            }
        },
        "entities.SwapProtocol": {
            "type": "object",
            "properties": {
                "factoryAddress": {
                    "type": "string"
                },
                "protocolName": {
                    "type": "string"
                },
                "swapRouter": {
                    "type": "string"
                }
            }
        },
        "entities.Token": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "wei": {
                    "type": "integer"
                }
            }
        },
        "entities.TokenPair": {
            "type": "object",
            "properties": {
                "token0": {
                    "$ref": "#/definitions/entities.Token"
                },
                "token1": {
                    "$ref": "#/definitions/entities.Token"
                }
            }
        },
        "entities.TradePair": {
            "type": "object",
            "properties": {
                "pool0": {
                    "$ref": "#/definitions/entities.TradePool"
                },
                "pool1": {
                    "$ref": "#/definitions/entities.TradePool"
                }
            }
        },
        "entities.TradePool": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "pair": {
                    "$ref": "#/definitions/entities.TokenPair"
                },
                "tradeProvider": {
                    "$ref": "#/definitions/entities.SwapProtocol"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "Flash Loan Arbitrage bot",
	Description:      "interacting with smart contract, collecting tokens, perform trade",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}