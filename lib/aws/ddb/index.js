let customDDB, aws;

function main() {
	if (!aws) aws = require("../index");
	return customDDB || new aws.sdk.DynamoDB();
}
main.set = (ddb) => customDDB = ddb;
main.revert = () => customDDB = null;
main.local = (endpoint = "http://localhost:8000") => {
	if (!aws) aws = require("../index");
	main.set(new aws.sdk.DynamoDB({
		endpoint
	}));
};

module.exports = main;
