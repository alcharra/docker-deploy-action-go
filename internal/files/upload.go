package files

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alcharra/docker-deploy-action-go/config"
	"github.com/alcharra/docker-deploy-action-go/internal/logs"
	"github.com/alcharra/docker-deploy-action-go/internal/ssh/client"
	"github.com/alcharra/docker-deploy-action-go/internal/ssh/scp"
)

func UploadFiles(cli *client.Client, cfg config.DeployConfig) []UploadedFile {
	logs.IsVerbose = cfg.Verbose

	var uploaded []UploadedFile
	var planned []UploadItem
	var flattenConflicts int
	seenFlattened := map[string]string{}

	deployFileName := filepath.Base(cfg.DeployFile)
	deployRemotePath := path.Join(cfg.ProjectPath, deployFileName)
	planned = append(planned, UploadItem{
		Source:      filepath.ToSlash(cfg.DeployFile),
		Destination: filepath.ToSlash(deployRemotePath),
		Note:        "(deploy file)",
		NoteColor:   logs.GrayColor,
	})

	for _, ef := range cfg.ExtraFiles {
		src := ef.Src
		dst := ef.Dst
		flatten := ef.Flatten

		matches, err := filepath.Glob(src)
		if err != nil {
			logs.Fatalf("Invalid glob pattern: %s", src)
		}
		if len(matches) == 0 {
			logs.Fatalf("No matches found for: %s", src)
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				logs.Fatalf("Cannot access '%s': %v", match, err)
			}

			if info.IsDir() {
				err := filepath.Walk(match, func(walkedPath string, walkedInfo os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if walkedInfo.IsDir() {
						return nil
					}

					localPath := filepath.ToSlash(walkedPath)
					var remotePath, note, color string
					relPath, _ := filepath.Rel(match, walkedPath)
					base := filepath.Base(localPath)

					if flatten {
						if dst != "" {
							if strings.HasSuffix(dst, "/") {
								remotePath = path.Join(cfg.ProjectPath, dst, base)
							} else {
								remotePath = path.Join(cfg.ProjectPath, dst)
							}
							note = "(flattened-custom)"
						} else {
							remotePath = path.Join(cfg.ProjectPath, base)
							note = "(flattened)"
						}
						if existing, ok := seenFlattened[base]; ok {
							flattenConflicts++
							logs.Fatalf("Flattening conflict: both '%s' and '%s' target '%s'", existing, localPath, base)
						}
						seenFlattened[base] = localPath
						color = logs.GrayColor
					} else if dst != "" {
						remotePath = path.Join(cfg.ProjectPath, dst, relPath)
						note = "(custom-dir)"
						color = logs.GrayColor
					} else {
						rel, err := filepath.Rel(".", localPath)
						if err != nil {
							logs.Fatalf("Failed to resolve relative path: %v", err)
						}
						remotePath = path.Join(cfg.ProjectPath, filepath.ToSlash(rel))
						note = "(preserved-dir)"
						color = logs.GrayColor
					}

					planned = append(planned, UploadItem{
						Source:      localPath,
						Destination: filepath.ToSlash(remotePath),
						Note:        note,
						NoteColor:   color,
					})
					return nil
				})
				if err != nil {
					logs.Fatalf("Failed to walk directory '%s': %v", match, err)
				}
				continue
			}

			localPath := filepath.ToSlash(match)
			var remotePath, note, color string
			base := filepath.Base(localPath)

			if flatten {
				if dst != "" {
					if strings.HasSuffix(dst, "/") {
						remotePath = path.Join(cfg.ProjectPath, dst, base)
					} else {
						remotePath = path.Join(cfg.ProjectPath, dst)
					}
					note = "(flattened-custom)"
				} else {
					remotePath = path.Join(cfg.ProjectPath, base)
					note = "(flattened)"
				}
				if existing, ok := seenFlattened[base]; ok {
					flattenConflicts++
					logs.Fatalf("Flattening conflict: both '%s' and '%s' target '%s'", existing, localPath, base)
				}
				seenFlattened[base] = localPath
				color = logs.GrayColor
			} else if dst != "" {
				if strings.HasSuffix(dst, "/") {
					remotePath = path.Join(cfg.ProjectPath, dst, filepath.Base(match))
				} else {
					remotePath = path.Join(cfg.ProjectPath, dst)
				}
				note = "(custom)"
				color = logs.GrayColor
			} else {
				rel, err := filepath.Rel(".", localPath)
				if err != nil {
					logs.Fatalf("Failed to resolve relative path: %v", err)
				}
				remotePath = path.Join(cfg.ProjectPath, filepath.ToSlash(rel))
				note = "(preserved)"
				color = logs.GrayColor
			}

			planned = append(planned, UploadItem{
				Source:      localPath,
				Destination: filepath.ToSlash(remotePath),
				Note:        note,
				NoteColor:   color,
			})
		}
	}

	if cfg.EnvVars != "" {
		planned = append(planned, UploadItem{
			Source:      ".env",
			Destination: filepath.ToSlash(path.Join(cfg.ProjectPath, ".env")),
			Note:        "(generated)",
			NoteColor:   logs.BlueColor,
		})
	}

	sort.SliceStable(planned, func(i, j int) bool {
		return planned[i].Source < planned[j].Source
	})

	maxSrcLen, maxDstLen := 0, 0
	for _, item := range planned {
		if len(item.Source) > maxSrcLen {
			maxSrcLen = len(item.Source)
		}
		if len(item.Destination) > maxDstLen {
			maxDstLen = len(item.Destination)
		}
	}

	logs.Step("\U0001F4C4 Planned uploads...")
	for _, item := range planned {
		src := fmt.Sprintf("%-*s", maxSrcLen, item.Source)
		dst := fmt.Sprintf("%-*s", maxDstLen, item.Destination)
		if item.Note != "" {
			logs.Substepf("\u2022 %s -> %s %s%s%s", src, dst, item.NoteColor, item.Note, logs.ResetColor)
		} else {
			logs.Substepf("\u2022 %s -> %s", src, dst)
		}
	}
	logs.Break()
	logs.Successf("%d files prepared for upload", len(planned))
	logs.Warnf("%d flattening conflicts", flattenConflicts)

	logs.Step("\U0001F4E6 Uploading files...")
	for _, item := range planned {
		if item.Source == ".env" && cfg.EnvVars != "" {
			logs.Verbose("Creating temporary .env file with inline variables")
			if err := os.WriteFile(".env", []byte(cfg.EnvVars), 0644); err != nil {
				logs.Fatalf("Failed to create .env file: %v", err)
			}
			defer os.Remove(".env")
		}

		logs.Verbosef("Uploading '%s' to '%s'", item.Source, item.Destination)
		if err := scp.UploadFileSCP(cli, item.Source, item.Destination); err != nil {
			logs.Fatalf("Failed to upload '%s': %v", item.Source, err)
		}
		logs.Successf("%s uploaded", filepath.Base(item.Source))

		uploaded = append(uploaded, UploadedFile{
			File:       item.Source,
			RemotePath: item.Destination,
		})
	}

	return uploaded
}
