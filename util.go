package main

// VolumeStatusToNumber converts volume status strings to corresponding numeric values
func VolumeStatusToNumber(status string) float64 {
	switch status {
	case "available":
			return 0
	case "error":
			return 1
	case "creating":
			return 2
	case "deleting":
			return 3
	case "in-use":
			return 4
	case "attaching":
			return 5
	case "detaching":
			return 6
	case "error_deleting":
			return 7
	case "maintenance":
			return 8
	default:
			return 99 // or any other value indicating an unknown status
	}
}

// SnapshotStatusToNumber converts snapshot status strings to corresponding numeric values
func SnapshotStatusToNumber(status string) float64 {
	switch status {
	case "available":
			return 0
		case "error":
			return 1
	case "creating":
			return 2
	case "deleting":
			return 3
	case "error_deleting":
			return 7
	case "maintenance":
			return 8
	default:
			return 99 // or any other value indicating an unknown status
	}
}