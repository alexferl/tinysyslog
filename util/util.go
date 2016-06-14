package util

var severities = map[int]string{
	0: "EMERGENCY",
	1: "ALERT",
	2: "CRITICAL",
	3: "ERROR",
	4: "WARNING",
	5: "NOTICE",
	6: "INFO",
	7: "DEBUG",
}

func SeverityNumToString(severity int) string {
	switch severity {
	case 0:
		return severities[0]
	case 1:
		return severities[1]
	case 2:
		return severities[2]
	case 3:
		return severities[3]
	case 4:
		return severities[4]
	case 5:
		return severities[5]
	case 6:
		return severities[6]
	case 7:
		return severities[7]
	default:
		return "DEFAULT"
	}
}
