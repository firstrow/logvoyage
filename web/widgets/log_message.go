package widgets

import (
	"encoding/json"
	"fmt"
)

// Prepares source to be rendered in log table
// Return it in next format:
// message {json}
// TODO: Bench. Improve speed.
func BuildLogLine(s map[string]interface{}) string {
	message := s["message"]
	delete(s, "datetime")
	delete(s, "message")
	if len(s) > 0 {
		j, err := json.Marshal(s)
		if err != nil {
			//TODO: Log it!
			return "Error rendering message"
		}
		return fmt.Sprintf("%s %s", message, string(j))
	}
	return fmt.Sprintf("%s", message)
}
