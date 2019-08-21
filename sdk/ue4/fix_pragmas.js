/**
 * Fix pragmas script
 * Adds extra pragmas to ignore warning when compiling in UE4
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

const PRAGMAS = `// Extra pramas to fix compilation in UE4.  DO NOT EDIT!
#pragma warning (disable : 4800) // forcing value to bool true or false
#pragma warning (disable : 4125) // decimal digit terminates octal escape sequence
#pragma warning (disable : 4647) // behavior change __is_pod has different value in previous version
#pragma warning (disable : 4668) // 'symbol' is not defined as a preprocessor macro, replacing with '0' for 'directives'
#pragma warning (disable : 4946) // reinterpret_cast used

`;

// Read the file
let contents = fs.readFileSync(fileValue, 'utf8');

// Process the contents
contents = `${PRAGMAS}${contents}`;

// Write the output
fs.writeFileSync(fileValue, contents, 'utf8');
