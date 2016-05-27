package cleanUp

import (
	"os"

	"github.com/cpg1111/maestro/config"
)

// Run runs the clean tasks
func Run(conf *config.CleanUp, clonePath *string) error {
	if len(conf.AdditionalCMDs) > 0 {
		// TODO run AdditionalCMDs
		return nil
	}
	if len(conf.Artifacts) > 0 {
		// TODO save artifacts
		return nil
	}
	if conf.InDaemon {
		// TODO send status to daemon
		return nil
	}
	os.RemoveAll(*clonePath)
	return nil
}
