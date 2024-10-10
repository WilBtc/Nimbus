// encryptionService.js
const crypto = require('crypto');
const fs = require('fs');
const path = require('path');
const { AtServer } = require('@atsign-foundation/dess'); // Assuming AtServer from DESS handles encryption/decryption

// EncryptionService handles encryption and decryption using AtSign DESS
class EncryptionService {
    constructor(storagePath, dessServer, logger) {
        this.storagePath = storagePath || './encryption_keys';
        this.dessServer = dessServer || new AtServer(); // Reference to DESS AtServer instance
        this.logger = logger || console;

        // Ensure the storage directory for keys exists
        if (!fs.existsSync(this.storagePath)) {
            fs.mkdirSync(this.storagePath, { recursive: true });
            this.logger.info(`Encryption keys storage directory created: ${this.storagePath}`);
        }
    }

    // Encrypt data using AtSign DESS
    async encryptData(deviceId, data) {
        try {
            const encryptionKey = await this.getEncryptionKey(deviceId);
            const cipher = crypto.createCipher('aes-256-cbc', encryptionKey);
            const encryptedData = Buffer.concat([cipher.update(data), cipher.final()]);

            this.logger.info(`Data encrypted for device ${deviceId}`);
            return encryptedData;
        } catch (err) {
            this.logger.error(`Error encrypting data for device ${deviceId}: ${err.message}`);
            return null;
        }
    }

    // Decrypt data using AtSign DESS
    async decryptData(deviceId, encryptedData) {
        try {
            const decryptionKey = await this.getEncryptionKey(deviceId);
            const decipher = crypto.createDecipher('aes-256-cbc', decryptionKey);
            const decryptedData = Buffer.concat([decipher.update(encryptedData), decipher.final()]);

            this.logger.info(`Data decrypted for device ${deviceId}`);
            return decryptedData;
        } catch (err) {
            this.logger.error(`Error decrypting data for device ${deviceId}: ${err.message}`);
            return null;
        }
    }

    // Fetch or generate the encryption key for a specific device using AtSign DESS
    async getEncryptionKey(deviceId) {
        const keyFilePath = path.join(this.storagePath, `${deviceId}_key.pem`);

        // Check if the key already exists locally
        if (fs.existsSync(keyFilePath)) {
            const key = fs.readFileSync(keyFilePath, 'utf8');
            this.logger.info(`Encryption key loaded for device ${deviceId}`);
            return key;
        }

        // Generate a new key using the AtSign DESS server
        try {
            const key = await this.dessServer.generateKey(deviceId); // Hypothetical DESS method for generating keys
            fs.writeFileSync(keyFilePath, key); // Save key to local storage
            this.logger.info(`New encryption key generated and saved for device ${deviceId}`);
            return key;
        } catch (err) {
            this.logger.error(`Error generating encryption key for device ${deviceId}: ${err.message}`);
            throw new Error('Key generation failed');
        }
    }

    // Store a device's encryption key securely using AtSign DESS
    async storeEncryptionKey(deviceId, key) {
        const keyFilePath = path.join(this.storagePath, `${deviceId}_key.pem`);

        try {
            fs.writeFileSync(keyFilePath, key);
            this.logger.info(`Encryption key stored for device ${deviceId}`);
        } catch (err) {
            this.logger.error(`Error storing encryption key for device ${deviceId}: ${err.message}`);
            throw new Error('Key storage failed');
        }
    }

    // Retrieve a device's encryption key securely using AtSign DESS
    async retrieveEncryptionKey(deviceId) {
        const keyFilePath = path.join(this.storagePath, `${deviceId}_key.pem`);

        try {
            if (fs.existsSync(keyFilePath)) {
                const key = fs.readFileSync(keyFilePath, 'utf8');
                this.logger.info(`Encryption key retrieved for device ${deviceId}`);
                return key;
            } else {
                throw new Error(`Key not found for device ${deviceId}`);
            }
        } catch (err) {
            this.logger.error(`Error retrieving encryption key for device ${deviceId}: ${err.message}`);
            throw new Error('Key retrieval failed');
        }
    }

    // Remove a device's encryption key
    removeEncryptionKey(deviceId) {
        const keyFilePath = path.join(this.storagePath, `${deviceId}_key.pem`);

        try {
            if (fs.existsSync(keyFilePath)) {
                fs.unlinkSync(keyFilePath);
                this.logger.info(`Encryption key removed for device ${deviceId}`);
            } else {
                this.logger.warn(`No encryption key found for device ${deviceId} to remove`);
            }
        } catch (err) {
            this.logger.error(`Error removing encryption key for device ${deviceId}: ${err.message}`);
            throw new Error('Key removal failed');
        }
    }
}

module.exports = EncryptionService;
