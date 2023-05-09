package flags

import (
	"errors"

	"github.com/updiver/dumper"
)

var (
	ErrInvalidCloneMode = errors.New("invalid clone mode")
)

type CloneMode string

const (
	CloneModeAllBranches   = "all-branches"
	CloneModeDefaultBranch = "default-branch"
	CloneModeInvalid       = "invalid"
)

func (cm CloneMode) String() string {
	return string(cm)
}

func (cm CloneMode) Valid() error {
	switch cm {
	case CloneModeAllBranches:
	case CloneModeDefaultBranch:
		return nil
	default:
		return ErrInvalidCloneMode
	}

	return nil
}

// Cmd related

func ApplyCloneMode(config *dumper.DumpRepositoryOptions, cloneMode CloneMode) {
	switch cloneMode {
	case CloneModeAllBranches:
		config.OnlyDefaultBranch = dumper.NegativeBoolRef()
		config.BranchRestrictions = &dumper.BranchRestrictions{
			SingleBranch: false,
			BranchName:   "",
		}
	case CloneModeDefaultBranch:
		config.OnlyDefaultBranch = dumper.PositiveBoolRef()
		config.BranchRestrictions = nil
	default:
		return
	}
}
