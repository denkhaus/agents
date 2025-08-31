package prompt

import "fmt"

// convertMapKeysToStrings recursively converts map[interface{}]interface{} keys to map[string]interface{} keys.
func convertMapKeysToStrings(in interface{}) (interface{}, error) {
	switch in := in.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range in {
			strKey, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("map key is not a string: %v", k)
			}
			convertedValue, err := convertMapKeysToStrings(v)
			if err != nil {
				return nil, err
			}
			m[strKey] = convertedValue
		}
		return m, nil
	case []interface{}:
		var l []interface{}
		for _, v := range in {
			convertedValue, err := convertMapKeysToStrings(v)
			if err != nil {
				return nil, err
			}
			l = append(l, convertedValue)
		}
		return l, nil
	default:
		return in, nil
	}
}
