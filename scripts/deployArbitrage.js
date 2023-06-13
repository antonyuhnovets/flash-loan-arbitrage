const hre = require("hardhat");

async function main() {
  console.log("deploying...");
  const Arbitrage = await hre.ethers.getContractFactory("Arbitrage");
  const arbitrage = await Arbitrage.deploy(
    // Goerli
    "0xC911B590248d127aD18546B186cC6B324e99F02c", // Address Provider
    "0x68b3465833fb72A70ecDF485E0e4C7bD8665Fc45", // Uniswap Router
    "0x1b02dA8Cb0d097eB8D57A175b88c7D8b47997506", // Sushiswap Router
    "0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6", // WETH
    "0xBa8DCeD3512925e52FE67b1b5329187589072A55"  // DAI (?)
  );
  
  await arbitrage.deployed();

  console.log("Flash loan contract deployed: ", arbitrage.address);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
