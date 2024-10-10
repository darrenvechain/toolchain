package main

import "flag"

var (
	thorUrlFlag            = flag.String("thor-url", "http://localhost:8669", "URL of the Thor node")
	mnemonicFlag           = flag.String("mnemonic", "denial kitchen pet squirrel other broom bar gas better priority spoil cross", "mnemonic to use to fund the transactions")
	mnemonicAccountsFlag   = flag.Int("mnemonic-accounts", 10, "number of accounts to use for the funding")
	extraAccountMultiplier = flag.Int("extra-account-multiplier", 5, "the multiple of mnemonic accounts to use for the transactions")
	requireFundingFlag     = flag.Bool("require-funding", false, "require funding of the accounts")
)
