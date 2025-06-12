package svc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/autonomouskoi/mageutil"
)

var (
	baseDir string
)

func init() {
	var err error
	if baseDir, err = os.Getwd(); err != nil {
		panic(err)
	}
	baseDir = filepath.Join(baseDir, "..", "svc")
}

func TinyGoProtos() error {
	akcoreDir := filepath.Join(baseDir, "..", "..", "akcore")
	protosDir := filepath.Join(akcoreDir, "svc", "pb")
	protoFiles, err := mageutil.DirGlob(protosDir, "*.proto")
	if err != nil {
		return fmt.Errorf("globbing %s: %w", protosDir, err)
	}

	for _, protoFile := range protoFiles {
		baseName := strings.TrimSuffix(protoFile, ".proto")
		outFile := baseName + ".pb.go"
		err := mageutil.TinyGoProto(
			filepath.Join(baseDir, outFile),
			filepath.Join(protosDir, protoFile),
			protosDir,
		)
		if err != nil {
			return fmt.Errorf("%s -> %s: %w",
				filepath.Join(protosDir, protoFile),
				filepath.Join(baseDir, outFile),
				err,
			)
		}
	}
	return err
}
