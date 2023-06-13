const hre = require("hardhat");

async function main() {
  console.log("deploying...");
  const FlashSwap = await hre.ethers.getContractFactory("UniswapV3FlashSwap");
  const flash = await FlashSwap.deploy();
  
  await flash.deployed();

  console.log("Flash loan contract deployed: ", flash.address);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
