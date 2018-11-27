package eth

import (
	"io"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/belljustin/stamp/internal/configs"
	"github.com/belljustin/stamp/pkg/stamper"
)

type EthManager struct {
	conn *ethclient.Client
}

func (em *EthManager) Connect(rawURL string) {
	conn, err := ethclient.Dial(rawURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client, %s: %v", rawURL, err)
	}

	em.conn = conn
}

func (em *EthManager) GetStamper(keyin io.Reader, password, contractAddr string) *stamper.Stamper {
	auth, err := bind.NewTransactor(keyin, password)
	if err != nil {
		log.Fatalf("Could not decrypt private key: %v", err)
	}

	stampStorage, err := stamper.NewStampStorage(common.HexToAddress(contractAddr), em.conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}

	return stamper.NewStamper(stampStorage, auth)
}

func InitStamper(config *configs.ContractConfig) *stamper.Stamper {
	ethManager := new(EthManager)
	ethManager.Connect(config.RawURL())
	keyin, err := os.Open(config.PrivateKey)
	if err != nil {
		log.Fatalf("Could not read private key file: %v", config.PrivateKey)
	}
	defer keyin.Close()
	return ethManager.GetStamper(keyin, config.Password, config.Address)
}
