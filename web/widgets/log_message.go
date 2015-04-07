package widgets

import (
	"encoding/json"
	"fmt"
)

// Prepares source to be rendered in logs table
// Return it in next format:
// message {json}
// TODO: Bench. Improve speed.
func BuildLogLine(s map[string]interface{}) string {
	var message string
	if _, ok := s["message"]; ok {
		message = fmt.Sprintf("%v", s["message"])
	} else {
		message = ""
	}
	delete(s, "datetime")
	delete(s, "message")
	// If records has additional json attributes
	if len(s) > 0 {
		j, err := json.Marshal(s)
		if err != nil {
			return "Error rendering message"
		}

		return fmt.Sprintf("%s %s", message, string(j))
	}
	return fmt.Sprintf("%s", message)
}
