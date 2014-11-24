package widgets

import (
	"fmt"
	"html/template"

	"github.com/Unknwon/com"
	"github.com/firstrow/logvoyage/common"
)

// Renders widget with user source groups.
// groups - user source groups
// selected - list if types in query string
func NewSourceGroups(groups []*common.SourceGroup, selected []string) template.HTML {
	var result string
	// Iterate over users source groups
	for _, gr := range groups {
		// If group has some types
		if len(gr.Types) > 0 {
			result += fmt.Sprintf(`<optgroup label="%s">`, gr.Name)
			for _, t := range gr.Types {
				s := ""
				// Check if type is in search query
				if com.IsSliceContainsStr(selected, t) {
					s = "selected"
				}
				result += fmt.Sprintf(`<option value="%s" %s>%s</option>`, t, s, t)
			}
			result += `</optgroup>`
		}
	}
	return template.HTML(result)
}
