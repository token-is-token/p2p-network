package protocol

import (
	"fmt"
	"strings"
)

const (
	MinSupportedVersion = "1.0.0"
	MaxSupportedVersion = "1.9.9"
	CurrentVersion      = "1.1.0"
)

var SupportedVersions = []string{"1.0.0", "1.1.0"}

type VersionInfo struct {
	Version   string   `json:"version"`
	Protocols []string `json:"protocols"`
}

func NewVersionInfo() *VersionInfo {
	return &VersionInfo{
		Version:   CurrentVersion,
		Protocols: SupportedVersions,
	}
}

func (v *VersionInfo) ProtocolID() string {
	return fmt.Sprintf("/llm-share/%s", v.Version)
}

func ParseVersion(protocolID string) (string, error) {
	prefix := "/llm-share/"
	if !strings.HasPrefix(protocolID, prefix) {
		return "", fmt.Errorf("invalid protocol ID format: %s", protocolID)
	}
	version := strings.TrimPrefix(protocolID, prefix)
	if version == "" {
		return "", fmt.Errorf("missing version in protocol ID: %s", protocolID)
	}
	return version, nil
}

func IsVersionSupported(version string) bool {
	for _, v := range SupportedVersions {
		if v == version {
			return true
		}
	}
	return false
}

func NegotiateVersion(remoteVersions []string) (string, error) {
	for _, remote := range remoteVersions {
		if IsVersionSupported(remote) {
			return remote, nil
		}
	}
	return "", fmt.Errorf("no compatible version found (remote: %v, supported: %v)", remoteVersions, SupportedVersions)
}

func CompareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")
	
	for i := 0; i < 3; i++ {
		var n1, n2 int
		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &n1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &n2)
		}
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
	}
	return 0
}

type VersionNegotiator struct {
	localVersion  *VersionInfo
	negotiatedVer string
}

func NewVersionNegotiator() *VersionNegotiator {
	return &VersionNegotiator{
		localVersion: NewVersionInfo(),
	}
}

func (n *VersionNegotiator) LocalVersion() *VersionInfo {
	return n.localVersion
}

func (n *VersionNegotiator) Negotiate(remoteInfo *VersionInfo) error {
	version, err := NegotiateVersion(remoteInfo.Protocols)
	if err != nil {
		return fmt.Errorf("version negotiation failed: %w", err)
	}
	n.negotiatedVer = version
	return nil
}

func (n *VersionNegotiator) NegotiatedVersion() string {
	return n.negotiatedVer
}

func (n *VersionNegotiator) IsCompatible(remoteVersion string) bool {
	return IsVersionSupported(remoteVersion)
}

const (
	MsgTypeVersionRequest  MessageType = 100
	MsgTypeVersionResponse MessageType = 101
)

type VersionRequest struct {
	Protocols []string `json:"protocols"`
}

type VersionResponse struct {
	SelectedVersion string `json:"selected_version"`
	Success         bool   `json:"success"`
	Error           string `json:"error,omitempty"`
}

func NewVersionRequest() *VersionRequest {
	return &VersionRequest{
		Protocols: SupportedVersions,
	}
}

func NewVersionResponse(selected string, success bool, errMsg string) *VersionResponse {
	return &VersionResponse{
		SelectedVersion: selected,
		Success:         success,
		Error:           errMsg,
	}
}
