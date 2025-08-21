package main

import . "github.com/enetx/g"

func main() {
	// Example 1: Scan to show cumulative sum
	SliceOf(1, 2, 3, 4).
		Iter().
		Scan(0, func(acc, val int) int {
			return acc + val
		}).
		Collect().
		Println() // Slice[0, 1, 3, 6, 10]

	// Example 2: Scan to build string progressively
	SliceOf("a", "b", "c", "d").
		Iter().
		Scan("", func(acc, val string) string {
			return acc + val
		}).
		Collect().
		Println() // Slice[, a, ab, abc, abcd]

	// Example 3: Scan to track running maximum
	SliceOf(3, 1, 4, 1, 5, 9, 2, 6).
		Iter().
		Scan(0, func(acc, val int) int {
			if val > acc {
				return val
			}
			return acc
		}).
		Collect().
		Println() // Slice[0, 3, 3, 4, 4, 5, 9, 9, 9]

	// Example 4: Scan to calculate factorial sequence
	SliceOf(1, 2, 3, 4, 5).
		Iter().
		Scan(1, func(acc, val int) int {
			return acc * val
		}).
		Collect().
		Println() // Slice[1, 1, 2, 6, 24, 120]

	// Example 5: Scan to track balance with transactions
	transactions := SliceOf(100, -30, -20, 50, -10)
	transactions.
		Iter().
		Scan(1000, func(balance, transaction int) int {
			return balance + transaction
		}).
		Collect().
		Println() // Slice[1000, 1100, 1070, 1050, 1100, 1090]

	// Example 6: Scan with custom type - building a path
	SliceOf("home", "user", "documents", "file.txt").
		Iter().
		Scan("", func(path, segment string) string {
			if path == "" {
				return "/" + segment
			}
			return path + "/" + segment
		}).
		Skip(1). // Skip the initial empty value
		Collect().
		Println() // Slice[/home, /home/user, /home/user/documents, /home/user/documents/file.txt]

	// Example 7: Scan vs Reduce comparison
	numbers := SliceOf(1, 2, 3, 4)

	// Reduce gives only final result
	reduceResult := numbers.Iter().Reduce(func(a, b int) int {
		return a + b
	})
	String("Reduce result: ").Append(Int(reduceResult.UnwrapOr(0)).String()).Println()
	// Reduce result: 10

	// Scan gives all intermediate results
	String("Scan results: ").Println()
	numbers.Iter().
		Scan(0, func(acc, val int) int {
			return acc + val
		}).
		Collect().
		Println()
	// Scan results:
	// Slice[0, 1, 3, 6, 10]
}
