package exec

import (
	"errors"
	"fmt"
	utils "github.com/go-utils/utils/ssh"
	"github.com/spf13/cobra"
	"strings"
)

var cmdStr string
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if sshHost == "" {
			return errors.New("ssh-host can not be empty")
		}
		utils.SshManagerIns.AddHost(user, sshHost, sshPort)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if cmdStr == "" {
			if pods, err := GetNodePods(sshHost, filterStatus); err == nil {
				for key, pod := range pods {
					fmt.Println(key, "--", pod)
				}
			} else {
				return err
			}
		} else {
			fmt.Println("cmd: ", cmdStr)
			res, s2 := utils.SshManagerIns.Exec(sshHost, cmdStr)
			if s2 != nil {
				return s2
			}
			fmt.Println(res)
		}
		return nil
	},
}

var user string
var sshHost string
var sshPort int
var filterStatus string

func init() {
	//cobra.OnInitialize(initConfig)
	execCmd.PersistentFlags().StringVar(&cmdStr, "cmd", "", "")
	execCmd.PersistentFlags().StringVar(&user, "user", "root", "")
	execCmd.PersistentFlags().IntVar(&sshPort, "ssh-port", 22, "")
	execCmd.PersistentFlags().StringVar(&sshHost, "ip", "", "")
	execCmd.PersistentFlags().StringVar(&filterStatus, "filter-status", "", "")

	rootCmd.AddCommand(execCmd)
}

func GetNodePods(nodeIp, status string) (map[string]string, error) {
	cmd := fmt.Sprintf("kubectl get pods -A -o custom-columns=NAMESPACE:.metadata.namespace,NAME:.metadata.name,STATUS:.status.phase --field-selector spec.nodeName=%s|awk '{print $1 \"/\" $2 \"/\" $3}'", nodeIp)
	if status != "" {
		cmd = fmt.Sprintf("kubectl get pods -A -o custom-columns=NAMESPACE:.metadata.namespace,NAME:.metadata.name,STATUS:.status.phase --field-selector spec.nodeName=%s|grep %s|awk '{print $1 \"/\" $2 \"/\" $3}'", nodeIp, status)
	}
	res, s2 := utils.SshManagerIns.Exec(sshHost, cmd)
	if s2 != nil {
		return nil, s2
	}

	podStatusMap := map[string]string{}
	podLines := strings.Split(res, "\n")
	for _, line := range podLines {
		line = strings.TrimSpace(line)
		if line != "" {
			podElements := strings.Split(line, "/")
			podStatusMap[fmt.Sprintf("%s/%s", podElements[0], podElements[1])] = podElements[2]
		}
	}

	return podStatusMap, nil
}
