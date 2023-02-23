package utils

func NotNill(args ...interface{}) bool {
	status := true

	for _, it := range args {
		if it == nil {
			status = false
			break
		}
	}
	return status
}

func IsNill(args ...interface{}) bool {
	status := true

	for _, it := range args {

		if it != nil {
			status = false
			break
		}
	}

	return status
}
