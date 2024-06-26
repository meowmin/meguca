import options from ".";
import {importTemplate} from "../util";
import { page } from "../state";
import { sourcePath, serverNow } from "../posts";
import { fileTypes } from "../common"
import { handlers, message, connSM, connState, send } from "../connection"

type Data = {
	elapsed: number;
	playlist: Video[]
};

type Video = {
	sha1: string
	file_type: fileTypes
};

let playlist: Video[];
let lastStart = 0;

function render() {
	if (!playlist) {
		return
	}

	let cont = document.getElementById("megu-tv")
	if (!cont) {
		const modalOverlay = document.getElementById("modal-overlay");
		modalOverlay.prepend(importTemplate("megu-tv"));
	}
	cont = document.getElementById("megu-tv")

	if (options.workModeToggle) {
		cont.removeAttribute("style");
		return;
	}

	// Remove old videos and add new ones, while preserving existing one.
	// Should help caching.
	const existing: { [sha1: string]: HTMLVideoElement } = {};
	if(cont.children != null) {
		for (const ch of [...cont.children] as HTMLVideoElement[]) {
			if (ch.tagName === "VIDEO") {
				ch.remove();
				existing[ch.getAttribute("data-sha1")] = ch;
			}
		}
	}
	for (let i = 0; i < playlist.length; i++) {
		let el = existing[playlist[i].sha1];
		if (!el) {
			el = document.createElement("video");
			el.setAttribute("data-sha1", playlist[i].sha1);
			el.setAttribute("style", "max-width:30vw");
			el.setAttribute("preload", "auto")
			el.id = "megu-tv-video"
			el.onmouseenter = () => el.controls = true;
			el.onmouseleave = () => el.controls = false;
			el.src = sourcePath(playlist[i].sha1, playlist[i].file_type);
			el.volume = options.audioVolume / 100;
		}

		// Buffer videos about to play by playing them hidden and muted
		if (!i) {
			el.currentTime = serverNow() - lastStart;
			el.classList.remove("hidden");
			el.muted = false;
		} else {
			el.muted = true;
			el.classList.add("hidden");
		}
		cont.append(el);
		el.play();
	}
}

export function persistMessages() {
	handlers[message.meguTV] = (data: Data) => {
		lastStart = serverNow() - data.elapsed;
		playlist = data.playlist;
		if (options.meguTV) {
			render();
		}
	}

	// Subscribe to feed on connection
	connSM.on(connState.synced, subscribe);
}

function subscribe() {
	if (options.meguTV) {
		send(message.meguTV, null);
	}
}

export default function () {
	const el = document.getElementById("megu-tv");
	if (el || page.board === "all" || !page.thread) {
		return;
	}
	if (connSM.state === connState.synced) {
		subscribe();
	}
	render();

	// Handle toggling of the option
	options.onChange("meguTV", on => {
		if (on && page.board !== "all") {
			if (!document.getElementById("megu-tv")) {
				render();
			}
		} else {
			const el = document.getElementById("megu-tv");
			if (el) {
				el.remove();
			}
		}
	});

	options.onChange("workModeToggle", on => {
		const el = document.getElementById("megu-tv");
		if (el) {
			if (on) {
				for (const ch of [...el.children] as HTMLVideoElement[]) {
					ch.muted = true;
				}
				render();
			} else {
				render();
				el.setAttribute("style", "display: block");
				const ch = el.firstChild as HTMLVideoElement;
				ch.muted = false;
			}
		}
	})
}
