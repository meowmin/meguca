(function() {
	function loadScript(path) {
		var head = document.getElementsByTagName('head')[0];
		var script = document.createElement('script');
		script.type = 'text/javascript';
		script.src = path;
		head.appendChild(script);
		return script;
	}

	loadScript(webpackStaticJSFile)
	// 	.onload = function () {
	// 	require("client/static/main");
	// };
})();
