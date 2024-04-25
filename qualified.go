package q

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.
// You may assume that each input would have exactly one solution, and you may not use the same element twice.
func twoSum(nums []int, target int) []int {
	var (
		count = len(nums)
		memo  = make(map[int]int)
	)

	for i := 0; i <= count; i++ {
		remainder := target - nums[i]

		if v, ok := memo[remainder]; ok {
			return []int{v, i}
		}

		memo[nums[i]] = i
	}
	return nil
}

type point struct {
	x int
	y int
}

func pointArea(p1, p2, p3 point) float64 {
	a := (p1.x*(p2.y-p3.y) + p2.x*(p3.y-p1.y) + p3.x*(p1.y-p2.y)) / 2
	return math.Abs(float64(a))
}

func area(x1, y1, x2, y2, x3, y3 int) float64 {
	//Area A = [ x1(y2 – y3) + x2(y3 – y1) + x3(y1-y2)]/2

	a := (x1*(y2-y3) + x2*(y3-y1) + x3*(y1-y2)) / 2
	return math.Abs(float64(a))
}

func isInsidePoint(p1, p2, p3, p4 point) bool {
	var (
		abc = pointArea(p1, p2, p3)
		pbc = pointArea(p4, p2, p3)
		pac = pointArea(p1, p4, p3)
		pab = pointArea(p1, p2, p4)
	)
	return abc == (pbc + pac + pab)
}

func isInside(x1, y1, x2, y2, x3, y3, x, y int) bool {
	var (
		abc = area(x1, y1, x2, y2, x3, y3)
		pbc = area(x, y, x2, y2, x3, y3)
		pac = area(x1, y1, x, y, x3, y3)
		pab = area(x1, y1, x2, y2, x, y)
	)

	return abc == (pbc + pac + pab)
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}

	if n <= 3 {
		return true
	}

	if n%2 == 0 || n%3 == 0 {
		return false
	}

	i := 5

	for i*i <= n {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}

		i += 6
	}

	return true
}

func digitSum(n int) int {
	sum := 0
	for n > 0 {
		sum += n % 10
		n /= 10
	}
	return sum
}

func findClosestPairs(arr []int) [][]int {
	sort.Ints(arr)

	minDistance := arr[1] - arr[0]
	var closestPairs [][]int

	for i := 1; i < len(arr); i++ {
		if arr[i]-arr[i-1] < minDistance {
			minDistance = arr[i] - arr[i-1]
		}
	}

	for i := 1; i < len(arr); i++ {
		if arr[i]-arr[i-1] == minDistance {
			closestPairs = append(closestPairs, []int{arr[i-1], arr[i]})
		}
	}

	return closestPairs
}

func canRearrangeToPalindrome(s string) bool {
	charCount := make(map[rune]int)

	for _, char := range s {
		charCount[char]++
	}

	oddCount := 0

	for _, count := range charCount {
		if count%2 != 0 {
			oddCount++
		}
	}
	return oddCount <= 1
}

//func max(a, b int) int {
//	if a > b {
//		return a
//	}
//	return b
//}

func canEditToPalindrome(s string, i, j, k int) bool {
	substring := s[i-1 : j]

	editsRequired := max(0, len(substring)-sumCharFrequency(substring)/2)
	return editsRequired <= k
}

func sumCharFrequency(s string) int {
	charCount := make(map[rune]int)

	for _, char := range s {
		charCount[char]++
	}

	sum := 0
	for _, count := range charCount {
		sum += count
	}

	return sum
}

func isPalindrome(s string) bool {
	s = strings.ToLower(s)

	var builder strings.Builder

	for _, char := range s {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			builder.WriteRune(char)
		}
	}

	s = builder.String()

	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	reversedString := string(runes)
	return s == reversedString
}

//func divmod(a, b int) (int, int) {
//	return a / b, a % b
//}
//
//func MakeChange(amount int) map[byte]int {
//	coins := map[byte]int{'H': 50, 'Q': 25, 'D': 10, 'N': 5, 'P': 1}
//	change := make(map[byte]int)
//
//	for coin, val := range coins {
//		if amount == 0 {
//			break
//		}
//
//		quot, rem := divmod(amount, val)
//		if quot > 0 {
//			change[coin] = quot
//		}
//
//		amount = rem
//	}
//
//	var keys []byte
//	for k := range change {
//		keys = append(keys, k)
//	}
//
//	sort.Slice(keys, func(i, j int) bool {
//		return keys[i] < keys[j]
//	})
//
//	sortedChange := make(map[byte]int)
//	for _, k := range keys {
//		sortedChange[k] = change[k]
//	}
//
//	return sortedChange
//}

func divmod(a, b int) (int, int) {
	return a / b, a % b
}

func MakeChange(amount int) map[byte]int {
	coins := map[byte]int{'H': 50, 'Q': 25, 'D': 10, 'N': 5, 'P': 1}
	change := map[byte]int{}

	for coin, val := range coins {
		quot, rem := divmod(amount, val)
		if quot > 0 { // Only add non-zero counts to the map
			change[coin] = quot
		}
		amount = rem
	}

	sortedKeys := make([]byte, 0, len(change))
	for key := range change {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		return sortedKeys[i] < sortedKeys[j]
	})

	sortedChange := make(map[byte]int)
	for _, key := range sortedKeys {
		sortedChange[key] = change[key]
	}

	return sortedChange
}
func main() {
	//fmt.Println(twoSum([]int{2, 7, 11, 15}, 9)) // [0,1]
	//fmt.Println(twoSum([]int{3, 2, 4}, 6))      // [1,2]
	//fmt.Println(twoSum([]int{3, 3}, 6))         // [0,1]

	// P(10, 15)
	// A(0, 0), B(20, 0) and C(10, 30)

	//fmt.Println(isInside(0, 0, 20, 0, 10, 30, 10, 15))
	//fmt.Println(isInsidePoint(point{0, 0}, point{20, 0}, point{10, 30}, point{10, 15}))

	//for num := 1; num <= 1_000; num++ {
	//	if isPrime(num) && digitSum(num) < 10 {
	//		fmt.Println(num)
	//	}
	//}

	//arr := []int{4, 2, 1, 3, 5, 6}
	//closestPairs := findClosestPairs(arr)
	//
	//for _, pair := range closestPairs {
	//	fmt.Println(pair)
	//}

	//s := "abccda"
	//queries := [][]int{{1, 3, 1}, {2, 5, 2}, {1, 6, 4}}
	//
	//for _, query := range queries {
	//	i, j, k := query[0], query[1], query[2]
	//	if canEditToPalindrome(s, i, j, k) {
	//		fmt.Printf("Substring (%d, %d) can be edited at most %d times to form a palindrome\n", i, j, k)
	//	} else {
	//		fmt.Printf("Substring (%d, %d) cannot be edited at most %d times to form a palindrome\n", i, j, k)
	//	}
	//}

	fmt.Println(isPalindrome("A man, a plan, a canal, Panama")) // true
	fmt.Println(isPalindrome("race a car"))                     // false

}

// package challenge // import "qualified.io/challenge"
//import "sort"
//
//func divmod(a, b int) (int, int) {
//    return a / b, a % b
//}
//
//func MakeChange(amount int) map[byte]int {
//  coins := map[byte]int{'H': 50, 'Q': 25, 'D': 10, 'N': 5, 'P': 1}
//  sortedCoins := make(map[byte]int)
//
//  for coin, val := range coins {
//    sortedCoins[coin] = val
//  }
//  var sortedKeys []byte
//  for key := range sortedCoins {
//    sortedKeys = append(sortedKeys, key)
//  }
//  sort.Slice(sortedKeys, func(i, j int) bool {
//    return sortedCoins[sortedKeys[i]] > sortedCoins[sortedKeys[j]]
//  })
//
//  change := map[byte]int{}
//  for _, coin := range sortedKeys {
//    quot, rem := divmod(amount, sortedCoins[coin])
//    change[coin] = quot
//    amount = rem
//    if change[coin] == 0 {
//      delete(change, coin)
//    }
//  }
//  return change
//}

// Task
//Your team is working on a project for a retail chain which keeps their stores' weekly open hours listings as an array of data. The data is in the following format:
//
//openHours := []map[string]string{
//  {
//    "day":   "Monday",
//    "open":  "8:00 AM",
//    "close": "5:00 PM",
//  },
//  {
//    "day":   "Tuesday",
//    "open":  "8:00 AM",
//    "close": "5:00 PM",
//  },
//  {
//    "day":   "Wednesday",
//    "open":  "8:00 AM",
//    "close": "6:00 PM",
//  },
//  {
//    "day":   "Thursday",
//    "open":  "8:00 AM",
//    "close": "5:00 PM",
//  },
//  {
//    "day":   "Friday",
//    "open":  "8:00 AM",
//    "close": "5:00 PM",
//  },
//  {
//    "day":   "Saturday",
//    "open":  "8:00 AM",
//    "close": "4:00 PM",
//  },
//}
//However, the company's website needs the data to be transformed to a grouped format for displaying to visitors. The grouped format is as follows:
//
//groupedOpenHours := []map[string]string{
//  {
//    "days":  "Monday-Tuesday",
//    "open":  "8:00 AM",
//    "close": "5:00 PM",
//  },
//  {
//    "days":  "Wednesday",
//    "open":  "8:00 AM",
//    "close": "6:00 PM",
//  },
//  {
//    "days":  "Thursday-Friday",
//    "open":  "8:00 AM",
//    "close": "5:00 PM",
//  },
//  {
//    "days":  "Saturday",
//    "open":  "8:00 AM",
//    "close": "4:00 PM",
//  },
//  {
//    "days":  "Sunday",
//    "open":  "",
//    "close": "",
//  },
//}
//In the output above, any consecutive days sharing the same open and close times have been compressed into the same map. Whenever consecutive days sharing open and close times are encountered, the first and last day in the range are concatenated with a hyphen for the "day" key.
//
//Your task is to write code to perform the transformation. The function you'll complete is GroupOpenHours(openHours []map[string]string) []map[string]string. The function should return the transformed map slice in the above format.
//
//The output slice should be in order of the days of the week. Consider Monday as the beginning of the week and Sunday as the end. No range that bridges the gap between Sunday-Monday should be created (but a range from Monday-Sunday is valid whenever the entire week has the same open/closed hours or the input is empty.
//
//As shown above, any missing days of the week should be added to the returned slice as maps with "open": "" and "close": "" entries. When the missing dates consist of consecutive ranges, they should use the same hyphenated grouped "days" key format as open days would be.
//
//The openHours parameter will always be well-formed but may be empty and unsorted. There will never be more than 7 maps in the slice, and each map is guaranteed to have only "day", "open" and "close" keys present with string values formatted as in the structure shown above. All "day" keys are guaranteed to be unique in openHours and correctly capitalized, valid days of the week.
//
//Please do not mutate the openHours parameter.
