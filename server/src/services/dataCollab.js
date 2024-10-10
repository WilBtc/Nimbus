// dataCollab.js
const fs = require('fs');
const path = require('path');
const { EventEmitter } = require('events');

// DataCollab handles collaborative file access, updates, and sync between multiple users or devices.
class DataCollab extends EventEmitter {
    constructor(collabDirectory, logger) {
        super();
        this.collabDirectory = collabDirectory || './collaborative_files'; // Default directory for shared files
        this.logger = logger || console;

        // Ensure the collaboration directory exists
        if (!fs.existsSync(this.collabDirectory)) {
            fs.mkdirSync(this.collabDirectory, { recursive: true });
            this.logger.info(`Collaboration directory created: ${this.collabDirectory}`);
        }

        // Set to manage locks on files for collaborative access
        this.fileLocks = new Set();
    }

    // Create or update a collaborative file
    createOrUpdateFile(fileName, content, author) {
        const filePath = path.join(this.collabDirectory, fileName);
        const timestamp = new Date().toISOString();
        const updateDetails = `Updated by ${author} at ${timestamp}\n`;

        try {
            // If the file is locked by another user, prevent the update
            if (this.fileLocks.has(fileName)) {
                const errorMsg = `File ${fileName} is currently locked for collaboration. Try again later.`;
                this.logger.warn(errorMsg);
                this.emit('fileLocked', fileName);
                return { success: false, message: errorMsg };
            }

            // Write or append the update details and content to the file
            fs.appendFileSync(filePath, updateDetails + content + '\n');
            this.logger.info(`File ${fileName} updated by ${author}`);
            this.emit('fileUpdated', fileName, author, timestamp);
            return { success: true, message: `File ${fileName} updated successfully by ${author}` };
        } catch (err) {
            this.logger.error(`Error updating file ${fileName}: ${err.message}`);
            this.emit('error', err);
            return { success: false, message: err.message };
        }
    }

    // Read the contents of a collaborative file
    readFile(fileName) {
        const filePath = path.join(this.collabDirectory, fileName);
        try {
            if (!fs.existsSync(filePath)) {
                const errorMsg = `File ${fileName} not found.`;
                this.logger.warn(errorMsg);
                this.emit('fileNotFound', fileName);
                return { success: false, message: errorMsg };
            }

            const content = fs.readFileSync(filePath, 'utf-8');
            this.logger.info(`File ${fileName} read successfully.`);
            this.emit('fileRead', fileName);
            return { success: true, content };
        } catch (err) {
            this.logger.error(`Error reading file ${fileName}: ${err.message}`);
            this.emit('error', err);
            return { success: false, message: err.message };
        }
    }

    // Lock a file to prevent concurrent updates from other users
    lockFile(fileName) {
        if (this.fileLocks.has(fileName)) {
            const errorMsg = `File ${fileName} is already locked.`;
            this.logger.warn(errorMsg);
            this.emit('fileLocked', fileName);
            return { success: false, message: errorMsg };
        }

        this.fileLocks.add(fileName);
        this.logger.info(`File ${fileName} locked for collaboration.`);
        this.emit('fileLocked', fileName);
        return { success: true, message: `File ${fileName} locked.` };
    }

    // Unlock a file to allow updates from other users
    unlockFile(fileName) {
        if (!this.fileLocks.has(fileName)) {
            const errorMsg = `File ${fileName} is not locked.`;
            this.logger.warn(errorMsg);
            this.emit('fileNotLocked', fileName);
            return { success: false, message: errorMsg };
        }

        this.fileLocks.delete(fileName);
        this.logger.info(`File ${fileName} unlocked for collaboration.`);
        this.emit('fileUnlocked', fileName);
        return { success: true, message: `File ${fileName} unlocked.` };
    }

    // Get a list of collaborative files
    listCollaborativeFiles() {
        try {
            const files = fs.readdirSync(this.collabDirectory);
            this.logger.info(`Collaborative files listed successfully.`);
            return { success: true, files };
        } catch (err) {
            this.logger.error(`Error listing collaborative files: ${err.message}`);
            this.emit('error', err);
            return { success: false, message: err.message };
        }
    }

    // Delete a collaborative file
    deleteFile(fileName, author) {
        const filePath = path.join(this.collabDirectory, fileName);

        try {
            if (!fs.existsSync(filePath)) {
                const errorMsg = `File ${fileName} not found.`;
                this.logger.warn(errorMsg);
                this.emit('fileNotFound', fileName);
                return { success: false, message: errorMsg };
            }

            fs.unlinkSync(filePath);
            this.logger.info(`File ${fileName} deleted by ${author}.`);
            this.emit('fileDeleted', fileName, author);
            return { success: true, message: `File ${fileName} deleted by ${author}` };
        } catch (err) {
            this.logger.error(`Error deleting file ${fileName}: ${err.message}`);
            this.emit('error', err);
            return { success: false, message: err.message };
        }
    }
}

module.exports = DataCollab;
