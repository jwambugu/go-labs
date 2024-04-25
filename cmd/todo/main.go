//package main
//
//import (
//	"fmt"
//)
//
////func GroupOpenHours(openHours []map[string]string) []map[string]string {
////	if len(openHours) == 0 {
////		return []map[string]string{{"days": "Monday-Sunday", "open": "", "close": ""}}
////	}
////
////	groupedOpenHours := make([]map[string]string, 0)
////
////	var currentStart, currentEnd string
////	var currentOpen, currentClose string
////	var isOpen bool
////	for i, oh := range openHours {
////		if i == 0 {
////			currentStart = oh["day"]
////			currentEnd = oh["day"]
////			currentOpen = oh["open"]
////			currentClose = oh["close"]
////			isOpen = currentOpen != "" && currentClose != ""
////			continue
////		}
////
////		prev := openHours[i-1]
////		isPrevOpen := prev["open"] != "" && prev["close"] != ""
////		isCurrOpen := oh["open"] != "" && oh["close"] != ""
////
////		if isOpen == isPrevOpen && isOpen == isCurrOpen &&
////			(currentOpen == oh["open"] && currentClose == oh["close"]) {
////			currentEnd = oh["day"]
////		} else {
////			groupedOpenHours = append(groupedOpenHours, map[string]string{
////				"days":  formatDays(currentStart, currentEnd),
////				"open":  currentOpen,
////				"close": currentClose,
////			})
////			currentStart = oh["day"]
////			currentEnd = oh["day"]
////			currentOpen = oh["open"]
////			currentClose = oh["close"]
////			isOpen = currentOpen != "" && currentClose != ""
////		}
////	}
////
////	// Append the last entry
////	groupedOpenHours = append(groupedOpenHours, map[string]string{
////		"days":  formatDays(currentStart, currentEnd),
////		"open":  currentOpen,
////		"close": currentClose,
////	})
////
////	// Add missing days
////	for _, day := range []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"} {
////		if !containsDay(openHours, day) {
////			groupedOpenHours = append(groupedOpenHours, map[string]string{"days": day, "open": "", "close": ""})
////		}
////	}
////
////	// Merge consecutive closed days into a single entry
////	groupedOpenHours = mergeClosedDays(groupedOpenHours)
////
////	return groupedOpenHours
////}
////
////func containsDay(openHours []map[string]string, day string) bool {
////	for _, oh := range openHours {
////		if oh["day"] == day {
////			return true
////		}
////	}
////	return false
////}
////
////func formatDays(start, end string) string {
////	if start == end {
////		return start
////	}
////	return fmt.Sprintf("%s-%s", start, end)
////}
////
////func mergeClosedDays(openHours []map[string]string) []map[string]string {
////	merged := make([]map[string]string, 0)
////
////	var currentStart, currentEnd string
////	var isOpen bool
////	for i, oh := range openHours {
////		if i == 0 {
////			currentStart = oh["days"]
////			currentEnd = oh["days"]
////			isOpen = oh["open"] != "" && oh["close"] != ""
////			continue
////		}
////
////		prev := openHours[i-1]
////		isPrevOpen := prev["open"] != "" && prev["close"] != ""
////		isCurrOpen := oh["open"] != "" && oh["close"] != ""
////
////		if isOpen == isPrevOpen && isOpen == isCurrOpen && oh["open"] == "" && oh["close"] == "" {
////			currentEnd = oh["days"]
////		} else {
////			merged = append(merged, map[string]string{
////				"days":  formatDays(currentStart, currentEnd),
////				"open":  "",
////				"close": "",
////			})
////			currentStart = oh["days"]
////			currentEnd = oh["days"]
////			isOpen = oh["open"] != "" && oh["close"] != ""
////		}
////	}
////
////	// Append the last entry
////	merged = append(merged, map[string]string{
////		"days":  formatDays(currentStart, currentEnd),
////		"open":  "",
////		"close": "",
////	})
////
////	return merged
////}
//
//func main() {
//	openHours := []map[string]string{
//		{"day": "Monday", "open": "9:00 AM", "close": "4:00 PM"},
//		{"day": "Tuesday", "open": "9:00 AM", "close": "4:00 PM"},
//		{"day": "Wednesday", "open": "9:00 AM", "close": "4:00 PM"},
//		{"day": "Thursday", "open": "9:00 AM", "close": "4:00 PM"},
//	}
//	groupedOpenHours := GroupOpenHours(openHours)
//	for _, oh := range groupedOpenHours {
//		fmt.Println(oh)
//	}
//
//	fmt.Println("==============")
//
//	openHours = []map[string]string{
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
//	groupedOpenHours = GroupOpenHours(openHours)
//	for _, oh := range groupedOpenHours {
//		fmt.Println(oh)
//	}
//}
//func GroupOpenHours(openHours []map[string]string) []map[string]string {
//	if len(openHours) == 0 {
//		return []map[string]string{{"days": "Monday-Sunday", "open": "", "close": ""}}
//	}
//
//	groupedOpenHours := make([]map[string]string, 0)
//
//	var currentStart, currentEnd string
//	var currentOpen, currentClose string
//	var isOpen bool
//	for i, oh := range openHours {
//		if i == 0 {
//			currentStart = oh["day"]
//			currentEnd = oh["day"]
//			currentOpen = oh["open"]
//			currentClose = oh["close"]
//			isOpen = currentOpen != "" && currentClose != ""
//			continue
//		}
//
//		prev := openHours[i-1]
//		isPrevOpen := prev["open"] != "" && prev["close"] != ""
//		isCurrOpen := oh["open"] != "" && oh["close"] != ""
//
//		if isOpen == isPrevOpen && isOpen == isCurrOpen &&
//			(currentOpen == oh["open"] && currentClose == oh["close"]) {
//			currentEnd = oh["day"]
//		} else {
//			groupedOpenHours = append(groupedOpenHours, map[string]string{
//				"days":  formatDays(currentStart, currentEnd),
//				"open":  currentOpen,
//				"close": currentClose,
//			})
//			currentStart = oh["day"]
//			currentEnd = oh["day"]
//			currentOpen = oh["open"]
//			currentClose = oh["close"]
//			isOpen = currentOpen != "" && currentClose != ""
//		}
//	}
//
//	groupedOpenHours = append(groupedOpenHours, map[string]string{
//		"days":  formatDays(currentStart, currentEnd),
//		"open":  currentOpen,
//		"close": currentClose,
//	})
//
//	for _, day := range []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"} {
//		if !containsDay(openHours, day) {
//			groupedOpenHours = append(groupedOpenHours, map[string]string{"days": day, "open": "", "close": ""})
//		}
//	}
//
//	groupedOpenHours = mergeClosedDays(groupedOpenHours)
//
//	return groupedOpenHours
//}
//
//func containsDay(openHours []map[string]string, day string) bool {
//	for _, oh := range openHours {
//		if oh["day"] == day {
//			return true
//		}
//	}
//	return false
//}
//
//func formatDays(start, end string) string {
//	if start == end {
//		return start
//	}
//	return fmt.Sprintf("%s-%s", start, end)
//}
//
//func mergeClosedDays(openHours []map[string]string) []map[string]string {
//	merged := make([]map[string]string, 0)
//
//	var currentStart, currentEnd string
//	var isOpen bool
//	for i, oh := range openHours {
//		if i == 0 {
//			currentStart = oh["days"]
//			currentEnd = oh["days"]
//			isOpen = oh["open"] != "" && oh["close"] != ""
//			continue
//		}
//
//		prev := openHours[i-1]
//		isPrevOpen := prev["open"] != "" && prev["close"] != ""
//		isCurrOpen := oh["open"] != "" && oh["close"] != ""
//
//		if isOpen == isPrevOpen && isOpen == isCurrOpen && oh["open"] == "" && oh["close"] == "" {
//			currentEnd = oh["days"]
//		} else {
//			merged = append(merged, map[string]string{
//				"days":  formatDays(currentStart, currentEnd),
//				"open":  prev["open"],
//				"close": prev["close"],
//			})
//			currentStart = oh["days"]
//			currentEnd = oh["days"]
//			isOpen = oh["open"] != "" && oh["close"] != ""
//		}
//	}
//
//	merged = append(merged, map[string]string{
//		"days":  formatDays(currentStart, currentEnd),
//		"open":  openHours[len(openHours)-1]["open"],
//		"close": openHours[len(openHours)-1]["close"],
//	})
//
//	return merged
//}

package main

import (
	"fmt"
)

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

	// Append the last entry
	groupedOpenHours = append(groupedOpenHours, map[string]string{
		"days":  formatDays(currentStart, currentEnd),
		"open":  currentOpen,
		"close": currentClose,
	})

	// Add missing days
	for _, day := range []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"} {
		if !containsDay(openHours, day) {
			groupedOpenHours = append(groupedOpenHours, map[string]string{"days": day, "open": "", "close": ""})
		}
	}

	// Merge consecutive closed days into a single entry
	groupedOpenHours = mergeClosedDays(groupedOpenHours)

	return groupedOpenHours
}

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

	// Append the last entry
	merged = append(merged, map[string]string{
		"days":  formatDays(currentStart, currentEnd),
		"open":  openHours[len(openHours)-1]["open"],
		"close": openHours[len(openHours)-1]["close"],
	})

	return merged
}

func main() {
	openHours := []map[string]string{
		{"day": "Monday", "open": "9:00 AM", "close": "4:00 PM"},
		{"day": "Tuesday", "open": "9:00 AM", "close": "4:00 PM"},
		{"day": "Wednesday", "open": "9:00 AM", "close": "4:00 PM"},
		{"day": "Thursday", "open": "9:00 AM", "close": "4:00 PM"},
	}
	groupedOpenHours := GroupOpenHours(openHours)
	for _, oh := range groupedOpenHours {
		fmt.Println(oh)
	}
}
