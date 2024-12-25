package cmd

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	keysFolder = "keys/"
)

func init() {
	rootCmd.AddCommand(GenerateNewValidatorKeys)
	rootCmd.AddCommand(PrintBLSKeysOfValidator)
}

var GenerateNewValidatorKeys = &cobra.Command{
	Use:   "generate-new-validator-keys",
	Short: "Generate new validator keys",
	Long:  `Generate new validator keys of give amount (Usage: generate-new-validator-keys 10 ;; will generate 10 validator keys)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintHeader("🧱 Generating new validator keys")

		if len(args) == 0 {
			return fmt.Errorf("amount of validator keys is required")
		}

		// get node index id from args
		count, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid amount of validator keys: %w", err)
		}

		for i := 0; i < count; i++ {
			nodeIndexId := fmt.Sprintf("%d", i)

			// check if the folder exists, rules `data/validator_<node_index_id>`
			folderPath := fmt.Sprintf("%s/validator_%s/", keysFolder, nodeIndexId)
			// exit if the folder exists
			if _, err := os.Stat(folderPath); !os.IsNotExist(err) {
				return fmt.Errorf("folder %s already exists", folderPath)
			}

			// create the folder
			err := os.MkdirAll(folderPath, 0755)
			if err != nil {
				return fmt.Errorf("failed to create folder %s: %w", folderPath, err)
			}

			// generate the keys
			err = GenerateCredsIfNotExists(folderPath)
			if err != nil {
				return fmt.Errorf("failed to generate creds: %w", err)
			}
		}

		for i := 0; i < count; i++ {
			nodeIndexId := fmt.Sprintf("%d", i)
			folderPath := fmt.Sprintf("%s/validator_%s/", keysFolder, nodeIndexId)
			err = printBLSKeysOfValidator(folderPath)
			if err != nil {
				return fmt.Errorf("failed to print bls keys of validator: %w", err)
			}
		}

		return nil
	},
}

func printBLSKeysOfValidator(folderPath string) error {
	nodeId, proofOfPossession, err := NodeInfoFromCreds(folderPath)
	if err != nil {
		return fmt.Errorf("failed to get node info from creds: %w", err)
	}

	publicKey := "0x" + hex.EncodeToString(proofOfPossession.PublicKey[:])
	pop := "0x" + hex.EncodeToString(proofOfPossession.ProofOfPossession[:])

	fmt.Printf("🔑 BLS keys of %s\n", folderPath)
	fmt.Printf("Node ID: %s\n", nodeId)
	fmt.Printf("Public Key: %s\n", publicKey)
	fmt.Printf("Proof of Possession: %s\n\n", pop)

	stakerKeyBase64, stakerCertBase64, signerKeyBase64, err := GetCredsBase64(folderPath)
	if err != nil {
		return fmt.Errorf("failed to get creds base64: %w", err)
	}

	fmt.Printf("- AVALANCHEGO_STAKING_TLS_KEY_FILE_CONTENT=%s\n", stakerKeyBase64)
	fmt.Printf("- AVALANCHEGO_STAKING_TLS_CERT_FILE_CONTENT=%s\n", stakerCertBase64)
	fmt.Printf("- BLS_KEY_BASE64=%s\n\n\n", signerKeyBase64)

	return nil
}

var PrintBLSKeysOfValidator = &cobra.Command{
	Use:   "print-bls-key",
	Short: "Print BLS keys of validator",
	Long:  `Print BLS keys of validator`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("node index id is required")
		}

		// get node index id from args
		nodeIndexId := args[0]

		folderPath := fmt.Sprintf("%s/validator_%s/", keysFolder, nodeIndexId)

		return printBLSKeysOfValidator(folderPath)
	},
}
