// file gotask/gocli/icli.go
package gocli

// - Extracts into a temporary directory.
// - Tries to locate the binary named as cliName.
// - If not found, searches for any executable file.
// - Moves the binary to the install directory and sets executable permission.

// case tar.TypeReg:
// 			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
// 			if err != nil {
// 				return err
// 			}
// 			if _, err := io.Copy(outFile, tr); err != nil {
// 				outFile.Close()
// 				return err
// 			}
// 			outFile.Close()
