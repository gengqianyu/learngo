package common

import "command/view"

func View(templateName string) view.View {
	return view.CreateView(templateName)
}
