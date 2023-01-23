package e2e_test

import (
	"context"
	"log"
	"math/rand"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/edgefarm/provider-nats/cluster/local/e2e/pkg/stream"
	utils "github.com/edgefarm/provider-nats/cluster/local/e2e/pkg/utils"
)

func TestE2e(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	RegisterFailHandler(Fail)
	err := utils.CreateFramework(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to create framework: %v", err)
	}
	RunSpecs(t, "E2e Suite")
}
