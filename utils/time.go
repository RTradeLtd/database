package utils

import "time"

/*
Time related functionality
*/

// CalculateGarbageCollectDate is used to calculate the date at which a file
// will be removed from our inventory. To prepare data for input use `strconv.Atoi(fmt.Sprintf("%v", .....")"`
func CalculateGarbageCollectDate(holdTimeInMonths int) time.Time {
	return time.Now().UTC().AddDate(0, holdTimeInMonths, 0)
}
