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
var UsedContexts = []string{"Worship", "Shopping", "Work", "School", "Campus", "Sport", "Party"}
var ClothesCategories = []string{"upper_body", "bottom_body", "head", "foot", "hand"}
var ClothesTypes = []string{"hat", "pants", "shirt", "jacket", "shoes", "socks", "scarf", "gloves", "shorts", "skirt", "dress", "blouse", "sweater", "hoodie", "tie", "belt", "coat", "underwear", "swimsuit", "vest", "t-shirt", "jeans", "leggings", "boots", "sandals", "sneakers", "raincoat", "poncho", "cardigan"}
var ClothesGenders = []string{"male", "female", "unisex"}
var ClothesMadeFroms = []string{"cotton", "wool", "silk", "linen", "polyester", "denim", "leather", "nylon", "rayon", "synthetic", "cloth"}
var ClothesSizes = []string{"S", "M", "L", "XL", "XXL"}
