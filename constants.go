package main

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