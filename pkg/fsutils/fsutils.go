package fsutils

import (
	"io"
	"os"
)

func CopyFile(source string, dest string) error {
	s, err := os.Open(source)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer d.Close()

	if _, err = io.Copy(d, s); err != nil {
		return err
	}

	err = d.Sync()
	if err != nil {
		return err
	}

	return nil
}
