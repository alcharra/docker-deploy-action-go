package scp

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/alcharra/docker-deploy-action-go/internal/ssh/client"
)

func UploadFileSCP(cli *client.Client, localPath, remotePath string) error {
	if cli == nil {
		return fmt.Errorf("SSH client is not initialised")
	}

	dir := path.Dir(remotePath)

	session, err := cli.NewSession()
	if err != nil {
		return fmt.Errorf("unable to initialise SSH session for directory creation: %w", err)
	}
	if err := session.Run(fmt.Sprintf("mkdir -p %q", dir)); err != nil {
		session.Close()
		return fmt.Errorf("unable to create remote directory '%s': %w", dir, err)
	}
	session.Close()

	session, err = cli.NewSession()
	if err != nil {
		return fmt.Errorf("unable to initialise SSH session for file transfer: %w", err)
	}
	defer session.Close()

	srcFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("unable to open local file '%s': %w", localPath, err)
	}
	defer srcFile.Close()

	info, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("unable to retrieve file info for '%s': %w", localPath, err)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("unable to open stdin pipe for SCP transfer: %w", err)
	}

	errCh := make(chan error, 1)
	go func() {
		defer stdin.Close()

		if _, err := fmt.Fprintf(stdin, "C%#o %d %s\n",
			info.Mode().Perm(),
			info.Size(),
			filepath.Base(srcFile.Name()),
		); err != nil {
			errCh <- fmt.Errorf("unable to write file header for '%s': %w", srcFile.Name(), err)
			return
		}

		if _, err := io.Copy(stdin, srcFile); err != nil {
			errCh <- fmt.Errorf("unable to copy file content for '%s': %w", srcFile.Name(), err)
			return
		}

		if _, err := fmt.Fprint(stdin, "\x00"); err != nil {
			errCh <- fmt.Errorf("unable to send transfer completion signal for '%s': %w", srcFile.Name(), err)
			return
		}

		errCh <- nil
	}()

	cmd := fmt.Sprintf("scp -t %s", dir)
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("SCP transfer failed for '%s': %w", localPath, err)
	}

	if err := <-errCh; err != nil {
		return fmt.Errorf("failed to write SCP payload for '%s': %w", localPath, err)
	}

	return nil
}
