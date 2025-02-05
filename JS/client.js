var net = require('net');
class Client {
    constructor(HOST, PORT) {
        this.socket = new net.Socket();
        this.socket.connect(PORT, HOST);
        this.socket.on('data', (data) => {
            console.log(data.toString());
        });
    }

    stop() {
        this.socket.end();
    }
}
module.exports = { Client };