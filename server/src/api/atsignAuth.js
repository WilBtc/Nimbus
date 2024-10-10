// atsignAuth.js
const AtClient = require('@atsign/at_client'); // Assuming this is the correct package for AtClient SDK
const crypto = require('crypto');
const fs = require('fs');
const path = require('path');
const EventEmitter = require('events');

// AtSignAuthManager handles authentication and encryption logic for AtSign communication
class AtSignAuthManager extends EventEmitter {
    constructor(atSign, storagePath, logger) {
        super();
        this.atSign = atSign;
        this.storagePath = storagePath || './atsign_storage';
        this.logger = logger || console;
        this.client = null; // AtClient instance
    }

    // Initialize AtSign authentication with provided credentials or load from storage
    async initialize(atsignConfig) {
        try {
            this.logger.info(`Initializing AtSign authentication for: ${this.atSign}`);

            // Load or authenticate using AtSign credentials
            this.client = await this._authenticateAtSign(atsignConfig);
            this.logger.info(`Successfully authenticated AtSign: ${this.atSign}`);
            this.emit('authenticated', this.atSign);
        } catch (err) {
            this.logger.error(`Failed to authenticate AtSign: ${this.atSign}, Error: ${err.message}`);
            this.emit('error', err);
        }
    }

    // Private method to authenticate AtSign and create AtClient instance
    async _authenticateAtSign(atsignConfig) {
        // Load credentials from config or storage
        const keysFilePath = path.join(this.storagePath, `${this.atSign}_keys.json`);
        let keys = null;

        if (fs.existsSync(keysFilePath)) {
            this.logger.info(`Loading AtSign keys from storage for ${this.atSign}`);
            keys = JSON.parse(fs.readFileSync(keysFilePath));
        } else {
            this.logger.info(`Using provided credentials for AtSign: ${this.atSign}`);
            keys = {
                cramSecret: atsignConfig.cramSecret,
                encryptionPrivateKey: atsignConfig.encryptionPrivateKey,
                encryptionPublicKey: atsignConfig.encryptionPublicKey,
            };

            // Save keys to storage for future use
            fs.writeFileSync(keysFilePath, JSON.stringify(keys));
        }

        // Authenticate the AtSign with AtClient
        const atClient = new AtClient({
            atSign: this.atSign,
            rootDomain: atsignConfig.rootDomain,
            keys: keys,
        });

        await atClient.authenticate();
        return atClient;
    }

    // Encrypt data using the AtSign's encryption keys
    encryptData(data) {
        if (!this.client || !this.client.keys.encryptionPublicKey) {
            this.logger.error('Encryption failed: No valid public key or client not authenticated.');
            return null;
        }

        const publicKey = Buffer.from(this.client.keys.encryptionPublicKey, 'base64');
        const encryptedData = crypto.publicEncrypt(publicKey, Buffer.from(data));
        this.logger.info(`Data encrypted using AtSign public key: ${this.atSign}`);
        return encryptedData.toString('base64');
    }

    // Decrypt data using the AtSign's private key
    decryptData(encryptedData) {
        if (!this.client || !this.client.keys.encryptionPrivateKey) {
            this.logger.error('Decryption failed: No valid private key or client not authenticated.');
            return null;
        }

        const privateKey = Buffer.from(this.client.keys.encryptionPrivateKey, 'base64');
        const decryptedData = crypto.privateDecrypt(privateKey, Buffer.from(encryptedData, 'base64'));
        this.logger.info(`Data decrypted using AtSign private key: ${this.atSign}`);
        return decryptedData.toString();
    }

    // Sign data using the AtSign's private key (for message integrity)
    signData(data) {
        if (!this.client || !this.client.keys.encryptionPrivateKey) {
            this.logger.error('Signing failed: No valid private key or client not authenticated.');
            return null;
        }

        const privateKey = Buffer.from(this.client.keys.encryptionPrivateKey, 'base64');
        const sign = crypto.createSign('SHA256');
        sign.update(data);
        sign.end();
        const signature = sign.sign(privateKey, 'base64');
        this.logger.info(`Data signed using AtSign private key: ${this.atSign}`);
        return signature;
    }

    // Verify the signature of data using the AtSign's public key
    verifySignature(data, signature) {
        if (!this.client || !this.client.keys.encryptionPublicKey) {
            this.logger.error('Signature verification failed: No valid public key or client not authenticated.');
            return false;
        }

        const publicKey = Buffer.from(this.client.keys.encryptionPublicKey, 'base64');
        const verify = crypto.createVerify('SHA256');
        verify.update(data);
        verify.end();
        const isValid = verify.verify(publicKey, signature, 'base64');
        this.logger.info(`Signature verification result for AtSign ${this.atSign}: ${isValid}`);
        return isValid;
    }

    // Sign out the AtSign and remove keys from storage
    signOut() {
        const keysFilePath = path.join(this.storagePath, `${this.atSign}_keys.json`);
        if (fs.existsSync(keysFilePath)) {
            fs.unlinkSync(keysFilePath);
            this.logger.info(`AtSign keys removed from storage for ${this.atSign}`);
        }
        this.client = null;
        this.logger.info(`AtSign signed out: ${this.atSign}`);
        this.emit('signedOut', this.atSign);
    }
}

module.exports = AtSignAuthManager;
