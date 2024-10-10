// noPorts.js
const noports = require('noports'); // Assuming this is the correct package for NoPorts
const EventEmitter = require('events');
const fs = require('fs');
const path = require('path');

// NoPortsManager handles the initialization and management of NoPorts tunnels
class NoPortsManager extends EventEmitter {
    constructor(tunnelConfig, storagePath, logger) {
        super();
        this.tunnelConfig = tunnelConfig || {};  // Configuration for NoPorts tunnels
        this.storagePath = storagePath || './noports_storage';
        this.logger = logger || console;
        this.tunnels = new Map();  // Map to store active tunnels
    }

    // Initialize a NoPorts tunnel with the provided configuration
    initializeTunnel(tunnelId, customConfig = {}) {
        const config = { ...this.tunnelConfig, ...customConfig };
        this.logger.info(`Initializing NoPorts tunnel with ID: ${tunnelId}`);

        try {
            // Create and start a NoPorts tunnel
            const tunnel = noports.connect(config);
            tunnel.on('open', () => {
                this.logger.info(`Tunnel ${tunnelId} opened successfully.`);
                this.tunnels.set(tunnelId, tunnel);
                this.emit('tunnelOpened', tunnelId);
            });

            tunnel.on('close', () => {
                this.logger.info(`Tunnel ${tunnelId} closed.`);
                this.tunnels.delete(tunnelId);
                this.emit('tunnelClosed', tunnelId);
            });

            tunnel.on('error', (err) => {
                this.logger.error(`Error in tunnel ${tunnelId}: ${err.message}`);
                this.emit('error', err);
            });

            return tunnel;
        } catch (err) {
            this.logger.error(`Failed to initialize tunnel ${tunnelId}: ${err.message}`);
            this.emit('error', err);
        }
    }

    // Close a specific NoPorts tunnel by its ID
    closeTunnel(tunnelId) {
        const tunnel = this.tunnels.get(tunnelId);
        if (!tunnel) {
            this.logger.error(`Tunnel with ID ${tunnelId} not found.`);
            return;
        }

        tunnel.close((err) => {
            if (err) {
                this.logger.error(`Error closing tunnel ${tunnelId}: ${err.message}`);
                this.emit('error', err);
                return;
            }
            this.logger.info(`Tunnel ${tunnelId} closed successfully.`);
            this.emit('tunnelClosed', tunnelId);
        });
    }

    // Close all active tunnels
    closeAllTunnels() {
        this.tunnels.forEach((tunnel, tunnelId) => {
            this.closeTunnel(tunnelId);
        });
        this.logger.info('All NoPorts tunnels have been closed.');
    }

    // Save the tunnel configuration to a file for persistence
    saveTunnelConfig(tunnelId, config) {
        const configPath = path.join(this.storagePath, `${tunnelId}_config.json`);
        fs.writeFileSync(configPath, JSON.stringify(config, null, 2));
        this.logger.info(`Tunnel configuration saved for ${tunnelId} at ${configPath}`);
    }

    // Load a tunnel configuration from a file
    loadTunnelConfig(tunnelId) {
        const configPath = path.join(this.storagePath, `${tunnelId}_config.json`);
        if (!fs.existsSync(configPath)) {
            this.logger.error(`Tunnel configuration file not found for ${tunnelId}`);
            return null;
        }

        const config = JSON.parse(fs.readFileSync(configPath, 'utf-8'));
        this.logger.info(`Tunnel configuration loaded for ${tunnelId}`);
        return config;
    }

    // List all active tunnels
    listActiveTunnels() {
        const activeTunnels = Array.from(this.tunnels.keys());
        this.logger.info(`Active tunnels: ${activeTunnels.join(', ')}`);
        return activeTunnels;
    }

    // Restart a tunnel by closing and reinitializing it
    restartTunnel(tunnelId) {
        this.logger.info(`Restarting tunnel ${tunnelId}...`);

        const config = this.loadTunnelConfig(tunnelId);
        if (!config) {
            this.logger.error(`Failed to restart tunnel ${tunnelId}: Configuration not found.`);
            return;
        }

        this.closeTunnel(tunnelId);
        setTimeout(() => {
            this.initializeTunnel(tunnelId, config);
        }, 1000);  // Wait for 1 second before restarting
    }
}

module.exports = NoPortsManager;
