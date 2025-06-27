package config

var ResponseMessages = map[string]string{
	"post":        "created",
	"put":         "updated",
	"hard delete": "permanentally deleted",
	"soft delete": "deleted",
	"recover":     "recovered",
	"get":         "fetched",
	"login":       "login",
	"sign out":    "signed out",
	"empty":       "not found",
}
var DictionaryTypes = []string{
	"wash_type", "clothes_type", "clothes_category", "used_context", "clothes_gender", "clothes_made_from", "clothes_size",
}
var Days = []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
