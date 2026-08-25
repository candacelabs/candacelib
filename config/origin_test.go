package config_test

import (
	"github.com/candacelabs/candacelib/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("private origins",
	func(origin string, valid bool) {
		err := config.ValidatePrivateOrigin(origin)
		if valid {
			Expect(err).NotTo(HaveOccurred())
			return
		}
		Expect(err).To(HaveOccurred())
	},
	Entry("loopback", "http://127.0.0.1:4096", true),
	Entry("private DNS", "http://opencode:4096", true),
	Entry("tailnet", "https://100.80.127.79", true),
	Entry("public host", "https://opencode.example.com", false),
	Entry("path", "http://opencode:4096/session", false),
	Entry("credentials", "http://user:secret@opencode:4096", false),
	Entry("unsupported scheme", "ssh://opencode:4096", false),
)
