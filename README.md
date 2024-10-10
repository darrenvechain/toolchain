# toolchain simulator

Simple script to simulate toolchain traffic on vechain thor.

## Usage

Flags:
- `thor-url`: URL of the thor node to connect to.
- `mnemonic`: Mnemonic of the wallet to use for funding.
- `mnemonic-accounts`: Number of accounts to use for funding. Eg. `10` for thor solo
- `require-funding`: `true` to fund the extra accounts before running the simulation.
- `extra-account-multiplier`: Multiplier for the number of extra accounts to use. Multiply the number of mnemonic accounts by this number to get the total number of accounts to use for the simulation.
accounts to spam transactions.

### Example

- This example has accounts already funded.
- It will fund a further 10 accounts.
- 20 accounts will be used in the simulation.

```bash
go run . \
    --thor-url=https://testnet.vechain.org 
    --mnemonic="denial kitchen pet squirrel other broom bar gas better priority spoil cross" \
    --mnemonic-accounts=10 \
    --require-funding=true \
    --extra-account-multiplier=2
```
