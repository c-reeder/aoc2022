package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type FileType int

const (
	FileTypeRegular FileType = iota
	FileTypeDirectory
)

type File interface {
	FileType() FileType
}

type RegularFile struct {
	Size uint64
}

func (f *RegularFile) FileType() FileType {
	return FileTypeRegular
}

type Directory struct {
	Children map[string]File
	Parent   *Directory
}

func (f *Directory) FileType() FileType {
	return FileTypeDirectory
}

type FileSystem struct {
	DiskSize   uint64
	Root       *Directory
	CurrentDir *Directory
}

func NewFileSystem(diskSize uint64) FileSystem {
	root := Directory{
		Children: map[string]File{},
	}
	root.Parent = &root
	return FileSystem{
		DiskSize:   diskSize,
		Root:       &root,
		CurrentDir: &root,
	}
}

func (fs *FileSystem) ChangeDirectory(dirName string) {
	if dirName == "/" {
		fs.CurrentDir = fs.Root
		return
	}
	if dirName == ".." {
		fs.CurrentDir = fs.CurrentDir.Parent
		return
	}
	v, ok := fs.CurrentDir.Children[dirName]
	if !ok && dirName != ".." {
		log.Fatal(ErrDirectoryDoesNotExist)
	}
	dir, ok := v.(*Directory)
	if !ok {
		log.Fatal(ErrCouldNotChangeIntoDir)
	}
	fs.CurrentDir = dir
}

func (fs *FileSystem) CreateRegularFile(name string, size uint64) {
	fs.AddDirectoryEntry(name, &RegularFile{
		Size: size,
	})
}
func (fs *FileSystem) CreateDirectory(name string) {
	fs.AddDirectoryEntry(name, &Directory{
		Children: map[string]File{},
		Parent:   fs.CurrentDir,
	})
}

func (fs *FileSystem) AddDirectoryEntry(name string, entry File) {
	if _, ok := fs.CurrentDir.Children[name]; ok {
		log.Fatal(ErrDuplicateFileName)
	}
	fs.CurrentDir.Children[name] = entry

}

var (
	CdRgx   = regexp.MustCompile(`^\$ +cd +([^ ]+)$`)
	LsRgx   = regexp.MustCompile(`^\$ +ls$`)
	DirRgx  = regexp.MustCompile(`^dir +([^ ]+)$`)
	FileRgx = regexp.MustCompile(`^([0-9]+) +([^ ]+)$`)

	ErrCouldNotChangeIntoDir = errors.New("could not change into directory")
	ErrDirectoryDoesNotExist = errors.New("directory does not exist")
	ErrDuplicateFileName     = errors.New("tried to create duplicate file")
	ErrUnknownFileType       = errors.New("unknown file type")
)

func main() {
	// Parse flags
	{
		partBFlag := flag.Bool("b", false, "To switch to part b")
		flag.Parse()
		if partBFlag != nil && *partBFlag {
		}

	}

	// Check args
	if len(flag.Args()) != 1 {
		log.Fatal("Expected 1 argument containing file name!")
	}

	// Open file
	file, err := os.Open(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fs := NewFileSystem(70000000)

	// Iterate over lines and build up file system
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if matches := CdRgx.FindStringSubmatch(line); len(matches) == 2 {
			fs.ChangeDirectory(matches[1])
		} else if matches := LsRgx.FindStringSubmatch(line); len(matches) == 1 {
			// No-op
		} else if matches := FileRgx.FindStringSubmatch(line); len(matches) == 3 {
			size, err := strconv.ParseUint(matches[1], 10, 64)
			if err != nil {
				log.Fatal("could not parse filesize: %v\n", matches[1])
			}
			fs.CreateRegularFile(matches[2], size)
		} else if matches := DirRgx.FindStringSubmatch(line); len(matches) == 2 {
			if err != nil {
				log.Fatal("could not parse filesize: %v\n", matches[1])
			}
			fs.CreateDirectory(matches[1])
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("------------------")

	// Part A
	const maxSize = 100000
	totalSumUnderMax, totalUsed := PartA(&fs, maxSize)
	fmt.Printf("Sum of directories under maximum %v is: %v\n", maxSize, totalSumUnderMax)

	// Part B
	dirName, dirSize := PartB(&fs, totalUsed)
	fmt.Printf("Directory to remove is %v with a size of %v\n", dirName, dirSize)

}

// PartA finds sum of the sizes of all directories under maxSize
// Also returns totalUsed space to jump-start PartB
func PartA(fs *FileSystem, maxSize uint64) (totalSumUnderMax uint64, totalUsed uint64) {
	totalUsed = AddSizeToSumIfUnderMax(fs.Root, maxSize, &totalSumUnderMax)
	return totalSumUnderMax, totalUsed
}

// PartB finds smallest directory to remove to get the amount of remaining space
func PartB(fs *FileSystem, totalUsed uint64) (dirName string, dirSize uint64) {
	unused := fs.DiskSize - totalUsed
	neededForUpdate := 30000000 - unused
	flatDirSizes := map[string]uint64{}
	AddToFlatMapIfOverMin("/", fs.Root, neededForUpdate, flatDirSizes)

	// Iterate over directories with size greater than min
	// to find the smallest one
	for currDirName, currSize := range flatDirSizes {
		if dirSize == 0 || currSize < dirSize {
			dirName = currDirName
			dirSize = currSize
		}
	}
	return dirName, dirSize
}

func AddSizeToSumIfUnderMax(dir *Directory, maxSize uint64, sum *uint64) uint64 {
	var localSum uint64
	for _, childFile := range dir.Children {
		switch v := childFile.(type) {
		case *Directory:
			dirSum := AddSizeToSumIfUnderMax(v, maxSize, sum)
			localSum += dirSum
		case *RegularFile:
			localSum += v.Size
		default:
			log.Fatal(ErrUnknownFileType)
		}
	}
	if localSum < maxSize {
		*sum += localSum
	}
	return localSum
}

func AddToFlatMapIfOverMin(dirName string, dir *Directory, minSize uint64, flatMap map[string]uint64) uint64 {
	var localSum uint64
	for childName, childFile := range dir.Children {
		switch v := childFile.(type) {
		case *Directory:
			dirSum := AddToFlatMapIfOverMin(childName, v, minSize, flatMap)
			localSum += dirSum
		case *RegularFile:
			localSum += v.Size
		default:
			log.Fatal(ErrUnknownFileType)
		}
	}
	if localSum > minSize {
		flatMap[dirName] = localSum
	}
	return localSum
}
