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
		isOutlier := rand.Float64() < 0.05

		// Select neighborhood
		neighborhood := neighborhoods[rand.Intn(len(neighborhoods))]

		// Generate square footage
		var squareFootage int
		if isOutlier {
			// Outlier square footage: very small or very large
			if rand.Float64() < 0.5 {
				squareFootage = rand.Intn(600) + 500 // Small: 500-1100 sq ft
			} else {
				squareFootage = rand.Intn(1500) + 2500 // Large: 2500-4000 sq ft
			}
		} else {
			// Normal range: 900-2200 sq ft
			squareFootage = rand.Intn(1300) + 900
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
			if squareFootage < 1000 {
				bedrooms = rand.Intn(2) + 1 // 1-2 bedrooms for small houses
			} else if squareFootage < 1500 {
				bedrooms = rand.Intn(2) + 2 // 2-3 bedrooms for medium houses
			} else if squareFootage < 2000 {
				bedrooms = rand.Intn(2) + 3 // 3-4 bedrooms for larger houses
			} else {
				bedrooms = rand.Intn(3) + 4 // 4-6 bedrooms for very large houses
			}
		}

		// Generate bathrooms based on bedrooms and square footage
		var bathrooms float64
		if isOutlier && rand.Float64() < 0.3 {
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
		ageRange := [2]int{0, 70} // Most homes 0-70 years old
		if isOutlier && rand.Float64() < 0.3 {
			// Some outliers are historic homes
			yearBuilt = rand.Intn(50) + 1900 // Built between 1900-1950
		} else {
			yearBuilt = currentYear - rand.Intn(ageRange[1]-ageRange[0]+1) - ageRange[0]
		}
		age = currentYear - yearBuilt

		// Renovation age (years since last major renovation)
		if age < 5 {
			renovationAge = 0 // New homes don't need renovation
		} else if isOutlier && rand.Float64() < 0.3 {
			renovationAge = age // Never renovated
		} else {
			// Most homes are renovated every 10-30 years
			renovationAge = rand.Intn(min(age, 30)) + 1
		}

		// Lot size in acres (0.1 to 1.0 typical)
		lotSize := 0.1 + rand.Float64()*0.9
		if isOutlier && rand.Float64() < 0.3 {
			// Some outliers have very large lots
			lotSize = 1.0 + rand.Float64()*9.0 // 1-10 acres
		}

		// Property area in square feet (lot size * 43,560 sq ft per acre)
		propertyArea := int(lotSize * 43560)

		// Garage spaces
		var garageSpaces int
		if isOutlier && rand.Float64() < 0.3 {
			// Outlier garages
			if rand.Float64() < 0.5 {
				garageSpaces = 0 // No garage
			} else {
				garageSpaces = 4 + rand.Intn(3) // Large 4-6 car garage
			}
		} else {
			// Normal distribution of garage spaces based on house size
			if squareFootage < 1200 {
				garageSpaces = rand.Intn(2) // 0-1 spaces for small houses
			} else if squareFootage < 2000 {
				garageSpaces = 1 + rand.Intn(2) // 1-2 spaces for medium houses
			} else {
				garageSpaces = 2 + rand.Intn(2) // 2-3 spaces for large houses
			}
		}

		// Total rooms (including bedrooms, excluding bathrooms)
		var rooms int
		if isOutlier && rand.Float64() < 0.3 {
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
		if isOutlier && rand.Float64() < 0.3 {
			halfBaths = rand.Intn(4) // 0-3 half baths
		} else {
			// Typically 0-2 half baths
			if squareFootage < 1500 {
				halfBaths = rand.Intn(2) // 0-1 for smaller homes
			} else {
				halfBaths = rand.Intn(3) // 0-2 for larger homes
			}
		}

		// Number of stories
		var stories float64
		if isOutlier && rand.Float64() < 0.3 {
			// Outlier stories
			if rand.Float64() < 0.5 {
				stories = 1.0 // Single story despite large size
			} else {
				stories = 3.0 + rand.Float64() // 3+ stories
			}
		} else {
			// Normal distribution based on size
			if squareFootage < 1200 {
				stories = 1.0 // Single story for smaller homes
			} else if squareFootage < 2000 {
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
		if isOutlier && rand.Float64() < 0.3 {
			basement = rand.Intn(2) // Random for outliers
		} else {
			// More likely in older homes or larger homes
			basementProb := 0.3 + 0.01*float64(age) + 0.0002*float64(squareFootage)
			if basementProb > 0.9 {
				basementProb = 0.9
			}
			if rand.Float64() < basementProb {
				basement = 1
			}
		}

		// Pool more common in expensive neighborhoods
		if isOutlier && rand.Float64() < 0.3 {
			pool = rand.Intn(2) // Random for outliers
		} else {
			poolProb := 0.05 + 0.1*neighborhood.avgPrice + 0.0001*float64(squareFootage)
			if poolProb > 0.5 {
				poolProb = 0.5 // Pools are somewhat rare
			}
			if rand.Float64() < poolProb {
				pool = 1
			}
		}

		// Fireplace more common in older homes and larger homes
		if isOutlier && rand.Float64() < 0.3 {
			fireplace = rand.Intn(2) // Random for outliers
		} else {
			fireplaceProb := 0.2 + 0.005*float64(age) + 0.0002*float64(squareFootage)
			if fireplaceProb > 0.8 {
				fireplaceProb = 0.8
			}
			if rand.Float64() < fireplaceProb {
				fireplace = 1
			}
		}

		// Central air more common in newer homes
		if isOutlier && rand.Float64() < 0.3 {
			centralAir = rand.Intn(2) // Random for outliers
		} else {
			centralAirProb := 0.95 - 0.01*float64(age)
			if centralAirProb < 0.5 {
				centralAirProb = 0.5 // Even old homes may have central air added
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
		if isOutlier && rand.Float64() < 0.3 {
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
		if isOutlier && rand.Float64() < 0.3 {
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
		if isOutlier && rand.Float64() < 0.3 {
			// Some outliers have very high or zero HOA
			if rand.Float64() < 0.5 {
				hoaFees = 0 // No HOA
			} else {
				hoaFees = hoaRange[1] + rand.Float64()*500 // Very high HOA
			}
		}

		// Generate price based on all factors
		// Base price factor: about $120 per sq ft + $30k per bedroom + $20k per bathroom
		basePricePerSqFt := 0.12
		bedroomBonus := 30.0
		bathroomBonus := 20.0

		// Calculate base price from main factors
		basePrice := float64(squareFootage)*basePricePerSqFt +
			float64(bedrooms)*bedroomBonus +
			bathrooms*bathroomBonus

		// Apply neighborhood factor
		basePrice *= neighborhood.avgPrice

		// Additional price factors
		priceFactors := 1.0

		// Adjust for age (newer houses worth more)
		ageDiscount := 0.0
		if float64(age)*0.005 < 0.3 {
			ageDiscount = 1.0 - float64(age)*0.005
		} else {
			ageDiscount = 1.0 - 0.3
		}
		priceFactors *= ageDiscount

		// Renovation reduces age penalty
		if renovationAge < 10 {
			priceFactors *= 1.0 + (0.1 * (1.0 - float64(renovationAge)/10.0))
		}

		// Premium features
		if basement == 1 {
			priceFactors *= 1.05 // +5% for basement
		}
		if pool == 1 {
			priceFactors *= 1.07 // +7% for pool
		}
		if fireplace == 1 {
			priceFactors *= 1.03 // +3% for fireplace
		}
		if centralAir == 1 {
			priceFactors *= 1.04 // +4% for central air
		}

		// Quality and condition adjustments
		priceFactors *= 0.85 + float64(conditionScore)*0.03 // Up to +15% for condition

		// School rating premium
		priceFactors *= 1.0 + float64(schoolRating-5)*0.02 // ±10% based on schools

		// Crime rate discount
		priceFactors *= 1.0 - (crimeRate-5.0)*0.01 // ±5% based on crime

		// Lot size premium
		priceFactors *= 1.0 + (lotSize-0.5)*0.1 // Bigger lots worth more

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
			price = price * (0.9 + rand.Float64()*0.2)
		}

		// Calculate taxes based on final price (1-2% of home value)
		taxRate := 0.01 + rand.Float64()*0.01
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
