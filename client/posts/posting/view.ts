import PostView from "../view"
import FormModel from "./model"
import { Post } from "../model"
import { boardConfig } from "../../state"
import { setAttrs, importTemplate, atBottom, scrollToBottom } from "../../util"
import { postSM, postEvent, postState } from "."
import UploadForm from "./upload"
import identity from "./identity"
import lang from "../../lang";

// Element at the bottom of the thread to keep the fixed reply form from
// overlapping any other posts, when scrolled till bottom
let bottomSpacer: HTMLElement

// Post creation and update view
export default class FormView extends PostView {
    public model: FormModel
    public input: HTMLTextAreaElement
    private observer: MutationObserver
    private previousHeight: number
    public upload: UploadForm

    constructor(model: Post) {
        super(model, null)
        this.renderInputs()
        this.initDraft()
    }

    // Render extra input fields for inputting text and optionally uploading
    // images
    private renderInputs() {
        this.input = document.createElement("textarea")
        setAttrs(this.input, {
            id: "text-input",
            name: "body",
            rows: "1",
            maxlength: "2000",
        })
        this.el.append(importTemplate("post-controls"))
        this.resizeInput()

        this.input.addEventListener("input", e => {
            e.stopImmediatePropagation()
            this.onInput()
        })
        this.onClick({
            "input[name=\"done\"]": postSM.feeder(postEvent.done),
        })
        this.updateDoneButton();

        if (!boardConfig.textOnly) {
            this.upload = new UploadForm(this.model,
                this.el.querySelector(".upload-container"),
                ()=>{this.onAttachTiktokButton()});
        }

        const bq = this.el.querySelector("blockquote")
        bq.innerHTML = ""
        bq.append(this.input)
        requestAnimationFrame(() =>
            this.input.focus());
    }

    public onAttachTiktokButton() {
        this.upload.attachTiktokButton.hidden = true;
        const template = importTemplate("attach-tiktok-form")
        this.el.append(template);
        const templateEl = this.el.querySelector(".attach-tiktok-form")
        const input = <HTMLInputElement>templateEl.querySelector(".attach-tiktok-form-row1 > input")
        // const lastHDSetting = localStorage.getItem('attach-hd-tiktoks');
        // const hdCheck = (<HTMLInputElement>templateEl.querySelector(".attach-tiktok-form-params input"))
        // hdCheck.checked = JSON.parse(lastHDSetting)

        templateEl.querySelector(".attach-tiktok-cancel").addEventListener("click",()=>{
            templateEl.remove()
            this.upload.attachTiktokButton.hidden = false;
        })
        templateEl.querySelector(".attach-tiktok-paste").addEventListener("click",()=>{
            navigator.clipboard.readText().then(text => {input.value = text})
        })
        templateEl.querySelector(".attach-tiktok-attach").addEventListener("click",()=>{
            const inputVal = input.value
            const rotation = templateEl.querySelector("select").value
            this.model.attachTiktok(inputVal,true,rotation)
            templateEl.remove()
            this.upload.attachTiktokButton.hidden = false;
        })
    }

    // Render a temporary view of the identity fields, so the user can see what
    // credentials he is about to post with
    public renderIdentity() {
        let { name, auth } = identity,
            trip = ""
        const i = name.indexOf("#")
        if (i !== -1) {
            trip = "?"
            name = name.slice(0, i)
        }

        this.el.querySelector(".name").classList.remove("admin")
        this.model.name = name.trim()
        this.model.trip = trip
        this.model.auth = auth ? -1 : 0 // Force question marks
        this.model.sage = identity.sage
        this.renderName()
    }

    // Initialize extra elements for a draft unallocated post
    private initDraft() {
        bottomSpacer = document.getElementById("bottom-spacer")
        this.el.classList.add("reply-form")
        this.el.querySelector("header").classList.add("temporary")
        this.renderIdentity()

        // Keep this post and bottomSpacer the same height
        this.observer = new MutationObserver(() =>
            this.resizeSpacer())
        this.observer.observe(this.el, {
            childList: true,
            attributes: true,
            characterData: true,
            subtree: true,
        })

        document.getElementById("thread-container").append(this.el)
        this.resizeSpacer()
        this.setEditing(true);
    }

    // Resize bottomSpacer to the same top position as this post
    private resizeSpacer() {
        // Not a reply
        if (!bottomSpacer) {
            return
        }

        const { height } = this.el.getBoundingClientRect()
        // Avoid needless writes
        if (this.previousHeight === height) {
            return
        }
        this.previousHeight = height
        bottomSpacer.style.height = `calc(${height}px - 2.1em)`
    }

    // Handle input events on this.input
    public onInput() {
        if (!this.input) {
            return;
        }
        this.model.parseInput(this.input.value);
        this.resizeInput();
        this.renderCounter();
    }

    // Resize textarea to content width and adjust height
    private resizeInput() {
        const el = this.input,
            s = el.style
        s.width = "0px"
        s.height = "0px"
        el.wrap = "off"
        // Make the line slightly larger, so there is enough space for the next
        // character. This prevents wrapping on type.
        s.width = Math.max(260, el.scrollWidth + 5) + "px"
        el.wrap = "soft"
        s.height = Math.max(16, el.scrollHeight) + "px"
    }

    // Trim input from the end by the supplied length
    public trimInput(length: number) {
        this.input.value = this.input.value.slice(0, -length)
    }

    // Replace the current body and set the cursor to the input's end.
    // commit sets, if the onInput method should be run.
    public replaceText(body: string, pos: number, commit: boolean) {
        const el = this.input
        el.value = body

        if (commit) {
            this.onInput()
        } else {
            this.resizeInput()
        }

        requestAnimationFrame(() => {
            el.focus()
            el.setSelectionRange(pos, pos)

            // Because Firefox refocuses the clicked <a>
            requestAnimationFrame(() =>
                el.focus())
        })
    }

    // Transform form into a generic post. Removes any dangling form controls
    // and frees up references.
    public cleanUp() {
        if (this.upload) {
            this.upload.cancel();
        }
        this.el.classList.remove("reply-form")
        const pc = this.el.querySelector("#post-controls")
        if (pc) {
            pc.remove()
        }
        if (bottomSpacer) {
            bottomSpacer.style.height = ""
            if (atBottom) {
                scrollToBottom()
            }
        }
        if (this.observer) {
            this.observer.disconnect()
        }
        bottomSpacer
            = this.observer
            = this.upload
            = null
    }

    // Special override for displaying the post, as if it was committed to
    // server. Increases perceived response speed.
    public closePost() {
        let oldBody: string;
        if (this.model.inputBody) {
            oldBody = this.model.body
            this.model.body = this.model.inputBody
            this.model.inputBody = null
        }
        const claudeExists = /#claude\s\S.*/.test(this.model.body)
        const claude = this.model.claude_state
        if (!claudeExists) {
            this.setEditing(false)
        }
        else if (claude != null && (claude.status == "error"|| claude.status == "done" )){
            this.setEditing(false)
        }
        this.reparseBody();
        if (oldBody) {
            this.model.body = oldBody
        }
        const attachForm = this.el.querySelector(".attach-tiktok-form")
        if (attachForm) {
            attachForm.remove()
        }
    }

    // Clean up on form removal
    public remove() {
        super.remove()
        this.cleanUp()
    }

    // Lock the post form after a critical error occurs
    public renderError() {
        this.el.classList.add("erred")
        this.input.setAttribute("contenteditable", "false")
    }

    // Transition into allocated post
    public renderAlloc() {
        this.id = this.el.id = "p" + this.model.id
        this.el.querySelector("header").classList.remove("temporary")
        this.renderHeader()
    }

    // Toggle the spoiler input checkbox
    public toggleSpoiler() {
        if (this.model.image && postSM.state === postState.alloc) {
            this.upload.hideSpoilerToggle();
            this.model.commitSpoiler();
        } else {
            const el = this.inputElement("spoiler");
            el.checked = !el.checked;
        }
    }

    // Insert image into an open post
    public insertImage() {
        this.renderImage(false);
        this.resizeInput();
        this.upload.hideMaskToggle();
        this.upload.hideButton();

        if (postSM.state !== postState.alloc) {
            return;
        }
        if (this.model.image.spoiler) {
            this.upload.hideSpoilerToggle();
        } else {
            this.inputElement("spoiler").addEventListener("change",
                this.toggleSpoiler.bind(this), { passive: true });
        }
    }

    // Update the display of the Done button according to postSM state
    public updateDoneButton() {
        const el = this.inputElement("done");
        if (!el) {
            return;
        }
        let text = lang.ui["done"];
        let disable = false;
        switch (postSM.state) {
            case postState.halted:
                disable = true;
                break;
            case postState.draft:
                text = lang.ui["cancel"];
                break;
            case postState.alloc:
                break;
            case postState.allocating:
            case postState.erred:
                disable = true;
                break;
        }
        el.disabled = disable;
        el.value = text;
    }

    // Update the character counter on an editing post
    private renderCounter() {
        const el = this.el.querySelector("#char-count") as HTMLDivElement;
        if (!el) {
            return;
        }
        const cnt = this.input.value.length;
        if (cnt < 1000) {
            el.style.display = "none";
            return;
        }
        el.style.display = "flex";
        (el.firstChild as HTMLDivElement).innerText = "" + cnt;
        if (cnt > 1900) {
            el.classList.add("admin");
        }
        else {
            el.classList.remove("admin");
        }
    }
}
