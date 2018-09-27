package utils_test

import (
	"testing"

	"github.com/RTradeLtd/database/utils"
)

var (
	testIpfsMultiAddrString = "/ip4/192.168.1.242/tcp/4001/ipfs/QmXivHtDyAe8nS7cbQiS7ri9haUM2wGvbinjKws3a4EstT"
	testP2PMultiAddrString  = "/ip4/192.168.1.242/tcp/4001/ipfs/QmXivHtDyAe8nS7cbQiS7ri9haUM2wGvbinjKws3a4EstT"
)

func TestGenerateMultiAddrAndParsePeerID(t *testing.T) {
	addr, err := utils.GenerateMultiAddrFromString(testP2PMultiAddrString)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = utils.ParsePeerIDFromIPFSMultiAddr(addr); err != nil {
		t.Fatal(err)
	}

	addr, err = utils.GenerateMultiAddrFromString(testIpfsMultiAddrString)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = utils.ParsePeerIDFromIPFSMultiAddr(addr); err != nil {
		t.Fatal(err)
	}
}
