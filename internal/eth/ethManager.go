package eth

import (
	"io"
	"log"
	"math"
	"os"
	"time"

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
	conn, err := exponentialBackoff(rawURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client, %s: %v", rawURL, err)
	}

	em.conn = conn
}

func exponentialBackoff(rawURL string) (*ethclient.Client, error) {
	conn, err := ethclient.Dial(rawURL)
	if err == nil {
		return conn, nil
	}
	for i := uint(1); i < 6; i++ {
		d := math.Exp2(float64(i))
		log.Printf("Could not connect to Ethereum client. Will retry in %f seconds", d)
		time.Sleep(time.Duration(d) * time.Second)
		if err == nil {
			return conn, nil
		}
	}
	return nil, err
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
