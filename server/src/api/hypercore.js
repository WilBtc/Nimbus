// hypercore.js
const Hypercore = require('hypercore');
const crypto = require('crypto');
const fs = require('fs');
const path = require('path');
const EventEmitter = require('events');

// HypercoreManager manages file chunking, verification, and distribution using Hypercore
class HypercoreManager extends EventEmitter {
    constructor(storagePath, logger) {
        super();
        this.storagePath = storagePath || './hypercore_storage';
        this.logger = logger || console;
        this.feeds = new Map(); // Store active Hypercore feeds
    }

    // Initialize a Hypercore feed for a given key or create a new feed
    initializeFeed(feedKey) {
        const key = feedKey ? Buffer.from(feedKey, 'hex') : null;
        const feed = new Hypercore(path.join(this.storagePath, key ? key.toString('hex') : ''), key, {
            valueEncoding: 'binary',
        });

        feed.ready(() => {
            const feedKeyHex = feed.key.toString('hex');
            this.logger.info(`Hypercore feed initialized with key: ${feedKeyHex}`);
            this.feeds.set(feedKeyHex, feed);
            this.emit('ready', feedKeyHex);
        });

        // Handle feed errors
        feed.on('error', (err) => {
            this.logger.error(`Error in Hypercore feed: ${err.message}`);
            this.emit('error', err);
        });

        return feed;
    }

    // Append data to the specified feed
    appendToFeed(feedKey, data) {
        const feed = this.feeds.get(feedKey);
        if (!feed) {
            this.logger.error(`Feed with key ${feedKey} not found`);
            return;
        }

        feed.append(data, (err) => {
            if (err) {
                this.logger.error(`Error appending data to feed ${feedKey}: ${err.message}`);
                this.emit('error', err);
                return;
            }

            this.logger.info(`Data appended to feed ${feedKey}`);
            this.emit('dataAppended', feedKey, data);
        });
    }

    // Retrieve data from the feed by index
    getFeedData(feedKey, index) {
        const feed = this.feeds.get(feedKey);
        if (!feed) {
            this.logger.error(`Feed with key ${feedKey} not found`);
            return;
        }

        feed.get(index, (err, data) => {
            if (err) {
                this.logger.error(`Error retrieving data from feed ${feedKey} at index ${index}: ${err.message}`);
                this.emit('error', err);
                return;
            }

            this.logger.info(`Data retrieved from feed ${feedKey} at index ${index}`);
            this.emit('dataRetrieved', feedKey, index, data);
        });
    }

    // Chunk and distribute a file by appending its chunks to the feed
    distributeFile(feedKey, filePath) {
        const feed = this.feeds.get(feedKey);
        if (!feed) {
            this.logger.error(`Feed with key ${feedKey} not found`);
            return;
        }

        const fileStream = fs.createReadStream(filePath);
        fileStream.on('data', (chunk) => {
            this.appendToFeed(feedKey, chunk);
        });

        fileStream.on('end', () => {
            this.logger.info(`File ${filePath} distributed through feed ${feedKey}`);
            this.emit('fileDistributed', feedKey, filePath);
        });

        fileStream.on('error', (err) => {
            this.logger.error(`Error reading file ${filePath}: ${err.message}`);
            this.emit('error', err);
        });
    }

    // Verify the integrity of a feed by checking its hashes
    verifyFeed(feedKey) {
        const feed = this.feeds.get(feedKey);
        if (!feed) {
            this.logger.error(`Feed with key ${feedKey} not found`);
            return;
        }

        feed.audit((err, report) => {
            if (err) {
                this.logger.error(`Error verifying feed ${feedKey}: ${err.message}`);
                this.emit('error', err);
                return;
            }

            this.logger.info(`Feed verification completed for feed ${feedKey}:`, report);
            this.emit('feedVerified', feedKey, report);
        });
    }

    // Close a feed and release resources
    closeFeed(feedKey) {
        const feed = this.feeds.get(feedKey);
        if (!feed) {
            this.logger.warn(`Feed with key ${feedKey} not found`);
            return;
        }

        feed.close((err) => {
            if (err) {
                this.logger.error(`Error closing feed ${feedKey}: ${err.message}`);
                this.emit('error', err);
                return;
            }

            this.feeds.delete(feedKey);
            this.logger.info(`Feed closed with key ${feedKey}`);
            this.emit('feedClosed', feedKey);
        });
    }
}

module.exports = HypercoreManager;
