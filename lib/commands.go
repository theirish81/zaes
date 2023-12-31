package lib

import (
	"crypto/rand"
	"fmt"
	"github.com/augani/zipper"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path"
	"path/filepath"
)

// Params are the possible input parameters
type Params struct {
	password       string
	nonInteractive bool
	eraseSource    bool
}

// params are the param instance
var params = Params{}

// EncryptCmd is the encryption command
var EncryptCmd = &cobra.Command{
	Use:   "encrypt [path]",
	Short: "Encrypts a directory into a file, and optionally securely erase the source",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		px := args[0]
		if !checkFileExistenceAndType(px, true) {
			PrintError("the provided directory does not exist or it is not a directory")
			os.Exit(1)
		}
		px = cleanupDirPath(px)
		if params.password == "" {
			PrintError("password is required")
			_ = cmd.Help()
			os.Exit(1)
		}
		parentDir := filepath.Dir(px)
		fileName := filepath.Base(px)
		zipFileName := zipFilePath(parentDir, fileName)
		zaesFileName := zaesFilePath(parentDir, fileName)
		if checkExists(zipFileName) {
			PrintError("a zip file with the given name already exists")
			os.Exit(1)
		}
		if checkExists(zaesFileName) {
			PrintError("a zaes file with the given name already exists")
			os.Exit(1)
		}

		if _, err := zipper.ZipIt(px, parentDir, fileName); err != nil {
			PrintError("error creating file: %s", err.Error())
			os.Exit(1)
		}

		params.password = extendPassword(params.password)

		gcm, err := NewGCM(params.password)
		if err != nil {
			PrintError("cipher GCM err: %s", err.Error())
			os.Exit(1)
		}
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			PrintError("nonce  err: %s", err.Error())
			os.Exit(1)
		}
		data, err := os.ReadFile(zipFileName)
		if err != nil {
			PrintError("could not read zip file: %s", err.Error())
		}
		cipherText := gcm.Seal(nonce, nonce, data, nil)
		err = os.WriteFile(zaesFileName, cipherText, 0777)
		if err != nil {
			PrintError("could not write zaes file: %s", err.Error())
		}
		err = WipeFile(zipFileName)
		if err != nil {
			PrintError("could not remove temp zip file: %s", err.Error())
			os.Exit(1)
		}
		if params.eraseSource && prompt(fmt.Sprintf("WARN: the directory %s will be SECURELY ERASED. Continue?", px), params.nonInteractive) {
			if err = WipeDir(px); err != nil {
				PrintError("could not remove directory: %s ", err.Error())
				os.Exit(1)
			}
		}
	},
}

// DecryptCmd is the decryption command
var DecryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypts a Zaes file into the original directory, and optionally securely erase the source",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		px := args[0]
		if !checkFileExistenceAndType(px, false) {
			PrintError("the provided directory does not exist or it is not a file")
			os.Exit(1)
		}
		if params.password == "" {
			PrintError("password is required")
			_ = cmd.Help()
			os.Exit(1)
		}
		params.password = extendPassword(params.password)
		parentDir := filepath.Dir(px)
		fileName := filepath.Base(px)
		extension := filepath.Ext(fileName)
		fileName = fileName[0 : len(fileName)-len(extension)]
		zipFileName := zipFilePath(parentDir, fileName)
		zaesFileName := zaesFilePath(parentDir, fileName)
		destinationDirName := path.Join(parentDir, fileName)
		if checkExists(zipFileName) {
			PrintError("a zip file with the given name already exists")
			os.Exit(1)
		}
		if checkExists(destinationDirName) {
			PrintError("a directory with the given name already exists")
			os.Exit(1)
		}
		gcm, err := NewGCM(params.password)

		if err != nil {
			PrintError("cipher GCM err: %s", err.Error())
			os.Exit(1)
		}

		cipherText, nonce, err := ReadCypherText(gcm, px)
		if err != nil {
			PrintError("cannot read CypherText: %s", err.Error())
			os.Exit(1)
		}
		plainText, err := gcm.Open(nil, nonce, cipherText, nil)
		if err != nil {
			PrintError("could not decrypt CypherText: %s", err.Error())
			os.Exit(1)
		}
		err = os.WriteFile(zipFileName, plainText, 0777)
		if err != nil {
			PrintError("could not write zip file: %s", err.Error())
			os.Exit(1)
		}
		_, err = zipper.UnZipIt(zipFileName, path.Join(parentDir, fileName))
		if err != nil {
			PrintError("could not unzip file: %s", err.Error())
			os.Exit(1)
		}
		err = WipeFile(zipFileName)
		if err != nil {
			PrintError("could not remove zip file: %s", err.Error())
			os.Exit(1)
		}
		if params.eraseSource && prompt(fmt.Sprintf("WARN: the file %s will be SECURELY ERASED. Continue?", zaesFileName), params.nonInteractive) {
			err = WipeFile(zaesFileName)
			if err != nil {
				PrintError("could not remove zaes file: %s", err.Error())
				os.Exit(1)
			}
		}
	},
}

// WipeCmd the wipe command
var WipeCmd = &cobra.Command{
	Use:   "wipe [path]",
	Short: "Securely erases a file or a directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !checkExists(args[0]) {
			PrintError("the provided path does not exist: %s", args[0])
			os.Exit(1)
		}
		if !prompt(fmt.Sprintf("WARN: the file %s will be SECURELY ERASED. Continue?", args[0]), params.nonInteractive) {
			os.Exit(0)
		}
		info, err := os.Stat(args[0])
		if err != nil {
			PrintError("could not access the provided path: %s", err.Error())
			os.Exit(1)
		}
		if info.IsDir() {
			if err := WipeDir(args[0]); err != nil {
				PrintError("error while wiping directory: %s", err.Error())
				os.Exit(1)
			}
		} else {
			if err := WipeFile(args[0]); err != nil {
				PrintError("error wiping file: %s", err.Error())
			}
		}
	},
}

// init is the initialization function
func init() {
	EncryptCmd.PersistentFlags().StringVarP(&params.password, "password", "p", "", "The password to encrypt the archive")
	EncryptCmd.PersistentFlags().BoolVarP(&params.nonInteractive, "non-interactive", "y", false, "If activated, no interactive warning will be issued")
	EncryptCmd.PersistentFlags().BoolVarP(&params.eraseSource, "erase-source", "e", false, "If activated, the source file(s) gets securely erased after completion")
	_ = EncryptCmd.MarkPersistentFlagRequired("password")

	DecryptCmd.PersistentFlags().StringVarP(&params.password, "password", "p", "", "The password to decrypt the archive")
	DecryptCmd.PersistentFlags().BoolVarP(&params.nonInteractive, "non-interactive", "y", false, "If activated, no interactive warning will be issued")
	DecryptCmd.PersistentFlags().BoolVarP(&params.eraseSource, "erase-source", "e", false, "If activated, the source file(s) gets securely erased after completion")
	_ = DecryptCmd.MarkPersistentFlagRequired("password")

	WipeCmd.PersistentFlags().BoolVarP(&params.nonInteractive, "non-interactive", "y", false, "If activated, no interactive warning will be issued")

}
