package cluster

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/dunstall/gorqlite"
	log "github.com/sirupsen/logrus"
)

type ToxiproxyNode struct {
	id string
}

func NewToxiproxyNode(id string, targetPort uint16, proxyPort uint16) (ToxiproxyNode, error) {
	lg, err := newLogger(fmt.Sprintf("toxiproxy-create-%s", id))
	if err != nil {
		return ToxiproxyNode{}, gorqlite.WrapError(err, "failed to create toxiproxy node")
	}

	cmd := exec.Command(
		"toxiproxy-cli",
		"create",
		id,
		"--listen", fmt.Sprintf("0.0.0.0:%d", proxyPort),
		"--upstream", fmt.Sprintf("localhost:%d", targetPort),
	)
	log.WithFields(log.Fields{
		"cmd": strings.Join(cmd.Args, " "),
	}).Debug("running command")
	cmd.Stdout = lg
	cmd.Stderr = lg
	if err := cmd.Run(); err != nil {
		return ToxiproxyNode{}, gorqlite.WrapError(err, "failed to create toxiproxy node")
	}

	log.WithFields(log.Fields{
		"id":          id,
		"target_port": targetPort,
		"proxy_port":  proxyPort,
		"log":         lg.Path,
	}).Info("started proxy")

	return ToxiproxyNode{
		id: id,
	}, nil
}

func (n *ToxiproxyNode) Close() error {
	lg, err := newLogger(fmt.Sprintf("toxiproxy-delete-%s", n.id))
	if err != nil {
		return gorqlite.WrapError(err, "failed to close toxiproxy node")
	}

	cmd := exec.Command("toxiproxy-cli", "delete", n.id)
	log.WithFields(log.Fields{
		"cmd": strings.Join(cmd.Args, " "),
	}).Debug("running command")
	cmd.Stdout = lg
	cmd.Stderr = lg
	if err := cmd.Run(); err != nil {
		return gorqlite.WrapError(err, "failed to close toxiproxy node")
	}

	log.WithFields(log.Fields{
		"id":  n.id,
		"log": lg.Path,
	}).Info("stopped proxy")

	return nil
}
