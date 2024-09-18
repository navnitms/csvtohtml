package converter

func getStringOption(options map[string]interface{}, key, defaultValue string) string {
	if val, ok := options[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func getIntOption(options map[string]interface{}, key string, defaultValue int) int {
	if val, ok := options[key]; ok {
		if i, ok := val.(int); ok {
			return i
		}
	}
	return defaultValue
}