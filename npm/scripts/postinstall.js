const https = require('https');
const fs = require('fs');
const path = require('path');
const os = require('os');
const packageJson = require('../package.json');

const version = packageJson.version;
const platform = process.platform;
const arch = process.arch;

let releaseName = '';
if (platform === 'win32') {
  if (arch === 'x64') releaseName = 'nipo-windows-amd64.exe';
} else if (platform === 'darwin') {
  if (arch === 'x64') releaseName = 'nipo-darwin-amd64';
  else if (arch === 'arm64') releaseName = 'nipo-darwin-arm64';
} else if (platform === 'linux') {
  if (arch === 'x64') releaseName = 'nipo-linux-amd64';
  else if (arch === 'arm64') releaseName = 'nipo-linux-arm64';
}

if (!releaseName) {
  console.error(`Error: Unsupported platform/architecture: ${platform}/${arch}`);
  process.exit(1);
}

const url = `https://github.com/ngtuonghy/nipo/releases/download/v${version}/${releaseName}`;

const destDir = path.join(os.homedir(), '.nipo', 'bin');
if (!fs.existsSync(destDir)) {
  fs.mkdirSync(destDir, { recursive: true });
}

const destFile = path.join(destDir, platform === 'win32' ? 'nipo.exe' : 'nipo');

console.log(`Downloading Nipo core binary for ${platform}-${arch} from ${url}...`);

function downloadFile(url, dest, callback) {
  const file = fs.createWriteStream(dest);
  https.get(url, (response) => {
    if (response.statusCode === 302 || response.statusCode === 301) {
      file.close();
      fs.unlinkSync(dest);
      downloadFile(response.headers.location, dest, callback);
      return;
    }

    if (response.statusCode !== 200) {
      file.close();
      fs.unlinkSync(dest);
      callback(new Error(`Server returned status code ${response.statusCode}`));
      return;
    }

    response.pipe(file);
    file.on('finish', () => {
      file.close(callback);
    });
  }).on('error', (err) => {
    file.close();
    fs.unlink(dest, () => { });
    callback(err);
  });
}

downloadFile(url, destFile, (err) => {
  if (err) {
    console.error(`\x1b[31mError downloading binary: ${err.message}\x1b[0m`);
    console.error('Please make sure you have internet access.');
    process.exit(1);
  }

  if (platform !== 'win32') {
    try {
      fs.chmodSync(destFile, 0755);
    } catch (chmodErr) {
      console.warn(`Warning: Failed to set executable permissions: ${chmodErr.message}`);
    }
  }

  console.log('\x1b[32mSuccessfully installed Nipo core binary!\x1b[0m');
});
