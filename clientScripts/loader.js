// Selects and loads the client files and polyfills, if any. Use only ES5.

(function () {
	// Check if the client is an automated crawler
	var isBot,
		botStrings = [
			"bot", "googlebot", "crawler", "spider", "robot", "crawling"
		];

	for (var i = 0; i < botStrings.length; i++) {
		if (navigator.userAgent.indexOf(botStrings[i]) !== -1) {
			isBot = true;
			break;
		}
	}

	// Display mature content warning
	if (!isBot && config.mature && !localStorage.getItem("termsAccepted")) {
		var confirmText =
			"To access this website you understand and agree to the following:\n\n" +
			"1. The content of this website is for mature audiences only and may not be suitable for minors. If you are a minor or it is illegal for you to access mature images and language, do not proceed.\n\n" +
			"2. This website is presented to you AS IS, with no warranty, express or implied. By proceeding you agree not to hold the owner(s) of the website responsible for any damages from your use of the website, and you understand that the content posted is not owned or generated by the website, but rather by the website's users.";

		if (!confirm(confirmText)) {
			location.href = "http://www.gaiaonline.com/";
			return;
		}

		localStorage.setItem("termsAccepted", "true");
	}

	// Really old browser. Run in noscript mode.
	// if (!window.WebAssembly) {
	// 	var ns = document.getElementsByTagName("noscript");
	//
	// 	while (ns.length) { // Collection is live and changes with DOM updates
	// 		var el = ns[0],
	// 			cont = document.createElement("div");
	// 		cont.innerHTML = el.innerHTML;
	// 		el.parentNode.replaceChild(cont, el);
	// 	}
	//
	// 	var bc = document.getElementById("banner-center");
	// 	bc.classList.add("admin");
	// 	bc.innerHTML = "UPDATE YOUR FUCKING BROWSER";
	// 	return;
	// }

	// Remove prefixes on Web Crypto API for Safari
	if (!checkFunction("window.crypto.subtle.digest")) {
		window.crypto.subtle = window.crypto.webkitSubtle;
	}

	// TODO: Uncomment for WASM client rewrite
	// var wasm = /[\?&]wasm=true/.test(location.search);

	var head = document.getElementsByTagName('head')[0];
	loadClient();

	// Check if a browser API function is defined
	function checkFunction(func) {
		try {
			// See comment on line 134
			return typeof eval(func) === 'function';
		} catch (e) {
			return false;
		}
	}

	function loadScript(path) {
		var script = document.createElement('script');
		script.type = 'module';
		script.src = path;
		head.appendChild(script);
		return script;
	}

	function loadClient() {
		// Iterable NodeList
		if (!checkFunction('NodeList.prototype[Symbol.iterator]')) {
			NodeList.prototype[Symbol.iterator] =
				Array.prototype[Symbol.iterator];
		}

		// TODO: Uncomment for WASM client rewrite
		// if (wasm) {
		// 	window.Module = {};
		// 	fetch("/assets/wasm/main.wasm").then(function (res) {
		// 		return res.arrayBuffer();
		// 	}).then(function (bytes) {
		// 		// TODO: Parallel downloads of main.js and main.wasm
		// 		var script = document.createElement('script');
		// 		script.src = "/assets/wasm/main.js";
		// 		Module.wasmBinary = bytes;
		// 		document.head.appendChild(script);
		// 	});
		// } else {
		loadScript(webpackMainJSFile)
		// 	.onload = function () {
		// 	require("main");
		// };
		// }

		if ('serviceWorker' in navigator && (
			location.protocol === "https:" ||
			location.hostname === "localhost"
		)) {
			navigator.serviceWorker
				.register("/assets/js/scripts/worker.js", { scope: "/" })
				.catch(function (err) {
					throw err;
				});
		}
	}
})();
