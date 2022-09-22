package drivers

import (
	"fmt"
	"os"
	"os/exec"
)

const CredentialFilePrefix = "mysqldumpcred-"

type Mysql struct {
	MysqlDumpBinaryPath      string
	SkipComments             bool
	UseExtendedInserts       bool
	UseSingleTransaction     bool
	SkipLockTables           bool
	DoNotUseColumnStatistics bool
	UseQuick                 bool
	DefaultCharacterSet      string
	SetGtidPurged            string
	CreateTables             bool
	*DBDumper
}

func NewMysqlDumper() *Mysql {
	return &Mysql{
		MysqlDumpBinaryPath:      "mysqldump",
		SkipComments:             true,
		UseExtendedInserts:       true,
		UseSingleTransaction:     false,
		SkipLockTables:           false,
		DoNotUseColumnStatistics: false,
		UseQuick:                 false,
		DefaultCharacterSet:      "",
		SetGtidPurged:            "AUTO",
		CreateTables:             true,
		DBDumper:                 NewDBDumper(),
	}
}

func (mysql *Mysql) getDumpCommandArgs() ([]string, error) {
	credentialsFileName, err := mysql.createCredentialFile()
	if err != nil {
		return nil, err
	}

	args := []string{"--defaults-extra-file=" + credentialsFileName + ""}

	if !mysql.CreateTables {
		args = append(args, "--no-create-info")
	}

	if mysql.DefaultCharacterSet != "" {
		args = append(args, "--default-character-set="+mysql.DefaultCharacterSet)
	}

	if mysql.UseExtendedInserts {
		args = append(args, "--extended-insert")
	} else {
		args = append(args, "--skip-extended-insert")
	}

	if mysql.UseSingleTransaction {
		args = append(args, "--single-transaction")
	}

	if mysql.SkipComments {
		args = append(args, "--skip-comments")
	}

	if mysql.SkipLockTables {
		args = append(args, "--skip-lock-tables")
	}

	if mysql.DoNotUseColumnStatistics {
		args = append(args, "--column-statistics=0")
	}

	if mysql.UseQuick {
		args = append(args, "--quick")
	}

	if mysql.SetGtidPurged != "AUTO" {
		args = append(args, "--set-gtid-purged="+mysql.SetGtidPurged)
	}

	args = append(args, mysql.DBName)

	return args, nil
}

func (mysql *Mysql) createCredentialFile() (string, error) {
	var fileName string

	contents := `[client]
user = %s
password = %s
port = %d
host = %s`

	contents = fmt.Sprintf(contents, mysql.Username, mysql.Password, mysql.Port, mysql.Host)

	file, err := os.CreateTemp("", CredentialFilePrefix)
	if err != nil {
		return fileName, fmt.Errorf("failed to create temp folder: %w", err)
	}

	defer file.Close()

	_, err = file.WriteString(contents)
	if err != nil {
		return fileName, fmt.Errorf("failed to write credentials to temp file: %w", err)
	}

	return file.Name(), nil
}

func (mysql *Mysql) Dump(dumpFile string) error {
	args, err := mysql.getDumpCommandArgs()

	if err != nil {
		return fmt.Errorf("failed to get dump command args %w", err)
	}

	mysqldumpBinaryPath, err := exec.LookPath(mysql.MysqlDumpBinaryPath)
	if err != nil {
		return fmt.Errorf("failed to find mysqldump executable %s %w", mysql.MysqlDumpBinaryPath, err)
	}

	cmd := exec.Command(mysqldumpBinaryPath, args...)

	dumpOutFile, err := os.Create(dumpFile)
	if err != nil {
		return fmt.Errorf("failed to create the dump file %w", err)
	}
	defer dumpOutFile.Close()
	cmd.Stdout = dumpOutFile

	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		return fmt.Errorf("failed to run dump command %w", err)
	}

	return nil
}
