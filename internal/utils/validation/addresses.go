package validation

import addressentity "domashka-backend/internal/entity/geo"

const (
	RussiaLatitudeMin  = 41.0  // Южная граница
	RussiaLatitudeMax  = 81.0  // Северная граница
	RussiaLongitudeMin = 19.0  // Западная граница
	RussiaLongitudeMax = 170.0 // Восточная граница
)

func IsAddressInRussia(addr addressentity.Address) bool {
	return addr.Latitude >= RussiaLatitudeMin && addr.Latitude <= RussiaLatitudeMax &&
		addr.Longitude >= RussiaLongitudeMin && addr.Longitude <= RussiaLongitudeMax
}
