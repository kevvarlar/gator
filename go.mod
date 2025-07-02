module github.com/kevvarlar/gator

go 1.24.4

require (
	github.com/google/uuid v1.6.0
	github.com/kevvarlar/gator/internal/config v0.0.0-20250701200540-509742bbd8a6
	github.com/kevvarlar/gator/internal/database v0.0.0
	github.com/lib/pq v1.10.9
)

replace (
	github.com/kevvarlar/gator/internal/config v0.0.0-20250701200540-509742bbd8a6 => ./internal/config
	github.com/kevvarlar/gator/internal/database v0.0.0 => ./internal/database
)
