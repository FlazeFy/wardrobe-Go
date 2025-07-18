package config

import "time"

// Template Message
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
	"conflict":    "already exist",
}

// Rules
var RedisTime = 10 * time.Minute

// Rules
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
var TrackSources = []string{"Web", "Mobile", "Telegram Bot", "Line Bot"}
var Colors = []string{"red", "blue", "green", "yellow", "purple", "orange", "pink", "brown", "black", "white"}
var WeatherHitFroms = []string{"Task Schedule", "Manual"}
var WeatherConditions = []string{"Thunderstorm", "Drizzle", "Rain", "Snow", "Mist", "Smoke", "Haze", "Dust", "Fog", "Sand", "Ash", "Squall", "Tornado", "Clear", "Clouds"}
var WashTypes = []string{"Laundry", "Self-Wash"}
var StatsClothesField = []string{"clothes_type", "clothes_category", "clothes_gender", "clothes_made_from", "clothes_size", "clothes_merk"}
var StatsClothesUsedsField = []string{"used_context"}
var StatsSchedulesField = []string{"day"}
var StatsWashField = []string{"wash_type"}
var StatsWeatherField = []string{"weather_condition", "weather_city"}

// Query
var QueryOrder = []string{"desc", "asc"}

// List Menu
type MenuItem struct {
	Label string
	Data  string
}

var MenuList = []MenuItem{
	{Label: "All Clothes", Data: "/Show All Clothes"},
	{Label: "Used History", Data: "/Show Used Clothes History"},
	{Label: "By Category", Data: "/Show Clothes By Category"},
	{Label: "Schedule", Data: "/Show Schedule"},
	{Label: "Wash History", Data: "/Show Wash History"},
	{Label: "Wishlist", Data: "/Show Wishlist"},
	{Label: "Apps History", Data: "/Show Apps History"},
	{Label: "All Outfit", Data: "/Show All Outfit"},
	{Label: "Use Clothes", Data: "/Used a Clothes"},
	{Label: "Random Outfit", Data: "/Randomize My Outfit"},
	{Label: "Most Used", Data: "/Show Most Used Clothes"},
	{Label: "Most Used Daily", Data: "/Show Most Used Clothes for Daily"},
	{Label: "Last History", Data: "/Show Last History"},
	{Label: "Exit Bot", Data: "/Exit Bot"},
}
