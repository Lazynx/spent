package analyzer

import "github.com/Lazynx/spent/internal/statement"

type Category string

const (
	Food          Category = "Food"
	Transport     Category = "Transport"
	Shopping      Category = "Shopping"
	Entertainment Category = "Entertainment"
	Utilities     Category = "Utilities"
	Health        Category = "Health"
	Other         Category = "Other"
)

type Categorizer interface {
	Categorize(t statement.Transaction) Category
}
