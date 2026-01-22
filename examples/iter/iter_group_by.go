package main

import (
	"time"

	. "github.com/enetx/g"
)

// Person represents an employee with basic information
type Person struct {
	Name       String
	Age        Int
	Department String
	Salary     Int
	HireDate   time.Time
}

// Transaction represents a financial transaction
type Transaction struct {
	ID     Int
	Amount Float
	Type   String // "income", "expense"
	Date   time.Time
}

// LogEntry represents a log entry with timestamp, level, and metadata
type LogEntry struct {
	Timestamp time.Time
	Level     String // "INFO", "WARN", "ERROR"
	Message   String
	Source    String
}

func main() {
	// Example 1: Grouping employees by department
	employees := SliceOf(
		Person{Name: "Alice", Age: 30, Department: "Engineering", Salary: 80000},
		Person{Name: "Andy", Age: 30, Department: "Engineering", Salary: 90000},
		Person{Name: "Bob", Age: 25, Department: "Engineering", Salary: 75000},
		Person{Name: "Charlie", Age: 35, Department: "Engineering", Salary: 90000},
		Person{Name: "Diana", Age: 28, Department: "Marketing", Salary: 60000},
		Person{Name: "Eve", Age: 32, Department: "Marketing", Salary: 65000},
		Person{Name: "Frank", Age: 29, Department: "Sales", Salary: 55000},
		Person{Name: "Grace", Age: 31, Department: "Sales", Salary: 58000},
		Person{Name: "Henry", Age: 27, Department: "Sales", Salary: 52000},
	)

	Println("=== Grouping employees by department ===")
	// Group consecutive employees with the same department
	departmentGroups := employees.Iter().GroupBy(func(a, b Person) bool {
		return a.Department.Eq(b.Department)
	}).Collect()

	for _, group := range departmentGroups {
		if !group.IsEmpty() {
			Println("Department {}: {} employees", group[0].Department, group.Len())
			for _, emp := range group {
				Println("  - {} (age: {}, salary: ${})", emp.Name, emp.Age, emp.Salary)
			}
		}
	}

	// Example 2: Grouping by age categories
	Println("\n=== Grouping by age categories ===")
	ageGroups := employees.Iter().GroupBy(func(a, b Person) bool {
		// Helper function to determine age category
		ageCategory := func(age Int) string {
			if age.Lt(30) {
				return "young"
			}
			return "experienced"
		}
		return ageCategory(a.Age) == ageCategory(b.Age)
	}).Collect()

	for _, group := range ageGroups {
		if !group.IsEmpty() {
			category := "young"
			if group[0].Age.Gte(30) {
				category = "experienced"
			}

			Println("Category '{}': {} people", category, group.Len())
			for _, emp := range group {
				Println("  - {} ({} years old)", emp.Name, emp.Age)
			}
		}
	}

	// Example 3: Grouping transactions by type
	now := time.Now()
	transactions := SliceOf(
		Transaction{ID: 1, Amount: 1000, Type: "income", Date: now.AddDate(0, 0, -5)},
		Transaction{ID: 2, Amount: 1500, Type: "income", Date: now.AddDate(0, 0, -4)},
		Transaction{ID: 3, Amount: -200, Type: "expense", Date: now.AddDate(0, 0, -3)},
		Transaction{ID: 4, Amount: -150, Type: "expense", Date: now.AddDate(0, 0, -2)},
		Transaction{ID: 5, Amount: -300, Type: "expense", Date: now.AddDate(0, 0, -1)},
		Transaction{ID: 6, Amount: 2000, Type: "income", Date: now},
	)

	Println("\n=== Grouping transactions by type ===")
	// Group consecutive transactions of the same type
	transactionGroups := transactions.Iter().GroupBy(func(a, b Transaction) bool {
		return a.Type.Eq(b.Type)
	}).Collect()

	for _, group := range transactionGroups {
		if !group.IsEmpty() {
			var total Float
			// Calculate total amount for each group
			for _, t := range group {
				total += t.Amount
			}

			Println("Type '{}': {} transactions, total amount: ${.RoundDecima(2)}", group[0].Type, group.Len(), total)
		}
	}

	// Example 4: Grouping logs by severity level
	logs := SliceOf(
		LogEntry{Timestamp: now.Add(-1 * time.Hour), Level: "INFO", Message: "Server started", Source: "main"},
		LogEntry{Timestamp: now.Add(-50 * time.Minute), Level: "INFO", Message: "User logged in", Source: "auth"},
		LogEntry{Timestamp: now.Add(-40 * time.Minute), Level: "WARN", Message: "High memory usage", Source: "monitor"},
		LogEntry{
			Timestamp: now.Add(-30 * time.Minute),
			Level:     "ERROR",
			Message:   "Database connection failed",
			Source:    "db",
		},
		LogEntry{Timestamp: now.Add(-20 * time.Minute), Level: "ERROR", Message: "API timeout", Source: "api"},
		LogEntry{Timestamp: now.Add(-10 * time.Minute), Level: "INFO", Message: "Backup completed", Source: "backup"},
	)

	Println("\n=== Grouping logs by level ===")
	// Group consecutive log entries with the same severity level
	logGroups := logs.Iter().GroupBy(func(a, b LogEntry) bool {
		return a.Level == b.Level
	}).Collect()

	for _, group := range logGroups {
		if !group.IsEmpty() {
			Println("Level {}: {} entries", group[0].Level, group.Len())
			for _, log := range group {
				Println("  - [{}] {}: {}",
					log.Timestamp.Format("15:04"), log.Source, log.Message)
			}
		}
	}

	// Example 5: Complex grouping - by salary trend (descending or equal)
	Println("\n=== Grouping by salary trends ===")
	salaryTrendGroups := employees.Iter().GroupBy(func(a, b Person) bool {
		// Group while salary doesn't increase (descending or equal trend)
		return a.Salary >= b.Salary
	}).Collect()

	for i, group := range salaryTrendGroups {
		Println("Group {} (salary trend):", i+1)
		for j, emp := range group {
			trend := ""
			if j > 0 {
				// Determine trend direction compared to previous employee
				if emp.Salary > group[j-1].Salary {
					trend = " ↗" // increasing
				} else if emp.Salary < group[j-1].Salary {
					trend = " ↘" // decreasing
				} else {
					trend = " →" // equal
				}
			}
			Println("  - {} - ${}{}", emp.Name, emp.Salary, trend)
		}
	}

	// Example 6: Grouping by first letter of name
	Println("\n=== Grouping by first letter of name ===")
	nameGroups := employees.Iter().GroupBy(func(a, b Person) bool {
		// Group consecutive employees whose names start with the same letter
		return a.Name[0] == b.Name[0]
	}).Collect()

	for _, group := range nameGroups {
		if !group.IsEmpty() {
			firstLetter := string(group[0].Name[0])
			Println("Names starting with '{}': {}", firstLetter, group.Len())
			// Extract names and join them with commas
			names := TransformSlice(group, func(person Person) String { return person.Name }).Join(", ")
			Println("  {}", names)
		}
	}
}
