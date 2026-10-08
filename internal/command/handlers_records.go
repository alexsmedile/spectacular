package command

import (
	"encoding/json"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/missionbundle"
)

func (r Runner) executeRecords(op operation, input operationInput, service missionbundle.Service) (any, error) {
	rest := input.rest
	dataPayload := input.dataPayload
	dryRun := input.dryRun
	fromFile := input.fromFile
	override := input.override
	ws := service.Workspace
	var value any
	var err error
	switch op {
	case opReviewRecord:
		missionRef := rest[0]
		if dataPayload != "" {
			var draft missionbundle.ReviewDraft
			if unmarshalErr := json.Unmarshal([]byte(dataPayload), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			value, err = service.RecordReviewDraft(missionRef, draft)
		} else if len(rest) > 1 && strings.HasPrefix(strings.TrimSpace(rest[1]), "{") {
			var draft missionbundle.ReviewDraft
			if unmarshalErr := json.Unmarshal([]byte(rest[1]), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			value, err = service.RecordReviewDraft(missionRef, draft)
		} else if len(rest) > 1 {
			stdin, readErr := r.stdinIfNeeded(rest[1])
			if readErr != nil {
				err = readErr
				break
			}
			value, err = service.RecordReview(missionRef, inputPath(r.Cwd, rest[1]), stdin)
		} else {
			draft := missionbundle.ReviewDraft{
				Type:   "ReviewDraft",
				Title:  "Mission review for " + missionRef,
				Status: "passed",
			}
			value, err = service.RecordReviewDraft(missionRef, draft)
		}
	case opHandoffRecord:
		missionRef := rest[0]
		sender := ""
		for i := 1; i < len(rest); i++ {
			if (rest[i] == "--by" || rest[i] == "--sender") && i+1 < len(rest) {
				sender = rest[i+1]
			}
		}
		if dataPayload != "" {
			var draft missionbundle.HandoffDraft
			if unmarshalErr := json.Unmarshal([]byte(dataPayload), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			value, err = service.RecordHandoffDraft(missionRef, draft, sender)
		} else if len(rest) > 1 && strings.HasPrefix(strings.TrimSpace(rest[1]), "{") {
			var draft missionbundle.HandoffDraft
			if unmarshalErr := json.Unmarshal([]byte(rest[1]), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			value, err = service.RecordHandoffDraft(missionRef, draft, sender)
		} else if len(rest) > 1 && rest[1] != "--by" && !strings.HasPrefix(rest[1], "--") {
			stdin, readErr := r.stdinIfNeeded(rest[1])
			if readErr != nil {
				err = readErr
				break
			}
			value, err = service.RecordHandoff(missionRef, inputPath(r.Cwd, rest[1]), sender, stdin)
		} else {
			draft := missionbundle.HandoffDraft{
				Type:  "HandoffDraft",
				Title: "Handoff for " + missionRef,
				Task:  "Continue mission implementation",
			}
			value, err = service.RecordHandoffDraft(missionRef, draft, sender)
		}
	case opEvidenceRecord:
		missionRef := rest[0]
		if dataPayload != "" {
			var draft missionbundle.EvidenceDraft
			if unmarshalErr := json.Unmarshal([]byte(dataPayload), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			value, err = service.RecordEvidenceDraft(missionRef, draft)
		} else if len(rest) > 1 && strings.HasPrefix(strings.TrimSpace(rest[1]), "{") {
			var draft missionbundle.EvidenceDraft
			if unmarshalErr := json.Unmarshal([]byte(rest[1]), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			value, err = service.RecordEvidenceDraft(missionRef, draft)
		} else {
			targetPath := ""
			if len(rest) > 1 && !strings.HasPrefix(rest[1], "--") {
				targetPath = rest[1]
			}
			stdin, readErr := r.stdinIfNeeded(targetPath)
			if readErr != nil {
				err = readErr
				break
			}
			var absTarget string
			if targetPath != "" {
				absTarget = inputPath(r.Cwd, targetPath)
			}
			var absFrom string
			if fromFile != "" {
				absFrom = inputPath(r.Cwd, fromFile)
			}
			value, err = service.RecordEvidence(missionRef, absTarget, stdin, absFrom)
		}
	case opProposalCheck:
		value, err = missionbundle.ValidateProposal(ws, rest[0])
	case opContractAmend:
		value, err = service.AmendContract(rest[0], rest[2], rest[4], override, dryRun)
	case opContractCreate:
		title := ""
		if len(rest) == 3 && rest[1] == "--title" {
			title = rest[2]
		}
		value, err = service.CreateContract(rest[0], title, ws.Config.Defaults.Operator)
	case opDecide:
		flags, parseErr := parseDecideArgs(rest)
		if parseErr != nil {
			err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", parseErr.Error(), nil)
			break
		}
		if dataPayload != "" {
			flags.jsonPayload = dataPayload
		}
		if flags.jsonPayload != "" {
			var draft missionbundle.DecisionDraft
			if unmarshalErr := json.Unmarshal([]byte(flags.jsonPayload), &draft); unmarshalErr != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+unmarshalErr.Error(), nil)
				break
			}
			actor := draft.Actor
			if actor == "" {
				actor = flags.actor
			}
			if actor == "" {
				actor = ws.Config.Defaults.Operator
			}
			if actor == "" {
				actor = "Alex"
			}
			draft.Actor = actor
			if draft.ActorRole == "" {
				draft.ActorRole = "owner"
			}
			value, err = service.RecordDecisionDraft(draft)
		} else if flags.file != "" {
			stdin, readErr := r.stdinIfNeeded(flags.file)
			if readErr != nil {
				err = readErr
				break
			}
			value, err = service.RecordDecision(inputPath(r.Cwd, flags.file), stdin)
		} else {
			if strings.TrimSpace(flags.title) == "" || strings.TrimSpace(flags.disposition) == "" || strings.TrimSpace(flags.rationale) == "" {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "direct decision recording requires --title, --disposition, and --rationale", nil)
				break
			}
			actor := flags.actor
			if actor == "" {
				actor = ws.Config.Defaults.Operator
			}
			if actor == "" {
				actor = "Alex"
			}
			draft := missionbundle.DecisionDraft{
				Type:        "DecisionDraft",
				Title:       flags.title,
				Disposition: flags.disposition,
				Rationale:   flags.rationale,
				Actor:       actor,
				ActorRole:   "owner",
				Question:    flags.question,
				Supersedes:  flags.supersedes,
				Scope:       flags.scope,
				Targets:     flags.targets,
			}
			value, err = service.RecordDecisionDraft(draft)
		}
	}
	return value, err
}
