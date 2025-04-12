const fs = require('fs');
const { execSync } = require('child_process');

// Get the release type from command line arguments
// npm run release [patch|minor|major]
const releaseType = process.argv[2];

// Validate release type
if (!releaseType || !['patch', 'minor', 'major'].includes(releaseType)) {
  console.error('Please specify a valid release type: npm run release [patch|minor|major]');
  process.exit(1);
}

// Function to get the latest version
function getCurrentVersion() {
  try {
    // Try to get version from the latest git tag
    const gitTags = execSync('git tag -l "v*" --sort=-v:refname').toString().trim().split('\n');
    if (gitTags.length > 0 && gitTags[0]) {
      return gitTags[0].substring(1); // Remove the 'v' prefix
    }
  } catch (error) {
    console.log('No git tags found or error getting tags, falling back to package.json');
  }

  // Fall back to package.json if no git tags
  try {
    const packageJson = JSON.parse(fs.readFileSync('./package.json', 'utf8'));
    return packageJson.version;
  } catch (error) {
    console.error('Error reading package.json:', error.message);
    process.exit(1);
  }
}

// Function to increment version based on release type
function incrementVersion(version, type) {
  const parts = version.split('.').map(Number);
  
  switch (type) {
    case 'major':
      parts[0]++;
      parts[1] = 0;
      parts[2] = 0;
      break;
    case 'minor':
      parts[1]++;
      parts[2] = 0;
      break;
    case 'patch':
      parts[2]++;
      break;
  }
  
  return parts.join('.');
}

// Get current version and increment it
const currentVersion = getCurrentVersion();
const newVersion = incrementVersion(currentVersion, releaseType);
const tagVersion = `v${newVersion}`;

console.log(`Current version: ${currentVersion}`);
console.log(`New version: ${newVersion} (${releaseType} release)`);

// Update package.json
const packageJson = JSON.parse(fs.readFileSync('./package.json', 'utf8'));
packageJson.version = newVersion;
fs.writeFileSync('./package.json', JSON.stringify(packageJson, null, 2) + '\n');
console.log(`Updated package.json to version ${newVersion}`);

// Get current branch name
const currentBranch = execSync('git rev-parse --abbrev-ref HEAD').toString().trim();
console.log(`Current branch: ${currentBranch}`);

// Git operations
try {
  // First, make sure the release.js script itself is committed
  try {
    const status = execSync('git status --porcelain scripts/release.js').toString();
    if (status) {
      console.log('Committing changes to release.js first...');
      execSync('git add scripts/release.js', {stdio: 'inherit'});
      execSync('git commit -m "Update release script"', {stdio: 'inherit'});
    }
  } catch (e) {
    console.log('No changes to release.js or error checking:', e.message);
  }

  // Now proceed with the version bump
  execSync('git add package.json', {stdio: 'inherit'});
  execSync(`git commit -m "Bump version to ${newVersion} (${releaseType})"`, {stdio: 'inherit'});
  execSync(`git push origin ${currentBranch}`, {stdio: 'inherit'});
  console.log(`Pushed version update to ${currentBranch} branch`);
  
  execSync(`git tag -a ${tagVersion} -m "Release ${tagVersion}"`, {stdio: 'inherit'});
  execSync('git push --tags', {stdio: 'inherit'});
  console.log(`Tagged and pushed ${tagVersion}`);
  
  console.log('\nRelease process completed successfully!');
  console.log(`GitHub Actions workflow should now be building release ${tagVersion}`);
} catch (error) {
  console.error('Error during git operations:', error.message);
  process.exit(1);
}