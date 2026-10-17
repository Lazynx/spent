package analyzer

import (
	"strings"

	"github.com/Lazynx/spent/internal/statement"
)

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
		mccRules: map[string]Category{
			"5411": Food,
			"5812": Food,
			"5814": Food,
			"4111": Transport,
			"4121": Transport,
			"5541": Transport,
			"5311": Shopping,
			"5651": Shopping,
			"5732": Shopping,
			"7832": Entertainment,
			"7841": Entertainment,
			"5815": Entertainment,
			"4900": Utilities,
			"4814": Utilities,
			"5912": Health,
			"8011": Health,
			"8062": Health,
		},
		keywords: []keywordRule{
			{"magnum", Food},
			{"glovo", Food},
			{"wolt", Food},
			{"coffee", Food},
			{"kfc", Food},
			{"yandex go", Transport},
			{"uber", Transport},
			{"indrive", Transport},
			{"onay", Transport},
			{"wildberries", Shopping},
			{"technodom", Shopping},
			{"sulpak", Shopping},
			{"zara", Shopping},
			{"netflix", Entertainment},
			{"spotify", Entertainment},
			{"steam", Entertainment},
			{"kinopark", Entertainment},
			{"kazakhtelecom", Utilities},
			{"beeline", Utilities},
			{"kcell", Utilities},
			{"altel", Utilities},
			{"alseco", Utilities},
			{"pharm", Health},
			{"clinic", Health},
			{"hospital", Health},
		},
	}
}

func (c *RuleCategorizer) Categorize(t statement.Transaction) Category {
	if cat, ok := c.mccRules[t.MCC]; ok {
		return cat
	}
	desc := normalize(t.Description)
	for _, r := range c.keywords {
		if strings.Contains(desc, r.keyword) {
			return r.category
		}
	}
	return Other
}
