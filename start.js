const { spawn } = require('child_process');
const fs = require('fs');
const path = require('path');

const colors = {
  reset: '\x1b[0m', bright: '\x1b[1m',
  green: '\x1b[32m', cyan: '\x1b[36m',
  yellow: '\x1b[33m', red: '\x1b[31m', magenta: '\x1b[35m'
};

console.log(`${colors.bright}${colors.magenta}==================================================${colors.reset}`);
console.log(`${colors.bright}${colors.cyan}    Omnichannel CRM & WebRTC PBX Platform Runner  ${colors.reset}`);
console.log(`${colors.bright}${colors.magenta}==================================================${colors.reset}\n`);

// 1. Check/copy .env
const envPath = path.join(__dirname, '.env');
const envExamplePath = path.join(__dirname, '.env.example');
if (!fs.existsSync(envPath)) {
  console.log(`${colors.yellow}[*] .env file not found. Copying from .env.example...${colors.reset}`);
  try {
    fs.copyFileSync(envExamplePath, envPath);
    console.log(`${colors.green}[+] Created .env file successfully.${colors.reset}`);
  } catch (err) {
    console.error(`${colors.red}[-] Failed to copy .env file: ${err.message}${colors.reset}`);
    process.exit(1);
  }
} else {
  console.log(`${colors.green}[+] Found existing .env configuration.${colors.reset}`);
}

// Helper: run a command, inherit stdio (user can type sudo password), return exit code
function runCmd(cmd, args, label) {
  return new Promise((resolve) => {
    console.log(`\n${colors.cyan}[*] ${label}${colors.reset}`);
    const proc = spawn(cmd, args, { stdio: 'inherit', cwd: __dirname });
    proc.on('close', (code) => resolve(code));
  });
}

async function main() {
  console.log(`${colors.yellow}[!] Sudo password may be required for Docker access.${colors.reset}`);

  // Step 1: Stop & remove any stale containers from previous runs
  await runCmd('sudo', ['docker', 'compose', 'down', '--remove-orphans'],
    'Stopping any existing containers (cleanup)...');

  // Step 2: Start all services fresh
  const code = await runCmd('sudo', ['docker', 'compose', 'up', '--build', '-d'],
    'Starting all services...');

  if (code === 0) {
    console.log(`\n${colors.bright}${colors.green}==================================================${colors.reset}`);
    console.log(`${colors.bright}${colors.green}🚀 ALL SERVICES STARTED SUCCESSFULLY!${colors.reset}`);
    console.log(`${colors.bright}${colors.green}==================================================${colors.reset}\n`);
    console.log(`  ${colors.bright}Frontend Dashboard:${colors.reset}  ${colors.cyan}http://localhost:3000${colors.reset}`);
    console.log(`  ${colors.bright}Backend REST API:${colors.reset}    ${colors.cyan}http://localhost:8080/health${colors.reset}`);
    console.log(`  ${colors.bright}RabbitMQ Panel:${colors.reset}      ${colors.cyan}http://localhost:25672${colors.reset}`);
    console.log(`  ${colors.bright}MinIO S3 Panel:${colors.reset}      ${colors.cyan}http://localhost:9001${colors.reset}`);
    console.log(`  ${colors.bright}Grafana Metrics:${colors.reset}     ${colors.cyan}http://localhost:3001${colors.reset}\n`);
    console.log(`${colors.yellow}View logs:  ${colors.bright}sudo docker compose logs -f${colors.reset}`);
    console.log(`${colors.yellow}Stop all:   ${colors.bright}sudo docker compose down${colors.reset}\n`);
  } else {
    console.error(`\n${colors.red}[-] Docker Compose failed with code ${code}. Check the logs above.${colors.reset}`);
  }
}

main();