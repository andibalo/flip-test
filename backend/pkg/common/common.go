package common

func IsCsvFile(filename string) bool {
	return len(filename) > 4 && filename[len(filename)-4:] == ".csv"
}
