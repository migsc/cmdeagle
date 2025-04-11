const fs = require('fs');
const { execSync } = require('child_process');

// When running "npm run release v0.11.2", the version will be in process.argv[2]
// npm run [script] [arg1] [arg2] ...
//                    ^
//                    process.argv[2]
const version = process.argv[2];

if (!version || !/^v\d+\.\d+\.\d+$/.test(version)) {
    console.error('Please provide a valid version number: npm run release v1.2.3');
    process.exit(1);
}

// Remove the 'v' prefix for package.json
const npmVersion = version.substring(1);

// Read and update package.json
const packageJson = JSON.parse(fs.readFileSync('./package.json', 'utf8'));
packageJson.version = npmVersion;

// Write back with proper formatting (2 spaces indent)
fs.writeFileSync('./package.json', JSON.stringify(packageJson, null, 2) + '\n');
console.log(`Updated package.json to version ${npmVersion}`);

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
    execSync(`git commit -m "Bump version to ${npmVersion}"`, {stdio: 'inherit'});
    execSync(`git push origin ${currentBranch}`, {stdio: 'inherit'});
    console.log(`Pushed version update to ${currentBranch} branch`);
    
    execSync(`git tag -a ${version} -m "Release ${version}"`, {stdio: 'inherit'});
    execSync('git push --tags', {stdio: 'inherit'});
    console.log(`Tagged and pushed ${version}`);
    
    console.log('\nRelease process completed successfully!');
    console.log(`GitHub Actions workflow should now be building release ${version}`);
} catch (error) {
    console.error('Error during git operations:', error.message);
    process.exit(1);
}