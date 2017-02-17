package pack

import (
	"os"
	"archive/zip"
	"io"
	"runtime"
	"path/filepath"
	"strings"
	"github.com/klauspost/compress/gzip"
	"archive/tar"
)

func createZip(filename string, files []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	zipper := zip.NewWriter(file)
	defer zipper.Close()
	for _, name := files {
		if err := writeFileToZip(zipper, name); err != nil {
			return err
		}
	}
	return nil
}

func writeFileToZip(zipper *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = sanitizedName(filename)
	writer, err := zipper.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}

func sanitizedName(filename string) string {
	if len(filename) > 1 && filename[1] == ':' &&
		runtime.GOOS == "windows"{
		filename = filename[2:]
	}
	filename = filepath.ToSlash(filename)
	filename = strings.TrimLeft(filename, "./")
	return strings.Replace(filename, "../", "", -1)
}

func createTar(filename string, files []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	var fileWriter io.WriteCloser = file
	if strings.HasSuffix(filename, ".gz") {
		fileWriter = gzip.NewWriter(file)
		defer fileWriter.Close()
	}
	writer := tar.NewWriter(fileWriter)
	defer writer.Close()
	for _, name := range files {
		if err := writeFileToTar(writer, name); err != nil {
			return err
		}
	}
	return nil
}

func writeFileToTar(writer *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	header := &tar.Header{
		Name: sanitizedName(filename),
		Mode:int64(stat.Mode()),
		Uid:os.Getuid(),
		Gid:os.Getuid(),
		Size:stat.Size(),
		ModTime:stat.ModTime(),
	}
	if err = writer.WriteHeader(header); err != nil {
		return err
	}
	_, err = io.Copy(writer, file)
	return err
}