// Package runpolicy owns pure Run lifecycle rules, without workspace effects.
package runpolicy

import (
	"fmt"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"strings"
)

var validRunTransitions = map[string][]string{
	"active":          {"paused", "blocked", "awaiting-review", "completed", "stopped"},
	"paused":          {"active", "blocked", "stopped"},
	"blocked":         {"active", "stopped"},
	"awaiting-review": {"active", "completed", "stopped"},
	"completed":       {}, // terminal
	"stopped":         {}, // terminal
}

func ValidateTransition(from, to string) error {
	allowed, ok := validRunTransitions[from]
	if !ok {
		return domain.NewRefusal(domain.RefusalInvalidTransition, "status", fmt.Sprintf("unknown origin run state %q", from), nil)
	}
	for _, a := range allowed {
		if a == to {
			return nil
		}
	}
	if len(allowed) == 0 {
		return domain.NewRefusal(domain.RefusalInvalidTransition, "status", fmt.Sprintf("state %q is terminal and cannot transition to %q", from, to), nil)
	}
	return domain.NewRefusal(domain.RefusalInvalidTransition, "status", fmt.Sprintf("illegal transition from %q to %q (allowed: %s)", from, to, strings.Join(allowed, ", ")), nil)
}
