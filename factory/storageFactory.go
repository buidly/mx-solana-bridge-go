package factory

import (
	"path"

	"github.com/multiversx/mx-chain-go/config"
	"github.com/multiversx/mx-chain-go/storage/factory"
	"github.com/multiversx/mx-chain-go/storage/storageunit"
	"github.com/multiversx/mx-solana-bridge-go/core"
)

// CreateUnitStorer based on the config and the working directory
func CreateUnitStorer(config config.StorageConfig, workingDir string) (core.Storer, error) {
	statusMetricsDbConfig := factory.GetDBFromConfig(config.DB)
	dbPath := path.Join(workingDir, config.DB.FilePath)
	statusMetricsDbConfig.FilePath = dbPath

	return storageunit.NewStorageUnitFromConf(
		factory.GetCacherFromConfig(config.Cache),
		statusMetricsDbConfig)
}
