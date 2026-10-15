package analyzer

import "github.com/Lazynx/spent/internal/statement"

const topMerchantsLimit = 5

func Analyze(s statement.Statement, c Categorizer) (Report, error) {
	return Report{}, errNotImplemented
}

func checkBalance(s statement.Statement) error {
	return errNotImplemented
}
