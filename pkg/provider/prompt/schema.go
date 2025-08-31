package prompt

import "fmt"

// JSONSchema is a custom type to handle unmarshaling of JSON schema from YAML.
type JSONSchema map[string]interface{}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (s *JSONSchema) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw interface{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	converted, err := convertMapKeysToStrings(raw)
	if err != nil {
		return err
	}

	convertedMap, ok := converted.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected a map[string]interface{}, got %T", converted)
	}
	*s = convertedMap
	return nil
}
