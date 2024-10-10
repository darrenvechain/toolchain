package contracts

//go:generate docker run -v ./:/sources ethereum/solc:0.8.19 -o /sources --abi --bin /sources/Toolchain.sol --overwrite
