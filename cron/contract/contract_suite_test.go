package contract_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCronContract(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cron Contract Suite")
}
