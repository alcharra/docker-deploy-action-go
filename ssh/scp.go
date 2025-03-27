package ssh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"path"
)

func (c *Client) UploadFileSCP(localPath, remotePath string) error {
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	srcFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %w", err)
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat local file: %w", err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %w", err)
	}

	go func() {
		defer stdin.Close()
		fmt.Fprintf(stdin, "C%#o %d %s\n", info.Mode().Perm(), info.Size(), filepath.Base(localPath))
		io.Copy(stdin, srcFile)
		fmt.Fprint(stdin, "\x00")
	}()

	cmd := fmt.Sprintf("scp -t %s", path.Dir(remotePath))
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("failed to run scp command: %w", err)
	}

	return nil
}
