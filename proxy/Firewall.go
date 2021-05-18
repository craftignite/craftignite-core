package proxy

import (
	"log"
	"os/exec"
)

func InstallRedirect() {
	log.Println("Enabling firewall redirect...")
	cmd := exec.Command("iptables", "-t", "nat", "-A", "PREROUTING", "-i", "eth0", "-p", "tcp", "--dport", "25565", "-j", "REDIRECT", "--to-port", "25566")
	v, err := cmd.Output()
	log.Println(v)
	if err != nil {
		log.Fatalln(err)
	}
}

func UninstallRedirect() {
	log.Println("Disabling firewall redirect...")
	cmd := exec.Command("iptables", "-t", "nat", "-D", "PREROUTING", "-i", "eth0", "-p", "tcp", "--dport", "25565", "-j", "REDIRECT", "--to-port", "25566")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalln(err)
	}
}
