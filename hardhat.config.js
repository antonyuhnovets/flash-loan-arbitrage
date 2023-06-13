require("@nomicfoundation/hardhat-toolbox");
require("dotenv").config();

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: {
    compilers: [
      {
        version: "0.5.0",
        settings: {},
      },
      {
        version: "0.7.6",
        settings: {},
      },
      {
        version: "0.8.17",
        settings: {},
      },
    ],
  },
  networks: {
    mumbai: {
      url: process.env.INFURA_MUMBAI_ENDPOINT,
      accounts: [process.env.PRIVATE_KEY],
      gasLimit: 300000000,
    },
    goerli: {
      url: process.env.INFURA_ETH_GOERLI_ENDPOINT,
      accounts: [process.env.PRIVATE_KEY],
      gasLimit: 300000000,
    },
    sepolita: {
      url: process.env.INFURA_ETH_SEPOLITA_ENDPOINT,
      accounts: [process.env.PRIVATE_KEY],
      gasLimit: 300000000,
    },
  },
  etherscan: {
    apiKey: {
      arbitrumGoerli: process.env.ARB_GOERLI_KEY,
    }
  }
};
