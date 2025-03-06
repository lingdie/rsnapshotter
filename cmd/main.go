//go:build linux

package main

import (
	"log"
	"net"
	"os"
	"path"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"

	snapshotsapi "github.com/containerd/containerd/api/services/snapshots/v1"
	"github.com/containerd/containerd/contrib/snapshotservice"
	"github.com/containerd/containerd/snapshots/overlay"

	"github.com/lingdie/rsnapshotter/pkg/consts"
	"github.com/lingdie/rsnapshotter/pkg/snapshotter"
)

func main() {
	app := &cli.App{
		Name:  "rsnapshotter",
		Usage: "Run a rsnapshotter",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "root-dir",
				Value: consts.DefaultRootDir,
			},
		},
		Action: func(ctx *cli.Context) error {
			rootDir := ctx.String("root-dir")
			sOpts := []overlay.Opt{}
			//!! must enable async remove and upperdir label
			sOpts = append(sOpts, overlay.AsynchronousRemove, overlay.WithUpperdirLabel)
			sn, err := snapshotter.NewSnapshotter(rootDir, sOpts...)
			if err != nil {
				return err
			}
			service := snapshotservice.FromSnapshotter(sn)
			rpc := grpc.NewServer()
			snapshotsapi.RegisterSnapshotsServer(rpc, service)
			socksPath := path.Join(rootDir, consts.SocksFileName)
			err = os.RemoveAll(socksPath)
			if err != nil {
				return err
			}
			l, err := net.Listen("unix", socksPath)
			if err != nil {
				return err
			}
			return rpc.Serve(l)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
