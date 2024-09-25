package utils

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

var SshManagerIns *SshManager = &SshManager{
	SshInfos:   make(map[string]SshConnInfo),
	SshClients: make(map[string]*ssh.Client),
}

type SshConnInfo struct {
	SshConfig *ssh.ClientConfig
	Port      int
	Host      string
	User      string
}

type SshManager struct {

	//ssh连接信息，key为host
	SshInfos map[string]SshConnInfo
	//ssh客户端，key为host
	SshClients map[string]*ssh.Client

	mutex sync.RWMutex
}

func (c *SshManager) checkHealth() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for ip, v := range c.SshClients {
		session, err := v.NewSession()
		if err == nil {
			session.Close()
			continue
		}
		sshInfo, ok := c.SshInfos[ip]
		if !ok {
			continue
		}
		addr := fmt.Sprintf("%s:%d", sshInfo.Host, sshInfo.Port)
		newClient, err := ssh.Dial("tcp", addr, sshInfo.SshConfig)
		if err != nil {
			continue
		}
		c.SshClients[ip] = newClient

	}
}

func (c *SshManager) Run() {
	for true {
		c.checkHealth()
		time.Sleep(5 * time.Second)
	}
}

func (c *SshManager) AddHost(user string, host string, port int) error {

	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.SshClients[host] != nil {
		return nil
	}
	config := &ssh.ClientConfig{
		Timeout:         20 * time.Second,
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	key, err := ioutil.ReadFile("/root/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}
	config.Auth = []ssh.AuthMethod{
		// Use the PublicKeys method for remote authentication.
		ssh.PublicKeys(signer),
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return errors.New(fmt.Sprintf("create ssh client fail, host:%s, port:%d, user:%s", host, port, user, err))
	}
	c.SshInfos[host] = SshConnInfo{
		User:      user,
		Host:      host,
		Port:      port,
		SshConfig: config,
	}
	c.SshClients[host] = sshClient
	return nil
}

func (c *SshManager) Exec(host string, cmd string) (string, error) {

	c.mutex.RLock()
	defer c.mutex.RUnlock()
	sshClient := c.SshClients[host]
	if nil == sshClient {
		return "", errors.New(fmt.Sprintf("host: %s is not connect", host))
	}
	session, err := sshClient.NewSession()
	if err != nil {
		return "", errors.New(fmt.Sprintf("host: %s create ssh session fail, error: %v", host, err))
	}
	defer session.Close()
	res, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", errors.New(fmt.Sprintf("host: %s exec cmd fail, cmd:%s, error: %v, res:%s", host, cmd, err, res))
	} else {
		return strings.Trim(string(res), " \n"), nil
	}
}

func (c *SshManager) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for _, v := range c.SshClients {
		v.Close()
	}
}
