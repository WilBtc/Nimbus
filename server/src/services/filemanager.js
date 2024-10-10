// fileManager.js
const fs = require('fs');
const path = require('path');
const EventEmitter = require('events');
const crypto = require('crypto');

// FileManager handles file uploads, downloads, and data sharing.
class FileManager extends EventEmitter {
    constructor(storageDir, logger) {
        super();
        this.storageDir = storageDir || './file_storage'; // Default directory for storing files
        this.logger = logger || console;

        // Ensure the storage directory exists
        if (!fs.existsSync(this.storageDir)) {
            fs.mkdirSync(this.storageDir, { recursive: true });
            this.logger.info(`Storage directory created: ${this.storageDir}`);
        }
    }

    // Upload a file and store it securely
    uploadFile(fileBuffer, fileName, deviceId) {
        const filePath = path.join(this.storageDir, fileName);
        try {
            fs.writeFileSync(filePath, fileBuffer);
            this.logger.info(`File uploaded: ${fileName} from device ${deviceId}`);
            this.emit('fileUploaded', fileName, deviceId);

            // Optionally, encrypt the file for secure storage
            this.encryptFile(filePath, deviceId);
        } catch (err) {
            this.logger.error(`Error uploading file ${fileName}: ${err.message}`);
            this.emit('error', err);
        }
    }

    // Download a file securely
    downloadFile(fileName, deviceId) {
        const filePath = path.join(this.storageDir, fileName);
        try {
            if (!fs.existsSync(filePath)) {
                throw new Error(`File ${fileName} not found`);
            }

            // Decrypt the file before sending it for download
            this.decryptFile(filePath, deviceId);
            const fileBuffer = fs.readFileSync(filePath);
            this.logger.info(`File downloaded: ${fileName} by device ${deviceId}`);
            this.emit('fileDownloaded', fileName, deviceId);
            return fileBuffer;
        } catch (err) {
            this.logger.error(`Error downloading file ${fileName}: ${err.message}`);
            this.emit('error', err);
            return null;
        }
    }

    // Delete a file securely
    deleteFile(fileName, deviceId) {
        const filePath = path.join(this.storageDir, fileName);
        try {
            if (!fs.existsSync(filePath)) {
                throw new Error(`File ${fileName} not found`);
            }

            fs.unlinkSync(filePath);
            this.logger.info(`File deleted: ${fileName} by device ${deviceId}`);
            this.emit('fileDeleted', fileName, deviceId);
        } catch (err) {
            this.logger.error(`Error deleting file ${fileName}: ${err.message}`);
            this.emit('error', err);
        }
    }

    // Encrypt a file before saving it to the storage
    encryptFile(filePath, deviceId) {
        try {
            const cipher = crypto.createCipher('aes-256-cbc', this.getEncryptionKey(deviceId));
            const input = fs.createReadStream(filePath);
            const output = fs.createWriteStream(filePath + '.enc');

            input.pipe(cipher).pipe(output);

            output.on('finish', () => {
                fs.unlinkSync(filePath); // Remove the original unencrypted file
                fs.renameSync(filePath + '.enc', filePath); // Rename encrypted file to original filename
                this.logger.info(`File encrypted: ${filePath} for device ${deviceId}`);
                this.emit('fileEncrypted', filePath, deviceId);
            });
        } catch (err) {
            this.logger.error(`Error encrypting file ${filePath}: ${err.message}`);
            this.emit('error', err);
        }
    }

    // Decrypt a file before sending it to a device
    decryptFile(filePath, deviceId) {
        try {
            const decipher = crypto.createDecipher('aes-256-cbc', this.getEncryptionKey(deviceId));
            const input = fs.createReadStream(filePath);
            const output = fs.createWriteStream(filePath + '.dec');

            input.pipe(decipher).pipe(output);

            output.on('finish', () => {
                fs.unlinkSync(filePath); // Remove the original encrypted file
                fs.renameSync(filePath + '.dec', filePath); // Rename decrypted file to original filename
                this.logger.info(`File decrypted: ${filePath} for device ${deviceId}`);
                this.emit('fileDecrypted', filePath, deviceId);
            });
        } catch (err) {
            this.logger.error(`Error decrypting file ${filePath}: ${err.message}`);
            this.emit('error', err);
        }
    }

    // Get encryption key for a specific device (placeholder implementation)
    getEncryptionKey(deviceId) {
        // Placeholder key generation logic; in production, you would securely generate/store keys
        return crypto.createHash('sha256').update(deviceId).digest('base64').substr(0, 32); // Generate a 32-byte key from device ID
    }

    // List all files stored in the file storage directory
    listFiles() {
        try {
            const files = fs.readdirSync(this.storageDir);
            this.logger.info(`Files in storage: ${files.join(', ')}`);
            return files;
        } catch (err) {
            this.logger.error(`Error listing files: ${err.message}`);
            this.emit('error', err);
            return [];
        }
    }

    // Share a file with another device securely
    shareFile(fileName, fromDeviceId, toDeviceId) {
        try {
            const filePath = path.join(this.storageDir, fileName);

            // Ensure file exists before sharing
            if (!fs.existsSync(filePath)) {
                throw new Error(`File ${fileName} not found for sharing`);
            }

            // Optionally, re-encrypt the file for the receiving device
            const fileBuffer = fs.readFileSync(filePath);
            const encryptedBuffer = this.encryptBuffer(fileBuffer, toDeviceId);

            this.logger.info(`File ${fileName} shared from device ${fromDeviceId} to ${toDeviceId}`);
            this.emit('fileShared', fileName, fromDeviceId, toDeviceId);
            return encryptedBuffer; // Return the encrypted buffer to be sent to the receiving device
        } catch (err) {
            this.logger.error(`Error sharing file ${fileName}: ${err.message}`);
            this.emit('error', err);
            return null;
        }
    }

    // Encrypt a buffer of data for a specific device
    encryptBuffer(buffer, deviceId) {
        try {
            const cipher = crypto.createCipher('aes-256-cbc', this.getEncryptionKey(deviceId));
            return Buffer.concat([cipher.update(buffer), cipher.final()]);
        } catch (err) {
            this.logger.error(`Error encrypting buffer for device ${deviceId}: ${err.message}`);
            return null;
        }
    }
}

module.exports = FileManager;
