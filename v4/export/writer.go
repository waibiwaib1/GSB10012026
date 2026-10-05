package export

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/pingcap/dumpling/v4/log"
	"go.uber.org/zap"
)

type Writer interface {
	WriteDatabaseMeta(ctx context.Context, db, createSQL string) error
	WriteTableMeta(ctx context.Context, db, table, createSQL string) error
	WriteTableData(ctx context.Context, ir TableDataIR) error
}

type SimpleWriter struct {
	cfg *Config
}

func NewSimpleWriter(config *Config) (Writer, error) {
	sw := &SimpleWriter{cfg: config}
	return sw, os.MkdirAll(config.OutputDirPath, 0755)
}

func (f *SimpleWriter) WriteDatabaseMeta(ctx context.Context, db, createSQL string) error {
	fileName := fmt.Sprintf("%s-schema-create.sql", db)
	filePath := path.Join(f.cfg.OutputDirPath, fileName)
	return writeMetaToFile(db, createSQL, filePath)
}

func (f *SimpleWriter) WriteTableMeta(ctx context.Context, db, table, createSQL string) error {
	fileName := fmt.Sprintf("%s.%s-schema.sql", db, table)
	filePath := path.Join(f.cfg.OutputDirPath, fileName)
	return writeMetaToFile(db, createSQL, filePath)
}

func (f *SimpleWriter) WriteTableData(ctx context.Context, ir TableDataIR) error {
	log.Zap().Debug("start dumping table...", zap.String("table", ir.TableName()))

	chunks := splitTableDataIntoChunks(ir, f.cfg.FileSize, f.cfg.StatementSize)
	chunkCount := 0
	fileName := fmt.Sprintf("%s.%s.sql", ir.DatabaseName(), ir.TableName())
	if f.cfg.FileSize != UnspecifiedSize {
		fileName = fmt.Sprintf("%s.%s.%d.sql", ir.DatabaseName(), ir.TableName(), chunkCount)
	}

	for {
		filePath := path.Join(f.cfg.OutputDirPath, fileName)
		fileWriter, tearDown := buildLazyFileWriter(filePath)
		intWriter := &InterceptStringWriter{StringWriter: fileWriter}
		err := WriteInsert(chunks, intWriter)
		tearDown()
		if err != nil {
			return err
		}

		if !intWriter.SomethingIsWritten {
			if chunkCount == 0 && f.cfg.Where != "" {
				if err := writeFilteredOutTableData(f.cfg.Where, filePath); err != nil {
					return err
				}
			}
			break
		}

		if f.cfg.FileSize == UnspecifiedSize {
			break
		}
		chunkCount += 1
		fileName = fmt.Sprintf("%s.%s.%d.sql", ir.DatabaseName(), ir.TableName(), chunkCount)
	}
	log.Zap().Debug("dumping table successfully",
		zap.String("table", ir.TableName()))
	return nil
}

func writeFilteredOutTableData(where, filePath string) error {
	fileWriter, tearDown, err := buildFileWriter(filePath)
	if err != nil {
		return err
	}
	defer tearDown()

	filteredOutComments := fmt.Sprintf("/* All data are filtered out by this condition: */\n/* WHERE (%s) */\n", where)
	return write(fileWriter, filteredOutComments)
}

func writeMetaToFile(target, metaSQL, path string) error {
	fileWriter, tearDown, err := buildFileWriter(path)
	if err != nil {
		return err
	}
	defer tearDown()

	return WriteMeta(&metaData{
		target:  target,
		metaSQL: metaSQL,
	}, fileWriter)
}
