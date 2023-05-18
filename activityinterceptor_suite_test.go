package interceptor_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestActivityinterceptor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Activityinterceptor Suite")
}
