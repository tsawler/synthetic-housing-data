package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Pricing constants
const (
	BASE_PRICE_PER_SQFT = 0.12 // Base price per square foot in thousands
	BEDROOM_BONUS       = 30.0 // Value added per bedroom in thousands
	BATHROOM_BONUS      = 20.0 // Value added per bathroom in thousands
)

// Age related constants
const (
	AGE_DISCOUNT_RATE     = 0.005 // Rate at which age reduces value (0.5% per year)
	MAX_AGE_DISCOUNT      = 0.3   // Maximum discount due to age (30%)
	RENOVATION_BONUS_MAX  = 0.1   // Maximum bonus for recent renovation (10%)
)

// Feature premium constants
const (
	BASEMENT_PREMIUM    = 1.05 // +5% for basement
	POOL_PREMIUM        = 1.07 // +7% for pool
	FIREPLACE_PREMIUM   = 1.03 // +3% for fireplace
	CENTRAL_AIR_PREMIUM = 1.04 // +4% for central air
)

// Quality adjustment constants
const (
	MIN_CONDITION_FACTOR   = 0.85  // Base factor for condition score
	CONDITION_FACTOR_RATE  = 0.03  // Rate at which condition score improves value
	SCHOOL_RATING_MIDPOINT = 5     // Midpoint for school ratings
	SCHOOL_FACTOR_RATE     = 0.02  // Rate at which school rating affects value
	CRIME_MIDPOINT         = 5.0   // Midpoint for crime rate
	CRIME_FACTOR_RATE      = 0.01  // Rate at which crime affects value
	AVG_LOT_SIZE           = 0.5   // Average lot size in acres
	LOT_SIZE_FACTOR_RATE   = 0.1   // Rate at which lot size affects value
)

// Normal variation constants
const (
	NORMAL_PRICE_MIN       = 0.9  // Minimum random price multiplier (90%)
	NORMAL_PRICE_VAR       = 0.2  // Price variation range (up to +20%)
	OUTLIER_CHANCE         = 0.05 // 5% chance of an outlier
	OUTLIER_FEATURE_CHANCE = 0.3  // 30% chance of outlier having unusual feature
)

// Tax related constants
const (
	MIN_TAX_RATE = 0.01 // Minimum property tax rate
	MAX_TAX_RATE = 0.02 // Maximum property tax rate
)

// Home age related constants
const (
	MIN_AGE            = 0  // Minimum home age
	MAX_AGE            = 70 // Maximum typical home age
	HISTORIC_HOME_MIN  = 1900 // Minimum year for historic homes
	HISTORIC_HOME_RANGE = 50  // Range of years for historic homes
)

// Home feature thresholds
const (
	SMALL_HOME_SQFT    = 1200 // Threshold for small homes
	MEDIUM_HOME_SQFT   = 2000 // Threshold for medium homes
	SMALL_HOME_MIN_SQFT = 900  // Minimum square footage for normal homes
	SMALL_HOME_RANGE   = 1300 // Range of square footage for normal homes
	
	SMALL_OUTLIER_MIN  = 500  // Minimum square footage for small outliers
	SMALL_OUTLIER_RANGE = 600 // Range for small outliers
	LARGE_OUTLIER_MIN  = 2500 // Minimum square footage for large outliers
	LARGE_OUTLIER_RANGE = 1500 // Range for large outliers
)

// Probability constants
const (
	CENTRAL_AIR_BASE_PROB = 0.95 // Base probability of central air
	CENTRAL_AIR_AGE_FACTOR = 0.01 // How much age reduces central air probability
	CENTRAL_AIR_MIN_PROB = 0.5    // Minimum probability of central air
	
	BASEMENT_BASE_PROB = 0.3     // Base probability of basement
	BASEMENT_AGE_FACTOR = 0.01   // How much age increases basement probability
	BASEMENT_SQFT_FACTOR = 0.0002 // How much square footage increases basement probability
	BASEMENT_MAX_PROB = 0.9      // Maximum probability of basement
	
	POOL_BASE_PROB = 0.05       // Base probability of pool
	POOL_PRICE_FACTOR = 0.1     // How neighborhood price affects pool probability
	POOL_SQFT_FACTOR = 0.0001   // How square footage affects pool probability
	POOL_MAX_PROB = 0.5         // Maximum probability of pool
	
	FIREPLACE_BASE_PROB = 0.2     // Base probability of fireplace
	FIREPLACE_AGE_FACTOR = 0.005  // How age affects fireplace probability
	FIREPLACE_SQFT_FACTOR = 0.0002 // How square footage affects fireplace probability
	FIREPLACE_MAX_PROB = 0.8      // Maximum probability of fireplace
)

func main() {
	// Define command line flags
	numEntries := flag.Int("n", 100, "number of entries to generate")
	outputFile := flag.String("o", "house_data.csv", "output file name")

	// Define all features (excluding price_thousands which will always be included)
	allFeatures := []string{
		"age",
		"basement",
		"bathrooms",
		"bedrooms",
		"central_air",
		"condition_score",
		"crime_rate",
		"distance_downtown",
		"energy_efficiency",
		"fireplace",
		"garage_spaces",
		"half_baths",
		"hoa_fees",
		"lot_size",
		"neighborhood",
		"pool",
		"property_area",
		"renovation_age",
		"rooms",
		"school_rating",
		"square_footage",
		"stories",
		"taxes",
		"walkability_score",
	}

	// Create a comma-separated list of all features for help text
	allFeaturesStr := strings.Join(allFeatures, ", ")

	// Add a new flag for features
	featuresFlag := flag.String("features", "",
		fmt.Sprintf("comma-separated list of features to include in the output. "+
			"If not specified, all features will be included.\n"+
			"Available features: %s", allFeaturesStr))

	// Custom usage function to improve help message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nAvailable Features:\n")
		for _, feature := range allFeatures {
			fmt.Fprintf(os.Stderr, "  %s\n", feature)
		}
	}

	flag.Parse()

	// Determine which features to include
	var selectedFeatures []string
	var selectedIndices []int

	if *featuresFlag == "" {
		// If no features specified, use all features
		selectedFeatures = allFeatures
		selectedIndices = make([]int, len(allFeatures))
		for i := range selectedIndices {
			selectedIndices[i] = i
		}
	} else {
		// Parse the features flag
		requestedFeatures := strings.Split(*featuresFlag, ",")

		// Trim whitespace from each feature name
		for i, feature := range requestedFeatures {
			requestedFeatures[i] = strings.TrimSpace(feature)
		}

		// Validate and collect selected features
		for _, requested := range requestedFeatures {
			found := false
			for i, available := range allFeatures {
				if requested == available {
					selectedFeatures = append(selectedFeatures, requested)
					selectedIndices = append(selectedIndices, i)
					found = true
					break
				}
			}
			if !found {
				fmt.Fprintf(os.Stderr, "Warning: Unknown feature '%s' - will be ignored\n", requested)
			}
		}

		// Check if any valid features were selected
		if len(selectedFeatures) == 0 {
			fmt.Fprintf(os.Stderr, "Error: No valid features selected. Use -h for help.\n")
			os.Exit(1)
		}
	}

	// Initialize random number generator
	rand.Seed(time.Now().UnixNano())

	// Create output file
	file, err := os.Create(*outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Prepare header with selected features plus price_thousands
	headerWithPrice := append(selectedFeatures, "price_thousands")

	// Write header
	if err := writer.Write(headerWithPrice); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing header: %v\n", err)
		os.Exit(1)
	}

	// Define neighborhood characteristics for consistency
	neighborhoods := []struct {
		name          string
		avgPrice      float64    // price multiplier
		schoolRating  int        // 1-10
		distanceRange [2]float64 // min and max distance to downtown
		crimeRate     [2]float64 // min and max crime rate
		walkScore     [2]int     // min and max walkability score
		hoaRange      [2]float64 // min and max HOA fees
	}{
		{"Downtown", 1.5, 6, [2]float64{0.1, 2.0}, [2]float64{5.0, 15.0}, [2]int{85, 100}, [2]float64{200, 600}},
		{"Midtown", 1.2, 7, [2]float64{2.0, 5.0}, [2]float64{3.0, 8.0}, [2]int{70, 90}, [2]float64{150, 400}},
		{"NorthSide", 1.3, 8, [2]float64{5.0, 10.0}, [2]float64{1.0, 4.0}, [2]int{50, 75}, [2]float64{100, 300}},
		{"SouthSide", 0.9, 5, [2]float64{4.0, 9.0}, [2]float64{4.0, 12.0}, [2]int{40, 65}, [2]float64{75, 250}},
		{"Eastside", 1.1, 6, [2]float64{6.0, 12.0}, [2]float64{2.0, 7.0}, [2]int{30, 60}, [2]float64{50, 200}},
		{"Westside", 1.0, 7, [2]float64{7.0, 15.0}, [2]float64{2.0, 6.0}, [2]int{25, 55}, [2]float64{25, 150}},
		{"SuburbsNorth", 1.4, 9, [2]float64{12.0, 20.0}, [2]float64{0.5, 2.0}, [2]int{15, 40}, [2]float64{150, 350}},
		{"SuburbsSouth", 1.25, 8, [2]float64{10.0, 18.0}, [2]float64{1.0, 3.0}, [2]int{10, 35}, [2]float64{125, 300}},
	}

	// Current year for age calculations
	// Get the current local time
	currentTime := time.Now()

	// Extract the year as an integer
	currentYear := currentTime.Year()

	// Generate and write data rows
	for i := 0; i < *numEntries; i++ {
		// Decide if this entry will be an outlier (approximately 5% chance)
		isOutlier := rand.Float64() < OUTLIER_CHANCE

		// Select neighborhood
		neighborhood := neighborhoods[rand.Intn(len(neighborhoods))]

		// Generate square footage
		var squareFootage int
		if isOutlier {
			// Outlier square footage: very small or very large
			if rand.Float64() < 0.5 {
				squareFootage = rand.Intn(SMALL_OUTLIER_RANGE) + SMALL_OUTLIER_MIN // Small: 500-1100 sq ft
			} else {
				squareFootage = rand.Intn(LARGE_OUTLIER_RANGE) + LARGE_OUTLIER_MIN // Large: 2500-4000 sq ft
			}
		} else {
			// Normal range: 900-2200 sq ft
			squareFootage = rand.Intn(SMALL_HOME_RANGE) + SMALL_HOME_MIN_SQFT
		}

		// Generate number of bedrooms based on square footage
		var bedrooms int
		if isOutlier {
			// Outlier bedrooms: unusually many or few for the square footage
			if rand.Float64() < 0.5 {
				bedrooms = 1 // Unusually few bedrooms
			} else {
				bedrooms = rand.Intn(3) + 5 // Unusually many bedrooms (5-7)
			}
		} else {
			// Normal relationship between square footage and bedrooms
			if squareFootage < SMALL_HOME_SQFT {
				bedrooms = rand.Intn(2) + 1 // 1-2 bedrooms for small houses
			} else if squareFootage < MEDIUM_HOME_SQFT {
				bedrooms = rand.Intn(2) + 2 // 2-3 bedrooms for medium houses
			} else if squareFootage < 2000 {
				bedrooms = rand.Intn(2) + 3 // 3-4 bedrooms for larger houses
			} else {
				bedrooms = rand.Intn(3) + 4 // 4-6 bedrooms for very large houses
			}
		}

		// Generate bathrooms based on bedrooms and square footage
		var bathrooms float64
		if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
			// Outlier bathrooms
			if rand.Float64() < 0.5 {
				bathrooms = 1.0 // Unusually few bathrooms
			} else {
				bathrooms = float64(bedrooms) + 1.5 // Unusually many bathrooms
			}
		} else {
			// Normal case: roughly proportional to bedrooms
			bathroomsBase := float64(bedrooms) * 0.75
			// Add 0, 0.5, or 1 bathroom
			bathrooms = bathroomsBase + float64(rand.Intn(3))*0.5
			// Minimum 1 bathroom
			if bathrooms < 1.0 {
				bathrooms = 1.0
			}
		}

		// Year built and age
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

		// Lot size in acres (0.1 to 1.0 typical)
		lotSize := 0.1 + rand.Float64()*0.9
		if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
			// Some outliers have very large lots
			lotSize = 1.0 + rand.Float64()*9.0 // 1-10 acres
		}

		// Property area in square feet (lot size * 43,560 sq ft per acre)
		propertyArea := int(lotSize * 43560)

		// Garage spaces
		var garageSpaces int
		if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
			// Outlier garages
			if rand.Float64() < 0.5 {
				garageSpaces = 0 // No garage
			} else {
				garageSpaces = 4 + rand.Intn(3) // Large 4-6 car garage
			}
		} else {
			// Normal distribution of garage spaces based on house size
			if squareFootage < SMALL_HOME_SQFT {
				garageSpaces = rand.Intn(2) // 0-1 spaces for small houses
			} else if squareFootage < MEDIUM_HOME_SQFT {
				garageSpaces = 1 + rand.Intn(2) // 1-2 spaces for medium houses
			} else {
				garageSpaces = 2 + rand.Intn(2) // 2-3 spaces for large houses
			}
		}

		// Total rooms (including bedrooms, excluding bathrooms)
		var rooms int
		if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
			// Outlier room count
			if rand.Float64() < 0.5 {
				rooms = bedrooms + 1 // Minimal other rooms
			} else {
				rooms = bedrooms + 5 + rand.Intn(5) // Many extra rooms
			}
		} else {
			// Normal: bedrooms plus 2-5 other rooms (living, dining, etc.)
			rooms = bedrooms + 2 + rand.Intn(4)
		}

		// Half baths
		var halfBaths int
		if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
			halfBaths = rand.Intn(4) // 0-3 half baths
		} else {
			// Typically 0-2 half baths
			if squareFootage < SMALL_HOME_SQFT {
				halfBaths = rand.Intn(2) // 0-1 for smaller homes
			} else {
				halfBaths = rand.Intn(3) // 0-2 for larger homes
			}
		}

		// Number of stories
		var stories float64
		if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
			// Outlier stories
			if rand.Float64() < 0.5 {
				stories = 1.0 // Single story despite large size
			} else {
				stories = 3.0 + rand.Float64() // 3+ stories
			}
		} else {
			// Normal distribution based on size
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
		}

		// Boolean features (0=No, 1=Yes)
		var basement, pool, fireplace, centralAir int

		// Basement more common in older homes and certain neighborhoods
		if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
			basement = rand.Intn(2) // Random for outliers
		} else {
			// More likely in older homes or larger homes
			basementProb := BASEMENT_BASE_PROB + BASEMENT_AGE_FACTOR*float64(age) + BASEMENT_SQFT_FACTOR*float64(squareFootage)
			if basementProb > BASEMENT_MAX_PROB {
				basementProb = BASEMENT_MAX_PROB
			}
			if rand.Float64() < basementProb {
				basement = 1
			}
		}

		// Pool more common in expensive neighborhoods
		if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
			pool = rand.Intn(2) // Random for outliers
		} else {
			poolProb := POOL_BASE_PROB + POOL_PRICE_FACTOR*neighborhood.avgPrice + POOL_SQFT_FACTOR*float64(squareFootage)
			if poolProb > POOL_MAX_PROB {
				poolProb = POOL_MAX_PROB // Pools are somewhat rare
			}
			if rand.Float64() < poolProb {
				pool = 1
			}
		}

		// Fireplace more common in older homes and larger homes
		if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
			fireplace = rand.Intn(2) // Random for outliers
		} else {
			fireplaceProb := FIREPLACE_BASE_PROB + FIREPLACE_AGE_FACTOR*float64(age) + FIREPLACE_SQFT_FACTOR*float64(squareFootage)
			if fireplaceProb > FIREPLACE_MAX_PROB {
				fireplaceProb = FIREPLACE_MAX_PROB
			}
			if rand.Float64() < fireplaceProb {
				fireplace = 1
			}
		}

		// Central air more common in newer homes
		if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
			centralAir = rand.Intn(2) // Random for outliers
		} else {
			centralAirProb := CENTRAL_AIR_BASE_PROB - CENTRAL_AIR_AGE_FACTOR*float64(age)
			if centralAirProb < CENTRAL_AIR_MIN_PROB {
				centralAirProb = CENTRAL_AIR_MIN_PROB // Even old homes may have central air added
			}
			if rand.Float64() < centralAirProb {
				centralAir = 1
			}
		}

		// School rating from neighborhood with slight variation
		schoolRating := neighborhood.schoolRating
		if rand.Float64() < 0.3 {
			// Add some variation to school rating
			schoolRating += rand.Intn(3) - 1 // -1, 0, or +1
			if schoolRating < 1 {
				schoolRating = 1
			} else if schoolRating > 10 {
				schoolRating = 10
			}
		}

		// Distance to downtown based on neighborhood
		distRange := neighborhood.distanceRange
		distanceDowntown := distRange[0] + rand.Float64()*(distRange[1]-distRange[0])
		if isOutlier && rand.Float64() < 0.2 {
			// Some outliers are unusually far or close
			distanceDowntown = rand.Float64() * 25.0
		}

		// Crime rate based on neighborhood
		crimeRange := neighborhood.crimeRate
		crimeRate := crimeRange[0] + rand.Float64()*(crimeRange[1]-crimeRange[0])
		if isOutlier && rand.Float64() < 0.2 {
			// Some outliers have unusual crime rates
			crimeRate = rand.Float64() * 20.0
		}

		// Walkability score based on neighborhood
		walkRange := neighborhood.walkScore
		walkabilityScore := walkRange[0] + rand.Intn(walkRange[1]-walkRange[0]+1)
		if isOutlier && rand.Float64() < 0.2 {
			// Some outliers have unusual walkability
			walkabilityScore = rand.Intn(101) // 0-100
		}

		// Condition score (1-10)
		var conditionScore int
		if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
			conditionScore = rand.Intn(11) // 0-10 for outliers
		} else {
			// Better condition for newer homes or recently renovated homes
			baseCondition := 10 - int(math.Min(float64(age), float64(renovationAge))*0.2)
			if baseCondition < 3 {
				baseCondition = 3 // Even old homes not too terrible
			}
			// Add some randomness
			conditionScore = baseCondition + rand.Intn(3) - 1
			if conditionScore < 1 {
				conditionScore = 1
			} else if conditionScore > 10 {
				conditionScore = 10
			}
		}

		// Energy efficiency (1-10)
		var energyEfficiency int
		if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
			energyEfficiency = rand.Intn(11) // 0-10 for outliers
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
			energyEfficiency = baseEfficiency + rand.Intn(3) - 1
			if energyEfficiency < 1 {
				energyEfficiency = 1
			} else if energyEfficiency > 10 {
				energyEfficiency = 10
			}
		}

		// Property taxes (roughly 1-2% of home value)
		var taxes float64
		// Will calculate after determining price

		// HOA fees
		var hoaFees float64
		hoaRange := neighborhood.hoaRange
		hoaFees = hoaRange[0] + rand.Float64()*(hoaRange[1]-hoaRange[0])
		if isOutlier && rand.Float64() < OUTLIER_FEATURE_CHANCE {
			// Some outliers have very high or zero HOA
			if rand.Float64() < 0.5 {
				hoaFees = 0 // No HOA
			} else {
				hoaFees = hoaRange[1] + rand.Float64()*500 // Very high HOA
			}
		}

		// Calculate base price from main factors
		basePrice := float64(squareFootage)*BASE_PRICE_PER_SQFT +
			float64(bedrooms)*BEDROOM_BONUS +
			bathrooms*BATHROOM_BONUS

		// Apply neighborhood factor
		basePrice *= neighborhood.avgPrice

		// Additional price factors
		priceFactors := 1.0

		// Adjust for age (newer houses worth more)
		ageDiscount := 0.0
		if float64(age)*AGE_DISCOUNT_RATE < MAX_AGE_DISCOUNT {
			ageDiscount = 1.0 - float64(age)*AGE_DISCOUNT_RATE
		} else {
			ageDiscount = 1.0 - MAX_AGE_DISCOUNT
		}
		priceFactors *= ageDiscount

		// Renovation reduces age penalty
		if renovationAge < 10 {
			priceFactors *= 1.0 + (RENOVATION_BONUS_MAX * (1.0 - float64(renovationAge)/10.0))
		}

		// Premium features
		if basement == 1 {
			priceFactors *= BASEMENT_PREMIUM
		}
		if pool == 1 {
			priceFactors *= POOL_PREMIUM
		}
		if fireplace == 1 {
			priceFactors *= FIREPLACE_PREMIUM
		}
		if centralAir == 1 {
			priceFactors *= CENTRAL_AIR_PREMIUM
		}

		// Quality and condition adjustments
		priceFactors *= MIN_CONDITION_FACTOR + float64(conditionScore)*CONDITION_FACTOR_RATE

		// School rating premium
		priceFactors *= 1.0 + float64(schoolRating-SCHOOL_RATING_MIDPOINT)*SCHOOL_FACTOR_RATE

		// Crime rate discount
		priceFactors *= 1.0 - (crimeRate-CRIME_MIDPOINT)*CRIME_FACTOR_RATE

		// Lot size premium
		priceFactors *= 1.0 + (lotSize-AVG_LOT_SIZE)*LOT_SIZE_FACTOR_RATE

		// Apply factors to base price
		price := basePrice * priceFactors

		// For outliers, make price unusually high or low
		if isOutlier && rand.Float64() < 0.4 {
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

		// Calculate taxes based on final price (1-2% of home value)
		taxRate := MIN_TAX_RATE + rand.Float64()*(MAX_TAX_RATE-MIN_TAX_RATE)
		taxes = price * taxRate

		// Round price to nearest thousand
		priceThousands := int(price + 0.5)

		// Create a complete data row with all features
		allDataValues := []string{
			strconv.Itoa(age),                     // age
			strconv.Itoa(basement),                // basement
			fmt.Sprintf("%.1f", bathrooms),        // bathrooms
			strconv.Itoa(bedrooms),                // bedrooms
			strconv.Itoa(centralAir),              // central_air
			strconv.Itoa(conditionScore),          // condition_score
			fmt.Sprintf("%.1f", crimeRate),        // crime_rate
			fmt.Sprintf("%.1f", distanceDowntown), // distance_downtown
			strconv.Itoa(energyEfficiency),        // energy_efficiency
			strconv.Itoa(fireplace),               // fireplace
			strconv.Itoa(garageSpaces),            // garage_spaces
			strconv.Itoa(halfBaths),               // half_baths
			fmt.Sprintf("%.0f", hoaFees),          // hoa_fees
			fmt.Sprintf("%.2f", lotSize),          // lot_size
			neighborhood.name,                     // neighborhood
			strconv.Itoa(pool),                    // pool
			strconv.Itoa(propertyArea),            // property_area
			strconv.Itoa(renovationAge),           // renovation_age
			strconv.Itoa(rooms),                   // rooms
			strconv.Itoa(schoolRating),            // school_rating
			strconv.Itoa(squareFootage),           // square_footage
			fmt.Sprintf("%.1f", stories),          // stories
			fmt.Sprintf("%.0f", taxes),            // taxes
			strconv.Itoa(walkabilityScore),        // walkability_score
		}

		// Create a new row with only the selected features
		selectedRow := make([]string, len(selectedIndices))
		for i, idx := range selectedIndices {
			selectedRow[i] = allDataValues[idx]
		}

		// Add price_thousands as the last column
		rowWithPrice := append(selectedRow, strconv.Itoa(priceThousands))

		// Write to CSV
		if err := writer.Write(rowWithPrice); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing data row: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Generated %d housing data entries with %d features (plus price_thousands) in %s\n",
		*numEntries, len(selectedFeatures), *outputFile)
}
