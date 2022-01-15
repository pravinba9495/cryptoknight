const Web3 = require("web3");
const fs = require('fs')

let provider = ""
let tokenAddress = "";
let walletAddress = "";

process.argv.forEach((arg) => {
  if(arg.includes("--token")){
    tokenAddress = arg.split("=")[1]
  }
  if(arg.includes("--address")){
    walletAddress = arg.split("=")[1]
  }
  if(arg.includes("--provider")){
    provider = arg.split("=")[1]
  }
})

const Web3Client = new Web3(new Web3.providers.HttpProvider(provider));
const minABI = [
  {
    constant: true,
    inputs: [{ name: "_owner", type: "address" }],
    name: "balanceOf",
    outputs: [{ name: "balance", type: "uint256" }],
    type: "function",
  },
];
const contract = new Web3Client.eth.Contract(minABI, tokenAddress);
contract.methods.balanceOf(walletAddress).call().then((balance) => {
  try {
    fs.writeFileSync(tokenAddress + '.txt', balance)
  } catch (err) {
    console.error(err)
  }
}).catch((error) => {
  try {
    fs.writeFileSync(tokenAddress + '.txt', 0)
  } catch (err) {
    console.error(err)
  }
})