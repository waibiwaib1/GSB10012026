let status = true; // Indicates if log events are being emitted to log providers or not

module.exports = {
	"pause": () => {
		status = false;
	},
	"resume": () => {
		status = true;
	},
	"status": () => {
		return status ? "active" : "paused";
	}
};
