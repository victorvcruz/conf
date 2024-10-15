# Configuration Loader

A simple Go package for loading configuration values from environment variables into struct fields using reflection. This package makes it easy to manage application configurations by defining environment variables alongside default values in struct tags.

## Overview

The `conf` package allows you to define your application's configuration in a structured way. You can specify environment variable names and default values using struct tags. If an environment variable is not set, the default value will be used.

### Features

- Load configuration from environment variables.
- Define default values for each configuration field.
- Support for various data types: strings, integers, booleans, floats, and string slices.

## Installation

To use this package, you can import it into your Go project:

```bash
go get github.com/victorvcruz/conf
```

## Usage

### Define Configuration Struct

You can define a configuration struct by specifying the environment variable name and the default value in the struct tags. Here's an example:

```go
package main

import (
	"log"
	"github.com/victorvcruz/conf"
)

type Config struct {
	DatabaseURL string   `conf:"DATABASE_URL,localhost:5432"`
	Debug       bool     `conf:"DEBUG,false"`
	Timeout     int      `conf:"TIMEOUT,30"`
	AllowedIPs  []string `conf:"ALLOWED_IPS,127.0.0.1;192.168.1.1"`
}

func main() {
	var cfg Config
	if err := conf.Load(&cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf("Loaded Config: %+v", cfg)
}
```

In the above example, if the `DATABASE_URL` environment variable is not set, the default value `'localhost:5432'` will be used. Similarly, if `DEBUG` is not set, it will default to `false`.

## Supported Data Types

- **String:** Loaded as a string value.
- **Integer:** Loaded as an integer. The default value can be specified as a string.
- **Boolean:** Loaded as a boolean. The default value can be either true or false.
- **Float:** Loaded as a float64. The default value can be specified as a string.
- **String Slice:** Loaded as a slice of strings, separated by semicolons (e.g., 127.0.0.1;192.168.1.1).