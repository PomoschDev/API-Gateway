// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/users": {
            "get": {
                "description": "Массив пользователей в базе данных",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Список всех пользователей",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/DatabaseServicev1.UsersResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/server.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/users/{id}": {
            "get": {
                "description": "Поиск пользователя по ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Поиск пользователя",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/DatabaseServicev1.CreateUserResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/server.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "DatabaseServicev1.Card": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "description": "* Дата создания сущности в базе данных",
                    "type": "string"
                },
                "cvv": {
                    "description": "* CVV код карты",
                    "type": "integer"
                },
                "date": {
                    "description": "* Дата до которой активна карта",
                    "type": "string"
                },
                "fullName": {
                    "description": "* ФИО с банковской карты",
                    "type": "string"
                },
                "id": {
                    "description": "* ID банковской карты в базе данных",
                    "type": "integer"
                },
                "number": {
                    "description": "* Номер карты",
                    "type": "string"
                },
                "updatedAt": {
                    "description": "* Дата последнего обновления сущности в базе данных",
                    "type": "string"
                },
                "userId": {
                    "description": "* ID пользователя которому принадлежит данная карта",
                    "type": "integer"
                }
            }
        },
        "DatabaseServicev1.CardCompany": {
            "type": "object",
            "properties": {
                "companyId": {
                    "description": "* ID компании которой принадлежит данная карта",
                    "type": "integer"
                },
                "createdAt": {
                    "description": "* Дата создания сущности в базе данных",
                    "type": "string"
                },
                "cvv": {
                    "description": "* CVV код карты",
                    "type": "integer"
                },
                "date": {
                    "description": "* Дата до которой активна карта",
                    "type": "string"
                },
                "fullName": {
                    "description": "* ФИО с банковской карты",
                    "type": "string"
                },
                "id": {
                    "description": "* ID банковской карты в базе данных",
                    "type": "integer"
                },
                "number": {
                    "description": "* Номер карты",
                    "type": "string"
                },
                "updatedAt": {
                    "description": "* Дата последнего обновления сущности в базе данных",
                    "type": "string"
                }
            }
        },
        "DatabaseServicev1.Company": {
            "type": "object",
            "properties": {
                "address": {
                    "description": "* Адрес офиса компании",
                    "type": "string"
                },
                "card": {
                    "description": "* Банковская карта компании",
                    "allOf": [
                        {
                            "$ref": "#/definitions/DatabaseServicev1.CardCompany"
                        }
                    ]
                },
                "createdAt": {
                    "description": "* Дата создания сущности в базе данных",
                    "type": "string"
                },
                "id": {
                    "description": "* ID компании в базе данных",
                    "type": "integer"
                },
                "inn": {
                    "description": "* ИНН юридического лица",
                    "type": "string"
                },
                "kpp": {
                    "description": "* КПП юридического лица",
                    "type": "string"
                },
                "okpo": {
                    "description": "* ОКПО предприятия/организации",
                    "type": "string"
                },
                "phone": {
                    "description": "* Номер телефона компании",
                    "type": "string"
                },
                "site": {
                    "description": "* Сайт компании",
                    "type": "string"
                },
                "title": {
                    "description": "* Название компании",
                    "type": "string"
                },
                "updatedAt": {
                    "description": "* Дата последнего обновления сущности в базе данных",
                    "type": "string"
                },
                "userId": {
                    "description": "* ID пользователя к которому относится данная компания",
                    "type": "integer"
                }
            }
        },
        "DatabaseServicev1.CreateUserResponse": {
            "type": "object",
            "properties": {
                "AvatarPath": {
                    "description": "* Локальный путь к аватару пользователя",
                    "type": "string"
                },
                "card": {
                    "description": "* Банковские карты пользователя",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/DatabaseServicev1.Card"
                    }
                },
                "company": {
                    "description": "* Компания пользователя, если он является юридическим лицом",
                    "allOf": [
                        {
                            "$ref": "#/definitions/DatabaseServicev1.Company"
                        }
                    ]
                },
                "createdAt": {
                    "description": "* Дата создания сущности записи в базе данных",
                    "type": "string"
                },
                "donations": {
                    "description": "* Пожертвования пользователя",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/DatabaseServicev1.Donations"
                    }
                },
                "email": {
                    "description": "* Email пользователя",
                    "type": "string"
                },
                "id": {
                    "description": "* ID пользователя",
                    "type": "integer"
                },
                "password": {
                    "description": "* Пароль пользователя (закодирован в MD5Hash)",
                    "type": "string"
                },
                "phone": {
                    "description": "* Номер телефона пользователя",
                    "type": "string"
                },
                "role": {
                    "description": "* Роль пользователя",
                    "type": "string"
                },
                "type": {
                    "description": "* Тип пользователя, 0 - физическое лицо, 1 - юридическое лицо",
                    "type": "integer"
                },
                "updatedAt": {
                    "description": "* Дата последнего обновления сущности в базе данных",
                    "type": "string"
                },
                "username": {
                    "description": "* Имя (никнейм) пользователя",
                    "type": "string"
                }
            }
        },
        "DatabaseServicev1.Donations": {
            "type": "object",
            "properties": {
                "amount": {
                    "description": "* Сумма пожертвования",
                    "type": "number"
                },
                "createdAt": {
                    "description": "* Дата создания сущности в базе данных",
                    "type": "string"
                },
                "id": {
                    "description": "* ID пожертвования в базе данных",
                    "type": "integer"
                },
                "title": {
                    "description": "* Название пожертвования (например \"На лекарства\")",
                    "type": "string"
                },
                "updatedAt": {
                    "description": "* Дата последнего обновления сущности в базе данных",
                    "type": "string"
                },
                "userId": {
                    "description": "* ID пользователя, которому принадлежит пожертвование",
                    "type": "integer"
                },
                "ward": {
                    "description": "* Подопечный этого пожертвования",
                    "allOf": [
                        {
                            "$ref": "#/definitions/DatabaseServicev1.Ward"
                        }
                    ]
                }
            }
        },
        "DatabaseServicev1.UsersResponse": {
            "type": "object",
            "properties": {
                "users": {
                    "description": "* Массив пользователей",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/DatabaseServicev1.CreateUserResponse"
                    }
                }
            }
        },
        "DatabaseServicev1.Ward": {
            "type": "object",
            "properties": {
                "AvatarPath": {
                    "description": "* Локальный путь к аватару подопечного",
                    "type": "string"
                },
                "createdAt": {
                    "description": "* Дата создания сущности в базе данных",
                    "type": "string"
                },
                "donationId": {
                    "description": "* ID пожертвования к которому относится данный подопечный",
                    "type": "integer"
                },
                "fullName": {
                    "description": "* Полное имя подопечного",
                    "type": "string"
                },
                "id": {
                    "description": "* ID подопечного в базе данных",
                    "type": "integer"
                },
                "necessary": {
                    "description": "* Необходимая сумма денег на необходимость",
                    "type": "number"
                },
                "title": {
                    "description": "* Дополнительный текст к подопечному",
                    "type": "string"
                },
                "updatedAt": {
                    "description": "* Дата последнего обновления сущности в базе данных",
                    "type": "string"
                },
                "want": {
                    "description": "* Необходимость подопечного (то в чем он нуждается, например \"Лекарства\")",
                    "type": "string"
                }
            }
        },
        "server.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
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
	Title:            "API Gateway",
	Description:      "Сервер маршрутизации",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}