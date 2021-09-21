// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Devashish Gupta",
            "email": "devagpta@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/login": {
            "post": {
                "description": "Login verification and first load of complete user data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "data retrieval"
                ],
                "summary": "Login API",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user identifier",
                        "name": "username",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.LoginResponse"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Pings the server for checking the health of the server",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "root"
                ],
                "summary": "Pings the server",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.PongResponse"
                        }
                    }
                }
            }
        },
        "/refresh/population": {
            "post": {
                "description": "Refresh endpoint to quickly refresh population data with the latest values",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "data retrieval"
                ],
                "summary": "Refresh population API",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user identifier",
                        "name": "username",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "planet identifier",
                        "name": "planet_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Population"
                        }
                    }
                }
            }
        },
        "/refresh/resources": {
            "post": {
                "description": "Refresh endpoint to quickly refresh mine data with the latest values",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "data retrieval"
                ],
                "summary": "Refresh mine API",
                "parameters": [
                    {
                        "type": "string",
                        "description": "user identifier",
                        "name": "username",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "planet identifier",
                        "name": "planet_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "mine identifier",
                        "name": "mine_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Resources"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.PongResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "pong"
                }
            }
        },
        "models.Clan": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "example": "Mind Krackers"
                },
                "role": {
                    "type": "string",
                    "example": "MEMBER"
                }
            }
        },
        "models.EmployedPopulation": {
            "type": "object",
            "properties": {
                "idle": {
                    "type": "integer",
                    "example": 4
                },
                "total": {
                    "type": "integer",
                    "example": 21
                }
            }
        },
        "models.Experience": {
            "type": "object",
            "properties": {
                "current": {
                    "type": "integer",
                    "example": 185
                },
                "level": {
                    "type": "integer",
                    "example": 4
                },
                "required": {
                    "type": "integer",
                    "example": 368
                }
            }
        },
        "models.LoginResponse": {
            "type": "object",
            "properties": {
                "home_sector": {
                    "$ref": "#/definitions/models.Sector"
                },
                "occupied_planets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.StaticPlanetData"
                    }
                },
                "profile": {
                    "$ref": "#/definitions/models.Profile"
                }
            }
        },
        "models.Mine": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string",
                    "example": "W101"
                },
                "max_limit": {
                    "type": "integer",
                    "example": 550
                },
                "mined": {
                    "type": "integer",
                    "example": 125
                },
                "mining_plant": {
                    "$ref": "#/definitions/models.MiningPlant"
                },
                "type": {
                    "type": "string",
                    "example": "WATER"
                }
            }
        },
        "models.MiningPlant": {
            "type": "object",
            "properties": {
                "building_id": {
                    "type": "string",
                    "example": "WMP101"
                },
                "building_level": {
                    "type": "integer",
                    "example": 3
                },
                "next_level": {
                    "$ref": "#/definitions/models.NextLevelAttributes"
                },
                "workers": {
                    "type": "integer",
                    "example": 12
                }
            }
        },
        "models.NextLevelAttributes": {
            "type": "object",
            "properties": {
                "current_mining_rate_per_worker": {
                    "type": "integer",
                    "example": 1
                },
                "current_workers_max_limit": {
                    "type": "integer",
                    "example": 40
                },
                "graphene_required": {
                    "type": "integer",
                    "example": 101
                },
                "max_mining_rate_per_worker": {
                    "type": "integer",
                    "example": 12
                },
                "max_workers_max_limit": {
                    "type": "integer",
                    "example": 130
                },
                "next_mining_rate_per_worker": {
                    "type": "integer",
                    "example": 1
                },
                "next_workers_max_limit": {
                    "type": "integer",
                    "example": 65
                },
                "shelio_required": {
                    "type": "integer",
                    "example": 0
                },
                "water_required": {
                    "type": "integer",
                    "example": 5
                }
            }
        },
        "models.OccupiedPlanet": {
            "type": "object",
            "properties": {
                "home": {
                    "type": "boolean",
                    "example": true
                },
                "mines": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Mine"
                    }
                },
                "planet_config": {
                    "type": "string",
                    "example": "Planet2.json"
                },
                "population": {
                    "$ref": "#/definitions/models.Population"
                },
                "position": {
                    "$ref": "#/definitions/models.PlanetPosition"
                },
                "resources": {
                    "$ref": "#/definitions/models.Resources"
                }
            }
        },
        "models.PlanetPosition": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string",
                    "example": "023:049:07"
                },
                "planet": {
                    "type": "integer",
                    "example": 7
                },
                "sector": {
                    "type": "integer",
                    "example": 49
                },
                "system": {
                    "type": "integer",
                    "example": 23
                }
            }
        },
        "models.Population": {
            "type": "object",
            "properties": {
                "generation_rate": {
                    "type": "integer",
                    "example": 3
                },
                "soldiers": {
                    "$ref": "#/definitions/models.EmployedPopulation"
                },
                "total": {
                    "type": "integer",
                    "example": 45
                },
                "unemployed": {
                    "type": "integer",
                    "example": 3
                },
                "workers": {
                    "$ref": "#/definitions/models.EmployedPopulation"
                }
            }
        },
        "models.Profile": {
            "type": "object",
            "properties": {
                "clan": {
                    "$ref": "#/definitions/models.Clan"
                },
                "experience": {
                    "$ref": "#/definitions/models.Experience"
                },
                "username": {
                    "type": "string",
                    "example": "devashish"
                }
            }
        },
        "models.Resource": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer",
                    "example": 23
                },
                "max_limit": {
                    "type": "number",
                    "example": 100
                },
                "reserved": {
                    "type": "integer",
                    "example": 14
                }
            }
        },
        "models.Resources": {
            "type": "object",
            "properties": {
                "graphene": {
                    "$ref": "#/definitions/models.Resource"
                },
                "shelio": {
                    "type": "integer",
                    "example": 23
                },
                "water": {
                    "$ref": "#/definitions/models.Resource"
                }
            }
        },
        "models.Sector": {
            "type": "object",
            "properties": {
                "occupied_planets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.OccupiedPlanet"
                    }
                },
                "position": {
                    "$ref": "#/definitions/models.SectorPosition"
                },
                "unoccupied_planets": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.UnoccupiedPlanet"
                    }
                }
            }
        },
        "models.SectorPosition": {
            "type": "object",
            "properties": {
                "_id": {
                    "type": "string",
                    "example": "023:049"
                },
                "sector": {
                    "type": "integer",
                    "example": 49
                },
                "system": {
                    "type": "integer",
                    "example": 23
                }
            }
        },
        "models.StaticPlanetData": {
            "type": "object",
            "properties": {
                "planet_config": {
                    "type": "string",
                    "example": "Planet2.json"
                },
                "position": {
                    "$ref": "#/definitions/models.PlanetPosition"
                },
                "same_sector": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "models.UnoccupiedPlanet": {
            "type": "object",
            "properties": {
                "planet_config": {
                    "type": "string",
                    "example": "Planet2.json"
                },
                "position": {
                    "$ref": "#/definitions/models.PlanetPosition"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0.0",
	Host:        "mmo-server-4xcaklgmnq-el.a.run.app",
	BasePath:    "",
	Schemes:     []string{"https"},
	Title:       "MMO Game Server",
	Description: "This is the server for new MMO Game based in space.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
