package conf

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// InvalidLoadError represents an error that occurs when trying to load configurations
// for an invalid type.
type InvalidLoadError struct {
	Type reflect.Type
}

// Error returns a descriptive message about the load error.
func (e *InvalidLoadError) Error() string {
	return fmt.Sprintf("conf: invalid load for type: %s", e.Type.String())
}

// Load populates the provided structure with configuration values
// obtained from environment variables.
//
// The function expects a pointer to a struct as an argument,
// and it uses field tags "conf" to map environment variables
// to fields in the struct. If the environment variable is not
// set, the default value specified in the tag is used.
//
// Example:
//
//	type Config struct {
//	    DatabaseURL string `conf:"DATABASE_URL,localhost:5432"`
//	    Debug       bool   `conf:"DEBUG,false"`
//	}
//
//	var cfg Config
//	if err := conf.Load(&cfg); err != nil {
//	    log.Fatalf("failed to load config: %v", err)
//	}
//
// In this example, if the environment variable DATABASE_URL is not set,
// the default value 'localhost:5432' will be used. If the variable
// DEBUG is not set, 'false' will be the value used.
//
// The structure must be a pointer so that the values can be
// set correctly. Any error loading the environment variables
// will result in an error returned by the function.
func Load(v any) error {
	if err := validateInput(v); err != nil {
		return err
	}

	return load(reflect.ValueOf(v).Elem())
}

// validateInput checks if the provided value is a non-nil pointer
// to a struct. If not, it returns an appropriate error.
func validateInput(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidLoadError{reflect.TypeOf(v)}
	}
	if rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("conf: expected a pointer to a struct, got: %s", rv.Elem().Kind())
	}
	return nil
}

// load iterates over the fields of the structure and loads the values
// from environment variables. If a field is a nested struct,
// the function is called recursively.
func load(v reflect.Value) error {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			if err := load(field); err != nil {
				return err
			}
			continue
		}

		confTag := fieldType.Tag.Get("conf")
		if confTag == "" {
			continue
		}

		parts := strings.Split(confTag, ",")
		envVar := parts[0]
		defaultValue := ""
		if len(parts) > 1 {
			defaultValue = parts[1]
		}

		// Get the environment variable value
		envValue := os.Getenv(envVar)
		if envValue == "" {
			envValue = defaultValue
		}

		if err := setField(field, envValue); err != nil {
			return err
		}
	}

	return nil
}

// setField sets the value of the field based on its type.
func setField(field reflect.Value, envValue string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(envValue)
	case reflect.Int:
		intValue, err := parseInt(envValue)
		if err != nil {
			return err
		}
		field.SetInt(int64(intValue))
	case reflect.Bool:
		boolValue, err := parseBool(envValue)
		if err != nil {
			return err
		}
		field.SetBool(boolValue)
	case reflect.Float64:
		floatValue, err := parseFloat(envValue)
		if err != nil {
			return err
		}
		field.SetFloat(floatValue)
	case reflect.Slice:
		if field.Type().Elem().Kind() == reflect.String {
			field.Set(reflect.ValueOf(strings.Split(envValue, ";")))
		} else {
			return fmt.Errorf("unsupported slice type: %s", field.Type())
		}
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}

// parseInt converts a string to an integer, returning an error if the conversion fails.
func parseInt(value string) (int, error) {
	var intValue int
	_, err := fmt.Sscanf(value, "%d", &intValue)
	return intValue, err
}

// parseBool converts a string to a boolean, returning an error if the conversion fails.
func parseBool(value string) (bool, error) {
	var boolValue bool
	_, err := fmt.Sscanf(value, "%t", &boolValue)
	return boolValue, err
}

// parseFloat converts a string to a float64, returning an error if the conversion fails.
func parseFloat(value string) (float64, error) {
	var floatValue float64
	_, err := fmt.Sscanf(value, "%f", &floatValue)
	return floatValue, err
}
