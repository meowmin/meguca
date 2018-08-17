import { ThreadData } from "../common"
import {
    extractConfigs, extractPost, reparseOpenPosts, extractPageData, hidePosts,
} from "./common"
import { findSyncwatches } from "../posts"
import { config } from "../state"
import { postSM, postState } from "../posts"

const counters = document.getElementById("thread-post-counters"),
    threads = document.getElementById("threads")
let postCtr = 0,
    imgCtr = 0,
    bumpTime = 0,
    isDeleted = false

// Render the HTML of a thread page
export default function () {
    extractConfigs()

    const { threads: data, backlinks } = extractPageData<ThreadData>(),
        { posts } = data

    data.posts = null

    postCtr = data.postCtr
    imgCtr = data.imageCtr
    bumpTime = data.bumpTime
    isDeleted = data.deleted
    renderPostCounter()

    extractPost(data, data.id, data.board, backlinks)

    for (let post of posts) {
        extractPost(post, data.id, data.board, backlinks)
    }
    hidePosts()
    reparseOpenPosts()
    findSyncwatches(threads)

    // Needs to be done, to  stop the FSM
    if (data.locked) {
        postSM.state = postState.threadLocked
    }
}

// Increment thread post counters and rerender the indicator in the banner
export function incrementPostCount(post: boolean, hasImage: boolean) {
    if (post) {
        postCtr++
        if (postCtr < 3000) {
            // An estimate, but good enough
            bumpTime = Math.floor(Date.now() / 1000)
        }
    }
    if (hasImage) {
        imgCtr++
    }
    renderPostCounter()
}

function renderPostCounter() {
    let text = ""
    if (postCtr) {
        text = `${postCtr} / ${imgCtr}`

        // Calculate estimated thread expiry time
        if (config.pruneThreads) {
            // Calculate expiry age
            const min = config.threadExpiryMin,
                max = config.threadExpiryMax
            let days = min + (-max + min) * (postCtr / 3000 - 1) ** 3
            if (isDeleted) {
                days /= 3
            }
            if (days < min) {
                days = min
            }

            // Subtract current bump time
            days -= (Date.now() / 1000 - bumpTime) / (3600 * 24)

            text += ` / `
            if (days > 1) {
                text += `${Math.round(days)}d`
            } else {
                text += `${Math.round(days / 24)}h`
            }
        }
    }
    counters.textContent = text
}
