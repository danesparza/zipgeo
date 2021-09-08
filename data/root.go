package data

// ZipGeo defines a zipcode -> lat,long lookup result
type ZipGeo struct {
	Zipcode   int     `json:"zipcode"`   // US zipcode
	Latitude  float64 `json:"latitude"`  // Latitude
	Longitude float64 `json:"longitude"` // Longitude
	Version   string  `json:"version"`   // Service version
}
