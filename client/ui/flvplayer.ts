import {importTemplate} from "../util";
import Mpegts from "mpegts.js";

let mpegts : typeof Mpegts = null;

async function loadMpegtsjs(){
    if (!mpegts) {
        mpegts = (await import("mpegts.js")).default;
    }
    return mpegts;
}


let playerOpen = false;
let player: Mpegts.Player;

export function openFlvPlayer() {
    const cont = document.getElementById("flv-player-cont")
    if (!cont) {
        const playerElement = importTemplate("flv-player")
        document.getElementById("modal-overlay").prepend(playerElement);
        document.getElementById("flv-close-button").addEventListener("click", closeFlvPlayer)
        document.getElementById("flv-reload-button").addEventListener("click", reloadPlayer)
    }
    playerOpen = true
}

async function reloadPlayer() {
    player.unload()
    player.load()
    player.play()
}

function closeFlvPlayer() {
    destroyPlayer()
    const cont = document.getElementById("flv-player-cont")
    cont.remove()
    playerOpen = false
}


function destroyPlayer() {
    if (typeof player !== "undefined" && player != null) {
        player.unload();
        player.detachMediaElement();
        player.destroy();
        player = null;
    }
}

export async function playLive(url: string) {
    await loadMpegtsjs()
    const playerConfig : Mpegts.Config = {
        enableWorker: false,
        liveBufferLatencyChasing: true,
        liveBufferLatencyMaxLatency: 2,
        liveBufferLatencyMinRemain: 1,
    }
    if (mpegts.getFeatureList().mseLivePlayback) {
        const videoElement = document.getElementById('flv-player');
        destroyPlayer()
        player = mpegts.createPlayer({
            type: 'flv',  // could also be mpegts, m2ts, flv
            isLive: true,
            url: url
        }, playerConfig);
        player.attachMediaElement(<HTMLMediaElement>videoElement);
        player.load();
        player.play();
    }
}


export async function playButtonClicked(url: string) {
    console.log("Button clicked, " + playerOpen);
    if (!playerOpen) {
        await openFlvPlayer()
    }
    await playLive(url)
}

export default function initFlvPlayer() {
    (window as any).playButtonClicked = playButtonClicked;
    (window as any).openFlvPlayer = openFlvPlayer;
    (window as any).playLive = playLive;
}

