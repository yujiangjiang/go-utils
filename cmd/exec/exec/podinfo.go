package exec

import (
	"encoding/json"
	"fmt"
	"github.com/go-utils/utils/exec"
	utils "github.com/go-utils/utils/ssh"
	"github.com/spf13/cobra"
	"strings"
)

type infoOption struct {
	namespace string
	podName   string
	host      string
}

var info infoOption

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		utils.SshManagerIns.AddHost("root", info.host, 22)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		statuses, err := GetPodContainersInfo(info)
		if err != nil {
			return err
		}
		for _, status := range statuses {
			fmt.Println(status)
		}

		fmt.Println(CurVersionIsV2())
		return nil
	},
}

func init() {
	infoCmd.PersistentFlags().StringVar(&info.podName, "name", "", "")
	infoCmd.PersistentFlags().StringVar(&info.namespace, "namespace", "", "")
	infoCmd.PersistentFlags().StringVar(&info.host, "host", "", "")
	rootCmd.AddCommand(infoCmd)
}

type ContainerStatus struct {
	Name    string                 `json:"name"`
	Ready   bool                   `json:"ready"`
	Started bool                   `json:"started"`
	State   map[string]interface{} `json:"state"`
}

func GetPodContainersInfo(info infoOption) ([]ContainerStatus, error) {
	cmd := fmt.Sprintf("kubectl get po %s -n%s -ojson|jq '.status.containerStatuses'", info.podName, info.namespace)
	res, s2 := utils.SshManagerIns.Exec(info.host, cmd)
	if s2 != nil {
		return nil, s2
	}
	fmt.Println(res)
	if res == "" || res == "[]" || res == "nil" {
		return nil, nil
	}
	var containerStatuses []ContainerStatus
	if err := json.Unmarshal([]byte(res), &containerStatuses); err != nil {
		return nil, err
	}

	return containerStatuses, nil
}

func CurVersionIsV2() bool {
	res, err := exec.Exec("bash", "-c", "cat /etc/systemd/magiccube/version")

	if err != "" {
		fmt.Errorf(err)
	}
	fmt.Println(res)
	return strings.HasPrefix(res, "V2")
}
