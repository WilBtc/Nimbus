// deviceRoutes.js
const express = require('express');
const router = express.Router();
const DeviceManager = require('../api/deviceManager'); // Assuming DeviceManager handles device connections
const logger = require('../utils/logger').getInstance(); // Assuming logger is initialized elsewhere

// Initialize the device manager
const deviceManager = new DeviceManager(logger);

// Route to connect a device
router.post('/connect', (req, res) => {
    const { deviceId } = req.body;

    if (!deviceId) {
        logger.warn('Device ID missing in connect request.');
        return res.status(400).json({ success: false, message: 'Device ID is required' });
    }

    // Attempt to connect the device
    try {
        const result = deviceManager.connectDevice(deviceId);
        logger.info(`Device ${deviceId} connected successfully.`);
        res.status(200).json({ success: true, message: `Device ${deviceId} connected`, result });
    } catch (err) {
        logger.error(`Error connecting device ${deviceId}: ${err.message}`);
        res.status(500).json({ success: false, message: err.message });
    }
});

// Route to disconnect a device
router.post('/disconnect', (req, res) => {
    const { deviceId } = req.body;

    if (!deviceId) {
        logger.warn('Device ID missing in disconnect request.');
        return res.status(400).json({ success: false, message: 'Device ID is required' });
    }

    // Attempt to disconnect the device
    try {
        const result = deviceManager.disconnectDevice(deviceId);
        logger.info(`Device ${deviceId} disconnected successfully.`);
        res.status(200).json({ success: true, message: `Device ${deviceId} disconnected`, result });
    } catch (err) {
        logger.error(`Error disconnecting device ${deviceId}: ${err.message}`);
        res.status(500).json({ success: false, message: err.message });
    }
});

// Route to list all connected devices
router.get('/list', (req, res) => {
    try {
        const devices = deviceManager.listConnectedDevices();
        logger.info('Retrieved list of connected devices.');
        res.status(200).json({ success: true, devices });
    } catch (err) {
        logger.error(`Error listing devices: ${err.message}`);
        res.status(500).json({ success: false, message: err.message });
    }
});

// Route to check the status of a specific device
router.get('/status/:deviceId', (req, res) => {
    const { deviceId } = req.params;

    try {
        const status = deviceManager.getDeviceStatus(deviceId);
        logger.info(`Device ${deviceId} status checked: ${status}`);
        res.status(200).json({ success: true, status });
    } catch (err) {
        logger.error(`Error checking status for device ${deviceId}: ${err.message}`);
        res.status(500).json({ success: false, message: err.message });
    }
});

module.exports = router;
