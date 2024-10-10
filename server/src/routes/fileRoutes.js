// fileRoutes.js
const express = require('express');
const router = express.Router();
const DataCollab = require('../api/dataCollab'); // Assuming DataCollab handles file sharing
const logger = require('../utils/logger').getInstance(); // Assuming logger is initialized elsewhere

// Initialize the data collaboration manager
const dataCollab = new DataCollab('./shared_files', logger);

// Route to upload or update a collaborative file
router.post('/upload', (req, res) => {
    const { fileName, content, author } = req.body;

    if (!fileName || !content || !author) {
        logger.warn('File upload request missing required fields.');
        return res.status(400).json({ success: false, message: 'FileName, content, and author are required' });
    }

    try {
        const result = dataCollab.createOrUpdateFile(fileName, content, author);
        res.status(200).json(result);
    } catch (err) {
        logger.error(`Error uploading file ${fileName}: ${err.message}`);
        res.status(500).json({ success: false, message: err.message });
    }
});

// Route to download a collaborative file
router.get('/download/:fileName', (req, res) => {
    const { fileName } = req.params;

    try {
        const result = dataCollab.readFile(fileName);
        if (result.success) {
            res.status(200).json(result);
        } else {
            res.status(404).json(result);
        }
    } catch (err) {
        logger.error(`Error downloading file ${fileName}: ${err.message}`);
        res.status(500).json({ success: false, message: err.message });
    }
});

// Route to list all collaborative files
router.get('/files', (req, res) => {
    try {
        const result = dataCollab.listCollaborativeFiles();
        res.status(200).json(result);
    } catch (err) {
        logger.error(`Error listing collaborative files: ${err.message}`);
        res.status(500).json({ success: false, message: err.message });
    }
});

// Route to delete a file
router.delete('/delete', (req, res) => {
    const { fileName, author } = req.body;

    if (!fileName || !author) {
        logger.warn('File delete request missing required fields.');
        return res.status(400).json({ success: false, message: 'FileName and author are required' });
    }

    try {
        const result = dataCollab.deleteFile(fileName, author);
        res.status(200).json(result);
    } catch (err) {
        logger.error(`Error deleting file ${fileName}: ${err.message}`);
        res.status(500).json({ success: false, message: err.message });
    }
});

module.exports = router;
