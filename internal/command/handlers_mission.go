package command

import (
	"encoding/json"
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/missionbundle"
)

func (r Runner) executeMission(op operation, input operationInput, service missionbundle.Service) (any, error) {
	rest := input.rest
	dataPayload := input.dataPayload
	allowMain := input.allowMain
	createBranch := input.createBranch
	statusFilter := input.statusFilter
	allMissions := input.allMissions
	dryRun := input.dryRun
	ws := service.Workspace
	var value any
	var err error
	switch op {
	case opMissionStart:
		var plan missionbundle.Plan
		var raw []byte
		if dataPayload != "" {
			plan, raw, err = missionbundle.ReadPlan("", []byte(dataPayload))
		} else if len(rest) > 0 && strings.HasPrefix(strings.TrimSpace(rest[0]), "{") {
			plan, raw, err = missionbundle.ReadPlan("", []byte(rest[0]))
		} else if len(rest) > 0 {
			path := inputPath(r.Cwd, rest[0])
			stdin, readErr := r.stdinIfNeeded(rest[0])
			if readErr != nil {
				err = readErr
				break
			}
			plan, raw, err = missionbundle.ReadPlan(path, stdin)
		} else {
			err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "requires a plan file, stdin (-), inline JSON, or --data", nil)
			break
		}
		if allowMain {
			plan.AllowMain = true
		}
		if createBranch {
			plan.CreateBranch = true
		}
		if err == nil {
			value, err = service.Start(plan, raw)
		}
	case opMissionList:
		value, err = service.ListMissions(statusFilter, allMissions)
	case opMissionShow:
		value, err = service.Show(rest[0])
	case opMissionCheck:
		targetRef := ""
		verifyMode := false
		for _, arg := range rest {
			if arg == "--verify" {
				verifyMode = true
			} else if targetRef == "" {
				targetRef = arg
			}
		}
		if targetRef == "" {
			err = domain.NewRefusal(domain.RefusalInvalidReference, "", "expected <mission-ref>", nil)
			break
		}
		if verifyMode {
			value, err = service.CheckWithVerify(targetRef)
		} else {
			value, err = service.Check(targetRef)
		}
	case opMissionAmendScope:
		addPaths := strings.Split(rest[2], ",")
		var cleaned []string
		for _, p := range addPaths {
			if t := strings.TrimSpace(p); t != "" {
				cleaned = append(cleaned, t)
			}
		}
		reason := ""
		if len(rest) == 7 && rest[5] == "--reason" {
			reason = rest[6]
		}
		value, err = service.AmendScope(rest[0], cleaned, rest[4], reason, dryRun)
	case opMissionClose:
		by := ""
		if len(rest) == 3 && rest[1] == "--by" {
			by = rest[2]
		}
		if by == "" {
			by = ws.Config.Defaults.Operator
		}
		if by == "" {
			by = "Alex"
		}
		value, err = service.CloseMission(rest[0], by)
	case opObjectiveShow:
		var objective missionbundle.Objective
		objective, _, err = service.Objective(rest[0])
		value = objective
	case opObjectivePromote:
		value, err = service.PromoteObjective(rest[0])
	case opObjectiveFinish:
		value, err = service.FinishObjective(rest[0])
	case opRunShow:
		var run missionbundle.Run
		run, _, err = service.Run(rest[0])
		value = run
	case opRunStart:
		value, err = service.StartRun(rest[0], rest[2])
	case opMissionComplete:
		by := ""
		if len(rest) == 3 && rest[1] == "--by" {
			by = rest[2]
		}
		if by == "" {
			by = ws.Config.Defaults.Operator
		}
		if by == "" {
			by = "Alex"
		}
		value, err = service.Complete(rest[0], by)
	case opRunTransition:
		targetRef := rest[0]
		if dataPayload != "" || (len(rest) > 0 && strings.HasPrefix(strings.TrimSpace(rest[len(rest)-1]), "{")) {
			rawJSON := dataPayload
			if rawJSON == "" {
				rawJSON = strings.TrimSpace(rest[len(rest)-1])
			}
			var payload struct {
				Target     string `json:"target"`
				To         string `json:"to"`
				By         string `json:"by"`
				Reason     string `json:"reason"`
				NextAction string `json:"next_action"`
			}
			if err := json.Unmarshal([]byte(rawJSON), &payload); err != nil {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "invalid JSON payload: "+err.Error(), nil)
				break
			}
			if payload.Target != "" {
				targetRef = payload.Target
			}
			actor := payload.By
			if actor == "" {
				actor = ws.Config.Defaults.Operator
			}
			if actor == "" {
				actor = "Alex"
			}
			value, err = service.TransitionRun(targetRef, payload.To, actor, payload.Reason, payload.NextAction)
		} else {
			toState := ""
			actor := ""
			reason := ""
			nextAction := ""
			for i := 1; i < len(rest); i++ {
				switch rest[i] {
				case "--to":
					if i+1 < len(rest) {
						toState = rest[i+1]
						i++
					}
				case "--by":
					if i+1 < len(rest) {
						actor = rest[i+1]
						i++
					}
				case "--reason":
					if i+1 < len(rest) {
						reason = rest[i+1]
						i++
					}
				case "--next-action":
					if i+1 < len(rest) {
						nextAction = rest[i+1]
						i++
					}
				}
			}
			if actor == "" {
				actor = ws.Config.Defaults.Operator
			}
			if actor == "" {
				actor = "Alex"
			}
			if toState == "" || reason == "" {
				err = domain.NewRefusal(domain.RefusalInvalidKnownField, "input", "run transition requires --to <state> and --reason <text>", nil)
				break
			}
			value, err = service.TransitionRun(targetRef, toState, actor, reason, nextAction)
		}
	}
	return value, err
}
