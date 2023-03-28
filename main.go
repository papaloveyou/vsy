package main

import (
	"context"
	"github.com/jacobsa/fuse"
	"github.com/papaloveyou/vsy/internal"
	"github.com/papaloveyou/vsy/tools"
	"log"
)

const (
	DEFAULT_DIR = "/var/lib/vsy"
)

func init() {
	err := tools.MkDir(DEFAULT_DIR)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	fs := internal.NewVsy()
	mfs, err := fuse.Mount(DEFAULT_DIR, fs, &fuse.MountConfig{
		FSName:                  tools.GenFsname(),
		DisableWritebackCaching: true,
		EnableNoOpenSupport:     true,
		EnableNoOpendirSupport:  true,
		Subtype:                 "vsy",
	})
	if err != nil {
		log.Fatalf("Mount: %v", err)
	}
	log.Println("Successfully Mounted", DEFAULT_DIR)
	// Wait for it to be unmounted.
	if err = mfs.Join(context.Background()); err != nil {
		log.Fatalf("Join: %v", err)
	}
	log.Println("Successfully exiting", DEFAULT_DIR)
}
