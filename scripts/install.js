#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const https = require('https');
const { execSync } = require('child_process');
const tar = require('tar');
const { createWriteStream } = require('fs');
const { promisify } = require('util');
const pipeline = promisify(require('stream').pipeline);
const fetch = require('node-fetch');
const crypto = require('crypto');

// Get the latest release version from GitHub API or use a specific version
async function getVersion() {
  try {
    // If a specific version is provided via environment variable, use that
    if (process.env.CMDEAGLE_VERSION) {
      return process.env.CMDEAGLE_VERSION;
    }
    
    // Otherwise, fetch the latest release from GitHub
    const response = await fetch('https://api.github.com/repos/migsc/cmdeagle/releases/latest');
    if (!response.ok) {
      throw new Error(`Failed to fetch latest release: ${response.statusText}`);
    }
    
    const data = await response.json();
    return data.tag_name.replace(/^v/, ''); // Remove 'v' prefix if present
  } catch (error) {
    console.error(`Error fetching version: ${error.message}`);
    // Fallback to a hardcoded version as last resort
    return '0.11.5'; // Update this with each release as a fallback
  }
}

// Map Node.js platform/arch to GoReleaser naming
const platformMap = {
  'darwin': 'Darwin',
  'linux': 'Linux',
  'win32': 'Windows'
};

const archMap = {
  'x64': 'x86_64',
  'arm64': 'arm64'
};

async function install() {
  try {
    const platform = process.platform;
    const arch = process.arch;
    
    // Check if platform/arch is supported
    if (!platformMap[platform]) {
      throw new Error(`Unsupported platform: ${platform}. Only macOS, Linux, and Windows are supported.`);
    }
    
    if (!archMap[arch]) {
      throw new Error(`Unsupported architecture: ${arch}. Only x64 and arm64 are supported.`);
    }
    
    // Get the version to download
    const version = await getVersion();
    console.log(`Installing cmdeagle version ${version}...`);
    
    // Construct download URL
    const osName = platformMap[platform];
    const archName = archMap[arch];
    const extension = platform === 'win32' ? 'zip' : 'tar.gz';
    const url = `https://github.com/migsc/cmdeagle/releases/download/v${version}/cmdeagle_${version}_${osName}_${archName}.${extension}`;
    
    console.log(`Downloading from ${url}...`);
    
    // Create bin directory if it doesn't exist
    const binDir = path.join(__dirname, '..', 'bin');
    if (!fs.existsSync(binDir)) {
      fs.mkdirSync(binDir, { recursive: true });
    }
    
    // Set binary name based on platform
    const binaryName = platform === 'win32' ? 'cmdeagle.exe' : 'cmdeagle';
    const binaryPath = path.join(binDir, binaryName);
    
    // Download and extract
    const tempFile = path.join(binDir, `cmdeagle-${version}.${extension}`);
    
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Failed to download binary: ${response.statusText}`);
    }
    
    // Save the downloaded file
    await pipeline(
      response.body,
      createWriteStream(tempFile)
    );
    
    // Extract the binary
    if (platform === 'win32') {
      // For Windows, use a simple extraction approach
      const AdmZip = require('adm-zip');
      const zip = new AdmZip(tempFile);
      const zipEntries = zip.getEntries();
      
      // Find the binary in the zip
      for (const entry of zipEntries) {
        if (entry.entryName.endsWith('.exe')) {
          zip.extractEntryTo(entry, binDir, false, true);
          // Rename if needed
          const extractedPath = path.join(binDir, entry.entryName.split('/').pop());
          if (extractedPath !== binaryPath) {
            fs.renameSync(extractedPath, binaryPath);
          }
          break;
        }
      }
    } else {
      // For Unix systems, use tar
      await tar.extract({
        file: tempFile,
        cwd: binDir,
        filter: (path) => path.endsWith(binaryName)
      });
      
      // Find and move the binary if it's in a subdirectory
      const files = fs.readdirSync(binDir);
      for (const file of files) {
        const filePath = path.join(binDir, file);
        if (fs.statSync(filePath).isDirectory()) {
          const nestedBinary = path.join(filePath, binaryName);
          if (fs.existsSync(nestedBinary)) {
            fs.renameSync(nestedBinary, binaryPath);
            fs.rmdirSync(filePath, { recursive: true });
          }
        }
      }
    }
    
    // Make binary executable on Unix systems
    if (platform !== 'win32') {
      fs.chmodSync(binaryPath, 0o755);
    }
    
    // Clean up
    fs.unlinkSync(tempFile);
    
    console.log(`Successfully installed cmdeagle ${version} to ${binaryPath}`);
  } catch (error) {
    console.error(`Installation failed: ${error.message}`);
    
    // Clean up any partial downloads
    const tempFile = path.join(__dirname, '..', 'bin', `cmdeagle-${await getVersion()}.${process.platform === 'win32' ? 'zip' : 'tar.gz'}`);
    if (fs.existsSync(tempFile)) {
      fs.unlinkSync(tempFile);
    }
    
    process.exit(1);
  }
}

// Run the installation
install();
