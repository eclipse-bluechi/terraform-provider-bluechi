/* SPDX-License-Identifier: LGPL-2.1-or-later */
package client

import (
	"fmt"
	"net"
	"os"
	"path"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func ignoreHostKeyCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	return nil
}

type SSHClient struct {
	Host                  string
	User                  string
	Password              string
	PKPath                string
	InsecureIgnoreHostKey bool

	conn        *ssh.Client
	connHasRoot bool
}

func (c *SSHClient) newSSHSession() (*ssh.Session, error) {
	if c == nil || c.conn == nil {
		return nil, fmt.Errorf("not connected")
	}

	return c.conn.NewSession()
}

func (c *SSHClient) determineRootPrivileges() error {
	session, err := c.newSSHSession()
	if err != nil {
		return err
	}
	defer session.Close()

	output, err := session.Output("whoami")
	if err != nil {
		return fmt.Errorf("failed to determine if root: (%s, %s)", err.Error(), string(output))
	}

	c.connHasRoot = (strings.TrimSpace(string(output)) == "root")
	return nil
}

func (c *SSHClient) runCommand(cmd string) ([]byte, error) {
	session, err := c.newSSHSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	if !c.connHasRoot {
		cmd = "sudo " + cmd
	}

	return session.Output(cmd)
}

func (c *SSHClient) isServiceInstalled(service string) (bool, error) {
	output, err := c.runCommand(fmt.Sprintf("systemctl list-unit-files %s", service))
	if err != nil {
		if serr, ok := err.(*ssh.ExitError); ok && serr.ExitStatus() == 1 {
			return false, nil
		}
		return false, fmt.Errorf("failed to list unit files: %s", string(output))
	}

	return strings.Contains(string(output), service), nil
}

func (c *SSHClient) determineOS() (string, error) {
	output, err := c.runCommand("cat /etc/os-release | grep -w ID=")
	if err != nil {
		return "", fmt.Errorf("failed to determine os: %s", string(output))
	}

	os := strings.ReplaceAll(string(output), "ID=", "")
	os = strings.ReplaceAll(os, "\"", "")
	return strings.TrimSpace(os), nil
}

func (c *SSHClient) Connect() error {
	var err error
	var authMethods []ssh.AuthMethod
	var hostkeyCallback ssh.HostKeyCallback

	if c.PKPath != "" {
		pkPath := c.PKPath
		// resolve home directory
		if strings.HasPrefix(pkPath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			pkPath = path.Join(homeDir, strings.Replace(pkPath, "~/", "", 1))
		}
		pKey, err := os.ReadFile(pkPath)
		if err != nil {
			return err
		}

		signer, err := ssh.ParsePrivateKey(pKey)
		if err != nil {
			return err
		}

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if c.Password != "" {
		authMethods = append(authMethods, ssh.Password(c.Password))
	}

	hostkeyCallback = ignoreHostKeyCallback
	if !c.InsecureIgnoreHostKey {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		knownHostsPath := fmt.Sprintf("%s/.ssh/known_hosts", homeDir)
		hostkeyCallback, err = knownhosts.New(knownHostsPath)
		if err != nil {
			return err
		}
	}

	conf := &ssh.ClientConfig{
		User:            c.User,
		HostKeyCallback: hostkeyCallback,
		Auth:            authMethods,
	}

	c.conn, err = ssh.Dial("tcp", c.Host, conf)
	if err != nil {
		return err
	}

	err = c.determineRootPrivileges()
	if err != nil {
		return err
	}

	return nil
}

func (c *SSHClient) Disconnect() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func (c *SSHClient) InstallBlueChi(installCtrl bool, installAgent bool) error {
	needsInstallCtrl := false
	needsInstallAgent := false

	if installCtrl {
		isInstalled, err := c.isServiceInstalled("bluechi-controller.service")
		if err != nil {
			return err
		}
		needsInstallCtrl = !isInstalled
	}
	if installAgent {
		isInstalled, err := c.isServiceInstalled("bluechi-agent.service")
		if err != nil {
			return err
		}
		needsInstallAgent = !isInstalled
	}

	if !needsInstallCtrl && !needsInstallAgent {
		return nil
	}

	os, err := c.determineOS()
	if err != nil {
		return err
	}

	if os == "autosd" || os == "centos" {
		packagesToInstall := ""
		if needsInstallCtrl {
			packagesToInstall += " bluechi-controller bluechi-ctl "
		}
		if needsInstallAgent {
			packagesToInstall += " bluechi-agent"
		}
		if packagesToInstall == "" {
			return nil
		}

		output, err := c.runCommand(fmt.Sprintf("dnf install -y %s", packagesToInstall))
		if err != nil {
			return fmt.Errorf("failed to install packages '%s': %s", packagesToInstall, output)
		}
	}

	return nil
}

func (c *SSHClient) CreateControllerConfig(file string, cfg BlueChiControllerConfig) error {
	cmd := fmt.Sprintf("bash -c 'echo \"%s\" > %s'", cfg.Serialize(), path.Join(BlueChiControllerConfdDirectory, file))
	output, err := c.runCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to create controller config file: %s", string(output))
	}

	return nil
}

func (c *SSHClient) RemoveControllerConfig(file string) error {
	cmd := fmt.Sprintf("rm -f %s", path.Join(BlueChiControllerConfdDirectory, file))
	output, err := c.runCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to remove controller config file: %s", string(output))
	}

	return nil
}

func (c *SSHClient) RestartBlueChiController() error {
	output, err := c.runCommand("systemctl start bluechi-controller")
	if err != nil {
		return fmt.Errorf("failed to restart controller service: %s", string(output))
	}

	return nil
}

func (c *SSHClient) StopBlueChiController() error {
	output, err := c.runCommand("systemctl stop bluechi-controller")
	if err != nil {
		return fmt.Errorf("failed to stop controller service: %s", string(output))
	}

	return nil
}

func (c *SSHClient) CreateAgentConfig(file string, cfg BlueChiAgentConfig) error {
	cmd := fmt.Sprintf("bash -c 'echo \"%s\" > %s'", cfg.Serialize(), path.Join(BlueChiAgentConfdDirectory, file))
	output, err := c.runCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to create agent config file: %s", string(output))
	}

	return nil
}

func (c *SSHClient) RemoveAgentConfig(file string) error {
	cmd := fmt.Sprintf("rm -f %s", path.Join(BlueChiAgentConfdDirectory, file))
	output, err := c.runCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to remove agent config file: %s", string(output))
	}

	return nil
}

func (c *SSHClient) RestartBlueChiAgent() error {
	output, err := c.runCommand("systemctl start bluechi-agent")
	if err != nil {
		return fmt.Errorf("failed to restart agent service: %s", string(output))
	}

	return nil
}

func (c *SSHClient) StopBlueChiAgent() error {
	output, err := c.runCommand("systemctl stop bluechi-agent")
	if err != nil {
		return fmt.Errorf("failed to stop agent service: %s", string(output))
	}

	return nil
}

func NewSSHClient(host string, user string, password string, pkPath string, insecureIgnoreHostKey bool) Client {
	return &SSHClient{
		Host:                  host,
		User:                  user,
		Password:              password,
		PKPath:                pkPath,
		InsecureIgnoreHostKey: insecureIgnoreHostKey,
	}
}
