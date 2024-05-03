import {connSM, connState, message, sendBinary} from "../connection";
import {escape} from "../util"
import {
    AddVideoEvent,
    ClearPlaylistEvent,
    ConnectedEvent,
    DumpEvent,
    GetTimeEvent,
    PauseEvent,
    PlayEvent,
    PlayItemEvent,
    RemoveVideoEvent,
    RewindEvent,
    SetNextItemEvent,
    SetRateEvent,
    SetTimeEvent,
    SkipVideoEvent,
    TogglePlaylistLockEvent,
    UpdatePlaylistEvent,
    WebSocketMessage,
} from "../typings/messages";
import {Player} from "./player";
import {ytPlayer} from "./ytplayer";

let player: Player;

export let playlistDiv: HTMLDivElement;
export let playlistOl: HTMLOListElement;
export let playerDiv: HTMLDivElement;
export let playlistStatus: HTMLElement;
export let vidEl: HTMLVideoElement;
export let watchStatus: HTMLElement;
export let currentSource: string;
export let watchDiv: HTMLElement;
let playerTimeInterval: NodeJS.Timeout | null = null;
let nekoTV = document.getElementById("banner-nekotv");
let isOpen : boolean;
let isPlaylistVisible = false;
let subscribeMessage = new Uint8Array([1,message.nekoTV]).buffer
let unsubMessage = new Uint8Array([0,message.nekoTV]).buffer
let isMuted : boolean;

export function initNekoTV() {
    if (!nekoTV) {
        return;
    }
    playlistDiv = document.getElementById('watch-playlist') as HTMLDivElement;
    playlistOl = document.getElementById('watch-playlist-entries') as HTMLOListElement;
    playerDiv = document.getElementById('watch-player') as HTMLDivElement;
    playlistStatus = document.getElementById('watch-playlist-status')!;
    vidEl = document.getElementById('watch-video') as HTMLVideoElement;
    watchStatus = document.getElementById('status-watch')!;
    watchDiv = document.getElementById("watch-panel");
    playerDiv.addEventListener("click",()=>{
        if (playlistDiv.style.display) {
            playlistDiv.style.display = ''
        } else {
            playlistDiv.style.display = 'block'
        }
    })
    let lastVal = localStorage.getItem('neko-tv')
    if (lastVal) {
        isOpen = lastVal === 't';
    } else {
        isOpen = true;
    }
    updateNekoTVIcon()
    connSM.on(connState.synced,subscribeToWatchFeed)
    nekoTV.addEventListener("click", () => {
        isOpen = !isOpen;
        localStorage.setItem('neko-tv', isOpen ? 't' : 'f');
        updateNekoTVIcon()
        togglePlayer()
    });

    let watchCloseButton = document.getElementById('watch-close-button');
    let watchMuteButton = document.getElementById('watch-mute-button');
    watchCloseButton.addEventListener('click',()=>{
        isOpen = false;
        localStorage.setItem('neko-tv', 'f');
        updateNekoTVIcon()
        togglePlayer()
    })
    lastVal = localStorage.getItem('neko-tv-mute')
    if (lastVal) {
        isMuted = lastVal === 't';
    } else {
        isMuted = false;
    }
    if(isMuted) {
        watchMuteButton.innerText = '􀊢'
        watchMuteButton.title = 'Unmute'
    }
    else {
        watchMuteButton.innerText = '􀊦'
        watchMuteButton.title = 'Mute'
    }
    watchMuteButton.addEventListener('click',()=> {
        isMuted = !isMuted;
        localStorage.setItem('neko-tv-mute', isMuted ? 't' : 'f');
        if (ytPlayer) {
            if(isMuted){
                ytPlayer.mute()
            }
            else {
                ytPlayer.unMute()
            }
        }
        if(isMuted) {
            watchMuteButton.innerText = '􀊢'
            watchMuteButton.title = 'Mute'
        }
        else {
            watchMuteButton.innerText = '􀊦'
            watchMuteButton.title = 'Unmute'
        }
    })
    player = new Player()

}

export function isNekoTVOpen() {
    return isOpen;
}

export function isNekoTVMuted() {
    return isMuted;
}

function updateNekoTVIcon(){
    if (isOpen) {
        nekoTV.innerText = '􀵨';
        nekoTV.title = 'NekoTV: Enabled'
    } else {
        nekoTV.innerText = '􁋞';
        nekoTV.title = 'NekoTV: Disabled'
    }

}
export function showWatchPanel() {
    watchDiv.style.display = 'block';
    watchDiv.classList.remove('hide-watch-panel');
}

export function hideWatchPanel() {
    watchDiv.classList.add('hide-watch-panel');
    watchDiv.style.display = 'none';
}
export function showPlaylist() {
    playlistDiv.style.display = 'block';
}

export function hidePlaylist() {
    playlistDiv.style.display = 'none';
    stopPlayerTimeInterval();
}

export function togglePlaylist() {
    isPlaylistVisible = !isPlaylistVisible;
    if (isPlaylistVisible) {
        showPlaylist();
    } else {
        hidePlaylist();
    }
}


export function updatePlayerTime() {
    if (!playlistOl || !playlistOl.firstElementChild || !ytPlayer || !ytPlayer.getCurrentTime) {
        console.log('Skipping updatePlayerTime');
        return;
    }

    let playerTime: number | undefined = null;
    playerTime = ytPlayer.getCurrentTime();

    if (playerTime === undefined) {
        console.error('Player time undefined');
        return;
    }

    playlistOl.children[player.getItemPos()].querySelector('.watch-player-time')!.innerHTML = `${secondsToTimeExact(playerTime)} / `;
}

function stopPlayerTimeInterval() {
    if (playerTimeInterval) {
        clearInterval(playerTimeInterval);
        playerTimeInterval = null;
    }
}

function handleConnectedEvent(connectedEvent: ConnectedEvent) {
    player.setItems(connectedEvent.videoList,connectedEvent.itemPos)
    handleSetTimeEvent(connectedEvent.getTime)
    updatePlaylist()
}

function handleAddVideoEvent(addVideoEvent: AddVideoEvent) {
    player.videoList.addItem(addVideoEvent.item, addVideoEvent.atEnd);
    if (player.videoList.length == 1) {
        player.setVideo(0);
    }
    updatePlaylist()
}

function handleRemoveVideoEvent(removeVideoEvent: RemoveVideoEvent) {
    player.removeItem(removeVideoEvent.url);
    if (player.isListEmpty()) {
        player.pause();
        hideWatchPanel();
    }
    updatePlaylist()
}

function handleSkipVideoEvent(skipVideoEvent: SkipVideoEvent) {
    player.skipItem(skipVideoEvent.url);
    if (player.isListEmpty()) player.pause();
    updatePlaylist()
}

function handlePauseEvent(pauseEvent: PauseEvent) {
    // player.setPauseIndicator(false);
    player.pause();
    player.setTime(pauseEvent.time);
}

function handlePlayEvent(playEvent: PlayEvent) {
    // player.setPauseIndicator(true);
    // const synchThreshold = player.settings.synchThreshold;
    const newTime = playEvent.time;
    const time = player.getTime();
    if (Math.abs(time - newTime) >= 1600) {
        player.setTime(newTime);
    }
    player.play();
}

function handleGetTimeEvent(getTimeEvent: GetTimeEvent) {
    console.log('Handling GetTimeEvent:', getTimeEvent);
    const paused = getTimeEvent.paused ?? false;
    const rate = getTimeEvent.rate ?? 1;

    if (player.getPlaybackRate() !== rate) {
        console.log('Updating playback rate to:', rate);
        player.setPlaybackRate(rate);
    }

    const synchThreshold = 1.6;
    const newTime = getTimeEvent.time;
    const time = player.getTime();

    console.log('Current time:', time);
    console.log('New time:', newTime);

    if (!player.isVideoLoaded()) {
        console.log('Video not loaded');
        // player.forceSyncNextTick = false;
    }
    if (player.getDuration() <= time + synchThreshold) {
        console.log('Video near end, skipping synchronization');
        return;
    }
    if (!paused) {
        console.log('Playing video');
        player.play();
    } else {
        console.log('Pausing video');
        player.pause();
    }
    // player.setPauseIndicator(!paused);
    if (Math.abs(time - newTime) < synchThreshold) {
        console.log('Time difference within threshold, skipping synchronization');
        return;
    }
    if (!paused) {
        console.log('Synchronizing time to:', newTime + 0.5);
        player.setTime(newTime + 0.5);
    } else {
        console.log('Synchronizing time to:', newTime);
        player.setTime(newTime);
    }
}

function handleSetTimeEvent(setTimeEvent: SetTimeEvent) {
    const synchThreshold = 1600;
    const newTime = setTimeEvent.time;
    const time = player.getTime();
    if (Math.abs(time - newTime) < synchThreshold) {
        return;
    }
    player.setTime(newTime);
}

function handleSetRateEvent(setRateEvent: SetRateEvent) {
    player.setPlaybackRate(setRateEvent.rate);
}

function handleRewindEvent(rewindEvent: RewindEvent) {
    player.setTime(rewindEvent.time + 0.5);
}

function handlePlayItemEvent(playItemEvent: PlayItemEvent) {
    player.setVideo(playItemEvent.pos);
}

function handleSetNextItemEvent(setNextItemEvent: SetNextItemEvent) {
    player.setNextItem(setNextItemEvent.pos);
}

function handleUpdatePlaylistEvent(updatePlaylistEvent: UpdatePlaylistEvent) {
    player.setItems(updatePlaylistEvent.videoList.items);
}

function handleTogglePlaylistLockEvent(togglePlaylistLockEvent: TogglePlaylistLockEvent) {
    // player.setPlaylistLock(togglePlaylistLockEvent.isOpen);
}

function handleDumpEvent(dumpEvent: DumpEvent) {
    // Implement the logic for handling the dump event if needed
}

function handleClearPlaylistEvent(clearPlaylistEvent: ClearPlaylistEvent) {
    player.clearItems();
    if (player.isListEmpty()) {
        player.pause();
    }
    updatePlaylist()
}

export function handleMessage(message: WebSocketMessage) {
    if (message.connectedEvent) {
        handleConnectedEvent(message.connectedEvent);
    } else if (message.addVideoEvent) {
        handleAddVideoEvent(message.addVideoEvent);
    } else if (message.removeVideoEvent) {
        handleRemoveVideoEvent(message.removeVideoEvent);
    } else if (message.skipVideoEvent) {
        handleSkipVideoEvent(message.skipVideoEvent);
    } else if (message.pauseEvent) {
        handlePauseEvent(message.pauseEvent);
    } else if (message.playEvent) {
        handlePlayEvent(message.playEvent);
    } else if (message.getTimeEvent) {
        handleGetTimeEvent(message.getTimeEvent);
    } else if (message.setTimeEvent) {
        handleSetTimeEvent(message.setTimeEvent);
    } else if (message.setRateEvent) {
        handleSetRateEvent(message.setRateEvent);
    } else if (message.rewindEvent) {
        handleRewindEvent(message.rewindEvent);
    } else if (message.playItemEvent) {
        handlePlayItemEvent(message.playItemEvent);
    } else if (message.setNextItemEvent) {
        handleSetNextItemEvent(message.setNextItemEvent);
    } else if (message.updatePlaylistEvent) {
        handleUpdatePlaylistEvent(message.updatePlaylistEvent);
    } else if (message.togglePlaylistLockEvent) {
        handleTogglePlaylistLockEvent(message.togglePlaylistLockEvent);
    } else if (message.dumpEvent) {
        handleDumpEvent(message.dumpEvent);
    } else if (message.clearPlaylistEvent) {
        handleClearPlaylistEvent(message.clearPlaylistEvent);
    } else {
        console.error("Invalid WebSocketMessage received");
    }
}
function truncateWithEllipsis(e, t) {
    return e.length <= t ? e : e.substring(0, t) + "…"
}
export function secondsToTimeExact(totalSeconds: number): string {
    totalSeconds = Math.floor(totalSeconds);

    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds - hours * 3600) / 60);
    const seconds = Math.round(totalSeconds - hours * 3600 - minutes * 60);

    let formattedTime: string;

    if (hours) {
        formattedTime = `${hours}:${padWithZero(minutes)}:${padWithZero(seconds)}`;
    } else if (minutes) {
        formattedTime = `${minutes}:${padWithZero(seconds)}`;
    } else {
        formattedTime = `0:${padWithZero(seconds)}`;
    }

    return formattedTime;
}

function padWithZero(value: number): string {
    return value < 10 ? `0${value}` : value.toString();
}
export function updatePlaylist() {
    if (player.isListEmpty()) {
        removePlayer()
        return;
    }
    if (!playerTimeInterval) {
        updatePlayerTime();
        playerTimeInterval = setInterval(updatePlayerTime, 1000);
    }

    // updatePlaylistStatus();
    showWatchPanel()

    const playlistItems: HTMLLIElement[] = [];

    const currentItemPos = player.getItemPos();

    for (let i = 0; i < player.videoList.items.length; i++) {
        const video = player.videoList.items[i];
        const li = document.createElement('li');
        li.classList.add('watch-playlist-entry');

        if (i === currentItemPos) {
            li.classList.add('selected');
        }

        let videoTerm = '';
        if (video.url && !video.url.startsWith('https')) {
            videoTerm = escape(truncateWithEllipsis(video.url, 25));
        }

        const videoTitle = escape(video.title);

        li.innerHTML = `
  <span class="watch-video-term">${videoTerm}</span>
  <a class="watch-video-title" target="_blank" href="${video.url}" title="${escape(video.title)}">
    ${videoTitle}
  </a>
  <span class="watch-video-time">
    <span class="watch-player-time"></span>
    <span class="watch-player-dur">${secondsToTimeExact(video.duration)}</span>
  </span>
`;

        playlistItems.push(li);
    }
    playlistOl.replaceChildren(...playlistItems);

    if (!isOpen) {
        isOpen = true;
    }
}

export function togglePlayer() {
    if (isOpen) {
        subscribeToWatchFeed();
    }
    else {
        unsubscribeFromWatchFeed();
    }
}

export function unsubscribeFromWatchFeed() {
    removePlayer();
    sendBinary(unsubMessage)
}

export function subscribeToWatchFeed() {
    if (isOpen) sendBinary(subscribeMessage)
}

export function removePlayer() {
    if (ytPlayer) {
        ytPlayer.stopVideo()
    }
    hideWatchPanel();
    // stopPlayerTimeInterval();
}