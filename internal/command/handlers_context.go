package command

import (
	"strings"

	"github.com/alexsmedile/spectacular/v2/internal/campaign"
	"github.com/alexsmedile/spectacular/v2/internal/charter"
	"github.com/alexsmedile/spectacular/v2/internal/domain"
	"github.com/alexsmedile/spectacular/v2/internal/guard"
	"github.com/alexsmedile/spectacular/v2/internal/missionbundle"
)

func (r Runner) executeContext(op operation, input operationInput, service missionbundle.Service) (any, error) {
	rest := input.rest
	ws := service.Workspace
	var value any
	var err error
	switch op {
	case opCampaignCheck:
		targetPath := ""
		asciiMode := false
		for _, arg := range rest {
			if arg == "--ascii" {
				asciiMode = true
			} else if targetPath == "" {
				targetPath = arg
			}
		}
		if targetPath == "" {
			err = domain.NewRefusal(domain.RefusalInvalidWorkspacePath, "", "expected <campaign-path>", nil)
			break
		}
		var c campaign.Check
		c, err = campaign.Validate(ws, targetPath)
		if err == nil {
			c.ASCIIMode = asciiMode
			value = c
		}
	case opCharter:
		targetRef := ""
		promptMode := false
		var extraSources []string
		for _, arg := range rest {
			if arg == "--prompt" {
				promptMode = true
			} else if targetRef == "" {
				targetRef = arg
			} else {
				extraSources = append(extraSources, arg)
			}
		}
		if targetRef == "" {
			err = domain.NewRefusal(domain.RefusalInvalidReference, "", "expected <mission-ref>/<objective-ref> (e.g. M17/O1)", nil)
			break
		}
		parts := strings.Split(targetRef, "/")
		if len(parts) != 2 {
			err = domain.NewRefusal(domain.RefusalInvalidReference, targetRef, "expected <mission-ref>/<objective-ref> (e.g. M17/O1)", nil)
			break
		}
		var c *charter.Charter
		c, err = charter.Compile(ws, parts[0], parts[1], extraSources, missionbundle.ReadView)
		if err == nil && c != nil {
			c.PromptMode = promptMode
			value = c
		}
	case opGuard:
		targetRef := ""
		watchMode := false
		execCmd := ""
		var cmdArgs []string
		inCmd := false
		for i := 0; i < len(rest); i++ {
			arg := rest[i]
			if inCmd {
				cmdArgs = append(cmdArgs, arg)
			} else if arg == "--" {
				inCmd = true
			} else if arg == "--watch" {
				watchMode = true
			} else if arg == "--exec" && i+1 < len(rest) {
				execCmd = rest[i+1]
				i++
			} else if targetRef == "" {
				targetRef = arg
			}
		}
		if execCmd != "" && len(cmdArgs) > 0 {
			err = domain.NewRefusal(domain.RefusalInvalidReference, targetRef, "cannot combine --exec with command arguments after --", nil)
			break
		}
		if targetRef == "" || (len(cmdArgs) == 0 && execCmd == "") {
			err = domain.NewRefusal(domain.RefusalInvalidReference, targetRef, "expected spectacular guard <mission-ref>/<objective-ref> [--watch] [--exec <command>] -- <command...>", nil)
			break
		}
		value, err = guard.Run(ws, targetRef, watchMode, execCmd, cmdArgs, missionbundle.ReadView)
	}
	return value, err
}
