package cmd

import (
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/config"
	"github.com/ava-labs/etna-devnet-resources/manual_etna_evm/helpers"
	"github.com/spf13/cobra"

	_ "embed"

	_ "embed"

	"github.com/ava-labs/avalanche-cli/pkg/vm"
	pluginEVM "github.com/ava-labs/coreth/plugin/evm"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/params"
	"github.com/ethereum/go-ethereum/common"
)

var (
	defaultPoAOwnerBalance = new(big.Int).Mul(vm.OneAvax, big.NewInt(10)) // 10 Native Tokens
)

//go:embed proxy_compiled/deployed_proxy_admin_bytecode.txt
var proxyAdminBytecodeHexString string

//go:embed proxy_compiled/deployed_transparent_proxy_bytecode.txt
var transparentProxyBytecodeHexString string

func init() {
	rootCmd.AddCommand(GenerateGenesisCmd)
}

var GenerateGenesisCmd = &cobra.Command{
	Use:   "generate-genesis",
	Short: "Generate genesis file for the L1",
	Long:  `Generate genesis file for the L1`,
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintHeader("🕸️  Generating genesis file")
		ownerKey, err := helpers.LoadSecp256k1PrivateKey(helpers.ValidatorManagerOwnerKeyPath)
		if err != nil {
			return fmt.Errorf("failed to load owner key: %s\n", err)
		}

		ethAddr := pluginEVM.PublicKeyToEthAddress(ownerKey.PublicKey())

		now := time.Now().Unix()

		feeConfig := commontype.FeeConfig{
			GasLimit:                 big.NewInt(12000000),
			TargetBlockRate:          2,
			MinBaseFee:               big.NewInt(25000000000),
			TargetGas:                big.NewInt(60000000),
			BaseFeeChangeDenominator: big.NewInt(36),
			MinBlockGasCost:          big.NewInt(0),
			MaxBlockGasCost:          big.NewInt(1000000),
			BlockGasCostStep:         big.NewInt(200000),
		}

		// zeroTime := uint64(0)

		genesis := core.Genesis{
			Config: &params.ChainConfig{
				BerlinBlock:         big.NewInt(0),
				ByzantiumBlock:      big.NewInt(0),
				ConstantinopleBlock: big.NewInt(0),
				EIP150Block:         big.NewInt(0),
				EIP155Block:         big.NewInt(0),
				EIP158Block:         big.NewInt(0),
				HomesteadBlock:      big.NewInt(0),
				IstanbulBlock:       big.NewInt(0),
				LondonBlock:         big.NewInt(0),
				MuirGlacierBlock:    big.NewInt(0),
				PetersburgBlock:     big.NewInt(0),
				FeeConfig:           feeConfig,
				ChainID:             big.NewInt(config.L1_CHAIN_ID),
			},
			Alloc: types.GenesisAlloc{
				ethAddr: {
					Balance: defaultPoAOwnerBalance,
				},
			},
			Difficulty: big.NewInt(0),
			GasLimit:   uint64(12000000),
			Timestamp:  uint64(now),
		}

		proxyAdminBytecodeHexString = strings.TrimSpace(strings.TrimPrefix(proxyAdminBytecodeHexString, "0x"))
		transparentProxyBytecodeHexString = strings.TrimSpace(strings.TrimPrefix(transparentProxyBytecodeHexString, "0x"))

		proxyAdminBytecode, err := hex.DecodeString(proxyAdminBytecodeHexString)
		if err != nil {
			return fmt.Errorf("failed to decode proxy admin bytecode: %s\n", err)
		}

		transparentProxyBytecode, err := hex.DecodeString(transparentProxyBytecodeHexString)
		if err != nil {
			return fmt.Errorf("failed to decode transparent proxy bytecode: %s\n", err)
		}

		genesis.Alloc[common.HexToAddress(config.ProxyAdminContractAddress)] = types.Account{
			Balance: big.NewInt(0),
			Code:    proxyAdminBytecode,
			Nonce:   1,
			Storage: map[common.Hash]common.Hash{
				common.HexToHash("0x0"): common.HexToHash(ethAddr.String()),
			},
		}

		genesis.Alloc[common.HexToAddress(config.ProxyContractAddress)] = types.Account{
			Balance: big.NewInt(0),
			Code:    transparentProxyBytecode,
			Nonce:   1,
			Storage: map[common.Hash]common.Hash{
				common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"): common.HexToHash(MustDeriveContractAddress(ethAddr, 1).String()),
				common.HexToHash("0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103"): common.HexToHash(config.ProxyAdminContractAddress),
			},
		}

		// Convert genesis to map to add warpConfig
		genesisMap := make(map[string]interface{})
		genesisBytes, err := json.Marshal(genesis)
		if err != nil {
			return fmt.Errorf("failed to marshal genesis to map: %s\n", err)
		}
		if err := json.Unmarshal(genesisBytes, &genesisMap); err != nil {
			return fmt.Errorf("failed to unmarshal genesis to map: %s\n", err)
		}

		// Add warpConfig to config
		configMap := genesisMap["config"].(map[string]interface{})
		configMap["warpConfig"] = map[string]interface{}{
			"blockTimestamp":               now,
			"quorumNumerator":              67,
			"requirePrimaryNetworkSigners": true,
		}

		prettyJSON, err := json.MarshalIndent(genesisMap, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal genesis: %s\n", err)
		}

		err = helpers.SaveText(helpers.L1GenesisPath, string(prettyJSON))
		if err != nil {
			return fmt.Errorf("failed to save genesis: %s\n", err)
		}

		log.Printf("Successfully wrote genesis to data/L1-genesis.json\n")

		return nil
	},
}
