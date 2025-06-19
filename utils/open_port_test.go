package utils_test

import (
	"net"
	"net/netip"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jghiloni/go-commonutils/v2/utils"
)

var _ = Describe("OpenPort", func() {
	var l *net.TCPListener
	var u *net.UDPConn

	BeforeEach(func() {
		defer GinkgoRecover()
		l, _ = net.ListenTCP("tcp", net.TCPAddrFromAddrPort(netip.MustParseAddrPort("127.0.0.1:34567")))
		u, _ = net.ListenUDP("udp", net.UDPAddrFromAddrPort(netip.MustParseAddrPort("127.0.0.1:12346")))
	})

	AfterEach(func() {
		l.Close()
		u.Close()
	})
	DescribeTable("iterations", func(errExpected bool, portExpected uint16, ports ...uint16) {
		port, err := utils.FindOpenLocalPort(50*time.Millisecond, ports...)
		errAssertion := Expect(err)
		checker := errAssertion.ShouldNot
		if errExpected {
			checker = errAssertion.Should
		}
		checker(HaveOccurred())

		if portExpected != 0 {
			Expect(port).Should(Equal(portExpected))
		}
	},
		Entry("tcp no optional args", false, uint16(1025)),
		Entry("tcp 1 optional arg", false, uint16(30000), uint16(30000)),
		Entry("tcp two optional args", false, uint16(34568), uint16(34567), uint16(34569)),
		Entry("tcp multipe args", false, uint16(12346), uint16(34567), uint16(12346), uint16(55555)),
		Entry("tcp error privileged", true, uint16(0), uint16(100)),
		Entry("tcp error no open ports", true, uint16(0), uint16(34567), uint16(34567)),
	)
})
