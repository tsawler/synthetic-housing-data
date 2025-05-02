package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	initRandom()

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

	// Generate the specified number of housing entries
	neighborhoods := loadNeighborhoods()
	currentYear := getCurrentYear()

	for i := 0; i < *numEntries; i++ {
		// Generate a housing entry
		entry := generateHousingEntry(neighborhoods, currentYear)
		
		// Create a complete data row with all features
		allDataValues := []string{
			strconv.Itoa(entry.Age),                      // age
			strconv.Itoa(entry.Basement),                 // basement
			fmt.Sprintf("%.1f", entry.Bathrooms),         // bathrooms
			strconv.Itoa(entry.Bedrooms),                 // bedrooms
			strconv.Itoa(entry.CentralAir),               // central_air
			strconv.Itoa(entry.ConditionScore),           // condition_score
			fmt.Sprintf("%.1f", entry.CrimeRate),         // crime_rate
			fmt.Sprintf("%.1f", entry.DistanceDowntown),  // distance_downtown
			strconv.Itoa(entry.EnergyEfficiency),         // energy_efficiency
			strconv.Itoa(entry.Fireplace),                // fireplace
			strconv.Itoa(entry.GarageSpaces),             // garage_spaces
			strconv.Itoa(entry.HalfBaths),                // half_baths
			fmt.Sprintf("%.0f", entry.HoaFees),           // hoa_fees
			fmt.Sprintf("%.2f", entry.LotSize),           // lot_size
			entry.Neighborhood.Name,                      // neighborhood
			strconv.Itoa(entry.Pool),                     // pool
			strconv.Itoa(entry.PropertyArea),             // property_area
			strconv.Itoa(entry.RenovationAge),            // renovation_age
			strconv.Itoa(entry.Rooms),                    // rooms
			strconv.Itoa(entry.SchoolRating),             // school_rating
			strconv.Itoa(entry.SquareFootage),            // square_footage
			fmt.Sprintf("%.1f", entry.Stories),           // stories
			fmt.Sprintf("%.0f", entry.Taxes),             // taxes
			strconv.Itoa(entry.WalkabilityScore),         // walkability_score
		}

		// Create a new row with only the selected features
		selectedRow := make([]string, len(selectedIndices))
		for i, idx := range selectedIndices {
			selectedRow[i] = allDataValues[idx]
		}

		// Add price_thousands as the last column
		rowWithPrice := append(selectedRow, strconv.Itoa(entry.PriceThousands))

		// Write to CSV
		if err := writer.Write(rowWithPrice); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing data row: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Generated %d housing data entries with %d features (plus price_thousands) in %s\n",
		*numEntries, len(selectedFeatures), *outputFile)
}