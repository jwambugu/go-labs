//package q
//
//import (
//	"github.com/stretchr/testify/require"
//	"sort"
//	"testing"
//)
//
//func TestMakeChange(t *testing.T) {
//
//	expected := map[byte]int{'Q': 1, 'D': 1, 'N': 1, 'P': 3}
//	t.Log(MakeChange(43), expected)
//
//	expected = map[byte]int{'H': 1, 'Q': 1, 'D': 1, 'N': 1, 'P': 1}
//
//	t.Log(MakeChange(91), expected)
//
//	//Expect().To(Equal(expected))
//	//var _ = Describe("MakeChange", func() {
//	//	It("should make the correct change on 43", func() {
//	//
//	//	})
//	//	It("should make the correct change on 91", func() {
//	//		expected := map[byte]int{'H': 1, 'Q': 1, 'D': 1, 'N': 1, 'P': 1}
//	//		Expect(MakeChange(91)).To(Equal(expected))
//	//	})
//	//})
//
//}
//
//func GroupOpenHours(openHours []map[string]string) []map[string]string {
//	// Create a map to store grouped open hours data
//	grouped := make(map[string]map[string]string)
//
//	// Loop through the original open hours data
//	for _, entry := range openHours {
//		day := entry["day"]
//		open := entry["open"]
//		close := entry["close"]
//
//		// Check if a group already exists for these open/close times
//		if existingGroup, ok := grouped[open+"-"+close]; ok {
//			// Existing group found, update start and end days
//			if existingGroup["days"] < day {
//				existingGroup["days"] = existingGroup["days"] + "-" + day
//			} else {
//				existingGroup["days"] = day + "-" + existingGroup["days"]
//			}
//		} else {
//			// Create a new group for these open/close times
//			grouped[open+"-"+close] = map[string]string{
//				"days":  day,
//				"open":  open,
//				"close": close,
//			}
//		}
//	}
//
//	// Convert grouped map to a slice and sort by day
//	result := make([]map[string]string, 0, len(grouped))
//	for _, entry := range grouped {
//		result = append(result, entry)
//	}
//	sort.Slice(result, func(i, j int) bool {
//		return getWeekdayOrder(result[i]["days"]) < getWeekdayOrder(result[j]["days"])
//	})
//
//	// Fill in missing days with empty open/close times
//	expectedDays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
//	currentDay := 0
//	for _, expectedDay := range expectedDays {
//		if currentDay < len(result) && result[currentDay]["days"] == expectedDay {
//			currentDay++
//		} else {
//			result = append(result, map[string]string{
//				"days":  expectedDay,
//				"open":  "",
//				"close": "",
//			})
//		}
//	}
//
//	return result
//}
//
//// Helper function to get weekday order (Monday = 0, Sunday = 6)
//func getWeekdayOrder(day string) int {
//	switch day {
//	case "Monday":
//		return 0
//	case "Tuesday":
//		return 1
//	case "Wednesday":
//		return 2
//	case "Thursday":
//		return 3
//	case "Friday":
//		return 4
//	case "Saturday":
//		return 5
//	case "Sunday":
//		return 6
//	default:
//		return -1 // Handle unexpected day format
//	}
//}
//
////package challenge_test
////
////import (
////. "github.com/onsi/ginkgo"
////. "github.com/onsi/gomega"
////. "qualified.io/challenge"
////)
////
////var _ = Describe("GroupOpenHours", func() {
////	It("should work on the provided example", func() {
////
////		Expect(actual).To(Equal(expected))
////	})
////
////	It("should handle ranges of closed days", func() {
////
////		Expect(actual).To(Equal(expected))
////	})
////})
//
//func TestGroup(t *testing.T) {
//	openHours := []map[string]string{
//		{
//			"day":   "Monday",
//			"open":  "8:00 AM",
//			"close": "5:00 PM",
//		},
//		{
//			"day":   "Tuesday",
//			"open":  "8:00 AM",
//			"close": "5:00 PM",
//		},
//		{
//			"day":   "Wednesday",
//			"open":  "8:00 AM",
//			"close": "6:00 PM",
//		},
//		{
//			"day":   "Thursday",
//			"open":  "8:00 AM",
//			"close": "5:00 PM",
//		},
//		{
//			"day":   "Friday",
//			"open":  "8:00 AM",
//			"close": "5:00 PM",
//		},
//		{
//			"day":   "Saturday",
//			"open":  "8:00 AM",
//			"close": "4:00 PM",
//		},
//	}
//	expected := []map[string]string{
//		{
//			"days":  "Monday-Tuesday",
//			"open":  "8:00 AM",
//			"close": "5:00 PM",
//		},
//		{
//			"days":  "Wednesday",
//			"open":  "8:00 AM",
//			"close": "6:00 PM",
//		},
//		{
//			"days":  "Thursday-Friday",
//			"open":  "8:00 AM",
//			"close": "5:00 PM",
//		},
//		{
//			"days":  "Saturday",
//			"open":  "8:00 AM",
//			"close": "4:00 PM",
//		},
//		{
//			"days":  "Sunday",
//			"open":  "",
//			"close": "",
//		},
//	}
//	actual := GroupOpenHours(openHours)
//
//	require.Equal(t, expected, actual)
//
//	//openHours = []map[string]string{
//	//	{
//	//		"day":   "Monday",
//	//		"open":  "9:00 AM",
//	//		"close": "4:00 PM",
//	//	},
//	//	{
//	//		"day":   "Tuesday",
//	//		"open":  "9:00 AM",
//	//		"close": "4:00 PM",
//	//	},
//	//	{
//	//		"day":   "Wednesday",
//	//		"open":  "9:00 AM",
//	//		"close": "4:00 PM",
//	//	},
//	//	{
//	//		"day":   "Thursday",
//	//		"open":  "9:00 AM",
//	//		"close": "4:00 PM",
//	//	},
//	//}
//	//expected = []map[string]string{
//	//	{
//	//		"days":  "Monday-Thursday",
//	//		"open":  "9:00 AM",
//	//		"close": "4:00 PM",
//	//	},
//	//	{
//	//		"days":  "Friday-Sunday",
//	//		"open":  "",
//	//		"close": "",
//	//	},
//	//}
//	//actual = GroupOpenHours(openHours)
//	//require.Equal(t, expected, actual)
//
//}

package q_test

import "fmt"

package challenge // import "qualified.io/challenge"

import "fmt"

func containsDay(openHours []map[string]string, day string) bool {
	for _, oh := range openHours {
		if oh["day"] == day {
			return true
		}
	}
	return false
}

func formatDays(start, end string) string {
	if start == end {
		return start
	}
	return fmt.Sprintf("%s-%s", start, end)
}

func mergeClosedDays(openHours []map[string]string) []map[string]string {
	merged := make([]map[string]string, 0)

	var currentStart, currentEnd string
	var isOpen bool
	for i, oh := range openHours {
		if i == 0 {
			currentStart = oh["days"]
			currentEnd = oh["days"]
			isOpen = oh["open"] != "" && oh["close"] != ""
			continue
		}

		prev := openHours[i-1]
		isPrevOpen := prev["open"] != "" && prev["close"] != ""
		isCurrOpen := oh["open"] != "" && oh["close"] != ""

		if isOpen == isPrevOpen && isOpen == isCurrOpen && oh["open"] == "" && oh["close"] == "" {
			currentEnd = oh["days"]
		} else {
			merged = append(merged, map[string]string{
				"days":  formatDays(currentStart, currentEnd),
				"open":  prev["open"],
				"close": prev["close"],
			})
			currentStart = oh["days"]
			currentEnd = oh["days"]
			isOpen = oh["open"] != "" && oh["close"] != ""
		}
	}

	merged = append(merged, map[string]string{
		"days":  formatDays(currentStart, currentEnd),
		"open":  openHours[len(openHours)-1]["open"],
		"close": openHours[len(openHours)-1]["close"],
	})

	return merged
}

func GroupOpenHours(openHours []map[string]string) []map[string]string {
	if len(openHours) == 0 {
		return []map[string]string{{"days": "Monday-Sunday", "open": "", "close": ""}}
	}

	groupedOpenHours := make([]map[string]string, 0)

	var currentStart, currentEnd string
	var currentOpen, currentClose string
	var isOpen bool
	for i, oh := range openHours {
		if i == 0 {
			currentStart = oh["day"]
			currentEnd = oh["day"]
			currentOpen = oh["open"]
			currentClose = oh["close"]
			isOpen = currentOpen != "" && currentClose != ""
			continue
		}

		prev := openHours[i-1]
		isPrevOpen := prev["open"] != "" && prev["close"] != ""
		isCurrOpen := oh["open"] != "" && oh["close"] != ""

		if isOpen == isPrevOpen && isOpen == isCurrOpen &&
			(currentOpen == oh["open"] && currentClose == oh["close"]) {
			currentEnd = oh["day"]
		} else {
			groupedOpenHours = append(groupedOpenHours, map[string]string{
				"days":  formatDays(currentStart, currentEnd),
				"open":  currentOpen,
				"close": currentClose,
			})
			currentStart = oh["day"]
			currentEnd = oh["day"]
			currentOpen = oh["open"]
			currentClose = oh["close"]
			isOpen = currentOpen != "" && currentClose != ""
		}
	}

	groupedOpenHours = append(groupedOpenHours, map[string]string{
		"days":  formatDays(currentStart, currentEnd),
		"open":  currentOpen,
		"close": currentClose,
	})

	for _, day := range []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"} {
		if !containsDay(openHours, day) {
			groupedOpenHours = append(groupedOpenHours, map[string]string{"days": day, "open": "", "close": ""})
		}
	}

	groupedOpenHours = mergeClosedDays(groupedOpenHours)

	return groupedOpenHours
}
