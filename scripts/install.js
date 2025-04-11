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

// Get package info
const packageJson = require('../package.json');
const version = packageJson.version;

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

// Determine current platform and architecture
const platform = process.platform;
const arch = process.arch;

if (!platformMap[platform]) {
  console.error(`Unsupported platform: ${platform}`);
  process.exit(1);
}

if (!archMap[arch]) {
  console.error(`Unsupported architecture: ${arch}`);
  process.exit(1);
}

const mappedPlatform = platformMap[platform];
const mappedArch = archMap[arch];
const extension = platform === 'win32' ? '.exe' : '';

// Create bin directory if it doesn't exist
const binDir = path.join(__dirname, 'bin');
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

// Set up paths
const binaryName = `cmdeagle${extension}`;
const binaryPath = path.join(binDir, binaryName);
const fileName = `cmdeagle_${version}_${mappedPlatform}_${mappedArch}.tar.gz`;
const downloadUrl = `https://github.com/migsc/cmdeagle/releases/download/v${version}/${fileName}`;
const checksumUrl = `https://github.com/migsc/cmdeagle/releases/download/v${version}/checksums.txt`;
const tempFile = path.join(__dirname, fileName);

console.log(`Downloading cmdeagle ${version} for ${mappedPlatform} ${mappedArch}...`);
console.log(`Download URL: ${downloadUrl}`);

async function downloadFile(url, dest) {
  const response = await fetch(url);
  
  if (!response.ok) {
    throw new Error(`Failed to download from ${url}: ${response.statusText}`);
  }
  
  const fileStream = createWriteStream(dest);
  await pipeline(response.body, fileStream);
}

async function verifyChecksum(filePath, checksumUrl) {
  try {
    // Download checksums file
    const response = await fetch(checksumUrl);
    if (!response.ok) {
      console.warn('Could not download checksums file for verification. Skipping verification.');
      return true;
    }
    
    const checksums = await response.text();
    
    // Calculate SHA256 of the downloaded file
    const fileBuffer = fs.readFileSync(filePath);
    const hashSum = crypto.createHash('sha256');
    hashSum.update(fileBuffer);
    const fileHash = hashSum.digest('hex');
    
    // Find the matching line in checksums file
    const expectedLine = checksums.split('\n').find(line => 
      line.includes(path.basename(filePath))
    );
    
    if (!expectedLine) {
      console.warn(`Could not find checksum for ${path.basename(filePath)}. Skipping verification.`);
      return true;
    }
    
    const expectedHash = expectedLine.split(/\s+/)[0];
    
    if (fileHash !== expectedHash) {
      console.error(`Checksum verification failed!`);
      console.error(`Expected: ${expectedHash}`);
      console.error(`Got: ${fileHash}`);
      return false;
    }
    
    console.log('Checksum verification passed.');
    return true;
  } catch (error) {
    console.warn(`Checksum verification error: ${error.message}`);
    console.warn('Continuing with installation...');
    return true; // Continue even if verification fails
  }
}

async function install() {
  try {
    // Download the tarball
    await downloadFile(downloadUrl, tempFile);
    
    // Verify checksum
    const isValid = await verifyChecksum(tempFile, checksumUrl);
    if (!isValid) {
      throw new Error('Checksum verification failed. Aborting installation.');
    }
    
    // Extract the tarball
    console.log('Extracting binary...');
    await tar.x({
      file: tempFile,
      cwd: binDir,
      strip: 1 // Remove the top-level directory from the archive
    });
    
    // Make the binary executable on Unix platforms
    if (platform !== 'win32') {
      fs.chmodSync(binaryPath, 0o755);
    }
    
    // Clean up
    fs.unlinkSync(tempFile);
    
    console.log(`Successfully installed cmdeagle ${version} to ${binaryPath}`);
  } catch (error) {
    console.error(`Installation failed: ${error.message}`);
    
    // Clean up any partial downloads
    if (fs.existsSync(tempFile)) {
      fs.unlinkSync(tempFile);
    }
    
    process.exit(1);
  }
}

// Run the installation
install();
