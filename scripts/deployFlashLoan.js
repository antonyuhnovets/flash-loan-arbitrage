const hre = require("hardhat");

async function main() {
  console.log("deploying...");
  const FlashLoan = await hre.ethers.getContractFactory("FlashLoan");
  const flashLoan = await FlashLoan.deploy(
    // ARB Goerli
    // "0x4EEE0BB72C2717310318f27628B3c8a708E4951C",
    // "0xE592427A0AEce92De3Edee1F18E0157C05861564"
    "0xC911B590248d127aD18546B186cC6B324e99F02c",
    "0xE592427A0AEce92De3Edee1F18E0157C05861564"
  );
  
  await flashLoan.deployed();

  console.log("Flash loan contract deployed: ", flashLoan.address);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
