package eth

import (
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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

func (em *EthManager) GetStamper(contractAddr string, privateKeyHex string) *stamper.Stamper {
	pk, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Could not convert from hex to private key")
	}
	auth := bind.NewKeyedTransactor(pk)

	stampStorage, err := stamper.NewStampStorage(common.HexToAddress(contractAddr), em.conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}

	return stamper.NewStamper(stampStorage, auth)
}

func InitStamper(config *configs.ContractConfig) *stamper.Stamper {
	ethManager := new(EthManager)
	ethManager.Connect(config.RawURL())
	return ethManager.GetStamper(config.Address, config.PrivateKey)
}
