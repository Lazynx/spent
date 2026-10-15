package analyzer

import "github.com/Lazynx/spent/internal/statement"

type keywordRule struct {
	keyword  string
	category Category
}

type RuleCategorizer struct {
	mccRules map[string]Category
	keywords []keywordRule
}

func NewRuleCategorizer() *RuleCategorizer {
	return &RuleCategorizer{
		mccRules: map[string]Category{},
		keywords: []keywordRule{},
	}
}

func (c *RuleCategorizer) Categorize(t statement.Transaction) Category {
	return Other
}
