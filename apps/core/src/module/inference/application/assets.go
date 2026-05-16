package application

import (
	"runtime"

	"github.com/authuser-dev/karacron/apps/core/data"
	"github.com/authuser-dev/karacron/apps/core/src/module/inference/domain"
)

func SelectCurrentAsset() (domain.BinaryAsset, bool) {
	for _, asset := range data.LocalBinaryMetadata.Assets {
		if asset.OS == runtime.GOOS && asset.Arch == runtime.GOARCH {
			return asset, true
		}
	}
	return domain.BinaryAsset{}, false
}
