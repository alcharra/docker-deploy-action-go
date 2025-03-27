package ssh

import "golang.org/x/crypto/ssh"

type Client struct {
	Host       string
	Port       string
	User       string
	PrivateKey string
	client     *ssh.Client
}
