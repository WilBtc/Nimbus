// deviceManager.js
const EventEmitter = require('events');
const net = require('net');
const crypto = require('crypto');
const fs = require('fs');
const path = require('path');

// DeviceManager handles device connections, interactions, and secure communication
class DeviceManager extends EventEmitter {
    constructor(storagePath, logger) {
        super();
        this.devices = new Map(); // Stores device connections
        this.storagePath = storagePath || './device_storage'; // Directory path for device data
        this.logger = logger || console;
    }

    // Add a new device to the device manager
    addDevice(deviceId, connection) {
        if (this.devices.has(deviceId)) {
            this.logger.warn(`Device with ID ${deviceId} is already connected.`);
            return;
        }

        this.devices.set(deviceId, connection);
        this.logger.info(`Device ${deviceId} connected.`);
        this.emit('deviceConnected', deviceId);

        connection.on('data', (data) => this.handleDeviceData(deviceId, data));
        connection.on('close', () => this.removeDevice(deviceId));
        connection.on('error', (err) => this.logger.error(`Error on device ${deviceId}: ${err.message}`));
    }

    // Remove a device from the device manager
    removeDevice(deviceId) {
        if (!this.devices.has(deviceId)) {
            this.logger.warn(`Device ${deviceId} is not connected.`);
            return;
        }

        this.devices.get(deviceId).destroy(); // Ensure the connection is closed
        this.devices.delete(deviceId);
        this.logger.info(`Device ${deviceId} disconnected.`);
        this.emit('deviceDisconnected', deviceId);
    }

    // Handle incoming data from a device
    handleDeviceData(deviceId, data) {
        this.logger.info(`Received data from device ${deviceId}: ${data.toString()}`);
        this.emit('dataReceived', deviceId, data);

        // Process the data (decrypt, parse, etc.)
        const decryptedData = this.decryptData(deviceId, data);
        if (decryptedData) {
            this.logger.info(`Decrypted data from device ${deviceId}: ${decryptedData.toString()}`);
        }
    }

    // Send data to a device
    sendData(deviceId, data) {
        const deviceConnection = this.devices.get(deviceId);
        if (!deviceConnection) {
            this.logger.error(`Failed to send data: Device ${deviceId} not connected.`);
            return;
        }

        const encryptedData = this.encryptData(deviceId, data);
        deviceConnection.write(encryptedData, (err) => {
            if (err) {
                this.logger.error(`Error sending data to device ${deviceId}: ${err.message}`);
                this.emit('error', err);
            } else {
                this.logger.info(`Data sent to device ${deviceId}.`);
                this.emit('dataSent', deviceId, encryptedData);
            }
        });
    }

    // Encrypt data for secure transmission
    encryptData(deviceId, data) {
        try {
            const publicKeyPath = path.join(this.storagePath, `${deviceId}_public.pem`);
            if (!fs.existsSync(publicKeyPath)) {
                this.logger.error(`Public key for device ${deviceId} not found.`);
                return null;
            }

            const publicKey = fs.readFileSync(publicKeyPath, 'utf8');
            const encryptedData = crypto.publicEncrypt(publicKey, Buffer.from(data));
            return encryptedData;
        } catch (err) {
            this.logger.error(`Error encrypting data for device ${deviceId}: ${err.message}`);
            return null;
        }
    }

    // Decrypt incoming data
    decryptData(deviceId, encryptedData) {
        try {
            const privateKeyPath = path.join(this.storagePath, `${deviceId}_private.pem`);
            if (!fs.existsSync(privateKeyPath)) {
                this.logger.error(`Private key for device ${deviceId} not found.`);
                return null;
            }

            const privateKey = fs.readFileSync(privateKeyPath, 'utf8');
            const decryptedData = crypto.privateDecrypt(privateKey, Buffer.from(encryptedData));
            return decryptedData;
        } catch (err) {
            this.logger.error(`Error decrypting data for device ${deviceId}: ${err.message}`);
            return null;
        }
    }

    // Load a device's public/private key pair from storage
    loadDeviceKeys(deviceId) {
        const publicKeyPath = path.join(this.storagePath, `${deviceId}_public.pem`);
        const privateKeyPath = path.join(this.storagePath, `${deviceId}_private.pem`);

        if (fs.existsSync(publicKeyPath) && fs.existsSync(privateKeyPath)) {
            const publicKey = fs.readFileSync(publicKeyPath, 'utf8');
            const privateKey = fs.readFileSync(privateKeyPath, 'utf8');
            this.logger.info(`Keys loaded for device ${deviceId}`);
            return { publicKey, privateKey };
        } else {
            this.logger.error(`Keys not found for device ${deviceId}`);
            return null;
        }
    }

    // Generate a new public/private key pair for a device
    generateDeviceKeys(deviceId) {
        const { publicKey, privateKey } = crypto.generateKeyPairSync('rsa', {
            modulusLength: 2048,
            publicKeyEncoding: { type: 'spki', format: 'pem' },
            privateKeyEncoding: { type: 'pkcs8', format: 'pem' },
        });

        const publicKeyPath = path.join(this.storagePath, `${deviceId}_public.pem`);
        const privateKeyPath = path.join(this.storagePath, `${deviceId}_private.pem`);

        fs.writeFileSync(publicKeyPath, publicKey);
        fs.writeFileSync(privateKeyPath, privateKey);
        this.logger.info(`New keys generated for device ${deviceId}`);

        return { publicKey, privateKey };
    }

    // List all connected devices
    listDevices() {
        const deviceIds = Array.from(this.devices.keys());
        this.logger.info(`Connected devices: ${deviceIds.join(', ')}`);
        return deviceIds;
    }
}

module.exports = DeviceManager;
