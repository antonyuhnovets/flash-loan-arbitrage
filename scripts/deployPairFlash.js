const hre = require("hardhat");

async function main() {
  console.log("deploying...");
  const PairFlash = await hre.ethers.getContractFactory("PairFlash");
  const flash = await PairFlash.deploy(
    "0xE592427A0AEce92De3Edee1F18E0157C05861564",
    "0x1F98431c8aD98523631AE4a59f267346ea31F984",
    "0x82aF49447D8a07e3bd95BD0d56f35241523fBab1"
    // Polygon
    // "0x9c3C9283D3e44854697Cd22D3Faa240Cfb032889"
    );
  
  await flash.deployed();

  console.log("Flash loan contract deployed: ", flash.address);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
