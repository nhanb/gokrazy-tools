package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gokrazy/internal/config"
	"github.com/gokrazy/internal/instanceflag"
	"github.com/gokrazy/tools/internal/packer"
	"github.com/spf13/cobra"
)

// overwriteCmd is gok overwrite.
var overwriteCmd = &cobra.Command{
	Use:   "overwrite",
	Short: "Build and deploy a gokrazy instance to a storage device",
	Long: `Build and deploy a gokrazy instance to a storage device.

You typically need to use the gok overwrite command only once,
when first deploying your gokrazy instance. Afterwards, you can
switch to the gok update command instead for updating over the network.

Examples:
  # Overwrite the contents of the SD card sdx with gokrazy instance scan2drive:
  % gok -i scan2drive overwrite --full=/dev/sdx
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().NArg() > 0 {
			fmt.Fprint(os.Stderr, `positional arguments are not supported

`)
			return cmd.Usage()
		}

		return overwriteImpl.run(cmd.Context(), args, cmd.OutOrStdout(), cmd.OutOrStderr())
	},
}

type overwriteImplConfig struct {
	full string
	boot string
	root string
	mbr  string
}

var overwriteImpl overwriteImplConfig

func init() {
	instanceflag.RegisterPflags(overwriteCmd.Flags())
	overwriteCmd.Flags().StringVarP(&overwriteImpl.full, "full", "", "", "write a full gokrazy device image to the specified device (e.g. /dev/sdx) or path (e.g. /tmp/gokrazy.img)")
	overwriteCmd.Flags().StringVarP(&overwriteImpl.boot, "boot", "", "", "write the gokrazy boot file system to the specified partition (e.g. /dev/sdx1) or path (e.g. /tmp/boot.fat)")
	overwriteCmd.Flags().StringVarP(&overwriteImpl.root, "root", "", "", "write the gokrazy root file system to the specified partition (e.g. /dev/sdx2) or path (e.g. /tmp/root.squashfs)")
	overwriteCmd.Flags().StringVarP(&overwriteImpl.mbr, "mbr", "", "", "write the gokrazy master boot record (MBR) to the specified device (e.g. /dev/sdx) or path (e.g. /tmp/mbr.img). only effective if -boot is specified, too")

}

func (r *overwriteImplConfig) run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	// TODO: call config generate hook

	cfg, err := config.ReadFromFile()
	if err != nil {
		return err
	}
	log.Printf("cfg: %+v", cfg)

	// gok overwrite is mutually exclusive with gok update
	cfg.InternalCompatibilityFlags.Update = ""

	cfg.InternalCompatibilityFlags.Overwrite = r.full
	cfg.InternalCompatibilityFlags.OverwriteBoot = r.boot
	cfg.InternalCompatibilityFlags.OverwriteRoot = r.root
	cfg.InternalCompatibilityFlags.OverwriteMBR = r.mbr

	if err := os.Chdir(config.InstancePath()); err != nil {
		return err
	}

	packer.Main(cfg)

	return nil
}
