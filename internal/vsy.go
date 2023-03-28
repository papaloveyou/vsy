package internal

import (
	"context"
	"fmt"
	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"io"
	"log"
	"runtime/debug"
	"strings"
)

type Vsy struct {
	fuseutil.NotImplementedFileSystem
	inodes map[fuseops.InodeID]*Inode
}

func NewVsy() fuse.Server {
	fs := &Vsy{
		inodes: make(map[fuseops.InodeID]*Inode),
	}
	root := initRootInode()
	fs.inodes[root.id] = root
	for _, dirent := range root.children {
		fs.inodes[dirent.Inode] = toInode(dirent)
	}
	return fuseutil.NewFileSystemServer(fs)
}

func (fs *Vsy) getInodeOrDie(id fuseops.InodeID) *Inode {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("stacktrace from panic: %v \n"+string(debug.Stack()), err)
			err = fuse.EIO
		}
	}()
	inode := fs.inodes[id]
	if inode == nil {
		panic(fmt.Sprintf("Unknown inode: %v", id))
	}
	return inode
}

func (fs *Vsy) StatFS(ctx context.Context, op *fuseops.StatFSOp) (err error) {
	const BLOCK_SIZE = 4096
	const TOTAL_SPACE = 1 * 1024 * 1024 * 1024 * 1024 // 1TiB
	const TOTAL_BLOCKS = TOTAL_SPACE / BLOCK_SIZE
	const INODES = 100_000_000 // 100 million

	op.BlockSize = BLOCK_SIZE
	op.Blocks = TOTAL_BLOCKS
	op.BlocksFree = TOTAL_BLOCKS
	op.BlocksAvailable = TOTAL_BLOCKS
	op.IoSize = 1 * 1024 * 1024 // 1MB
	op.Inodes = INODES
	op.InodesFree = INODES
	return
}

func (fs *Vsy) GetInodeAttributes(ctx context.Context, op *fuseops.GetInodeAttributesOp) error {
	inode := fs.getInodeOrDie(op.Inode)

	op.Attributes = inode.attributes
	return nil
}

func (fs *Vsy) LookUpInode(ctx context.Context, op *fuseops.LookUpInodeOp) error {
	parent := fs.getInodeOrDie(op.Parent)
	child, err := parent.findChildInode(op.Name)
	if err != nil {
		return err
	}

	// Copy over information.
	op.Entry.Child = child
	op.Entry.Attributes = fs.inodes[child].attributes
	return nil
}

func (fs *Vsy) ReadDir(ctx context.Context, op *fuseops.ReadDirOp) error {
	inode := fs.getInodeOrDie(op.Inode)
	if !inode.dir {
		return fuse.EIO
	}
	entries := inode.children
	// Grab the range of interest.
	if op.Offset > fuseops.DirOffset(len(entries)) {
		return fuse.EIO
	}
	entries = entries[op.Offset:]
	// Resume at the specified offset into the array.
	for _, e := range entries {
		n := fuseutil.WriteDirent(op.Dst[op.BytesRead:], e)
		if n == 0 {
			break
		}
		op.BytesRead += n
	}
	return nil
}

func (fs *Vsy) OpenFile(ctx context.Context, op *fuseops.OpenFileOp) error {
	// Allow opening any file.
	op.KeepPageCache = false
	op.UseDirectIO = true
	return nil
}

func (fs *Vsy) ReadFile(ctx context.Context, op *fuseops.ReadFileOp) error {
	var info string
	switch op.Inode {
	case meminfo:
		info = GetMeminfo()
	case cpuinfo:
		info = GetCpuinfo()
	default:
		log.Println("not supported")
		return fuse.ENOENT
	}
	if info == "" {
		return fuse.EIO
	}
	reader := strings.NewReader(info)
	var err error
	op.BytesRead, err = reader.ReadAt(op.Dst, op.Offset)
	// Special case: FUSE doesn't expect us to return io.EOF.
	if err == io.EOF {
		return nil
	}
	return nil
}
