# Housing Data Generator

A command-line tool for generating realistic synthetic housing data for data science and machine learning projects.

## Overview

This tool generates a CSV file containing synthetic housing data with realistic correlations between features. The data is suitable for regression analysis, feature importance studies, and training machine learning models for housing price prediction.

## Features

- Generate customizable number of housing data entries
- Select specific features to include in the output
- Realistic correlations between housing attributes
- Built-in neighborhood characteristics affecting pricing
- Random outliers for more realistic data distributions
- Configurable output file name

## Installation

### Prerequisites

- Go 1.16 or higher

### Building from source

1. Clone this repository
2. Build the executable:

```bash
go build -o house-data-gen
```

## Usage

Basic usage:

```bash
./house-data-gen
```

This will generate 100 entries with all available features and save to `house_data.csv`.

### Command Line Options

```
Usage of ./house-data-gen:
  -features string
        comma-separated list of features to include in the output.
        If not specified, all features will be included.
        Available features: age, basement, bathrooms, bedrooms, central_air, condition_score, crime_rate, distance_downtown, energy_efficiency, fireplace, garage_spaces, half_baths, hoa_fees, lot_size, neighborhood, pool, property_area, renovation_age, rooms, school_rating, square_footage, stories, taxes, walkability_score
  -n int
        number of entries to generate (default 100)
  -o string
        output file name (default "house_data.csv")

Available Features:
  age
  basement
  bathrooms
  bedrooms
  central_air
  condition_score
  crime_rate
  distance_downtown
  energy_efficiency
  fireplace
  garage_spaces
  half_baths
  hoa_fees
  lot_size
  neighborhood
  pool
  property_area
  renovation_age
  rooms
  school_rating
  square_footage
  stories
  taxes
  walkability_score
```

### Examples

Generate 500 entries:

```bash
./house-data-gen -n 500
```

Generate data with only specific features:

```bash
./house-data-gen -features "square_footage,bedrooms,bathrooms,age,neighborhood"
```

Save to a custom filename:

```bash
./house-data-gen -o my_housing_data.csv
```

Combine multiple options:

```bash
./house-data-gen -n 1000 -features "square_footage,bedrooms,bathrooms,school_rating,crime_rate" -o large_dataset.csv
```

## Output Data

The generated CSV always includes the selected features plus `price_thousands` (the house price in thousands of dollars).

### Feature Descriptions

| Feature | Description | Type | Range/Units |
|---------|-------------|------|-------------|
| age | Age of the house in years | int | 0-70+ |
| basement | Whether the house has a basement | int | 0=No, 1=Yes |
| bathrooms | Number of full bathrooms | float | 1.0+ |
| bedrooms | Number of bedrooms | int | 1-7 |
| central_air | Whether the house has central air conditioning | int | 0=No, 1=Yes |
| condition_score | Overall condition of the property | int | 1-10 (10=Excellent) |
| crime_rate | Crime rate in the area | float | 0.5-20.0 (higher=worse) |
| distance_downtown | Distance to downtown in miles | float | 0.1-25.0 |
| energy_efficiency | Energy efficiency rating | int | 1-10 (10=Excellent) |
| fireplace | Whether the house has a fireplace | int | 0=No, 1=Yes |
| garage_spaces | Number of garage parking spaces | int | 0-6 |
| half_baths | Number of half bathrooms (toilet + sink) | int | 0-3 |
| hoa_fees | Monthly homeowner association fees | float | 0-800+ |
| lot_size | Size of the property lot in acres | float | 0.1-10.0 |
| neighborhood | Name of the neighborhood | string | One of 8 predefined neighborhoods |
| pool | Whether the property has a pool | int | 0=No, 1=Yes |
| property_area | Total property area in square feet | int | Derived from lot_size |
| renovation_age | Years since last major renovation | int | 0-30+ |
| rooms | Total number of rooms (including bedrooms) | int | 2+ |
| school_rating | Rating of nearby schools | int | 1-10 (10=Excellent) |
| square_footage | Living area in square feet | int | 500-4000 |
| stories | Number of stories | float | 1.0-3.5 |
| taxes | Annual property taxes | float | Based on property value |
| walkability_score | Walkability rating of the location | int | 0-100 (100=Best) |
| price_thousands | House price in thousands of dollars | int | Generated based on features |

### Neighborhood Characteristics

The program includes 8 built-in neighborhoods, each with distinctive characteristics:

- Downtown: High prices, moderate schools, very close to downtown, higher crime, excellent walkability, high HOA fees
- Midtown: Above-average prices, good schools, close to downtown, moderate crime, very good walkability
- NorthSide: Above-average prices, very good schools, moderate distance, low crime, good walkability
- SouthSide: Below-average prices, average schools, moderate distance, moderate crime, moderate walkability
- Eastside: Average prices, average schools, moderate distance, low-moderate crime, moderate walkability
- Westside: Average prices, good schools, farther from downtown, low-moderate crime, moderate walkability
- SuburbsNorth: High prices, excellent schools, far from downtown, very low crime, poor walkability
- SuburbsSouth: Above-average prices, very good schools, far from downtown, low crime, poor walkability

## Data Generation Logic

The program generates realistic data with the following characteristics:

1. Base house prices are calculated from square footage, number of bedrooms, and bathrooms
2. Neighborhood factors significantly impact pricing and related attributes
3. Home features are realistically correlated (e.g., larger homes tend to have more bedrooms)
4. Approximately 5% of entries are outliers with unusual feature combinations
5. Property taxes are calculated as a percentage of the home value
6. Home condition and energy efficiency correlate with age and renovation status
7. Various premium features (pool, basement, etc.) add value to the property

## Use Cases

- Training and testing machine learning models for home price prediction
- Demonstrating feature importance in housing data
- Teaching data science concepts without requiring real estate data
- Benchmarking different regression algorithms
- Creating datasets with controllable characteristics

## License

MIT License
