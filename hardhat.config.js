require("@nomicfoundation/hardhat-toolbox");
require("dotenv").config();

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.10",
  networks: {
    goerli: {
      url: process.env.INFURA_MUMBAI_ENDPOINT,
      accounts: [process.env.PRIVATE_KEY],
      gasLimit: 300000000,
    },
  },
};
