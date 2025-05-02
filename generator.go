package main

import (
	"math"
	"math/rand"
)

// generateHousingEntry creates a complete housing entry with all features
func generateHousingEntry(neighborhoods []Neighborhood, currentYear int) HousingEntry {
	entry := HousingEntry{}
	
	// Decide if this entry will be an outlier (approximately 5% chance)
	entry.IsOutlier = rand.Float64() < OUTLIER_CHANCE

	// Select neighborhood
	entry.Neighborhood = &neighborhoods[rand.Intn(len(neighborhoods))]

	// Generate square footage
	entry.SquareFootage = generateSquareFootage(entry.IsOutlier)

	// Generate number of bedrooms based on square footage
	entry.Bedrooms = generateBedrooms(entry.SquareFootage, entry.IsOutlier)

	// Generate bathrooms based on bedrooms and square footage
	entry.Bathrooms = generateBathrooms(entry.Bedrooms, entry.IsOutlier)

	// Generate year built and age
	entry.YearBuilt, entry.Age, entry.RenovationAge = generateAgeData(currentYear, entry.IsOutlier)

	// Generate lot size and property area
	entry.LotSize = generateLotSize(entry.IsOutlier)
	entry.PropertyArea = int(entry.LotSize * 43560) // 43,560 sq ft per acre

	// Generate garage spaces
	entry.GarageSpaces = generateGarageSpaces(entry.SquareFootage, entry.IsOutlier)

	// Generate total rooms
	entry.Rooms = generateRooms(entry.Bedrooms, entry.IsOutlier)

	// Generate half baths
	entry.HalfBaths = generateHalfBaths(entry.SquareFootage, entry.IsOutlier)

	// Generate number of stories
	entry.Stories = generateStories(entry.SquareFootage, entry.IsOutlier)

	// Generate boolean features
	entry.Basement = generateBasement(entry.Age, entry.SquareFootage, entry.IsOutlier)
	entry.Pool = generatePool(*entry.Neighborhood, entry.SquareFootage, entry.IsOutlier)
	entry.Fireplace = generateFireplace(entry.Age, entry.SquareFootage, entry.IsOutlier)
	entry.CentralAir = generateCentralAir(entry.Age, entry.IsOutlier)

	// Generate school rating
	entry.SchoolRating = generateSchoolRating(*entry.Neighborhood)

	// Generate distance to downtown
	entry.DistanceDowntown = generateDistanceDowntown(*entry.Neighborhood, entry.IsOutlier)

	// Generate crime rate
	entry.CrimeRate = generateCrimeRate(*entry.Neighborhood, entry.IsOutlier)

	// Generate walkability score
	entry.WalkabilityScore = generateWalkabilityScore(*entry.Neighborhood, entry.IsOutlier)

	// Generate condition score
	entry.ConditionScore = generateConditionScore(entry.Age, entry.RenovationAge, entry.IsOutlier)

	// Generate energy efficiency
	entry.EnergyEfficiency = generateEnergyEfficiency(entry.Age, entry.RenovationAge, entry.IsOutlier)

	// Calculate price
	price := calculatePrice(entry)

	// Generate HOA fees
	entry.HoaFees = generateHoaFees(*entry.Neighborhood, entry.IsOutlier)

	// Calculate taxes
	taxRate := MIN_TAX_RATE + rand.Float64()*(MAX_TAX_RATE-MIN_TAX_RATE)
	entry.Taxes = price * taxRate

	// Round price to nearest thousand
	entry.PriceThousands = int(price + 0.5)

	return entry
}

// generateSquareFootage creates a realistic square footage value
func generateSquareFootage(isOutlier bool) int {
	if isOutlier {
		// Outlier square footage: very small or very large
		if rand.Float64() < 0.5 {
			return rand.Intn(SMALL_OUTLIER_RANGE) + SMALL_OUTLIER_MIN // Small: 500-1100 sq ft
		} else {
			return rand.Intn(LARGE_OUTLIER_RANGE) + LARGE_OUTLIER_MIN // Large: 2500-4000 sq ft
		}
	} else {
		// Normal range: 900-2200 sq ft
		return rand.Intn(SMALL_HOME_RANGE) + SMALL_HOME_MIN_SQFT
	}
}

// generateBedrooms determines a realistic number of bedrooms
func generateBedrooms(squareFootage int, isOutlier bool) int {
	if isOutlier {
		// Outlier bedrooms: unusually many or few for the square footage
		if rand.Float64() < 0.5 {
			return 1 // Unusually few bedrooms
		} else {
			return rand.Intn(3) + 5 // Unusually many bedrooms (5-7)
		}
	} else {
		// Normal relationship between square footage and bedrooms
		if squareFootage < SMALL_HOME_SQFT {
			return rand.Intn(2) + 1 // 1-2 bedrooms for small houses
		} else if squareFootage < MEDIUM_HOME_SQFT {
			return rand.Intn(2) + 2 // 2-3 bedrooms for medium houses
		} else if squareFootage < 2000 {
			return rand.Intn(2) + 3 // 3-4 bedrooms for larger houses
		} else {
			return rand.Intn(3) + 4 // 4-6 bedrooms for very large houses
		}
	}
}

// generateBathrooms determines a realistic number of bathrooms
func generateBathrooms(bedrooms int, isOutlier bool) float64 {
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		// Outlier bathrooms
		if rand.Float64() < 0.5 {
			return 1.0 // Unusually few bathrooms
		} else {
			return float64(bedrooms) + 1.5 // Unusually many bathrooms
		}
	} else {
		// Normal case: roughly proportional to bedrooms
		bathroomsBase := float64(bedrooms) * 0.75
		// Add 0, 0.5, or 1 bathroom
		bathrooms := bathroomsBase + float64(rand.Intn(3))*0.5
		// Minimum 1 bathroom
		if bathrooms < 1.0 {
			return 1.0
		}
		return bathrooms
	}
}

// generateAgeData determines year built, age, and renovation age
func generateAgeData(currentYear int, isOutlier bool) (int, int, int) {
	var yearBuilt, age, renovationAge int
	ageRange := [2]int{MIN_AGE, MAX_AGE} // Most homes 0-70 years old
	
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		// Some outliers are historic homes
		yearBuilt = rand.Intn(HISTORIC_HOME_RANGE) + HISTORIC_HOME_MIN // Built between 1900-1950
	} else {
		yearBuilt = currentYear - rand.Intn(ageRange[1]-ageRange[0]+1) - ageRange[0]
	}
	
	age = currentYear - yearBuilt

	// Renovation age (years since last major renovation)
	if age < 5 {
		renovationAge = 0 // New homes don't need renovation
	} else if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		renovationAge = age // Never renovated
	} else {
		// Most homes are renovated every 10-30 years
		renovationAge = rand.Intn(min(age, 30)) + 1
	}
	
	return yearBuilt, age, renovationAge
}

// generateLotSize determines a realistic lot size in acres
func generateLotSize(isOutlier bool) float64 {
	lotSize := 0.1 + rand.Float64()*0.9 // 0.1 to 1.0 typical
	
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		// Some outliers have very large lots
		lotSize = 1.0 + rand.Float64()*9.0 // 1-10 acres
	}
	
	return lotSize
}

// generateGarageSpaces determines a realistic number of garage spaces
func generateGarageSpaces(squareFootage int, isOutlier bool) int {
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		// Outlier garages
		if rand.Float64() < 0.5 {
			return 0 // No garage
		} else {
			return 4 + rand.Intn(3) // Large 4-6 car garage
		}
	} else {
		// Normal distribution of garage spaces based on house size
		if squareFootage < SMALL_HOME_SQFT {
			return rand.Intn(2) // 0-1 spaces for small houses
		} else if squareFootage < MEDIUM_HOME_SQFT {
			return 1 + rand.Intn(2) // 1-2 spaces for medium houses
		} else {
			return 2 + rand.Intn(2) // 2-3 spaces for large houses
		}
	}
}

// generateRooms determines a realistic total number of rooms
func generateRooms(bedrooms int, isOutlier bool) int {
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		// Outlier room count
		if rand.Float64() < 0.5 {
			return bedrooms + 1 // Minimal other rooms
		} else {
			return bedrooms + 5 + rand.Intn(5) // Many extra rooms
		}
	} else {
		// Normal: bedrooms plus 2-5 other rooms (living, dining, etc.)
		return bedrooms + 2 + rand.Intn(4)
	}
}

// generateHalfBaths determines a realistic number of half bathrooms
func generateHalfBaths(squareFootage int, isOutlier bool) int {
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		return rand.Intn(4) // 0-3 half baths
	} else {
		// Typically 0-2 half baths
		if squareFootage < SMALL_HOME_SQFT {
			return rand.Intn(2) // 0-1 for smaller homes
		} else {
			return rand.Intn(3) // 0-2 for larger homes
		}
	}
}

// generateStories determines a realistic number of stories
func generateStories(squareFootage int, isOutlier bool) float64 {
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		// Outlier stories
		if rand.Float64() < 0.5 {
			return 1.0 // Single story despite large size
		} else {
			return 3.0 + rand.Float64() // 3+ stories
		}
	} else {
		// Normal distribution based on size
		var stories float64
		if squareFootage < SMALL_HOME_SQFT {
			stories = 1.0 // Single story for smaller homes
		} else if squareFootage < MEDIUM_HOME_SQFT {
			if rand.Float64() < 0.7 {
				stories = 2.0 // Two stories most common for medium homes
			} else {
				stories = 1.0
			}
		} else {
			if rand.Float64() < 0.8 {
				stories = 2.0 // Most large homes are 2 stories
			} else {
				stories = 3.0 // Some are 3 stories
			}
		}
		// Add split-level possibility
		if rand.Float64() < 0.2 && stories < 3 {
			stories += 0.5 // Split level
		}
		return stories
	}
}

// generateBasement determines whether a home has a basement
func generateBasement(age int, squareFootage int, isOutlier bool) int {
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		return rand.Intn(2) // Random for outliers
	} else {
		// More likely in older homes or larger homes
		basementProb := BASEMENT_BASE_PROB + BASEMENT_AGE_FACTOR*float64(age) + BASEMENT_SQFT_FACTOR*float64(squareFootage)
		if basementProb > BASEMENT_MAX_PROB {
			basementProb = BASEMENT_MAX_PROB
		}
		if rand.Float64() < basementProb {
			return 1
		}
		return 0
	}
}

// generatePool determines whether a home has a pool
func generatePool(neighborhood Neighborhood, squareFootage int, isOutlier bool) int {
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		return rand.Intn(2) // Random for outliers
	} else {
		poolProb := POOL_BASE_PROB + POOL_PRICE_FACTOR*neighborhood.AvgPrice + POOL_SQFT_FACTOR*float64(squareFootage)
		if poolProb > POOL_MAX_PROB {
			poolProb = POOL_MAX_PROB // Pools are somewhat rare
		}
		if rand.Float64() < poolProb {
			return 1
		}
		return 0
	}
}

// generateFireplace determines whether a home has a fireplace
func generateFireplace(age int, squareFootage int, isOutlier bool) int {
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		return rand.Intn(2) // Random for outliers
	} else {
		fireplaceProb := FIREPLACE_BASE_PROB + FIREPLACE_AGE_FACTOR*float64(age) + FIREPLACE_SQFT_FACTOR*float64(squareFootage)
		if fireplaceProb > FIREPLACE_MAX_PROB {
			fireplaceProb = FIREPLACE_MAX_PROB
		}
		if rand.Float64() < fireplaceProb {
			return 1
		}
		return 0
	}
}

// generateCentralAir determines whether a home has central air conditioning
func generateCentralAir(age int, isOutlier bool) int {
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		return rand.Intn(2) // Random for outliers
	} else {
		centralAirProb := CENTRAL_AIR_BASE_PROB - CENTRAL_AIR_AGE_FACTOR*float64(age)
		if centralAirProb < CENTRAL_AIR_MIN_PROB {
			centralAirProb = CENTRAL_AIR_MIN_PROB // Even old homes may have central air added
		}
		if rand.Float64() < centralAirProb {
			return 1
		}
		return 0
	}
}

// generateSchoolRating determines a realistic school rating for the area
func generateSchoolRating(neighborhood Neighborhood) int {
	schoolRating := neighborhood.SchoolRating
	if rand.Float64() < 0.3 {
		// Add some variation to school rating
		schoolRating += rand.Intn(3) - 1 // -1, 0, or +1
		if schoolRating < 1 {
			schoolRating = 1
		} else if schoolRating > 10 {
			schoolRating = 10
		}
	}
	return schoolRating
}

// generateDistanceDowntown determines a realistic distance to downtown
func generateDistanceDowntown(neighborhood Neighborhood, isOutlier bool) float64 {
	distRange := neighborhood.DistanceRange
	distanceDowntown := distRange[0] + rand.Float64()*(distRange[1]-distRange[0])
	
	if isOutlier && rand.Float64() < 0.2 {
		// Some outliers are unusually far or close
		distanceDowntown = rand.Float64() * 25.0
	}
	
	return distanceDowntown
}

// generateCrimeRate determines a realistic crime rate for the area
func generateCrimeRate(neighborhood Neighborhood, isOutlier bool) float64 {
	crimeRange := neighborhood.CrimeRate
	crimeRate := crimeRange[0] + rand.Float64()*(crimeRange[1]-crimeRange[0])
	
	if isOutlier && rand.Float64() < 0.2 {
		// Some outliers have unusual crime rates
		crimeRate = rand.Float64() * 20.0
	}
	
	return crimeRate
}

// generateWalkabilityScore determines a realistic walkability score for the area
func generateWalkabilityScore(neighborhood Neighborhood, isOutlier bool) int {
	walkRange := neighborhood.WalkScore
	walkabilityScore := walkRange[0] + rand.Intn(walkRange[1]-walkRange[0]+1)
	
	if isOutlier && rand.Float64() < 0.2 {
		// Some outliers have unusual walkability
		walkabilityScore = rand.Intn(101) // 0-100
	}
	
	return walkabilityScore
}

// generateConditionScore determines a realistic condition score for the house
func generateConditionScore(age int, renovationAge int, isOutlier bool) int {
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		return rand.Intn(11) // 0-10 for outliers
	} else {
		// Better condition for newer homes or recently renovated homes
		baseCondition := 10 - int(math.Min(float64(age), float64(renovationAge))*0.2)
		if baseCondition < 3 {
			baseCondition = 3 // Even old homes not too terrible
		}
		// Add some randomness
		conditionScore := baseCondition + rand.Intn(3) - 1
		if conditionScore < 1 {
			conditionScore = 1
		} else if conditionScore > 10 {
			conditionScore = 10
		}
		return conditionScore
	}
}

// generateEnergyEfficiency determines a realistic energy efficiency score
func generateEnergyEfficiency(age int, renovationAge int, isOutlier bool) int {
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		return rand.Intn(11) // 0-10 for outliers
	} else {
		// Better efficiency for newer homes
		baseEfficiency := 10 - int(float64(age)*0.15)
		if baseEfficiency < 4 {
			baseEfficiency = 4 // Even old homes not terrible
		}
		// Better if recently renovated
		if renovationAge < 10 {
			baseEfficiency += 2
		}
		// Add some randomness
		energyEfficiency := baseEfficiency + rand.Intn(3) - 1
		if energyEfficiency < 1 {
			energyEfficiency = 1
		} else if energyEfficiency > 10 {
			energyEfficiency = 10
		}
		return energyEfficiency
	}
}

// generateHoaFees determines realistic HOA fees for the property
func generateHoaFees(neighborhood Neighborhood, isOutlier bool) float64 {
	hoaRange := neighborhood.HoaRange
	hoaFees := hoaRange[0] + rand.Float64()*(hoaRange[1]-hoaRange[0])
	
	if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
		// Some outliers have very high or zero HOA
		if rand.Float64() < 0.5 {
			hoaFees = 0 // No HOA
		} else {
			hoaFees = hoaRange[1] + rand.Float64()*500 // Very high HOA
		}
	}
	
	return hoaFees
}

// calculatePrice determines a realistic price for the house based on all features
func calculatePrice(entry HousingEntry) float64 {
	// Calculate base price from main factors
	basePrice := float64(entry.SquareFootage)*BASE_PRICE_PER_SQFT +
		float64(entry.Bedrooms)*BEDROOM_BONUS +
		entry.Bathrooms*BATHROOM_BONUS

	// Apply neighborhood factor
	basePrice *= entry.Neighborhood.AvgPrice

	// Additional price factors
	priceFactors := 1.0

	// Adjust for age (newer houses worth more)
	ageDiscount := 0.0
	if float64(entry.Age)*AGE_DISCOUNT_RATE < MAX_AGE_DISCOUNT {
		ageDiscount = 1.0 - float64(entry.Age)*AGE_DISCOUNT_RATE
	} else {
		ageDiscount = 1.0 - MAX_AGE_DISCOUNT
	}
	priceFactors *= ageDiscount

	// Renovation reduces age penalty
	if entry.RenovationAge < 10 {
		priceFactors *= 1.0 + (RENOVATION_BONUS_MAX * (1.0 - float64(entry.RenovationAge)/10.0))
	}

	// Premium features
	if entry.Basement == 1 {
		priceFactors *= BASEMENT_PREMIUM
	}
	if entry.Pool == 1 {
		priceFactors *= POOL_PREMIUM
	}
	if entry.Fireplace == 1 {
		priceFactors *= FIREPLACE_PREMIUM
	}
	if entry.CentralAir == 1 {
		priceFactors *= CENTRAL_AIR_PREMIUM
	}

	// Quality and condition adjustments
	priceFactors *= MIN_CONDITION_FACTOR + float64(entry.ConditionScore)*CONDITION_FACTOR_RATE

	// School rating premium
	priceFactors *= 1.0 + float64(entry.SchoolRating-SCHOOL_RATING_MIDPOINT)*SCHOOL_FACTOR_RATE

	// Crime rate discount
	priceFactors *= 1.0 - (entry.CrimeRate-CRIME_MIDPOINT)*CRIME_FACTOR_RATE

	// Lot size premium
	priceFactors *= 1.0 + (entry.LotSize-AVG_LOT_SIZE)*LOT_SIZE_FACTOR_RATE

	// Apply factors to base price
	price := basePrice * priceFactors

	// For outliers, make price unusually high or low
	if entry.IsOutlier && rand.Float64() < 0.4 {
		if rand.Float64() < 0.5 {
			// Unusually low price (60%-80% of expected)
			price = price * (0.6 + rand.Float64()*0.2)
		} else {
			// Unusually high price (120%-180% of expected)
			price = price * (1.2 + rand.Float64()*0.6)
		}
	} else {
		// Add some normal variation (90%-110% of expected)
		price = price * (NORMAL_PRICE_MIN + rand.Float64()*NORMAL_PRICE_VAR)
	}

	return price
}