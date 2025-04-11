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

// Git operations
try {
    execSync('git add package.json', {stdio: 'inherit'});
    execSync(`git commit -m "Bump version to ${npmVersion}"`, {stdio: 'inherit'});
    execSync('git push origin main', {stdio: 'inherit'});
    console.log('Pushed version update to main branch');
    
    execSync(`git tag -a ${version} -m "Release ${version}"`, {stdio: 'inherit'});
    execSync('git push --tags', {stdio: 'inherit'});
    console.log(`Tagged and pushed ${version}`);
    
    console.log('\nRelease process completed successfully!');
    console.log(`GitHub Actions workflow should now be building release ${version}`);
} catch (error) {
    console.error('Error during git operations:', error.message);
    process.exit(1);
}