module.exports = {
    networks: {
        development: {
            host: "127.0.0.1",
            port: 8545,
            network_id: "5777"
        },
        geth: {
            host: "127.0.0.1",
            port: 8545,
            network_id: "1999",
            gas: 4600000
        }
    }
};
