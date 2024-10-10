// autobase.js
const Autobase = require('autobase');
const Hypercore = require('hypercore');
const crypto = require('crypto');
const path = require('path');
const fs = require('fs');
const EventEmitter = require('events');

// AutobaseManager handles the setup of a multi-writer database for collaborative data sharing
class AutobaseManager extends EventEmitter {
    constructor(storagePath, logger) {
        super();
        this.storagePath = storagePath || './autobase_storage';
        this.logger = logger || console;
        this.inputs = []; // List of input feeds for the multi-writer setup
        this.outputs = []; // List of output feeds
    }

    // Initialize a new writer feed and add it to the Autobase setup
    initializeWriter(writerKey) {
        const writerPath = path.join(this.storagePath, 'writer_' + (writerKey || this._generateRandomKey()));
        const writerFeed = new Hypercore(writerPath, writerKey ? Buffer.from(writerKey, 'hex') : null, {
            valueEncoding: 'json',
        });

        writerFeed.ready(() => {
            const writerKeyHex = writerFeed.key.toString('hex');
            this.logger.info(`Writer feed initialized with key: ${writerKeyHex}`);
            this.inputs.push(writerFeed); // Add the writer feed to the inputs
            this.emit('writerReady', writerKeyHex);
        });

        writerFeed.on('error', (err) => {
            this.logger.error(`Error in writer feed: ${err.message}`);
            this.emit('error', err);
        });

        return writerFeed;
    }

    // Initialize an Autobase with the provided input and output feeds
    initializeAutobase(outputPath) {
        if (!fs.existsSync(outputPath)) {
            fs.mkdirSync(outputPath, { recursive: true });
        }

        const outputFeed = new Hypercore(outputPath, { valueEncoding: 'json' });
        outputFeed.ready(() => {
            const autobase = new Autobase({
                inputs: this.inputs,
                localInput: this.inputs[this.inputs.length - 1], // Set the last writer as the local input
                outputs: [outputFeed],
            });

            this.logger.info('Autobase initialized with multiple writers.');
            this.outputs.push(outputFeed);
            this.emit('autobaseReady', autobase);
        });

        outputFeed.on('error', (err) => {
            this.logger.error(`Error in output feed: ${err.message}`);
            this.emit('error', err);
        });
    }

    // Append data to the local writer feed
    appendToWriter(writerKey, data) {
        const writerFeed = this.inputs.find((feed) => feed.key.toString('hex') === writerKey);
        if (!writerFeed) {
            this.logger.error(`Writer feed with key ${writerKey} not found.`);
            return;
        }

        writerFeed.append(data, (err) => {
            if (err) {
                this.logger.error(`Error appending data to writer feed ${writerKey}: ${err.message}`);
                this.emit('error', err);
                return;
            }

            this.logger.info(`Data appended to writer feed ${writerKey}`);
            this.emit('dataAppended', writerKey, data);
        });
    }

    // Retrieve merged data from the Autobase output
    retrieveMergedData(outputIndex, callback) {
        if (this.outputs.length === 0) {
            this.logger.error('No output feeds found for Autobase.');
            return;
        }

        const outputFeed = this.outputs[0]; // Assume single output feed
        outputFeed.get(outputIndex, (err, data) => {
            if (err) {
                this.logger.error(`Error retrieving data from Autobase at index ${outputIndex}: ${err.message}`);
                this.emit('error', err);
                return;
            }

            this.logger.info(`Merged data retrieved from Autobase at index ${outputIndex}`);
            this.emit('dataRetrieved', outputIndex, data);
            callback(data);
        });
    }

    // Verify integrity of all writer feeds by comparing heads
    verifyWriters() {
        this.inputs.forEach((feed) => {
            feed.head((err, head) => {
                if (err) {
                    this.logger.error(`Error verifying writer feed ${feed.key.toString('hex')}: ${err.message}`);
                    this.emit('error', err);
                    return;
                }

                this.logger.info(`Writer feed ${feed.key.toString('hex')} head: ${JSON.stringify(head)}`);
                this.emit('writerVerified', feed.key.toString('hex'), head);
            });
        });
    }

    // Generate a random 32-byte key for new writer feeds
    _generateRandomKey() {
        return crypto.randomBytes(32).toString('hex');
    }

    // Close all writer and output feeds
    closeFeeds() {
        this.inputs.forEach((feed) => {
            feed.close((err) => {
                if (err) {
                    this.logger.error(`Error closing writer feed ${feed.key.toString('hex')}: ${err.message}`);
                    this.emit('error', err);
                    return;
                }

                this.logger.info(`Writer feed ${feed.key.toString('hex')} closed successfully.`);
            });
        });

        this.outputs.forEach((feed) => {
            feed.close((err) => {
                if (err) {
                    this.logger.error(`Error closing output feed: ${err.message}`);
                    this.emit('error', err);
                    return;
                }

                this.logger.info('Output feed closed successfully.');
            });
        });

        this.logger.info('All feeds closed successfully.');
    }
}

module.exports = AutobaseManager;
