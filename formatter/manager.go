package formatter

var formatterMap = make(map[string]IFormatter)

func RegisterFormatter(formatter IFormatter) {
	formatterMap[formatter.GetName()] = formatter
}

func GetAllFormatter() map[string]IFormatter {
	return formatterMap
}
