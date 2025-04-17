package common

import "fmt"

type DuplicationStrategyType string

const (
	DuplicationStrategySkip      DuplicationStrategyType = "skip"
	DuplicationStrategyOverwrite DuplicationStrategyType = "overwrite"
	DuplicationStrategyAbort     DuplicationStrategyType = "abort"
)

func (d *DuplicationStrategyType) String() string {
	return string(*d)
}

func (d *DuplicationStrategyType) Set(value string) error {
	switch value {
	case string(DuplicationStrategySkip), string(DuplicationStrategyOverwrite), string(DuplicationStrategyAbort):
		*d = DuplicationStrategyType(value)
		return nil
	default:
		return fmt.Errorf("invalid value for DuplicationStrategyType: %s", value)
	}
}

func (d *DuplicationStrategyType) Type() string {
	return "DuplicationStrategyType"
}

// DuplicationOnNonMirrorStrategyType defines the strategy for handling non-mirror repo duplication
type DuplicationOnNonMirrorStrategyType string

const (
	DuplicationOnNonMirrorStrategySkip      DuplicationOnNonMirrorStrategyType = "skip"
	DuplicationOnNonMirrorStrategyOverwrite DuplicationOnNonMirrorStrategyType = "overwrite"
	DuplicationOnNonMirrorStrategyAbort     DuplicationOnNonMirrorStrategyType = "abort"
)

func (d *DuplicationOnNonMirrorStrategyType) String() string {
	return string(*d)
}

func (d *DuplicationOnNonMirrorStrategyType) Set(value string) error {
	switch value {
	case string(DuplicationOnNonMirrorStrategySkip), string(DuplicationOnNonMirrorStrategyOverwrite), string(DuplicationOnNonMirrorStrategyAbort):
		*d = DuplicationOnNonMirrorStrategyType(value)
		return nil
	default:
		return fmt.Errorf("invalid value for DuplicationOnNonMirrorStrategyType: %s", value)
	}
}

func (d *DuplicationOnNonMirrorStrategyType) Type() string {
	return "DuplicationOnNonMirrorStrategyType"
}

// DeletedRepoStrategyType defines the strategy for handling deleted repos
type DeletedRepoStrategyType string

const (
	DeletedRepoStrategySkip   DeletedRepoStrategyType = "skip"
	DeletedRepoStrategyDelete DeletedRepoStrategyType = "delete"
	DeletedRepoStrategyAbort  DeletedRepoStrategyType = "abort"
)

func (d *DeletedRepoStrategyType) String() string {
	return string(*d)
}

func (d *DeletedRepoStrategyType) Set(value string) error {
	switch value {
	case string(DeletedRepoStrategySkip), string(DeletedRepoStrategyDelete), string(DeletedRepoStrategyAbort):
		*d = DeletedRepoStrategyType(value)
		return nil
	default:
		return fmt.Errorf("invalid value for DeletedRepoStrategyType: %s", value)
	}
}

func (d *DeletedRepoStrategyType) Type() string {
	return "DeletedRepoStrategyType"
}
