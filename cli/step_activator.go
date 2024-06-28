package cli

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/bitrise/log"
	"github.com/bitrise-io/bitrise/tools"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/command/git"
	"github.com/bitrise-io/go-utils/pointers"
	"github.com/bitrise-io/stepman/activator"
	stepmanModels "github.com/bitrise-io/stepman/models"
	"github.com/bitrise-io/stepman/stepid"
)

type stepActivator struct {
}

func newStepActivator() stepActivator {
	return stepActivator{}
}

func (a stepActivator) activateStep(
	stepIDData stepid.CanonicalID,
	isStepLibUpdated bool,
	stepDir string, // $TMPDIR/bitrise/step_src
	workDir string, // $TMPDIR/bitrise
	workflowStep *stepmanModels.StepModel,
	stepInfoPtr *stepmanModels.StepInfoModel,
	isSteplibOfflineMode bool,
) (stepYMLPth string, origStepYMLPth string, didStepLibUpdate bool, err error) {
	stepYMLPth = filepath.Join(workDir, "current_step.yml")
	stepmanLogger := log.NewLogger(log.GetGlobalLoggerOpts())

	if stepIDData.SteplibSource == "path" {
		log.Debugf("[BITRISE_CLI] - Local step found: (path:%s)", stepIDData.IDorURI)

		activatedStep, err := activator.ActivatePathRefStep(
			stepmanLogger,
			stepIDData,
			stepDir,
			workDir,
		)
		if err != nil {
			return "", "", false, fmt.Errorf("activate local step: %w", err)
		}

		stepYMLPth = activatedStep.StepYMLPath
		origStepYMLPth = activatedStep.OrigStepYMLPath
	} else if stepIDData.SteplibSource == "git" {
		log.Debugf("[BITRISE_CLI] - Remote step, with direct git uri: (uri:%s) (tag-or-branch:%s)", stepIDData.IDorURI, stepIDData.Version)
		repo, err := git.New(stepDir)
		if err != nil {
			return "", "", false, err
		}
		var cloneCmd *command.Model
		if stepIDData.Version == "" {
			cloneCmd = repo.Clone(stepIDData.IDorURI, "--depth=1")
		} else {
			cloneCmd = repo.CloneTagOrBranch(stepIDData.IDorURI, stepIDData.Version, "--depth=1")
		}
		if out, err := cloneCmd.RunAndReturnTrimmedCombinedOutput(); err != nil {
			if strings.HasPrefix(stepIDData.IDorURI, "git@") {
				log.Warnf(`Note: if the step's repository is an open source one,
you should probably use a "https://..." git clone URL,
instead of the "git@..." git clone URL which usually requires authentication
even if the repository is open source!`)
			}
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				return "", "", false, fmt.Errorf("command failed with exit status %d (%s): %w", exitErr.ExitCode(), cloneCmd.PrintableCommandArgs(), errors.New(out))
			}
			return "", "", false, err
		}

		if err := command.CopyFile(filepath.Join(stepDir, "step.yml"), stepYMLPth); err != nil {
			return "", "", false, err
		}
	} else if stepIDData.SteplibSource == "_" {
		log.Debugf("[BITRISE_CLI] - Steplib independent step, with direct git uri: (uri:%s) (tag-or-branch:%s)", stepIDData.IDorURI, stepIDData.Version)

		// StepLib independent steps are completely defined in the workflow
		stepYMLPth = ""
		if err := workflowStep.FillMissingDefaults(); err != nil {
			return "", "", false, err
		}

		repo, err := git.New(stepDir)
		if err != nil {
			return "", "", false, err
		}
		if err := repo.CloneTagOrBranch(stepIDData.IDorURI, stepIDData.Version).Run(); err != nil {
			return "", "", false, err
		}
	} else if stepIDData.SteplibSource != "" {
		stepInfo, didUpdate, err := activateStepLibStep(stepIDData, stepDir, stepYMLPth, isStepLibUpdated, isSteplibOfflineMode)
		didStepLibUpdate = didUpdate

		stepInfoPtr.ID = stepInfo.ID
		if stepInfoPtr.Step.Title == nil || *stepInfoPtr.Step.Title == "" {
			stepInfoPtr.Step.Title = pointers.NewStringPtr(stepInfo.ID)
		}
		stepInfoPtr.Version = stepInfo.Version
		stepInfoPtr.LatestVersion = stepInfo.LatestVersion
		stepInfoPtr.OriginalVersion = stepInfo.OriginalVersion
		stepInfoPtr.GroupInfo = stepInfo.GroupInfo

		if err != nil {
			return "", "", didStepLibUpdate, err
		}
	} else {
		return "", "", didStepLibUpdate, fmt.Errorf("invalid stepIDData: no SteplibSource or LocalPath defined (%v)", stepIDData)
	}

	return stepYMLPth, origStepYMLPth, didStepLibUpdate, nil
}

func activateStepLibStep(stepIDData stepid.CanonicalID, destination, stepYMLCopyPth string, isStepLibUpdated bool, isOfflineMode bool) (stepmanModels.StepInfoModel, bool, error) {
	didStepLibUpdate := false

	log.Debugf("[BITRISE_CLI] - Steplib (%s) step (id:%s) (version:%s) found, activating step", stepIDData.SteplibSource, stepIDData.IDorURI, stepIDData.Version)
	if err := tools.StepmanSetup(stepIDData.SteplibSource); err != nil {
		return stepmanModels.StepInfoModel{}, false, err
	}

	versionConstraint, err := stepmanModels.ParseRequiredVersion(stepIDData.Version)
	if err != nil {
		return stepmanModels.StepInfoModel{}, false,
			fmt.Errorf("activating step (%s) from source (%s) failed, invalid version specified: %s", stepIDData.IDorURI, stepIDData.SteplibSource, err)
	}
	if versionConstraint.VersionLockType == stepmanModels.InvalidVersionConstraint {
		return stepmanModels.StepInfoModel{}, false,
			fmt.Errorf("activating step (%s) from source (%s) failed, version constraint is invalid", stepIDData.IDorURI, stepIDData.SteplibSource)
	}

	isStepLibUpdateNeeded := (versionConstraint.VersionLockType == stepmanModels.Latest) ||
		(versionConstraint.VersionLockType == stepmanModels.MinorLocked) ||
		(versionConstraint.VersionLockType == stepmanModels.MajorLocked)
	if !isStepLibUpdated && isStepLibUpdateNeeded && !isOfflineMode {
		log.Print("Step uses latest version, updating StepLib...")
		if err := tools.StepmanUpdate(stepIDData.SteplibSource); err != nil {
			log.Warnf("Step version constraint is latest or version locked, but failed to update StepLib, err: %s", err)
		} else {
			didStepLibUpdate = true
		}
	}

	info, err := tools.StepmanStepInfo(stepIDData.SteplibSource, stepIDData.IDorURI, stepIDData.Version)
	if err != nil {
		if isStepLibUpdated || isOfflineMode {
			return stepmanModels.StepInfoModel{}, didStepLibUpdate, fmt.Errorf("stepman JSON steplib step info failed: %s", err)
		}

		// May StepLib should be updated
		log.Infof("Step info not found in StepLib (%s) -- Updating ...", stepIDData.SteplibSource)
		if err := tools.StepmanUpdate(stepIDData.SteplibSource); err != nil {
			return stepmanModels.StepInfoModel{}, didStepLibUpdate, err
		}

		didStepLibUpdate = true

		info, err = tools.StepmanStepInfo(stepIDData.SteplibSource, stepIDData.IDorURI, stepIDData.Version)
		if err != nil {
			return stepmanModels.StepInfoModel{}, didStepLibUpdate, fmt.Errorf("stepman JSON steplib step info failed: %s", err)
		}
	}

	if info.Step.Title == nil || *info.Step.Title == "" {
		info.Step.Title = pointers.NewStringPtr(info.ID)
	}
	info.OriginalVersion = stepIDData.Version

	if err := tools.StepmanActivate(stepIDData.SteplibSource, stepIDData.IDorURI, info.Version, destination, stepYMLCopyPth, isOfflineMode); err != nil {
		return stepmanModels.StepInfoModel{}, didStepLibUpdate, err
	}
	log.Debugf("[BITRISE_CLI] - Step activated: (ID:%s) (version:%s)", stepIDData.IDorURI, stepIDData.Version)

	return info, didStepLibUpdate, nil
}
