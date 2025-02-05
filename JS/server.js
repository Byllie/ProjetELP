var net = require('net');
class Server {
    constructor(HOST, PORT) {
        this.net = require('net');
        this.server = this.net.createServer((socket) => {
            console.log('Client connected');
            socket.on('data', (data) => {
                console.log(data.toString());
                socket.write('Hello from server!');
            });
            socket.on('end', () => {
                console.log('Client disconnected');
            });
        });
        this.server.listen(PORT, HOST);
    }
}
module.exports = { Server };