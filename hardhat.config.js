/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
    solidity: "0.8.17",
    networks: {
        hardhat: {
            allowUnlimitedContractSize: false,
            mining: {
                auto: true,
                interval: [3000,6000]
            }
        }
    }
};