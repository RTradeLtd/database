package utils_test

import (
	"testing"
	"time"

	"github.com/RTradeLtd/database/utils"
)

func TestCalculateGarbageCollectDate(t *testing.T) {
	esimatedTime := utils.CalculateGarbageCollectDate(5)
	if esimatedTime == time.Now() {
		t.Fatal("invalid time retrieved")
	}
}
