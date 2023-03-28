package internal

import (
	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"github.com/papaloveyou/vsy/tools"
	"os"
	"time"
)

const (
	root fuseops.InodeID = fuseops.RootInodeID + iota
	meminfo
	cpuinfo
)

type Inode struct {
	id         fuseops.InodeID
	name       *string
	attributes fuseops.InodeAttributes
	dir        bool
	children   []fuseutil.Dirent
}

func (inode *Inode) findChildInode(name string) (fuseops.InodeID, error) {
	l := len(inode.children)
	if l == 0 {
		return 0, fuse.ENOENT
	}
	for _, child := range inode.children {
		if child.Name == name {
			return child.Inode, nil
		}
	}
	return 0, fuse.ENOENT
}

func toInode(dirent fuseutil.Dirent) *Inode {
	now := time.Now()
	return &Inode{
		id:   dirent.Inode,
		name: tools.PString(dirent.Name),
		attributes: fuseops.InodeAttributes{
			Mode:  0444,
			Atime: now,
			Mtime: now,
			Ctime: now,
		},
	}
}

func initRootInode() *Inode {
	now := time.Now()
	return &Inode{
		id: root,
		attributes: fuseops.InodeAttributes{
			Size:  4096,
			Mode:  os.ModeDir,
			Atime: now,
			Mtime: now,
			Ctime: now,
		},
		dir: true,
		children: []fuseutil.Dirent{
			{
				Offset: 1,
				Inode:  meminfo,
				Name:   "meminfo",
				Type:   fuseutil.DT_File,
			},
			{
				Offset: 2,
				Inode:  cpuinfo,
				Name:   "cpuinfo",
				Type:   fuseutil.DT_File,
			},
		},
	}
}
