// Code generated by qtc from "index.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line index.html:1
package templates

//line index.html:1
import "encoding/json"

//line index.html:2
import "strings"

//line index.html:3
import "github.com/bakape/meguca/config"

//line index.html:4
import "github.com/bakape/meguca/lang"

//line index.html:5
import "github.com/bakape/meguca/common"

//line index.html:6
import "github.com/bakape/meguca/assets"

//line index.html:8
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line index.html:8
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line index.html:8
func streamrenderIndex(qw422016 *qt422016.Writer, pos common.ModerationLevel) {
//line index.html:9
	conf := config.Get()

//line index.html:10
	ln := lang.Get()

//line index.html:11
	confJSON, confHash := config.GetClient()

//line index.html:12
	boards := config.GetBoards()

//line index.html:12
	qw422016.N().S(`<!doctype html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width"><meta name="application-name" content="meguca"><meta name="description" content="Realtime imageboard"><link type="image/x-icon" rel="shortcut icon" id="favicon" href="/assets/favicons/default.ico"><title id="page-title">`)
//line index.html:23
	qw422016.N().S(`$$$</title><link rel="manifest" href="/assets/mobile/manifest.json">`)
//line index.html:29
	qw422016.N().S(`<link rel="stylesheet" href="/assets/css/base.css"><link rel="stylesheet" id="theme-css" href="/assets/css/$$$.css"><style id="user-background-style"></style>`)
//line index.html:35
	qw422016.N().S(`<noscript><link rel="stylesheet" href="/assets/css/noscript.css"></noscript>`)
//line index.html:41
	qw422016.N().S(`<script>var config =`)
//line index.html:43
	qw422016.N().Z(confJSON)
//line index.html:43
	qw422016.N().S(`,configHash = '`)
//line index.html:44
	qw422016.N().S(confHash)
//line index.html:44
	qw422016.N().S(`',`)
//line index.html:45
	boardJSON, _ := json.Marshal(boards)

//line index.html:45
	qw422016.N().S(`boards =`)
//line index.html:46
	qw422016.N().Z(boardJSON)
//line index.html:46
	qw422016.N().S(`,position =`)
//line index.html:47
	qw422016.N().D(int(pos))
//line index.html:47
	qw422016.N().S(`,`)
//line index.html:48
	videosJSON, _ := json.Marshal(assets.GetVideoNames())

//line index.html:48
	qw422016.N().S(`bgVideos =`)
//line index.html:49
	qw422016.N().Z(videosJSON)
//line index.html:49
	qw422016.N().S(`;if (localStorage.theme !== config.DefaultCSS) {document.getElementById('theme-css').href = '/assets/css/' + localStorage.theme + '.css'}</script>`)
//line index.html:56
	qw422016.N().S(`<template name="article">`)
//line index.html:58
	streamdeletedToggle(qw422016)
//line index.html:58
	qw422016.N().S(`<header class="spaced"><input type="checkbox" class="mod-checkbox hidden"><h3 hidden></h3><b class="name spaced"></b><img class="flag" hidden><time></time><nav><a>No.</a><a class="quote"></a></nav>`)
//line index.html:71
	streamcontrolLink(qw422016)
//line index.html:71
	qw422016.N().S(`</header><div class="post-container"><blockquote></blockquote></div></template><template name="figcaption"><figcaption class="spaced"><a class="image-toggle act" hidden></a><span class="spaced image-search-container">`)
//line index.html:81
	engines := [...][2]string{
		{"google", "G"},
		{"yandex", "Yd"},
		{"iqdb", "Iq"},
		{"saucenao", "Sn"},
		{"whatAnime", "Wa"},
		{"desustorage", "Ds"},
		{"exhentai", "Ex"},
	}

//line index.html:90
	for _, e := range engines {
//line index.html:90
		qw422016.N().S(`<a class="image-search`)
//line index.html:91
		qw422016.N().S(` `)
//line index.html:91
		qw422016.N().S(e[0])
//line index.html:91
		qw422016.N().S(`" target="_blank" rel="nofollow">`)
//line index.html:92
		qw422016.N().S(e[1])
//line index.html:92
		qw422016.N().S(`</a>`)
//line index.html:94
	}
//line index.html:94
	qw422016.N().S(`</span><span class="fileinfo"><span class="media-artist"></span><span class="media-title"></span><span hidden class="has-audio">♫</span><span class="media-length"></span><span class="filesize"></span><span class="dims"></span></span><a></a></figcaption></template><template name="figure"><figure><a target="_blank"><img></a></figure></template><template name="post-controls"><div id="post-controls"><input name="done" type="button" value="`)
//line index.html:116
	qw422016.N().S(ln.Common.UI["done"])
//line index.html:116
	qw422016.N().S(`"><span class="upload-container" hidden><button>`)
//line index.html:119
	qw422016.N().S(ln.Common.UI["uploadFile"])
//line index.html:119
	qw422016.N().S(`</button><span data-id="spoiler"><label><input type="checkbox" name="spoiler">`)
//line index.html:124
	qw422016.N().S(ln.Common.Posts["spoiler"])
//line index.html:124
	qw422016.N().S(`</label></span><input type="file" hidden name="image" accept="image/png, image/gif, image/jpeg, video/webm, video/ogg, audio/ogg, application/ogg, video/mp4, audio/mp4, audio/mp3, application/zip, application/x-7z-compressed, application/x-xz, application/x-gzip, audio/x-flac, text/plain, application/pdf, video/quicktime, audio/x-flac"></span></div></template><template name="notification"><div class="notification modal glass show"><b class="admin"><b></div></template><template name="sticky">`)
//line index.html:137
	streamrenderSticky(qw422016, true)
//line index.html:137
	qw422016.N().S(`</template><template name="locked">`)
//line index.html:140
	streamrenderLocked(qw422016, true)
//line index.html:140
	qw422016.N().S(`</template>`)
//line index.html:142
	if pos > common.NotLoggedIn {
//line index.html:142
		qw422016.N().S(`<template name="keyValue">`)
//line index.html:144
		streamkeyValueForm(qw422016, "", "")
//line index.html:144
		qw422016.N().S(`</template><template name="arrayItem">`)
//line index.html:147
		streamarrayItemForm(qw422016, "")
//line index.html:147
		qw422016.N().S(`</template>`)
//line index.html:149
	}
//line index.html:149
	qw422016.N().S(`</head><body><div id="user-background"></div><div class="overlay-container">`)
//line index.html:156
	qw422016.N().S(`<span id="banner" class="glass"><nav id="board-navigation"><noscript>[`)
//line index.html:161
	for i, b := range boards {
//line index.html:162
		if i != 0 {
//line index.html:163
			qw422016.N().S(` `)
//line index.html:163
			qw422016.N().S(`/`)
//line index.html:163
			qw422016.N().S(` `)
//line index.html:164
		}
//line index.html:164
		qw422016.N().S(`<a href="/`)
//line index.html:165
		qw422016.N().S(b)
//line index.html:165
		qw422016.N().S(`/">`)
//line index.html:166
		qw422016.N().S(b)
//line index.html:166
		qw422016.N().S(`</a>`)
//line index.html:168
	}
//line index.html:168
	qw422016.N().S(`]</noscript></nav>`)
//line index.html:174
	qw422016.N().S(`<b id="banner-center" class="spaced"></b>`)
//line index.html:178
	qw422016.N().S(`<span><b id="sync" class="banner-float svg-link noscript-hide" title="`)
//line index.html:180
	qw422016.N().S(ln.UI["sync"])
//line index.html:180
	qw422016.N().S(`"></b><b id="sync-counter" class="act hide-empty banner-float svg-link noscript-hide" title="`)
//line index.html:181
	qw422016.N().S(ln.UI["syncCount"])
//line index.html:181
	qw422016.N().S(`"></b><b id="thread-post-counters" class="act hide-empty banner-float svg-link noscript-hide" title="`)
//line index.html:182
	qw422016.N().S(ln.Common.UI["postsImages"])
//line index.html:182
	qw422016.N().S(`"></b><span id="banner-extensions" class="hide-empty banner-float svg-link noscript-hide"></span><a id="banner-feedback" href="mailto:`)
//line index.html:184
	qw422016.E().S(conf.FeedbackEmail)
//line index.html:184
	qw422016.N().S(`" target="_blank" class="banner-float svg-link noscript-hide" title="`)
//line index.html:184
	qw422016.N().S(ln.UI["feedback"])
//line index.html:184
	qw422016.N().S(`"><svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" viewBox="0 0 8 8"><path d="M0 0v1l4 2 4-2v-1h-8zm0 2v4h8v-4l-4 2-4-2z" transform="translate(0 1)" /></svg></a><a id="banner-FAQ" class="banner-float svg-link noscript-hide" title="`)
//line index.html:189
	qw422016.N().S(ln.UI["FAQ"])
//line index.html:189
	qw422016.N().S(`"><svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" viewBox="0 0 8 8"><path d="M3 0c-.55 0-1 .45-1 1s.45 1 1 1 1-.45 1-1-.45-1-1-1zm-1.5 2.5c-.83 0-1.5.67-1.5 1.5h1c0-.28.22-.5.5-.5s.5.22.5.5-1 1.64-1 2.5c0 .86.67 1.5 1.5 1.5s1.5-.67 1.5-1.5h-1c0 .28-.22.5-.5.5s-.5-.22-.5-.5c0-.36 1-1.84 1-2.5 0-.81-.67-1.5-1.5-1.5z" transform="translate(2)"/></svg></a><a id="banner-account" class="banner-float svg-link noscript-hide" title="`)
//line index.html:194
	qw422016.N().S(ln.UI["account"])
//line index.html:194
	qw422016.N().S(`"><svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" viewBox="0 0 8 8"><path d="m 2,2.681 c -1.31,0 -2,1.01 -2,2 0,0.99 0.69,2 2,2 0.79,0 1.42,-0.56 2,-1.22 0.58,0.66 1.19,1.22 2,1.22 1.31,0 2,-1.01 2,-2 0,-0.99 -0.69,-2 -2,-2 -0.81,0 -1.42,0.56 -2,1.22 C 3.42,3.241 2.79,2.681 2,2.681 Z m 0,1 c 0.42,0 0.88,0.47 1.34,1 -0.46,0.53 -0.92,1 -1.34,1 -0.74,0 -1,-0.54 -1,-1 0,-0.46 0.26,-1 1,-1 z m 4,0 c 0.74,0 1,0.54 1,1 0,0.46 -0.26,1 -1,1 -0.43,0 -0.89,-0.47 -1.34,-1 0.46,-0.53 0.91,-1 1.34,-1 z" id="path4" /></svg></a><a id="banner-identity" class="banner-float svg-link noscript-hide" title="`)
//line index.html:199
	qw422016.N().S(ln.UI["identity"])
//line index.html:199
	qw422016.N().S(`"><svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" viewBox="0 0 8 8"><path d="M4 0c-1.1 0-2 1.12-2 2.5s.9 2.5 2 2.5 2-1.12 2-2.5-.9-2.5-2-2.5zm-2.09 5c-1.06.05-1.91.92-1.91 2v1h8v-1c0-1.08-.84-1.95-1.91-2-.54.61-1.28 1-2.09 1-.81 0-1.55-.39-2.09-1z" /></svg></a><a id="banner-options" class="banner-float svg-link noscript-hide" title="`)
//line index.html:204
	qw422016.N().S(ln.UI["options"])
//line index.html:204
	qw422016.N().S(`"><svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" viewBox="0 0 8 8"><path d="M3.5 0l-.5 1.19c-.1.03-.19.08-.28.13l-1.19-.5-.72.72.5 1.19c-.05.1-.09.18-.13.28l-1.19.5v1l1.19.5c.04.1.08.18.13.28l-.5 1.19.72.72 1.19-.5c.09.04.18.09.28.13l.5 1.19h1l.5-1.19c.09-.04.19-.08.28-.13l1.19.5.72-.72-.5-1.19c.04-.09.09-.19.13-.28l1.19-.5v-1l-1.19-.5c-.03-.09-.08-.19-.13-.28l.5-1.19-.72-.72-1.19.5c-.09-.04-.19-.09-.28-.13l-.5-1.19h-1zm.5 2.5c.83 0 1.5.67 1.5 1.5s-.67 1.5-1.5 1.5-1.5-.67-1.5-1.5.67-1.5 1.5-1.5z"/></svg></a></span></span>`)
//line index.html:213
	qw422016.N().S(`<div id="modal-overlay" class="overlay">`)
//line index.html:217
	qw422016.N().S(`<div id="FAQ" class="modal glass">meguca is licensed under the`)
//line index.html:219
	qw422016.N().S(` `)
//line index.html:219
	qw422016.N().S(`<a href="https://www.gnu.org/licenses/agpl.html" target="_blank">GNU Affero General Public License</a><br>Source code repository:`)
//line index.html:224
	qw422016.N().S(` `)
//line index.html:224
	qw422016.N().S(`<a href="https://github.com/bakape/meguca" target="_blank">github.com/bakape/meguca</a><hr>`)
//line index.html:229
	qw422016.N().S(strings.Replace(conf.FAQ, "\n", "<br>", -1))
//line index.html:229
	qw422016.N().S(`</div>`)
//line index.html:233
	qw422016.N().S(`<div id="identity" class="modal glass">`)
//line index.html:235
	fields := specs["identity"]

//line index.html:236
	if pos > common.NotStaff {
//line index.html:237
		fields = make([]inputSpec, 1, len(fields)+1)

//line index.html:238
		fields[0] = staffTitleSpec

//line index.html:239
		fields = append(fields, specs["identity"]...)

//line index.html:240
	}
//line index.html:241
	streamtable(qw422016, fields)
//line index.html:241
	qw422016.N().S(`</div>`)
//line index.html:245
	qw422016.N().S(`<div id="account-panel" class="modal glass">`)
//line index.html:247
	if pos == common.NotLoggedIn {
//line index.html:247
		qw422016.N().S(`<div id="login-forms">`)
//line index.html:249
		f := ln.Forms

//line index.html:250
		streamtabButts(qw422016, []string{f["id"][0], f["register"][0]})
//line index.html:250
		qw422016.N().S(`<div class="tab-cont"><div class="tab-sel" data-id="0"><form id="login-form">`)
//line index.html:254
		streamtable(qw422016, specs["login"])
//line index.html:255
		streamcaptcha(qw422016, "all")
//line index.html:256
		streamsubmit(qw422016, false)
//line index.html:256
		qw422016.N().S(`</form></div><div data-id="1"><form id="registration-form">`)
//line index.html:261
		streamtable(qw422016, specs["register"])
//line index.html:262
		streamcaptcha(qw422016, "all")
//line index.html:263
		streamsubmit(qw422016, false)
//line index.html:263
		qw422016.N().S(`</form></div></div></div>`)
//line index.html:268
	} else {
//line index.html:268
		qw422016.N().S(`<div id="form-selection">`)
//line index.html:270
		for _, l := range [...]string{
			"logout", "logoutAll", "changePassword",
			"createBoard", "configureBoard", "deleteBoard",
			"assignStaff", "setBanners", "setLoading",
		} {
//line index.html:274
			qw422016.N().S(`<a id="`)
//line index.html:275
			qw422016.N().S(l)
//line index.html:275
			qw422016.N().S(`">`)
//line index.html:276
			qw422016.N().S(ln.UI[l])
//line index.html:276
			qw422016.N().S(`<br></a>`)
//line index.html:279
		}
//line index.html:280
		if pos == common.Admin {
//line index.html:280
			qw422016.N().S(`<a id="configureServer">`)
//line index.html:282
			qw422016.N().S(ln.UI["configureServer"])
//line index.html:282
			qw422016.N().S(`<br></a>`)
//line index.html:285
		}
//line index.html:285
		qw422016.N().S(`</div>`)
//line index.html:287
	}
//line index.html:287
	qw422016.N().S(`</div>`)
//line index.html:291
	qw422016.N().S(`<div id="options" class="modal glass">`)
//line index.html:293
	streamtabButts(qw422016, ln.Tabs)
//line index.html:293
	qw422016.N().S(`<div class="tab-cont">`)
//line index.html:295
	for i, sp := range optionSpecs {
//line index.html:295
		qw422016.N().S(`<div data-id="`)
//line index.html:296
		qw422016.N().D(i)
//line index.html:296
		qw422016.N().S(`"`)
//line index.html:296
		if i == 0 {
//line index.html:296
			qw422016.N().S(` `)
//line index.html:296
			qw422016.N().S(`class="tab-sel"`)
//line index.html:296
		}
//line index.html:296
		qw422016.N().S(`>`)
//line index.html:297
		streamoptions(qw422016, sp, ln)
//line index.html:301
		if i == 0 {
//line index.html:301
			qw422016.N().S(`<br><span class="spaced">`)
//line index.html:304
			for _, id := range [...]string{"export", "import", "hidden"} {
//line index.html:304
				qw422016.N().S(`<a id="`)
//line index.html:305
				qw422016.N().S(id)
//line index.html:305
				qw422016.N().S(`" title="`)
//line index.html:305
				qw422016.N().S(ln.Forms[id][1])
//line index.html:305
				qw422016.N().S(`">`)
//line index.html:306
				qw422016.N().S(ln.Forms[id][0])
//line index.html:306
				qw422016.N().S(`</a>`)
//line index.html:308
			}
//line index.html:308
			qw422016.N().S(`</span>`)
//line index.html:312
			qw422016.N().S(`<input type="file" id="importSettings" hidden>`)
//line index.html:314
		}
//line index.html:314
		qw422016.N().S(`</div>`)
//line index.html:316
	}
//line index.html:316
	qw422016.N().S(`</div></div>`)
//line index.html:319
	if pos > common.NotStaff {
//line index.html:319
		qw422016.N().S(`<div id="moderation-panel" class="modal glass"><form>`)
//line index.html:322
		if pos >= common.Moderator {
//line index.html:322
			qw422016.N().S(`<div id="ban-form" class="hidden">`)
//line index.html:324
			for _, id := range [...]string{"day", "hour", "minute"} {
//line index.html:324
				qw422016.N().S(`<input type="number" name="`)
//line index.html:325
				qw422016.N().S(id)
//line index.html:325
				qw422016.N().S(`" min="0" placeholder="`)
//line index.html:325
				qw422016.N().S(strings.Title(ln.Common.Plurals[id][1]))
//line index.html:325
				qw422016.N().S(`">`)
//line index.html:326
			}
//line index.html:326
			qw422016.N().S(`<br><input type="text" name="reason" required class="full-width" placeholder="`)
//line index.html:328
			qw422016.N().S(ln.Common.UI["reason"])
//line index.html:328
			qw422016.N().S(`" disabled><br>`)
//line index.html:330
			if pos == common.Admin {
//line index.html:330
				qw422016.N().S(`<label><input type="checkbox" name="global">`)
//line index.html:333
				qw422016.N().S(ln.UI["global"])
//line index.html:333
				qw422016.N().S(`</label>`)
//line index.html:335
			}
//line index.html:335
			qw422016.N().S(`</div>`)
//line index.html:337
		}
//line index.html:338
		if pos == common.Admin {
//line index.html:338
			qw422016.N().S(`<div id="purgePost-form" class="hidden"><input type="text" name="purge-reason" required class="full-width" placeholder="`)
//line index.html:340
			qw422016.N().S(ln.Common.UI["reason"])
//line index.html:340
			qw422016.N().S(`" disabled><br></div><div id="notification-form" class="hidden"><input type="text" name="notification" required class="full-width" placeholder="`)
//line index.html:344
			qw422016.N().S(ln.UI["text"])
//line index.html:344
			qw422016.N().S(`" style="min-width: 20em;" disabled><br></div>`)
//line index.html:347
		}
//line index.html:347
		qw422016.N().S(`<input type="checkbox" name="showCheckboxes"><select name="action">`)
//line index.html:350
		ids := append(make([]string, 0, 5), "deletePost", "deleteImage", "spoilerImage")

//line index.html:351
		if pos >= common.Moderator {
//line index.html:352
			ids = append(ids, "ban")

//line index.html:353
		}
//line index.html:354
		if pos == common.Admin {
//line index.html:355
			ids = append(ids, "purgePost", "notification")

//line index.html:356
		}
//line index.html:357
		for _, id := range ids {
//line index.html:357
			qw422016.N().S(`<option value="`)
//line index.html:358
			qw422016.N().S(id)
//line index.html:358
			qw422016.N().S(`">`)
//line index.html:359
			qw422016.N().S(ln.UI[id])
//line index.html:359
			qw422016.N().S(`</option>`)
//line index.html:361
		}
//line index.html:361
		qw422016.N().S(`</select><input type="button" value="`)
//line index.html:363
		qw422016.N().S(ln.UI["clear"])
//line index.html:363
		qw422016.N().S(`" name="clear">`)
//line index.html:364
		streamsubmit(qw422016, false)
//line index.html:364
		qw422016.N().S(`</form></div>`)
//line index.html:367
	}
//line index.html:367
	qw422016.N().S(`</div></div>`)
//line index.html:372
	qw422016.N().S(`<div class="overlay top-overlay" id="hover-overlay"></div><div id="captcha-overlay" class="overlay top-overlay"></div>`)
//line index.html:378
	qw422016.N().S(`<section id="threads">`)
//line index.html:382
	qw422016.N().S(`$$$</section>`)
//line index.html:387
	qw422016.N().S(`<script src="/assets/js/vendor/almond.js"></script><script id="lang-data" type="application/json">`)
//line index.html:390
	buf, _ := json.Marshal(ln.Common)

//line index.html:391
	qw422016.N().Z(buf)
//line index.html:391
	qw422016.N().S(`</script><script id="board-title-data" type="application/json">`)
//line index.html:394
	buf, _ = json.Marshal(config.GetBoardTitles())

//line index.html:395
	qw422016.N().Z(buf)
//line index.html:395
	qw422016.N().S(`</script><script src="/assets/js/scripts/loader.js"></script></body>`)
//line index.html:399
}

//line index.html:399
func writerenderIndex(qq422016 qtio422016.Writer, pos common.ModerationLevel) {
//line index.html:399
	qw422016 := qt422016.AcquireWriter(qq422016)
//line index.html:399
	streamrenderIndex(qw422016, pos)
//line index.html:399
	qt422016.ReleaseWriter(qw422016)
//line index.html:399
}

//line index.html:399
func renderIndex(pos common.ModerationLevel) string {
//line index.html:399
	qb422016 := qt422016.AcquireByteBuffer()
//line index.html:399
	writerenderIndex(qb422016, pos)
//line index.html:399
	qs422016 := string(qb422016.B)
//line index.html:399
	qt422016.ReleaseByteBuffer(qb422016)
//line index.html:399
	return qs422016
//line index.html:399
}