package policy

import (
	"time"
	"windops/internal/fault"
)

func EvaluateCampaignFitsBudget(ctx Context) (Result, error) {
	if err := requireTenant(ctx, "policy.campaign_budget"); err != nil {
		return Result{}, err
	}
	if ctx.Budget < 0 || ctx.Cost < 0 {
		return Result{}, fault.New(fault.CodeInvalid, "policy.campaign_budget", "budget and cost cannot be negative")
	}
	lineTotal := int64(0)
	for _, value := range ctx.Values {
		if value < 0 {
			return Result{}, fault.New(fault.CodeInvalid, "policy.campaign_budget", "line cost cannot be negative")
		}
		nextTotal := lineTotal + value
		if nextTotal < lineTotal {
			return Result{}, fault.New(fault.CodeInvalid, "policy.campaign_budget", "line cost total overflowed")
		}
		lineTotal = nextTotal
	}
	total := ctx.Cost + lineTotal
	if total > ctx.Budget {
		result := deny("budget_exceeded", "campaign cost exceeds approved budget")
		result.Quantity = total - ctx.Budget
		return result, nil
	}
	result := allow("campaign is within approved budget")
	result.Quantity = ctx.Budget - total
	return result, nil
}

var _ = time.Second
var _ = fault.CodeInvalid
