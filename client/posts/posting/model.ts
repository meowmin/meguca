import {handlers, message, send, sendBinary} from "../../connection"
import {Post} from "../model"
import {ImageData, PostData} from "../../common"
import FormView from "./view"
import {boardConfig, page, posts, storeMine, storeSeenPost} from "../../state"
import {postEvent, postSM, postState} from "."
import {extend, modPaste} from "../../util"
import {SpliceResponse} from "../../client"
import {FileData} from "./upload"
import {newAllocRequest} from "./identity"

// Form Model of an OP post
export default class FormModel extends Post {
	public inputBody = ""
	public view: FormView
	public allocatingImage: boolean = false;
	private static textEncoder = new TextEncoder();

	// Pass and ID, if you wish to hijack an existing model. To create a new
	// model pass zero.
	constructor() {
		// Initialize state
		super({
			id: 0,
			op: page.thread,
			editing: true,
			sage: false,
			sticky: false,
			locked: false,
			time: Math.floor(Date.now() / 1000),
			body: "",
			name: "",
			auth: 0,
			trip: "",
			moderation: [],
			state: {
				spoiler: false,
				quote: false,
				code: false,
				bold: false,
				italic: false,
				red: false,
				blue: false,
				haveSyncwatch: false,
				successive_newlines: 0,
				iDice: 0,
			},
		})
	}

	// Append a character to the model's body and reparse the line, if it's a
	// newline
	public append(code: number) {
		this.body += String.fromCodePoint(code)
	}

	// Append a string to the model's body and reparse the line, if it's a
	// newline
	public appendString(str : string) {
		this.body += str
	}

	// Remove the last character from the model's body
	public backspace() {
		this.body = this.body.slice(0, -1)
	}

	// Splice the last line of the body
	public splice(msg: SpliceResponse) {
		this.spliceText(msg)
	}

	// Compare new value to old and generate appropriate commands
	public parseInput(val: string): void {
		// These operations should only be performed on fresh allocations or
		// after the server has verified the allocation
		switch (postSM.state) {
			case postState.draft:
			case postState.alloc:
				break;
			default:
				return;
		}

		const old = this.inputBody
		val = this.trimInput(val, true);
		if (old === val) { // Everything already submitted
			return
		}

		const lenDiff = val.length - old.length;
		if (postSM.state === postState.draft) {
			this.requestAlloc(val, null)
		} else if (lenDiff === 1 && val.slice(0, -1) === old) {
			// Commit a character appendage to the end of the line to the server
			const char = val.slice(-1);
			this.inputBody += char
			const bytes = FormModel.textEncoder.encode(char);
			const newBytes = new Uint8Array(bytes.length + 1);
			newBytes.set(bytes, 0);
			newBytes[bytes.length] = message.append;
			// newBytes.set(bytes.length, message.append);
			this.sendBinary(newBytes);
			// sendBinary(newBytes)

		} else if (lenDiff === -1 && old.slice(0, -1) === val) {
			// Send a message about removing the last character of the line to
			// the server
			this.inputBody = this.inputBody.slice(0, -1)
			const backspaceMessage = new Uint8Array([message.backspace]);
			this.sendBinary(backspaceMessage);
		} else {
			this.commitSplice(val)
		}
	}

	// Trim input string, if it has too many lines
	private trimInput(val: string, write: boolean): string {
		if (val.length > 2000) {
			const extra = val.length - 2000;
			val = val.slice(0, 2000)
			if (write) {
				this.view.trimInput(extra);
			}
		}

		// Remove any lines past 30
		const lines = val.split("\n")
		if (lines.length - 1 > 100) {
			const trimmed = lines.slice(0, 100).join("\n")
			if (write) {
				this.view.trimInput(val.length - trimmed.length);
			}
			return trimmed;
		}

		return val;
	}


	private sendBinary(msg : ArrayBuffer) {
		if (postSM.state !== postState.halted) {
			sendBinary(msg)
		}
	}

	private send(type: message, msg: any,encode=true) {
		if (postSM.state !== postState.halted) {
			send(type, msg,encode)
		}
	}

	// Commit any other input change that is not an append or backspace
	private commitSplice(v: string) {
		// Convert to arrays of chars to deal with multibyte unicode chars
		const old = [...this.inputBody],
			val = [...v],
			start = diffIndex(old, val),
			till = diffIndex(
				old.slice(start).reverse(),
				val.slice(start).reverse(),
			)

		// create a new typed array
		const offsets = new Uint16Array([start,old.length - till - start]);
		const text = FormModel.textEncoder.encode(val.slice(start, -till || undefined).join(""));
		const newBytes = new Uint8Array(text.byteLength + 5);
		newBytes.set(new Uint8Array(offsets.buffer), 0);
		newBytes.set(text, offsets.byteLength);
		newBytes[offsets.byteLength + text.byteLength] = message.splice;
		this.sendBinary(newBytes);
		this.inputBody = v;
	}

	// Close the form and revert to regular post. Cancel also erases all post
	// contents.
	public commitClose() {
		this.parseInput(this.view.input.value)
		this.abandon()
		this.send(message.closePost, null)
	}

	// Turn post form into a regular post, because it has expired after a
	// period of posting ability loss
	public abandon() {
		this.view.cleanUp()
		this.closePost()
	}

	// Add a link to the target post in the input.
	public addReference(id: number, sel: string) {
		const pos = this.view.input.selectionEnd,
			old = this.view.input.value
		let s = '',
			b = false

		// Insert post link and preceding whitespace.
		switch (old.charAt(pos - 1)) {
			// If empty post or newline before cursor, tell
			// next switch to do a newline after the quote.
			case '':
			case '\n':
				b = true
			case ' ':
				s = `>>${id}`
				break
			default:
				s = sel ? `\n>>${id}` : ` >>${id}`
		}

		// Insert superceding whitespace after post link.
		switch (old.charAt(pos)) {
			// Remember the boolean from the last switch? If true, or
			// selection is true, add a newline and reset boolean.
			case '':
			case ' ':
			case '\n':
				s += (b || sel) ? '\n' : ''
				b = false
				break
			default:
				b = true
				s += sel ? '\n' : ' '
		}

		// If we do have a selection of text, then quote all lines.
		if (sel) {
			for (const line of sel.split('\n')) {
				s += `>${line}\n`
			}

			// If the boolean from earlier is still true, add a newline.
			s += b ? '\n' : ''
		}

		this.view.replaceText(
			old.slice(0, pos) + s + old.slice(pos),
			// If the boolean from earlier is still true, correct cursor position.
			pos + s.length - (b ? 1 : 0),
			// Don't commit a quote, if it is the only input in a post
			postSM.state !== postState.draft || old.length !== 0
		)
	}

	// Paste text to the text body
	public paste(sel: string) {
		const start = this.view.input.selectionStart,
			end = this.view.input.selectionEnd,
			old = this.view.input.value
		const p = modPaste(old, sel, end)

		if (!p) {
			return
		}

		if (p.body.length > 2000) {
			p.body = this.trimInput(p.body, false);
			p.pos = 2000;
		} else if (start != end) {
			p.body = old.slice(0, start) + p.body + old.slice(end)
			p.pos -= (end - start)
		} else {
			p.body = old.slice(0, end) + p.body + old.slice(end)
		}

		this.view.replaceText(p.body, p.pos,
			postSM.state !== postState.draft || old.length !== 0)
	}

	// Returns a function, that handles a message from the server, containing
	// the ID of the allocated post.
	private receiveID(): (id: number) => void {
		return (id: number) => {
			this.id = id
			this.op = page.thread
			this.seenOnce = true
			postSM.feed(postEvent.alloc)
			storeSeenPost(this.id, this.op)
			storeMine(this.id, this.op)
			posts.add(this)
			delete handlers[message.postID]
		}
	}
	public attachTiktok(input: string, hd: boolean, rotation: string){
		if(postSM.state == postState.draft) {
			this.allocatingImage = true;
			this.requestAlloc(this.trimInput(this.view.input.value, true),
				null);
			this.view.updateDoneButton()
		}
		const strArray = new TextEncoder().encode(input);
		const bufferSize = strArray.length + 3;
		const buffer = new ArrayBuffer(bufferSize);
		const view = new DataView(buffer);

		for (let i = 0; i < strArray.length; i++) {
			view.setUint8(i, strArray[i]);
		}
		if (rotation === "90") {
			view.setUint8(strArray.length, 1);
		} else if (rotation === "180") {
			view.setUint8(strArray.length, 2);
		} else if (rotation === "270") {
			view.setUint8(strArray.length, 3);
		} else {
			view.setUint8(strArray.length, 0);
		}
        view.setUint8(strArray.length + 1, hd ? 1 : 0);
        view.setUint8(strArray.length + 2, message.attachTiktok)
		sendBinary(buffer)
	}

	// Request allocation of a draft post to the server
	private requestAlloc(body: string, image: FileData | null) {
		const req = newAllocRequest();
		req["open"] = true;
		if (body) {
			req["body"] = this.inputBody = body;
		}
		if (image) {
			req["image"] = image;
		}

		send(message.insertPost, req);
		postSM.feed(postEvent.sentAllocRequest);
		handlers[message.postID] = this.receiveID();
	}

	// Handle draft post allocation
	public onAllocation(data: PostData) {
		extend(this, data);
		this.view.renderAlloc();
		if (this.image) {
			this.insertImage(this.image);
		}
		if (postSM.state !== postState.alloc) {
			this.propagateLinks();
		}
	}

	// Upload the file and request its allocation
	public async uploadFile(file: File) {
		if (!boardConfig.textOnly && !this.image) {
			const pr = this.view.upload.uploadFile(file);
			this.view.input.focus();
			this.handleUploadResponse(await pr);
		}
	}

	public async uploadFileHash(hash: string) {
		if (!boardConfig.textOnly && !this.image) {
			const pr = this.view.upload.uploadFileHash(hash);
			this.view.input.focus();
			this.handleUploadResponse(await pr);
			this.view.input.setSelectionRange(0, 0);
			this.view.input.value ="👈 ~~Reposted from MeguTV~~\n" +  this.view.input.value
			this.view.onInput()
			// this.paste()
			// this.view.replaceText("👈 ~~Reposted from MeguTV~~\n" + this.body, p.pos,
			// 	postSM.state !== postState.draft || old.length !== 0)
		}
	}

	private handleUploadResponse(data: FileData | null) {
		// Upload failed, canceled or image added while thumbnailing
		if (!data || this.image || this.allocatingImage) {
			return
		}

		switch (postSM.state) {
			case postState.draft:
				this.allocatingImage = true;
				this.requestAlloc(this.trimInput(this.view.input.value, true),
					data);
				break;
			case postState.allocating:
				// Will allocate post soon check back every 200 ms
				setTimeout(this.handleUploadResponse.bind(this, data), 200);
				break;
			case postState.alloc:
				this.allocatingImage = true;
				send(message.insertImage, data)
				break;
		}
	}

	// Retry to upload a file after it previously failed
	public async retryUpload() {
		if (this.view.upload) {
			this.allocatingImage = false;
			this.handleUploadResponse(await this.view.upload.retry());
		}
	}

	// Insert the uploaded image into the model
	public insertImage(img: ImageData) {
		if (!this.editing) {
			// Stolen image
			super.insertImage(img);
			return;
		}

		this.image = img
		this.view.insertImage()
	}

	// Spoiler an already allocated image
	public commitSpoiler() {
		this.send(message.spoiler, null)
	}
}

// Find the first differing character in 2 character arrays
function diffIndex(a: string[], b: string[]): number {
	for (let i = 0; i < a.length; i++) {
		if (a[i] !== b[i]) {
			return i
		}
	}
	return a.length
}
