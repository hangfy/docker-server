package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var appPath = "/docker/app"

func NewCmdServer() *cobra.Command {
	c := cobra.Command{
		Use:   "server",
		Short: "Run bubble server",
	}
	startCmd.Flags().StringP("service", "s", "", "Name of the service (required)")
	err := startCmd.MarkFlagRequired("service")
	if err != nil {
		fmt.Println(err)
	}
	stopCmd.Flags().StringP("service", "s", "", "Name of the service (required)")
	err = stopCmd.MarkFlagRequired("service")
	if err != nil {
		fmt.Println(err)
	}

	restartCmd.Flags().StringP("service", "s", "", "Name of the service (required)")
	err = restartCmd.MarkFlagRequired("service")
	if err != nil {
		fmt.Println(err)
	}

	c.AddCommand(startCmd)
	c.AddCommand(stopCmd)
	c.AddCommand(restartCmd)
	return &c
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Docker service",
	Long:  "Start the specified Docker service using the provided service name.",
	Example: `
  docker-ctl start --service my_service`,
	Run: func(cmd *cobra.Command, args []string) {
		serviceName, _ := cmd.Flags().GetString("service")
		startService(serviceName)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a Docker service",
	Long:  "Stop the specified Docker service using the provided service name.",
	Example: `
  docker-ctl stop --service my_service`,
	Run: func(cmd *cobra.Command, args []string) {
		serviceName, _ := cmd.Flags().GetString("service")
		stopService(serviceName)
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a Docker service",
	Long:  `Stop, remove, and start the specified Docker service using the provided service name.`,
	Example: `
  docker-ctl restart --service my_service`,
	Run: func(cmd *cobra.Command, args []string) {
		serviceName, _ := cmd.Flags().GetString("service")
		restartService(serviceName)
	},
}

func RemoveFile(serverName string) error {
	entries, err := os.ReadDir(appPath)
	if err != nil {
		return errors.Wrap(err, "read dir error")
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		newName := serverName + ext
		//文件名以serverName-开头
		if strings.HasPrefix(serverName+"-", entry.Name()) {
			oldFile := filepath.Join(appPath, entry.Name())
			newFile := filepath.Join(appPath, newName)
			if fileExists(newFile) {
				err := os.Remove(newFile)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("remove file %s error", newFile))
				}
			}
			err := os.Rename(oldFile, newFile)
			if err != nil {
				return errors.Wrap(err, "rename error")
			}
			return nil
		}
	}
	return nil
}

func startService(serviceName string) {
	fmt.Printf("Starting service: %s\n", serviceName)
	// 检查操作
	err := RemoveFile(serviceName)
	if err != nil {
		fmt.Printf("Error removing file: %v\n", err)
		os.Exit(1)
	}
	// 执行启动服务的命令
	executeDockerCommand(serviceName, "up", "-d")
}

func stopService(serviceName string) {
	fmt.Printf("Stopping service: %s\n", serviceName)
	// 检查操作
	// 执行停止服务的命令
	executeDockerCommand(serviceName, "stop")
}

func restartService(serviceName string) {
	fmt.Printf("Restarting service: %s\n", serviceName)
	// 检查操作
	// 先停止服务，再移除容器，最后重新启动服务
	stopService(serviceName)
	executeDockerCommand(serviceName, "rm", "-f")
	startService(serviceName)
}

func executeDockerCommand(serviceName string, dockerArgs ...string) {
	args := append([]string{"--compatibility", "-p", "config", "-f", "/docker/config/docker-compose-config.yaml"}, dockerArgs...)
	args = append(args, serviceName)

	cmd := exec.Command("docker-compose", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing docker-compose command: %v\n", err)
		os.Exit(1)
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
