package main

import (
	. "github.com/enetx/g"
)

func main() {
	// Example 1: FilterMap for config validation
	configs := NewMap[String, String]()
	configs.Set("host", "localhost")
	configs.Set("port", "8080")
	configs.Set("debug", "invalid")
	configs.Set("timeout", "30")

	validConfigs := configs.Iter().
		FilterMap(func(k, v String) Option[Pair[String, String]] {
			// Keep only port and host configs with validation suffix
			if k == "port" || k == "host" {
				return Some(Pair[String, String]{Key: k, Value: v + "_validated"})
			}

			return None[Pair[String, String]]()
		}).
		Collect()

	validConfigs.Println() // Map[host:localhost_validated port:8080_validated]

	// Example 2: FilterMap for user age filtering
	users := NewMap[String, Int]()
	users.Set("alice", 25)
	users.Set("bob", 17)
	users.Set("charlie", 30)
	users.Set("diana", 16)

	adults := users.Iter().
		FilterMap(func(name String, age Int) Option[Pair[String, Int]] {
			// Keep only users who are adults (18+)
			if age >= 18 {
				return Some(Pair[String, Int]{Key: name, Value: age})
			}
			return None[Pair[String, Int]]()
		}).
		Collect()

	adults.Println() // Map[alice:25 charlie:30]

	// Example 3: FilterMap for URL validation
	urls := NewMap[String, String]()
	urls.Set("google", "https://google.com")
	urls.Set("invalid", "not-a-url")
	urls.Set("github", "https://github.com")
	urls.Set("empty", "")

	validUrls := urls.Iter().
		FilterMap(func(name, url String) Option[Pair[String, String]] {
			// Keep only valid HTTPS URLs
			if String(url).StartsWith("https://") {
				return Some(Pair[String, String]{Key: name + "_secure", Value: url})
			}
			return None[Pair[String, String]]()
		}).
		Collect()

	validUrls.Println() // Map[github_secure:https://github.com google_secure:https://google.com]
}
