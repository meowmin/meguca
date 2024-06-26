// Keyboard shortcuts and such

import options from "../options"
import { FormModel, postSM, postEvent, toggleExpandAll, expandThreadForm } from "../posts"
import { page } from "../state"
import { trigger } from "../util"
import {toggleNekoTV} from "../nekotv/nekotv";

// Bind keyboard event listener to the document
export default () =>
	document.addEventListener("keydown", handleShortcut)

function handleShortcut(event: KeyboardEvent) {
	let caught = false

	const anyModifier = event.altKey || event.metaKey || event.ctrlKey || event.shiftKey;
	const inInput = 'selectionStart' in event.target
	let altGr = event.getModifierState && event.getModifierState("AltGraph")
	if (navigator.platform.includes("Mac")) {
		altGr = false;
	}

	if (!anyModifier && !inInput) {
		caught = true
		switch (event.key) {
			case "w":
			case "ArrowLeft":
				if (options.arrowKeysNavigate) {
					navigatePost(true)
				}
				break
			case "s":
			case "ArrowRight":
				if (options.arrowKeysNavigate) {
					navigatePost(false)
				}
				break
			case "q":
				if (page.thread) {
					postSM.feed(postEvent.open)
				} else {
					expandThreadForm()
				}
				break
			default:
				caught = false
		}
		switch (event.key) {
			case options.volumeUp.toUpperCase():
				options.audioVolume = Math.min(100, options.audioVolume + 10)
				caught = true
				break
			case options.volumeDown.toUpperCase():
				options.audioVolume = Math.max(0, options.audioVolume - 10)
				caught = true
				break
		}
	}


	if (event.altKey && !altGr) {
		caught = true

		switch (event.which) {
			case options.newPost:
				if (page.thread) {
					postSM.feed(postEvent.open)
				} else {
					expandThreadForm()
				}
				break
			case options.done:
				postSM.feed(postEvent.done)
				break
			case options.toggleSpoiler:
				const m = trigger("getPostModel") as FormModel
				if (m) {
					m.view.toggleSpoiler()
				}
				break
			case options.galleryMode:
				options.galleryModeToggle = !options.galleryModeToggle
				break
			case options.expandAll:
				toggleExpandAll()
				break
			case options.workMode:
				options.workModeToggle = !options.workModeToggle
				break
			case options.meguTVShortcut:
				options.meguTV = !options.meguTV
				break
			case options.nekoTVShortcut:
				toggleNekoTV()
				break
			case 38:
				navigateUp()
				break
			default:
				caught = false
		}

	}

	if (caught) {
		event.stopImmediatePropagation()
		event.preventDefault()
	}
}

// Navigate one level up the board tree, if possible
function navigateUp() {
	let url: string
	if (page.thread) {
		url = `/${page.board}/`
	} else if (page.board !== "all") {
		url = "/all/"
	}
	if (url) {
		// Convert to absolute URL
		const a = document.createElement("a")
		a.href = url
		location.href = url
	}
}

const postSelector = "article[id^=p]"

function getArticleClosestToCenter(articles: Element[]): Element {
	const windowHeight = window.innerHeight;
	const windowWidth = window.innerWidth;
	const centerX = windowWidth / 2;
	const centerY = windowHeight / 2;

	let closestArticle: Element = articles[0];
	let minDistance = Infinity;

	for (let i = 0; i < articles.length; i++) {
		const article = articles[i];
		const rect = article.getBoundingClientRect();
		const articleCenterX = rect.left + rect.width / 2;
		const articleCenterY = rect.top + rect.height / 2;

		const distanceX = Math.abs(centerX - articleCenterX);
		const distanceY = Math.abs(centerY - articleCenterY);
		const distance = Math.sqrt(distanceX * distanceX + distanceY * distanceY);

		if (distance < minDistance) {
			minDistance = distance;
			closestArticle = article;
		}
	}

	return closestArticle;
}

// move focus to next or previous visible post in document order.
// starts with first post if none is selected via current url fragment
function navigatePost(reverse: boolean) {
	const all: Element[] = Array.from(document.querySelectorAll(postSelector))
	let current: Element = document.querySelector(postSelector + ":target");

	if (!current) {
		current = getArticleClosestToCenter(all);
	}
	let currentIdx = all.indexOf(current)

	while (current) {
		currentIdx = reverse ? currentIdx - 1 : currentIdx + 1
		current = all[currentIdx]
		if (current && window.getComputedStyle(current).display != "none") {
			break
		}
	}

	if (current) {
		window.location.hash = current.id
	}
}
