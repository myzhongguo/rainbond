// RAINBOND, Application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

// thanks @lextoumbourou.

package initiate

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/eapache/channels"
	"github.com/goodrain/rainbond/cmd/node/option"
	"github.com/goodrain/rainbond/node/kubecache"

	"github.com/Sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	commentChar    string = "#"
	startOfSection        = "# Generate by Rainbond. DO NOT EDIT"
	endOfSection          = "# End of Section"
	eol                   = "\n"
)

var (
	// ErrRegistryAddressNotFound record not found error, happens when haven't find any matched registry address.
	ErrRegistryAddressNotFound = errors.New("registry address not found")
)

// HostManager is responsible for writing the resolution of the private image repository domain name to /etc/hosts.
type HostManager interface {
	Start()
	Stop()
}

// NewHostManager creates a new HostManager.
func NewHostManager(cfg *option.Conf, kubecli kubecache.KubeClient, registryUpdateCh *channels.RingChannel) HostManager {
	return &hostManager{
		registryUpdateCh: registryUpdateCh,
		stopCh:           make(chan struct{}),
		cfg:              cfg,
		kubecli:          kubecli,
	}
}

type hostManager struct {
	registryUpdateCh *channels.RingChannel
	stopCh           chan struct{}

	cfg     *option.Conf
	kubecli kubecache.KubeClient
}

func (h *hostManager) Start() {
	for {
		select {
		case <-h.registryUpdateCh.Out():
			// ignore all error. If an error occurs, do not write the hosts file
			ep, err := h.getEndpointForRegistry()
			if err != nil {
				logrus.Warningf("finding endpoint for registry: %v", err)
				break
			}
			if err := h.CleanupAndFlush(ep, h.cfg.ImageRepositoryHost); err != nil {
				logrus.Warningf("clean and flush /etc/hosts: %v", err)
			}
		case <-h.stopCh:
			return
		}
	}
}

func (h *hostManager) Stop() {
	close(h.stopCh)
}

func (h *hostManager) getEndpointForRegistry() (string, error) {
	selector, _ := labels.Parse("name=rbd-hub")
	pods, err := h.kubecli.GetPodsBySelector(h.cfg.RbdNamespace, selector)
	if err != nil {
		return "", err
	}
	for _, po := range pods {
		for _, cdt := range po.Status.Conditions {
			if cdt.Type == corev1.PodReady && cdt.Status == corev1.ConditionTrue {
				return po.Status.HostIP, nil
			}
		}
	}

	return "", ErrRegistryAddressNotFound
}

// CleanupAndFlush cleanup old content and write new ones.
func (h *hostManager) CleanupAndFlush(ip, domain string) error {
	hosts, err := NewHosts()
	if err != nil {
		return fmt.Errorf("error creating hosts: %v", err)
	}

	if err := hosts.Cleanup(); err != nil {
		return fmt.Errorf("error cleanup hosts: %v", err)
	}

	lines := []string{
		startOfSection,
		ip + " " + domain,
		endOfSection,
	}
	hosts.AddLines(lines...)

	return hosts.Flush()
}

// HostsLine represents a single line in the hosts file.
type HostsLine struct {
	IP    string
	Hosts []string
	Raw   string
	Err   error
}

// NewHostsLine returns a new instance of ```HostsLine```.
func NewHostsLine(raw string) HostsLine {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return HostsLine{Raw: raw}
	}

	output := HostsLine{Raw: raw}
	if !output.IsComment() {
		rawIP := fields[0]
		if net.ParseIP(rawIP) == nil {
			output.Err = errors.New(fmt.Sprintf("Bad hosts line: %q", raw))
		}

		output.IP = rawIP
		output.Hosts = fields[1:]
	}

	return output
}

// IsComment returns ```true``` if the line is a comment.
func (l HostsLine) IsComment() bool {
	trimLine := strings.TrimSpace(l.Raw)
	isComment := strings.HasPrefix(trimLine, commentChar)
	return isComment
}

// Hosts represents a hosts file.
type Hosts struct {
	Path  string
	Lines []HostsLine
}

// NewHosts return a new instance of ``Hosts``.
func NewHosts() (Hosts, error) {
	hosts := Hosts{Path: "/etc/hosts"}

	err := hosts.load()
	if err != nil {
		return hosts, err
	}

	return hosts, nil
}

// load the hosts file into ```l.Lines```.
// ```Load()``` is called by ```NewHosts()``` and ```Hosts.Flush()``` so you
// generally you won't need to call this yourself.
func (h *Hosts) load() error {
	var lines []HostsLine

	file, err := os.Open(h.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := NewHostsLine(scanner.Text())
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	h.Lines = lines

	return nil
}

// Add an entry to the hosts file.
func (h *Hosts) Add(ip string, hosts ...string) {
	position := h.getIPPosition(ip)
	if position == -1 {
		endLine := NewHostsLine(buildRawLine(ip, hosts))
		// Ip line is not in file, so we just append our new line.
		h.Lines = append(h.Lines, endLine)
	} else {
		// Otherwise, we replace the line in the correct position
		newHosts := h.Lines[position].Hosts
		for _, addHost := range hosts {
			if itemInSlice(addHost, newHosts) {
				continue
			}

			newHosts = append(newHosts, addHost)
		}
		endLine := NewHostsLine(buildRawLine(ip, newHosts))
		h.Lines[position] = endLine
	}
}

// AddLines adds entries to the hosts file.
func (h *Hosts) AddLines(lines ...string) {
	for _, line := range lines {
		h.Lines = append(h.Lines, NewHostsLine(line))
	}
}

// Cleanup remove entries created by rainbond from the hosts file.
func (h *Hosts) Cleanup() error {
	start := h.getStartPosition()
	if start == -1 {
		return nil
	}
	end := h.getEndPosition(start)
	if end == -1 {
		return fmt.Errorf("wrong hosts file, found start of section, but no end of section")
	}
	end += start + 1
	if end == len(h.Lines) {
		h.Lines = h.Lines[:start]
		return nil
	}

	pre := h.Lines[:start]
	post := h.Lines[end+1:]
	h.Lines = append(pre, post...)
	return nil
}

// Flush any changes made to hosts file.
func (h Hosts) Flush() error {
	file, err := os.Create(h.Path)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(file)

	for _, line := range h.Lines {
		_, _ = fmt.Fprintf(w, "%s%s", line.Raw, eol)
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return h.load()
}

func (h Hosts) getStartPosition() int {
	for i := range h.Lines {
		line := h.Lines[i]
		if line.Raw == startOfSection {
			return i
		}
	}

	return -1
}

func (h Hosts) getEndPosition(startPos int) int {
	newLines := h.Lines[startPos+1:]
	for i := range newLines {
		line := newLines[i]
		if line.Raw == endOfSection {
			return i
		}
	}

	return -1
}

func (h Hosts) getIPPosition(ip string) int {
	for i := range h.Lines {
		line := h.Lines[i]
		if !line.IsComment() && line.Raw != "" {
			if line.IP == ip {
				return i
			}
		}
	}

	return -1
}

func itemInSlice(item string, list []string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}

	return false
}

func buildRawLine(ip string, hosts []string) string {
	output := ip
	for _, host := range hosts {
		output = fmt.Sprintf("%s %s", output, host)
	}

	return output
}
