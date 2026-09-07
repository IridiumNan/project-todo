package store

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/IridiumNan/project-todo/internal/models"
)

type StoreToml struct{}

const defaultTomlSep = "@@@@@GAP@@@@@\n"

// LoadMetadata load the metadata and return the slices
// The metadata file is a toml file with sep = @@@@@@GAP@@@@@@
// This function will not close the dataSrc
func (st *StoreToml) LoadMetadata(dataSrc io.Reader) (works []*models.Work, err error) {
	byteData, err := io.ReadAll(dataSrc)
	if err != nil {
		return nil, fmt.Errorf("error when read byte data from dataSrc, err: %s", err)
	}

	workParts := bytes.Split(byteData, []byte(defaultTomlSep))

	for idx := range workParts {
		var work models.Work
		err = toml.Unmarshal(workParts[idx], &work)
		if err != nil {
			slog.Error("while unmarshal toml metadata", "err", err, "str", string(workParts[idx]))
		}

		works = append(works, &work)
	}

	return
}

// AppendMetadata function append data into dataDst directly, not backup original data
func (st *StoreToml) AppendMetadata(dataDst io.Writer, newWorks []*models.Work) (err error) {
	byteData := st.buildTomlByteData(newWorks)

	_, err = dataDst.Write(byteData)
	if err != nil {
		return fmt.Errorf("error when write data into the dataDst, err: %s", err.Error())
	}

	return
}

// AppendMetadataToFile append new works into data file
// The dstPath should be a toml file path
func (st *StoreToml) AppendMetadataToFile(dstPath string, newWorks []*models.Work) (err error) {
	dstFile, err := os.OpenFile(dstPath, os.O_RDONLY, 0o644)
	if err != nil {
		return fmt.Errorf("error when open dst file path, path: %s, err: %s", dstPath, err)
	}

	tmpFile, err := os.CreateTemp("/tmp", "project-todo-data-*")
	if err != nil {
		return fmt.Errorf("error when create tmp file, err: %s", err)
	}
	_, err = io.Copy(tmpFile, dstFile)
	if err != nil {
		return fmt.Errorf("error when copy original data into tmp file, tmp file path: %s, dst file path: %s, err: %s", tmpFile.Name(), dstFile.Name(), err)
	}
	dstFile.Close()

	err = st.AppendMetadata(tmpFile, newWorks)

	tmpFile.Close()
	if err != nil {
		return fmt.Errorf("error when append new toml data into tmp file, file path: %s, err: %s", tmpFile.Name(), err)
	}

	err = os.Rename(tmpFile.Name(), dstFile.Name())
	if err != nil {
		return fmt.Errorf("error when replace old data file, tmp file path %s, new file path: %s, err: %s", tmpFile.Name(), dstFile.Name(), err)
	}

	return
}

// DumpMetadataToFile write to a tmp file first then remove it into old dataPath
// This function will replace old file with a updated file
func (st *StoreToml) DumpMetadataToFile(dstPath string, allWorks []*models.Work) (err error) {
	tmpFile, err := os.CreateTemp("/tmp", "project-todo-data-*")
	if err != nil {
		return fmt.Errorf("error when create tmp file, err: %s", err)
	}

	err = st.DumpMetadata(tmpFile, allWorks)
	if err != nil {
		return fmt.Errorf("error when write data into tmp file, err: %s", err)
	}

	tmpFile.Close()

	err = os.Rename(tmpFile.Name(), dstPath)
	if err != nil {
		return fmt.Errorf("error when replace old data file, tmp file path %s, new file path: %s, err: %s", tmpFile.Name(), dstPath, err)
	}

	return nil
}

// DumpMetadata dump build the dump the [models.Work] structure into [io.Writer] directly
func (st *StoreToml) DumpMetadata(dataDst io.Writer, allWorks []*models.Work) (err error) {
	byteData := st.buildTomlByteData(allWorks)

	_, err = dataDst.Write(byteData)
	if err != nil {
		return fmt.Errorf("error when write byte data, err: %s", err)
	}

	return
}

// buildTomlByteData build the []byte from []*[models.Work]
// Then you can write the type into [io.Writer]
func (st *StoreToml) buildTomlByteData(allWorks []*models.Work) (byteData []byte) {
	for idx := range allWorks {
		data, currErr := toml.Marshal(allWorks[idx])
		if currErr != nil {
			slog.Error("error when Marshal Work into byte data", "err", currErr, "models.Work", allWorks[idx])
			continue
		}
		byteData = append(byteData, data...)
		byteData = append(byteData, []byte(defaultTomlSep)...)
	}

	return
}
