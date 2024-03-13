package utils

import (
	"encoding/json"
	"golang.org/x/xerrors"
	"io/ioutil"
	"path/filepath"
)

const MetaFile = "sectorstore.json"

type SectorstoreConfig struct {
	ID string

	// A high weight means data is more likely to be stored in this path
	Weight uint64 // 0 = readonly

	// Intermediate data for the sealing process will be stored here
	CanSeal bool

	// Finalized sectors that will be proved over time will be stored here
	CanStore bool

	// MaxStorage specifies the maximum number of bytes to use for sector storage
	// (0 = unlimited)
	MaxStorage uint64

	// List of storage groups this path belongs to
	Groups []string

	// List of storage groups to which data from this path can be moved. If none
	// are specified, allow to all
	AllowTo []string

	// AllowTypes lists sector file types which are allowed to be put into this
	// path. If empty, all file types are allowed.
	//
	// Valid values:
	// - "unsealed"
	// - "sealed"
	// - "cache"
	// - "update"
	// - "update-cache"
	// Any other value will generate a warning and be ignored.
	AllowTypes []string

	// DenyTypes lists sector file types which aren't allowed to be put into this
	// path.
	//
	// Valid values:
	// - "unsealed"
	// - "sealed"
	// - "cache"
	// - "update"
	// - "update-cache"
	// Any other value will generate a warning and be ignored.
	DenyTypes []string
}

func (cfg *SectorstoreConfig) WriteFile(path string) error {
	if !(cfg.CanStore || cfg.CanSeal) {
		return xerrors.Errorf("must specify at least one of --store or --seal")
	}

	b, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return xerrors.Errorf("marshaling storage config: %w", err)
	}

	if err := ioutil.WriteFile(filepath.Join(path, MetaFile), b, 0644); err != nil {
		return xerrors.Errorf("persisting storage metadata (%s): %w", filepath.Join(path, MetaFile), err)
	}
	return nil
}
