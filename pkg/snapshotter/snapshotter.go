package snapshotter

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/snapshots"
	"github.com/containerd/containerd/snapshots/overlay"
	meta "github.com/containerd/containerd/snapshots/storage"
	"github.com/lingdie/rsnapshotter/pkg/helper"
	"github.com/lingdie/rsnapshotter/pkg/storage"
)

type Snapshotter struct {
	sn snapshots.Snapshotter
	ms *meta.MetaStore
	sg *storage.Storage
}

func NewSnapshotter(root string, opts ...overlay.Opt) (snapshots.Snapshotter, error) {
	sn, err := overlay.NewSnapshotter(root, opts...)
	if err != nil {
		return nil, err
	}
	ms, err := meta.NewMetaStore(filepath.Join(root, "rsnapshotter-metadata.db"))
	if err != nil {
		return nil, err
	}
	sg := storage.NewStorage()

	return &Snapshotter{
		sn: sn,
		ms: ms,
		sg: sg,
	}, nil
}

func (s *Snapshotter) Stat(ctx context.Context, key string) (snapshots.Info, error) {
	return s.sn.Stat(ctx, key)
}

func (s *Snapshotter) Update(ctx context.Context, info snapshots.Info, fieldpaths ...string) (snapshots.Info, error) {
	return s.sn.Update(ctx, info, fieldpaths...)
}

func (s *Snapshotter) Usage(ctx context.Context, key string) (snapshots.Usage, error) {
	return s.sn.Usage(ctx, key)
}

func (s *Snapshotter) Walk(ctx context.Context, fn snapshots.WalkFunc, filters ...string) error {
	return s.sn.Walk(ctx, fn, filters...)
}

func (s *Snapshotter) Prepare(ctx context.Context, key, parent string, opts ...snapshots.Opt) ([]mount.Mount, error) {
	mounts, err := s.sn.Prepare(ctx, key, parent, opts...)
	if err != nil {
		return nil, err
	}
	// get storage key from snapshot opt
	storageKey := helper.GetStorageFromSnapshotOpt(opts...)
	if storageKey != "" {
		// do sync snapshot data from storage to local, and record the snapshot key in meta store
		s.syncSnapshot(ctx, key, storageKey)
	}
	return mounts, nil
}

func (s *Snapshotter) View(ctx context.Context, key, parent string, opts ...snapshots.Opt) ([]mount.Mount, error) {
	return s.sn.View(ctx, key, parent, opts...)
}

func (s *Snapshotter) Commit(ctx context.Context, name, key string, opts ...snapshots.Opt) error {
	return s.sn.Commit(ctx, name, key, opts...)
}

func (s *Snapshotter) Remove(ctx context.Context, key string) error {
	return s.sn.Remove(ctx, key)
}

func (s *Snapshotter) Close() error {
	return s.sn.Close()
}

func (s *Snapshotter) Mounts(ctx context.Context, key string) ([]mount.Mount, error) {
	return s.sn.Mounts(ctx, key)
}

func (s *Snapshotter) Cleanup(ctx context.Context) error {
	// !! important: need to save snapshot data from storage before delete,
	// !! and we need make sure this process is atomic and thread safe
	// todo: get snapshot from ctx with meta store
	// todo: get snapshot info from meta store
	// todo: get snapshot data from storage
	// todo: save snapshot data to storage
	// todo: delete snapshot data from storage, use s.sn.Cleanup()
	// !! check if the s.sn.Cleanup() will not delete the snapshot data before we save it to storage
	// todo: delete snapshot info from meta store
	// todo: delete snapshot from overlay
	return nil
}

// syncSnapshot sync snapshot data from storage to local, and record the snapshot key in meta store
// todo: need to check if the snapshot data is already in local, maybe need to delete the local snapshot data
// todo: implement the syncSnapshot function
func (s *Snapshotter) syncSnapshot(ctx context.Context, snapshotKey string, storageKey string) error {
	// todo: implement: save snapshot info in meta store
	ctx, tx, err := s.ms.TransactionContext(ctx, false)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// todo: maybe we can save snapshot upperdir path in meta store
	// get snapshot info
	mounts, err := s.sn.Mounts(ctx, snapshotKey)
	if err != nil {
		return err
	}
	upperDir := mounts[0].Source
	// get snapshot data from storage and write to upperdir
	ioReader, err := s.sg.Read(ctx, storageKey)
	if err != nil {
		return err
	}
	defer ioReader.Close()
	// write to upperdir
	ioWriter, err := os.Create(upperDir)
	if err != nil {
		return err
	}
	defer ioWriter.Close()
	_, err = io.Copy(ioWriter, ioReader)
	if err != nil {
		return err
	}
	return nil
}
