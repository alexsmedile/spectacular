package command

import (
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/missionbundle"
)

type operationInput struct {
	rest         []string
	dataPayload  string
	allowMain    bool
	createBranch bool
	statusFilter string
	allMissions  bool
	dryRun       bool
	fromFile     string
	override     string
}

func (r Runner) execute(op operation, input operationInput, service missionbundle.Service) (any, error) {
	switch op {
	case opMissionStart, opMissionList, opMissionShow, opMissionCheck, opMissionAmendScope, opMissionClose, opObjectiveShow, opObjectivePromote, opObjectiveFinish, opRunShow, opRunStart, opMissionComplete, opRunTransition:
		return r.executeMission(op, input, service)
	case opReviewRecord, opHandoffRecord, opEvidenceRecord, opProposalCheck, opContractAmend, opContractCreate, opDecide:
		return r.executeRecords(op, input, service)
	case opCampaignCheck, opCharter, opGuard:
		return r.executeContext(op, input, service)
	default:
		return nil, domain.NewRefusal(domain.RefusalInvalidKnownField, "command", "operation has no workspace handler", nil)
	}
}
