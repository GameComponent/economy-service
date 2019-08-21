/**
 * Convert script
 * This script will remove unsupported code from a protobuf file
 * The output is used to generate UE4 plugin code
 *
 * This program required two start arguments:
 * --input should link to an input .proto file
 * --output is the place where the final proto file will be writen to
 */

const fs = require('fs');
const args = process.argv;

const inputIndex = args.indexOf('--input');
const outputIndex = args.indexOf('--output');

// Check if the start arguments are present
if (inputIndex == -1) {
	console.log('No path to input file given. Use --input.')
	process.exit(1);
}

if (outputIndex == -1) {
	console.log('No path to output file given. Use --output.')
	process.exit(1);
}

const inputValue = args[inputIndex + 1];
const outputValue = args[outputIndex + 1];

// Validate the inputs
if (!inputValue) {
	console.log('No path to input file given. Use --input.')
	process.exit(1);
}

if (!outputValue) {
	console.log('No path to output file given. Use --output.')
	process.exit(1);
}

// Match settings
const MATCH_GRPC_GATEWAY_OPTIONS_START = /option(\s?)\(grpc.gateway.protoc_gen_swagger.options.openapiv2_swagger\)/g;
const MATCH_IMPORT_TIMESTAMP = 'import "google/protobuf/timestamp.proto";';
const MATCH_IMPORT_STRUCT = 'import "google/protobuf/struct.proto";';
const MATCH_IMPORTS = /import[^;]*;/g;
const MATCH_GOOGLE_API_HTTP_OPTIONS = /option(\s?)\(google.api.http\)[^{]*\{*[\w\W]*?};+/g;
const MATCH_GOOGLE_PROTOBUF_TIMESTAMP = 'google.protobuf.Timestamp';
const MATCH_GOOGLE_PROTOBUF_STRUCT = 'google.protobuf.Struct';
const MATCH_GOOGLE_PROTOBUF_VALUE = 'google.protobuf.Value';

// Inline settings
const INLINE_TIMESTAMP = `
message Timestamp {
  int64 seconds = 1;
  int32 nanos = 2;
}
`;

const INLINE_STUCT = `
message Struct {
  map<string, Value> fields = 1;
}

message Value {
  oneof kind {
    NullValue null_value = 1;
    double number_value = 2;
    string string_value = 3;
    bool bool_value = 4;
    Struct struct_value = 5;
    ListValue list_value = 6;
  }
}

enum NullValue {
  NULL_VALUE = 0;
}

message ListValue {
  repeated Value values = 1;
}
`;

const INLINE_GOOGLE_PROTOBUF_TIMESTAMP = 'Timestamp';
const INLINE_GOOGLE_PROTOBUF_STRUCT = 'Struct';
const INLINE_GOOGLE_PROTOBUF_VALUE = 'Value';

// Read the file
let contents = fs.readFileSync(inputValue, 'utf8');

// Remove until the bracket closes
let open = 0;
let startLine = -1;
let endLine = -1;

const lines = contents.split('\n');
lines.forEach((line, index) => {
	if (endLine > -1) return;
	if (line.match(MATCH_GRPC_GATEWAY_OPTIONS_START)) {
		startLine = index;
	}
	if (startLine == -1) return

	open += line.split('{').length - 1;
	open -= line.split('}').length - 1;

	if (open === 0 && startLine !== index) {
		endLine = index + 1;
	}
});

if (startLine > -1 && endLine > -1) {
	contents = [...lines.slice(0, startLine), ...lines.slice(endLine)].join('\n');
}

// Process the contents
contents = contents.split(MATCH_IMPORT_STRUCT).join(INLINE_STUCT);
contents = contents.split(MATCH_IMPORT_TIMESTAMP).join(INLINE_TIMESTAMP);
contents = contents.split(MATCH_IMPORTS).join('');
contents = contents.split(MATCH_GOOGLE_API_HTTP_OPTIONS).join('');
contents = contents.split(MATCH_GOOGLE_PROTOBUF_STRUCT).join(INLINE_GOOGLE_PROTOBUF_STRUCT);
contents = contents.split(MATCH_GOOGLE_PROTOBUF_VALUE).join(INLINE_GOOGLE_PROTOBUF_VALUE);
contents = contents.split(MATCH_GOOGLE_PROTOBUF_TIMESTAMP).join(INLINE_GOOGLE_PROTOBUF_TIMESTAMP);

// Write the output
fs.writeFileSync(outputValue, contents, 'utf8');
