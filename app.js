const express = require('express');
const bodyParser = require('body-parser');
const path = require('path');
const hyperswarm = require('./api/hyperswarm'); // Hyperswarm setup
const atsignAuth = require('./api/atsignAuth'); // Atsign authentication
const noPorts = require('./api/noPorts'); // NoPorts tunneling
const hypercore = require('./api/hypercore'); // Hypercore data handling
const autobase = require('./api/autobase'); // Autobase for collaborative data
const deviceRoutes = require('./routes/deviceRoutes'); // Device routes
const fileRoutes = require('./routes/fileRoutes'); // File management routes
const logger = require('./utils/logger'); // Custom logger

const app = express();
const port = process.env.PORT || 3000;

// Middleware
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

// Static files
app.use(express.static(path.join(__dirname, 'public')));

// Routes
app.use('/api/devices', deviceRoutes);
app.use('/api/files', fileRoutes);

// Initialize core services
async function initializeServices() {
  try {
    // Initialize Hyperswarm for peer discovery
    await hyperswarm.setupHyperswarm();

    // Initialize Atsign Authentication
    await atsignAuth.initializeAuth();

    // Initialize NoPorts tunneling
    await noPorts.initializeTunnel();

    // Initialize Hypercore for data storage
    await hypercore.setupHypercore('./storage');

    // Initialize Autobase for collaborative data sharing
    await autobase.setupAutobase();

    logger.info('All services initialized successfully.');
  } catch (error) {
    logger.error('Error initializing services:', error);
    process.exit(1);
  }
}

// Start the server
app.listen(port, () => {
  logger.info(`Nimbus server is running on port ${port}`);
  initializeServices(); // Initialize core services
});

module.exports = app;
