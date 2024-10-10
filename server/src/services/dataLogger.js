// dataLogger.js
const fs = require('fs');
const path = require('path');
const { EventEmitter } = require('events');

// DataLogger is responsible for logging real-time interactions and tracking updates.
class DataLogger extends EventEmitter {
    constructor(logFilePath, logRotationInterval, logger) {
        super();
        this.logFilePath = logFilePath || './logs/data.log'; // Default log file path
        this.logRotationInterval = logRotationInterval || 24; // Default rotation interval in hours
        this.logger = logger || console;

        // Ensure the log directory exists
        const logDir = path.dirname(this.logFilePath);
        if (!fs.existsSync(logDir)) {
            fs.mkdirSync(logDir, { recursive: true });
            this.logger.info(`Log directory created: ${logDir}`);
        }

        // Start log rotation
        this.startLogRotation();
    }

    // Log real-time interactions
    logInteraction(deviceId, action, details = '') {
        const timestamp = new Date().toISOString();
        const logEntry = `${timestamp} - Device: ${deviceId} - Action: ${action} - Details: ${details}\n`;

        try {
            fs.appendFileSync(this.logFilePath, logEntry);
            this.logger.info(`Interaction logged for device ${deviceId}: ${action}`);
            this.emit('interactionLogged', deviceId, action, details);
        } catch (err) {
            this.logger.error(`Error logging interaction for device ${deviceId}: ${err.message}`);
            this.emit('error', err);
        }
    }

    // Log system updates or changes in state
    logUpdate(updateType, details) {
        const timestamp = new Date().toISOString();
        const logEntry = `${timestamp} - Update: ${updateType} - Details: ${details}\n`;

        try {
            fs.appendFileSync(this.logFilePath, logEntry);
            this.logger.info(`System update logged: ${updateType}`);
            this.emit('updateLogged', updateType, details);
        } catch (err) {
            this.logger.error(`Error logging system update: ${err.message}`);
            this.emit('error', err);
        }
    }

    // Start log rotation based on the specified interval
    startLogRotation() {
        setInterval(() => {
            const rotatedLogFile = this.getRotatedLogFileName();
            try {
                // Rename the current log file to the rotated version
                fs.renameSync(this.logFilePath, rotatedLogFile);
                this.logger.info(`Log file rotated: ${rotatedLogFile}`);
                this.emit('logRotated', rotatedLogFile);

                // Create a new empty log file
                fs.writeFileSync(this.logFilePath, '');
                this.logger.info('New log file created after rotation.');
            } catch (err) {
                this.logger.error(`Error rotating log file: ${err.message}`);
                this.emit('error', err);
            }
        }, this.logRotationInterval * 60 * 60 * 1000); // Convert hours to milliseconds
    }

    // Generate a filename for the rotated log file
    getRotatedLogFileName() {
        const timestamp = new Date().toISOString().replace(/:/g, '-'); // Replace colons to avoid issues in filenames
        const ext = path.extname(this.logFilePath);
        const baseName = path.basename(this.logFilePath, ext);
        const dirName = path.dirname(this.logFilePath);
        return path.join(dirName, `${baseName}-${timestamp}${ext}`);
    }

    // Retrieve recent log entries (last N lines)
    getRecentLogs(numLines = 10) {
        try {
            const data = fs.readFileSync(this.logFilePath, 'utf-8');
            const lines = data.trim().split('\n');
            const recentLogs = lines.slice(-numLines).join('\n');
            this.logger.info(`Retrieved ${numLines} recent log entries.`);
            return recentLogs;
        } catch (err) {
            this.logger.error(`Error retrieving recent log entries: ${err.message}`);
            this.emit('error', err);
            return null;
        }
    }

    // Clear log file content
    clearLogs() {
        try {
            fs.writeFileSync(this.logFilePath, '');
            this.logger.info('Log file cleared.');
            this.emit('logsCleared');
        } catch (err) {
            this.logger.error(`Error clearing log file: ${err.message}`);
            this.emit('error', err);
        }
    }
}

module.exports = DataLogger;
