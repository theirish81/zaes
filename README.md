# Zaes
Zaes is a simple CLI utility that allows you to:
* ZIP and AES-encrypt files and entire directories
* Securely wipe files from your hard drive (single pass)

## Commands

### encrypt
```shell
zaes encrypt [path]
```
Zips and encrypts the provided path. The byproduct of the process is a file retaining the file name, with the `.zaes`
extension.

**Parameters**
* `-p --password` (string, required): the password to encrypt the archive
* `-e --erase-source` (boolean): securely wipes the source file once the encryption is complete
* `-y --non-interactive` (boolean): will automatically "yes" any prompt that may arise

### decrypt
```shell
zaes decrypt [path]
```
Decrypts and unzips the provided path to a `.zaes` file.

**Parameters**
* `-p --password` (string, required): the password to decrypt the archive
* `-e --erase-source` (boolean): securely wipes the source `.zaes` file once the decryption is complete
* `-y --non-interactive` (boolean): will automatically "yes" any prompt that may arise

### wipe
```shell
zaes wipe [path]
```
Securely wipes the provided path, whether it's a file or directory.

**Parameters**
* `-y --non-interactive` (boolean): will automatically "yes" any prompt that may arise