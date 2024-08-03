package difference

const (
	added   = "+"
	removed = "-"
	common  = ""
)

func isChangeSet(currentSign string, upcomingSign string) bool {
	if currentSign == added && upcomingSign == removed {
		return true
	} else if upcomingSign == added && currentSign == removed {
		return true
	} else {
		return false
	}
}
