// fileManager.js
document.addEventListener('DOMContentLoaded', function () {
    const fileListSection = document.getElementById('file-list');
    const uploadForm = document.getElementById('upload-form');
    
    // Function to fetch the list of files
    async function fetchFiles() {
      try {
        const response = await fetch('/api/files');
        if (!response.ok) {
          throw new Error('Failed to fetch files');
        }
  
        const data = await response.json();
        renderFileList(data.files);
      } catch (error) {
        console.error(error);
      }
    }
  
    // Function to render the file list
    function renderFileList(files) {
      fileListSection.innerHTML = ''; // Clear previous content
      files.forEach(file => {
        const fileItem = document.createElement('div');
        fileItem.classList.add('file-item');
        fileItem.innerHTML = `
          <strong>File:</strong> ${file.fileName} <br>
          <strong>Author:</strong> ${file.author} <br>
          <strong>Uploaded:</strong> ${file.uploadDate} <br>
          <button onclick="downloadFile('${file.fileName}')">Download</button>
          <button onclick="deleteFile('${file.fileName}', '${file.author}')">Delete</button>
        `;
        fileListSection.appendChild(fileItem);
      });
    }
  
    // Function to upload a file
    uploadForm.addEventListener('submit', async function (event) {
      event.preventDefault();
  
      const fileInput = document.getElementById('fileInput');
      const fileAuthor = document.getElementById('fileAuthor').value;
  
      if (!fileInput.files.length) {
        return alert('Please select a file to upload.');
      }
  
      const formData = new FormData();
      formData.append('file', fileInput.files[0]);
      formData.append('author', fileAuthor);
  
      try {
        const response = await fetch('/api/files/upload', {
          method: 'POST',
          body: formData,
        });
  
        if (!response.ok) {
          throw new Error('Failed to upload file');
        }
  
        alert('File uploaded successfully!');
        fetchFiles(); // Refresh the file list after upload
      } catch (error) {
        console.error(error);
      }
    });
  
    // Function to download a file
    window.downloadFile = async function (fileName) {
      try {
        const response = await fetch(`/api/files/download/${fileName}`);
        if (!response.ok) {
          throw new Error('Failed to download file');
        }
  
        const blob = await response.blob();
        const link = document.createElement('a');
        link.href = window.URL.createObjectURL(blob);
        link.download = fileName;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      } catch (error) {
        console.error(error);
      }
    };
  
    // Function to delete a file
    window.deleteFile = async function (fileName, author) {
      const confirmDelete = confirm(`Are you sure you want to delete the file "${fileName}"?`);
      if (!confirmDelete) return;
  
      try {
        const response = await fetch('/api/files/delete', {
          method: 'DELETE',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ fileName, author }),
        });
  
        if (!response.ok) {
          throw new Error('Failed to delete file');
        }
  
        alert('File deleted successfully!');
        fetchFiles(); // Refresh the file list after deletion
      } catch (error) {
        console.error(error);
      }
    };
  
    // Initial load: Fetch the list of files
    fetchFiles();
  });
  