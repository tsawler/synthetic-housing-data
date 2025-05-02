package main

// Neighborhood represents a city neighborhood with its characteristics
type Neighborhood struct {
	Name          string    // Neighborhood name
	AvgPrice      float64   // Price multiplier
	SchoolRating  int       // Rating 1-10
	DistanceRange [2]float64 // Min and max distance to downtown
	CrimeRate     [2]float64 // Min and max crime rate
	WalkScore     [2]int     // Min and max walkability score
	HoaRange      [2]float64 // Min and max HOA fees
}

// HousingEntry represents a single housing record with all its features
type HousingEntry struct {
	// Basic info
	SquareFootage    int
	Bedrooms         int
	Bathrooms        float64
	Age              int
	RenovationAge    int
	YearBuilt        int
	
	// Property characteristics
	LotSize          float64
	PropertyArea     int
	GarageSpaces     int
	Rooms            int
	HalfBaths        int
	Stories          float64
	
	// Boolean features (0=No, 1=Yes)
	Basement         int
	Pool             int
	Fireplace        int
	CentralAir       int
	
	// Quality metrics
	SchoolRating     int
	DistanceDowntown float64
	CrimeRate        float64
	WalkabilityScore int
	ConditionScore   int
	EnergyEfficiency int
	
	// Financial aspects
	Taxes            float64
	HoaFees          float64
	PriceThousands   int
	
	// Reference to the neighborhood
	Neighborhood     *Neighborhood
	
	// Whether this is an outlier entry
	IsOutlier        bool
}