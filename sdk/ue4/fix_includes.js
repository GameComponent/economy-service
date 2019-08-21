/**
 * Fix includes script
 * Uses correct file extension for Windows builds
 *
 * This program required two start arguments:
 * --file should link to an input .proto file
 */

const fs = require('fs');
const args = process.argv;

const fileIndex = args.indexOf('--file');

// Check if the start arguments are present
if (fileIndex == -1) {
  console.log('No path to input file given. Use --file.')
  process.exit(1);
}

const fileValue = args[fileIndex + 1];

// Validate the inputs
if (!fileValue) {
  console.log('No path to input file given. Use --file.')
  process.exit(1);
}

const REPLACE_PB_INCLUDE = '#include "Wrappers/economy_service/economy_service.pb.hpp"';
const REPLACEMENT_PB_INCLUDE = '#include "Wrappers/economy_service/economy_service.pb.h"';
const REPLACE_GRPC_INCLUDE = '#include "Wrappers/economy_service/economy_service.grpc.pb.hpp"';
const REPLACEMENT_GRPC_INCLUDE = '#include "Wrappers/economy_service/economy_service.grpc.pb.h"';

// Read the file
let contents = fs.readFileSync(fileValue, 'utf8');

// Process the contents
contents = contents.split(REPLACE_PB_INCLUDE).join(REPLACEMENT_PB_INCLUDE);
contents = contents.split(REPLACE_GRPC_INCLUDE).join(REPLACEMENT_GRPC_INCLUDE);

// Write the output
fs.writeFileSync(fileValue, contents, 'utf8');
