package main

import (
	"time"

	. "github.com/enetx/g"
	"github.com/enetx/g/cmp"
)

func main() {
	Println("=== Real-World Parallel Iterator Use Cases ===\n")

	// Use Case 1: Log Processing Pipeline
	Println("1. Log Processing Pipeline:")
	logEntries := SliceOf[String](
		"2024-01-15 10:30:15 INFO User login: alice@company.com",
		"2024-01-15 10:30:16 ERROR Database connection failed: timeout",
		"2024-01-15 10:30:17 INFO User login: bob@company.com",
		"2024-01-15 10:30:18 WARN High memory usage: 85%",
		"2024-01-15 10:30:19 INFO API request: /users/profile",
		"2024-01-15 10:30:20 ERROR API rate limit exceeded",
		"2024-01-15 10:30:21 INFO User logout: alice@company.com",
		"2024-01-15 10:30:22 DEBUG Cache hit: user_123",
		"2024-01-15 10:30:23 ERROR Failed to send email notification",
		"2024-01-15 10:30:24 INFO System backup completed",
	)

	start := time.Now()

	// Process logs: extract alerts, enrich with metadata, sample for monitoring
	alerts := logEntries.Iter().
		Parallel(3).
		// Extract error and warning logs
		FilterMap(func(log String) Option[String] {
			time.Sleep(5 * time.Millisecond) // Simulate log parsing
			if log.Contains("ERROR") || log.Contains("WARN") {
				// Extract timestamp and message
				parts := log.Split(" ").Collect()
				if parts.Len() >= 4 {
					timestamp := parts[0].Append(" ").Append(parts[1])
					level := parts[2]
					message := log.Split(level.Append(" ")).Collect()[1]
					return Some(String("ALERT[").Append(timestamp).Append("]: ").Append(message))
				}
			}
			return None[String]()
		}).
		// Sample every 2nd alert for dashboard
		StepBy(2).
		Collect()

	logDuration := time.Since(start)
	Println("   Found {} alerts in {}", alerts.Len(), logDuration)
	alerts.Iter().ForEach(func(alert String) {
		Println("   {}", alert)
	})
	Println("")

	// Use Case 2: E-commerce Order Processing
	Println("2. E-commerce Order Processing:")

	type Order struct {
		ID     int
		UserID int
		Amount float64
		Status String
		Items  Slice[String]
	}

	orders := SliceOf(
		Order{1001, 501, 149.99, "pending", SliceOf[String]("laptop_case", "mouse")},
		Order{1002, 502, 25.50, "shipped", SliceOf[String]("book")},
		Order{1003, 501, 899.00, "pending", SliceOf[String]("laptop", "warranty")},
		Order{1004, 503, 15.99, "delivered", SliceOf[String]("pen", "notebook")},
		Order{1005, 504, 299.99, "pending", SliceOf[String]("monitor", "cables")},
		Order{1006, 502, 45.00, "cancelled", SliceOf[String]("keyboard")},
	)

	start = time.Now()

	// Find high-value pending orders
	highValueOrders := orders.Iter().
		Parallel(3).
		// Only process pending orders over $100
		FilterMap(func(order Order) Option[Order] {
			time.Sleep(8 * time.Millisecond) // Simulate order validation
			if order.Status == "pending" && order.Amount > 100.0 {
				return Some(order)
			}
			return None[Order]()
		}).
		Collect()

	// Find most expensive pending order
	expensiveOrder := orders.Iter().
		Parallel(2).
		FilterMap(func(order Order) Option[Order] {
			if order.Status == "pending" {
				return Some(order)
			}
			return None[Order]()
		}).
		MaxBy(func(a, b Order) cmp.Ordering {
			return cmp.Cmp(a.Amount, b.Amount)
		})

	orderDuration := time.Since(start)
	Println("   Processed {} high-value orders in {}", highValueOrders.Len(), orderDuration)
	Println("   Sample orders:")
	highValueOrders.Iter().Take(3).ForEach(func(order Order) {
		itemsList := ""
		order.Items.Iter().ForEach(func(item String) {
			if itemsList != "" {
				itemsList += ", "
			}
			itemsList += string(item)
		})
		Println("     Order {}: ${} - Items: {}", order.ID, order.Amount, itemsList)
	})
	if expensiveOrder.IsSome() {
		order := expensiveOrder.Some()
		Println("   Most expensive pending: Order {} - ${}", order.ID, order.Amount)
	}
	Println("")

	// Use Case 3: Data Analytics - Website Performance
	Println("3. Website Performance Analytics:")

	type PageMetric struct {
		URL          String
		LoadTime     int // milliseconds
		UserAgent    String
		ResponseCode int
	}

	metrics := SliceOf(
		PageMetric{"/home", 250, "Chrome", 200},
		PageMetric{"/products", 800, "Firefox", 200},
		PageMetric{"/home", 180, "Safari", 200},
		PageMetric{"/checkout", 1200, "Chrome", 500}, // Error
		PageMetric{"/products", 300, "Chrome", 200},
		PageMetric{"/profile", 150, "Firefox", 200},
		PageMetric{"/home", 220, "Chrome", 200},
		PageMetric{"/api/data", 2000, "Chrome", 503}, // Error
		PageMetric{"/products", 450, "Safari", 200},
		PageMetric{"/checkout", 600, "Firefox", 200},
	)

	start = time.Now()

	// Analyze performance issues
	slowMetrics := metrics.Iter().
		Parallel(4).
		// Identify slow pages and errors
		FilterMap(func(metric PageMetric) Option[PageMetric] {
			time.Sleep(2 * time.Millisecond) // Simulate analysis
			if metric.LoadTime > 500 || metric.ResponseCode >= 400 {
				return Some(metric)
			}
			return None[PageMetric]()
		}).
		Collect()

	// Find slowest page
	slowestPage := metrics.Iter().
		Parallel(3).
		MaxBy(func(a, b PageMetric) cmp.Ordering {
			return cmp.Cmp(a.LoadTime, b.LoadTime)
		})

	// Sample every 3rd metric for detailed analysis
	sampleMetrics := metrics.Iter().
		Parallel(2).
		StepBy(3).
		Collect()

	analyticsDuration := time.Since(start)
	Println("   Found {} slow/error metrics in {}", slowMetrics.Len(), analyticsDuration)
	slowMetrics.Iter().ForEach(func(metric PageMetric) {
		severity := "SLOW"
		if metric.ResponseCode >= 400 {
			severity = "ERROR"
		}
		Println("     {}: {} - {}ms ({})", severity, metric.URL, metric.LoadTime, metric.ResponseCode)
	})
	if slowestPage.IsSome() {
		page := slowestPage.Some()
		Println("   Slowest page: {} ({}ms)", page.URL, page.LoadTime)
	}
	Println("   Sample metrics: {} entries", sampleMetrics.Len())
	Println("")

	// Use Case 4: Financial Data Processing
	Println("4. Financial Transaction Analysis:")

	type Transaction struct {
		ID        int
		AccountID int
		Amount    float64
		Type      String // "debit" or "credit"
		Merchant  String
	}

	transactions := SliceOf(
		Transaction{1, 101, 50.00, "debit", "Coffee Shop"},
		Transaction{2, 102, 1200.00, "credit", "Salary Deposit"},
		Transaction{3, 101, 25.99, "debit", "Online Store"},
		Transaction{4, 103, 75.50, "debit", "Gas Station"},
		Transaction{5, 101, 800.00, "credit", "Tax Refund"},
		Transaction{6, 102, 150.00, "debit", "Grocery Store"},
		Transaction{7, 104, 2000.00, "credit", "Investment Return"},
		Transaction{8, 103, 45.00, "debit", "Restaurant"},
		Transaction{9, 101, 300.00, "debit", "Utilities"},
		Transaction{10, 105, 5000.00, "credit", "Bonus Payment"},
	)

	start = time.Now()

	// Risk analysis: large transactions, unusual patterns
	riskTransactions := transactions.Iter().
		Parallel(4).
		// Identify risky transactions
		FilterMap(func(tx Transaction) Option[Transaction] {
			time.Sleep(3 * time.Millisecond) // Simulate risk calculation
			// Flag high-value or unusual transactions
			if tx.Amount > 1000.00 || (tx.Type == "credit" && tx.Amount > 500.00) ||
				(tx.Type == "debit" && tx.Amount > 300.00) {
				return Some(tx)
			}
			return None[Transaction]()
		}).
		// Sample for manual review
		StepBy(2).
		Collect()

	// Find largest transaction by type
	largestCredit := transactions.Iter().
		Parallel(2).
		FilterMap(func(tx Transaction) Option[Transaction] {
			if tx.Type == "credit" {
				return Some(tx)
			}
			return None[Transaction]()
		}).
		MaxBy(func(a, b Transaction) cmp.Ordering {
			return cmp.Cmp(a.Amount, b.Amount)
		})

	financialDuration := time.Since(start)
	Println("   Identified {} risky transactions in {}", riskTransactions.Len(), financialDuration)
	Println("   Risk transactions:")
	riskTransactions.Iter().Take(5).ForEach(func(tx Transaction) {
		riskType := "HIGH_VALUE"
		if tx.Amount > 1000.00 {
			riskType = "VERY_HIGH_VALUE"
		}
		Println("     {}: Transaction {} - ${} {} from {}", riskType, tx.ID, tx.Amount, tx.Type, tx.Merchant)
	})
	if largestCredit.IsSome() {
		tx := largestCredit.Some()
		Println("   Largest credit: ${} from {}", tx.Amount, tx.Merchant)
	}

	totalTime := time.Now().
		Sub(start.Add(-logDuration).Add(-orderDuration).Add(-analyticsDuration).Add(-financialDuration))
	Println("\n=== Performance Summary ===")
	Println("Log processing:      {}", logDuration)
	Println("Order processing:    {}", orderDuration)
	Println("Analytics:           {}", analyticsDuration)
	Println("Financial analysis:  {}", financialDuration)
	Println("Total execution:     {}", totalTime)
	Println("\nAll use cases demonstrate the power of parallel iterators for")
	Println("real-world data processing pipelines! 🚀")

	/* Expected Output:
	=== Real-World Parallel Iterator Use Cases ===

	1. Log Processing Pipeline:
	   Found 2 alerts in 13ms
	   ALERT[2024-01-15 10:30:16]: Database connection failed: timeout
	   ALERT[2024-01-15 10:30:20]: API rate limit exceeded

	2. E-commerce Order Processing:
	   Processed 6 high-value items in 28ms
	   Sample items:
	     Order_1001_User_501_Item_laptop_case_Value_149.99
	     Order_1001_User_501_Item_mouse_Value_149.99
	     Order_1003_User_501_Item_laptop_Value_899
	   Most expensive pending: Order 1003 - $899

	3. Website Performance Analytics:
	   Found 4 performance issues in 8ms
	     SLOW_/products_800ms_200
	     ERROR_/checkout_1200ms_500
	     ERROR_/api/data_2000ms_503
	     SLOW_/products_450ms_200
	   Slowest page: /api/data (2000ms)
	   Sample metrics: Slice[/home_250ms, /checkout_1200ms, /profile_150ms, /checkout_600ms]

	4. Financial Transaction Analysis:
	   Generated 11 risk assessment items in 13ms
	   Sample risk items:
	     TX_1_50
	     HIGH_VALUE_2
	     LARGE_DEPOSIT_2
	     TX_3_25.99
	     TX_4_75.5
	   Largest credit: $5000 from Bonus Payment

	=== Performance Summary ===
	Log processing:      13ms
	Order processing:    28ms
	Analytics:           8ms
	Financial analysis:  13ms
	Total execution:     62ms

	All use cases demonstrate the power of parallel iterators for
	real-world data processing pipelines! 🚀
	*/
}
