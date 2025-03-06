package helper

import (
	"github.com/containerd/containerd/snapshots"
	"github.com/lingdie/rsnapshotter/pkg/consts"
)

func GetStorageFromSnapshotOpt(opts ...snapshots.Opt) string {
	tempSnapshot := &snapshots.Info{}
	for _, opt := range opts {
		opt(tempSnapshot)
	}
	return tempSnapshot.Labels[consts.LabelDevboxContainerSnapshot]
}
