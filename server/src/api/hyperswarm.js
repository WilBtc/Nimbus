// hyperswarm.js
const Hyperswarm = require('hyperswarm');
const crypto = require('crypto');
const EventEmitter = require('events');

// HyperswarmManager manages the peer discovery and connection process
class HyperswarmManager extends EventEmitter {
    constructor(logger) {
        super();
        this.swarm = new Hyperswarm();
        this.logger = logger || console;
        this.peers = new Map(); // Store connected peers with peer ID as the key
    }

    // Initialize Hyperswarm with a topic and start discovering peers
    initialize(topic) {
        const discoveryKey = this._generateDiscoveryKey(topic);
        this.logger.info(`Initializing Hyperswarm with discovery key: ${discoveryKey.toString('hex')}`);

        // Join the swarm for peer discovery
        this.swarm.join(discoveryKey, {
            lookup: true, // Find and connect to peers
            announce: true // Announce self as a peer
        });

        // Handle new peer connections
        this.swarm.on('connection', (socket, details) => {
            const peerId = details.peer ? details.peer.host + ':' + details.peer.port : 'unknown';
            this.logger.info(`New peer connected: ${peerId}`);
            this._handleNewPeer(peerId, socket);
        });

        // Handle errors
        this.swarm.on('error', (err) => {
            this.logger.error(`Hyperswarm error: ${err.message}`);
            this.emit('error', err); // Emit error event
        });
    }

    // Close all peer connections and shut down the swarm
    async close() {
        this.logger.info('Closing Hyperswarm connections...');
        await this.swarm.destroy();
        this.peers.forEach((socket, peerId) => {
            socket.end();
            this.logger.info(`Closed connection with peer: ${peerId}`);
        });
        this.peers.clear();
        this.logger.info('Hyperswarm shut down completed.');
    }

    // Internal function to handle new peer connections
    _handleNewPeer(peerId, socket) {
        this.peers.set(peerId, socket);
        socket.on('data', (data) => {
            this.logger.info(`Data received from peer ${peerId}: ${data.toString()}`);
            this.emit('data', peerId, data); // Emit data event
        });

        socket.on('close', () => {
            this.logger.info(`Peer connection closed: ${peerId}`);
            this.peers.delete(peerId);
        });

        socket.on('error', (err) => {
            this.logger.error(`Error with peer ${peerId}: ${err.message}`);
        });

        // Emit a new peer connection event
        this.emit('connection', peerId, socket);
    }

    // Generate a discovery key from a topic using SHA-256
    _generateDiscoveryKey(topic) {
        return crypto.createHash('sha256').update(topic).digest();
    }

    // Send data to a connected peer
    sendDataToPeer(peerId, data) {
        const socket = this.peers.get(peerId);
        if (socket) {
            socket.write(data);
            this.logger.info(`Sent data to peer ${peerId}: ${data.toString()}`);
        } else {
            this.logger.warn(`Peer ${peerId} not found`);
        }
    }

    // Broadcast data to all connected peers
    broadcastData(data) {
        this.peers.forEach((socket, peerId) => {
            socket.write(data);
            this.logger.info(`Broadcasted data to peer ${peerId}: ${data.toString()}`);
        });
    }
}

module.exports = HyperswarmManager;
