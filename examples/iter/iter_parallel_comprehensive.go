package main

import (
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	Println("=== Comprehensive Parallel Iterator Pipeline ===\n")

	// Simulating a data processing pipeline with all new parallel methods
	type Transaction struct {
		ID        int
		Amount    float64
		Category  String
		UserID    int
		Timestamp int64
	}

	// Sample data - financial transactions
	transactions := SliceOf(
		Transaction{1, 150.50, "food", 101, 1640995200},
		Transaction{2, 25.00, "transport", 102, 1640995300},
		Transaction{3, 300.75, "shopping", 101, 1640995400},
		Transaction{4, 45.25, "food", 103, 1640995500},
		Transaction{5, 89.99, "entertainment", 102, 1640995600},
		Transaction{6, 12.50, "transport", 104, 1640995700},
		Transaction{7, 200.00, "shopping", 103, 1640995800},
		Transaction{8, 67.80, "food", 101, 1640995900},
		Transaction{9, 15.75, "transport", 105, 1641000000},
		Transaction{10, 125.30, "entertainment", 104, 1641000100},
		Transaction{11, 75.25, "food", 102, 1641000200},
		Transaction{12, 450.00, "shopping", 105, 1641000300},
	)

	Println("=== Processing {} transactions with parallel pipeline ===", transactions.Len())
	totalStart := time.Now()

	// Step 1: FlatMap - Expand each transaction into multiple analysis records
	Println("1. FlatMap: Expanding transactions into analysis records...")
	start := time.Now()

	analysisRecords := transactions.Iter().
		Parallel(4).
		FlatMap(func(t Transaction) SeqSlice[Transaction] {
			// Simulate complex analysis
			time.Sleep(20 * time.Millisecond)

			records := NewSlice[Transaction]()

			// Base transaction record
			records.Push(t)

			// Add risk analysis record for high amounts
			if t.Amount > 200 {
				riskRecord := t
				riskRecord.ID = t.ID + 10000 // Different ID for risk analysis
				riskRecord.Category = t.Category.Append("_risk")
				records.Push(riskRecord)
			}

			// Add user analysis record for frequent users
			if t.UserID <= 102 {
				userRecord := t
				userRecord.ID = t.ID + 20000 // Different ID for user analysis
				userRecord.Category = t.Category.Append("_user")
				records.Push(userRecord)
			}

			return records.Iter()
		}).
		Collect()

	duration1 := time.Since(start)
	Println("   Generated {} analysis records in {}", analysisRecords.Len(), duration1)
	Println("   Sample records: {} entries", analysisRecords.Iter().Take(5).Collect().Len())
	Println("")

	// Step 2: FilterMap - Process high-value transactions
	Println("2. FilterMap: Processing high-value transactions...")
	start = time.Now()

	processedTransactions := transactions.Iter().
		Parallel(3).
		FilterMap(func(t Transaction) Option[Transaction] {
			// Simulate transaction validation
			time.Sleep(15 * time.Millisecond)

			// Only process transactions above $50 in specific categories
			if t.Amount > 50.0 && (t.Category == "food" || t.Category == "shopping" || t.Category == "entertainment") {
				// Apply processing fee
				processed := t
				processed.Amount = t.Amount * 1.2
				return Some(processed)
			}
			return None[Transaction]()
		}).
		Collect()

	duration2 := time.Since(start)
	Println("   Processed {} qualified transactions in {}", processedTransactions.Len(), duration2)
	processedTransactions.Iter().Take(3).ForEach(func(tx Transaction) {
		Println("   - ID {}: ${} ({}) from user {}", tx.ID, tx.Amount, tx.Category, tx.UserID)
	})
	Println("")

	// Step 3: StepBy - Sample every 3rd transaction for audit
	Println("3. StepBy: Sampling transactions for audit (every 3rd)...")
	start = time.Now()

	auditSample := transactions.Iter().
		Parallel(3).
		Inspect(func(t Transaction) {
			// Simulate audit preparation
			time.Sleep(10 * time.Millisecond)
		}).
		StepBy(3).
		Collect()

	duration3 := time.Since(start)
	Println("   Sampled {} transactions for audit in {}", auditSample.Len(), duration3)
	auditSample.Iter().ForEach(func(tx Transaction) {
		Println("   - Audit ID {}: ${} ({})", tx.ID, tx.Amount, tx.Category)
	})
	Println("")

	// Step 4: MaxBy/MinBy - Find extreme values
	Println("4. MaxBy/MinBy: Finding transaction extremes...")
	start = time.Now()

	// Find highest value transaction
	maxTransaction := transactions.Iter().
		Parallel(3).
		Inspect(func(t Transaction) {
			// Simulate complex valuation
			time.Sleep(12 * time.Millisecond)
		}).
		MaxBy(func(a, b Transaction) cmp.Ordering {
			return cmp.Cmp(a.Amount, b.Amount)
		})

	// Find lowest value transaction
	minTransaction := transactions.Iter().
		Parallel(3).
		Inspect(func(t Transaction) {
			time.Sleep(12 * time.Millisecond)
		}).
		MinBy(func(a, b Transaction) cmp.Ordering {
			return cmp.Cmp(a.Amount, b.Amount)
		})

	// Find most recent transaction
	latestTransaction := transactions.Iter().
		Parallel(3).
		Inspect(func(t Transaction) {
			time.Sleep(12 * time.Millisecond)
		}).
		MaxBy(func(a, b Transaction) cmp.Ordering {
			return cmp.Cmp(a.Timestamp, b.Timestamp)
		})

	duration4 := time.Since(start)
	if maxTransaction.IsSome() {
		max := maxTransaction.Some()
		Println("   Highest value: ${} (ID: {}, Category: {})", max.Amount, max.ID, max.Category)
	}
	if minTransaction.IsSome() {
		min := minTransaction.Some()
		Println("   Lowest value: ${} (ID: {}, Category: {})", min.Amount, min.ID, min.Category)
	}
	if latestTransaction.IsSome() {
		latest := latestTransaction.Some()
		Println("   Most recent: ID {} at timestamp {}", latest.ID, latest.Timestamp)
	}
	Println("   Analysis completed in {}", duration4)
	Println("")

	// Step 5: Complex pipeline combining all methods
	Println("5. Combined Pipeline: All methods in sequence...")
	start = time.Now()

	finalResults := transactions.Iter().
		Parallel(4).
		// First, filter high-value transactions
		FilterMap(func(t Transaction) Option[Transaction] {
			time.Sleep(5 * time.Millisecond)
			// Keep transactions over $100
			if t.Amount > 100.0 {
				return Some(t)
			}
			return None[Transaction]()
		}).
		// Sample every 2nd transaction
		StepBy(2).
		// Find maximum amount
		MaxBy(func(a, b Transaction) cmp.Ordering {
			return cmp.Cmp(a.Amount, b.Amount)
		})

	duration5 := time.Since(start)
	Println("   Complex pipeline result: {}", finalResults)
	if finalResults.IsSome() {
		result := finalResults.Some()
		Println("   Highest sampled amount: ${} from category {}", result.Amount, result.Category)
	}
	Println("   Pipeline completed in {}", duration5)

	totalDuration := time.Since(totalStart)
	Println("\n=== Pipeline Summary ===")
	Println("Total processing time: {}", totalDuration)
	Println("Individual steps:")
	Println("  1. FlatMap expansion: {}", duration1)
	Println("  2. FilterMap processing: {}", duration2)
	Println("  3. StepBy sampling: {}", duration3)
	Println("  4. MaxBy/MinBy analysis: {}", duration4)
	Println("  5. Combined pipeline: {}", duration5)

	// Performance comparison note
	Println("\nNote: Without parallelism, this would take approximately:")
	estimatedSequential := time.Duration(
		int(transactions.Len()) * (20 + 15 + 10 + 12*3 + 5) * int(time.Millisecond), // Sum of all sleep times
	)
	Println("  Sequential estimate: {}", estimatedSequential)
	Println("  Actual parallel time: {}", totalDuration)
	Println("  Approximate speedup: {}x",
		Float(estimatedSequential.Milliseconds())/Float(totalDuration.Milliseconds()))
}
