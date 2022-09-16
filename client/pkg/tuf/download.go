package tuf

import (
	"io"
	"os"
	"path/filepath"

	tufClient "github.com/theupdateframework/go-tuf/client"
	tufUtil "github.com/theupdateframework/go-tuf/util"
)

func (c Client) DownloadFile(targetName, dest string, destMode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	f, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, destMode)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	file := destinationFile{f}
	if err := c.Download(targetName, &file); err != nil {
		return err
	}

	return nil
}

type destinationFile struct {
	*os.File
}

func (t *destinationFile) Delete() error {
	_ = t.Close()
	return os.Remove(t.Name())
}

func (c Client) Download(targetName string, destination tufClient.Destination) error {
	return c.Client.Download(tufUtil.NormalizeTarget(targetName), destination)
}

func (c Client) DownloadMeta(name string) ([]byte, error) {
	ioReader, _, err := c.RemoteStore.GetMeta(name)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(ioReader)
}
