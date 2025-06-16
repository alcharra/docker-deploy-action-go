package tests

import "github.com/alcharra/docker-deploy-action-go/config"

func extraFilesToEnv(files []config.ExtraFile) []string {
	var lines []string
	for _, f := range files {
		line := ""
		if f.Flatten {
			line += "flatten "
		}
		line += f.Src
		if f.Dst != "" {
			line += ":" + f.Dst
		}
		lines = append(lines, line)
	}
	return lines
}
