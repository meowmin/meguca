// Code generated by qtc from "article.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line article.html:1
package templates

//line article.html:1
import "fmt"

//line article.html:2
import "strconv"

//line article.html:3
import "strings"

//line article.html:4
import "github.com/bakape/meguca/common"

//line article.html:5
import "github.com/bakape/meguca/lang"

//line article.html:6
import "github.com/bakape/meguca/imager/assets"

//line article.html:7
import "github.com/bakape/meguca/util"

//line article.html:9
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line article.html:9
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line article.html:9
func streamrenderArticle(qw422016 *qt422016.Writer, p common.Post, c articleContext) {
//line article.html:10
	id := strconv.FormatUint(p.ID, 10)

//line article.html:11
	ln := lang.Get()

//line article.html:11
	qw422016.N().S(`<article id="p`)
//line article.html:12
	qw422016.N().S(id)
//line article.html:12
	qw422016.N().S(`"`)
//line article.html:12
	qw422016.N().S(` `)
//line article.html:12
	streampostClass(qw422016, p, c.op)
//line article.html:12
	qw422016.N().S(`>`)
//line article.html:13
	streamdeletedToggle(qw422016)
//line article.html:13
	qw422016.N().S(`<header class="spaced"><input type="radio" name="mod-checkbox" class="mod-checkbox hidden">`)
//line article.html:16
	streamrenderSticky(qw422016, c.sticky)
//line article.html:17
	streamrenderLocked(qw422016, c.locked)
//line article.html:18
	if c.subject != "" {
//line article.html:19
		if c.board != "" {
//line article.html:19
			qw422016.N().S(`<b class="board">/`)
//line article.html:21
			qw422016.N().S(c.board)
//line article.html:21
			qw422016.N().S(`/</b>`)
//line article.html:23
		}
//line article.html:23
		qw422016.N().S(`<h3>「`)
//line article.html:25
		qw422016.E().S(c.subject)
//line article.html:25
		qw422016.N().S(`」</h3>`)
//line article.html:27
	}
//line article.html:27
	qw422016.N().S(`<b class="name spaced`)
//line article.html:28
	if p.Auth != common.NotStaff {
//line article.html:28
		qw422016.N().S(` `)
//line article.html:28
		qw422016.N().S(`admin`)
//line article.html:28
	}
//line article.html:28
	if p.Sage {
//line article.html:28
		qw422016.N().S(` `)
//line article.html:28
		qw422016.N().S(`sage`)
//line article.html:28
	}
//line article.html:28
	qw422016.N().S(`">`)
//line article.html:29
	if p.Name != "" || p.Trip == "" {
//line article.html:29
		qw422016.N().S(`<span>`)
//line article.html:31
		if p.Name != "" {
//line article.html:32
			qw422016.E().S(p.Name)
//line article.html:33
		} else {
//line article.html:34
			qw422016.N().S(ln.Common.Posts["anon"])
//line article.html:35
		}
//line article.html:35
		qw422016.N().S(`</span>`)
//line article.html:37
	}
//line article.html:38
	if p.Trip != "" {
//line article.html:38
		qw422016.N().S(`<code>!`)
//line article.html:40
		qw422016.E().S(p.Trip)
//line article.html:40
		qw422016.N().S(`</code>`)
//line article.html:42
	}
//line article.html:43
	if p.Auth != common.NotStaff {
//line article.html:43
		qw422016.N().S(`<span>##`)
//line article.html:45
		qw422016.N().S(` `)
//line article.html:45
		qw422016.N().S(ln.Common.Posts[p.Auth.String()])
//line article.html:45
		qw422016.N().S(`</span>`)
//line article.html:47
	}
//line article.html:47
	qw422016.N().S(`</b>`)
//line article.html:49
	if p.Flag != "" {
//line article.html:50
		title, ok := countryMap[p.Flag]

//line article.html:51
		if !ok {
//line article.html:52
			title = p.Flag

//line article.html:53
		}
//line article.html:54
		if strings.HasPrefix(p.Flag, "us-") {
//line article.html:55
			title2, ok2 := countryMap["us"]

//line article.html:56
			if !ok2 {
//line article.html:57
				title2 = "us"

//line article.html:58
			}
//line article.html:58
			qw422016.N().S(`<img class="flag" src="/assets/flags/us.svg" title="`)
//line article.html:59
			qw422016.N().S(title2)
//line article.html:59
			qw422016.N().S(`">`)
//line article.html:60
		}
//line article.html:60
		qw422016.N().S(`<img class="flag" src="/assets/flags/`)
//line article.html:61
		qw422016.N().S(p.Flag)
//line article.html:61
		qw422016.N().S(`.svg" title="`)
//line article.html:61
		qw422016.N().S(title)
//line article.html:61
		qw422016.N().S(`">`)
//line article.html:62
	}
//line article.html:62
	qw422016.N().S(`<time>`)
//line article.html:64
	qw422016.N().S(formatTime(p.Time))
//line article.html:64
	qw422016.N().S(`</time><nav>`)
//line article.html:67
	url := "#p" + id

//line article.html:68
	if c.index {
//line article.html:69
		url = util.ConcatStrings("/all/", id, "?last=100", url)

//line article.html:70
	}
//line article.html:70
	qw422016.N().S(`<a href="`)
//line article.html:71
	qw422016.N().S(url)
//line article.html:71
	qw422016.N().S(`">No.</a><a class="quote">`)
//line article.html:75
	qw422016.N().S(id)
//line article.html:75
	qw422016.N().S(`</a></nav>`)
//line article.html:78
	if c.index && c.subject != "" {
//line article.html:78
		qw422016.N().S(`<span>`)
//line article.html:80
		streamexpandLink(qw422016, "all", id)
//line article.html:81
		streamlast100Link(qw422016, "all", id)
//line article.html:81
		qw422016.N().S(`</span>`)
//line article.html:83
	}
//line article.html:84
	streamcontrolLink(qw422016)
//line article.html:85
	if c.op == p.ID {
//line article.html:86
		streamthreadWatcherToggle(qw422016, p.ID)
//line article.html:87
	}
//line article.html:87
	qw422016.N().S(`</header>`)
//line article.html:89
	var src string

//line article.html:90
	if p.Image != nil {
//line article.html:91
		img := *p.Image

//line article.html:92
		src = assets.SourcePath(img.FileType, img.SHA1)

//line article.html:92
		qw422016.N().S(`<figcaption class="spaced"><a class="image-toggle act" hidden></a><span class="spaced image-search-container">`)
//line article.html:96
		streamimageSearch(qw422016, c.root, img)
//line article.html:96
		qw422016.N().S(`</span><span class="fileinfo">`)
//line article.html:99
		if img.Audio {
//line article.html:99
			qw422016.N().S(`<span>♫</span>`)
//line article.html:103
		}
//line article.html:104
		if img.Length != 0 {
//line article.html:104
			qw422016.N().S(`<span>`)
//line article.html:106
			l := img.Length

//line article.html:107
			if l < 60 {
//line article.html:108
				qw422016.N().S(fmt.Sprintf("0:%02d", l))
//line article.html:109
			} else {
//line article.html:110
				min := l / 60

//line article.html:111
				qw422016.N().S(fmt.Sprintf("%02d:%02d", min, l-min*60))
//line article.html:112
			}
//line article.html:112
			qw422016.N().S(`</span>`)
//line article.html:114
		}
//line article.html:114
		qw422016.N().S(`<span>`)
//line article.html:116
		qw422016.N().S(readableFileSize(img.Size))
//line article.html:116
		qw422016.N().S(`</span>`)
//line article.html:118
		if img.Dims != [4]uint16{} {
//line article.html:118
			qw422016.N().S(`<span>`)
//line article.html:120
			qw422016.N().S(strconv.FormatUint(uint64(img.Dims[0]), 10))
//line article.html:120
			qw422016.N().S(`x`)
//line article.html:122
			qw422016.N().S(strconv.FormatUint(uint64(img.Dims[1]), 10))
//line article.html:122
			qw422016.N().S(`</span>`)
//line article.html:124
		}
//line article.html:125
		if img.Artist != "" {
//line article.html:125
			qw422016.N().S(`<span>`)
//line article.html:127
			qw422016.E().S(img.Artist)
//line article.html:127
			qw422016.N().S(`</span>`)
//line article.html:129
		}
//line article.html:130
		if img.Title != "" {
//line article.html:130
			qw422016.N().S(`<span>`)
//line article.html:132
			qw422016.E().S(img.Title)
//line article.html:132
			qw422016.N().S(`</span>`)
//line article.html:134
		}
//line article.html:134
		qw422016.N().S(`</span>`)
//line article.html:136
		name := imageName(img.FileType, img.Name)

//line article.html:136
		qw422016.N().S(`<a href="`)
//line article.html:137
		qw422016.N().S(assets.RelativeSourcePath(img.FileType, img.SHA1))
//line article.html:137
		qw422016.N().S(`" download="`)
//line article.html:137
		qw422016.N().S(name)
//line article.html:137
		qw422016.N().S(`">`)
//line article.html:138
		qw422016.N().S(name)
//line article.html:138
		qw422016.N().S(`</a>`)
//line article.html:140
		tokid := getTokID(name)

//line article.html:141
		if tokid != nil {
//line article.html:142
			timeago := "Posted " + relativeTime(*tokid)

//line article.html:142
			qw422016.N().S(`<a class="sourcelink" href="https://www.tiktok.com/share/video/`)
//line article.html:143
			qw422016.N().S(*tokid)
//line article.html:143
			qw422016.N().S(`" title="`)
//line article.html:143
			qw422016.N().S(timeago)
//line article.html:143
			qw422016.N().S(`" target="_blank" rel="noopener noreferrer">source</a>`)
//line article.html:144
		}
//line article.html:144
		qw422016.N().S(`</figcaption>`)
//line article.html:146
	}
//line article.html:146
	qw422016.N().S(`<div class="post-container">`)
//line article.html:148
	if p.Image != nil {
//line article.html:149
		img := *p.Image

//line article.html:149
		qw422016.N().S(`<figure><a target="_blank" href="`)
//line article.html:151
		qw422016.N().S(src)
//line article.html:151
		qw422016.N().S(`">`)
//line article.html:152
		switch {
//line article.html:153
		case img.ThumbType == common.NoFile:
//line article.html:154
			var file string

//line article.html:155
			switch img.FileType {
//line article.html:156
			case common.WEBM, common.MP4, common.MP3, common.OGG, common.FLAC:
//line article.html:157
				file = "audio"

//line article.html:158
			default:
//line article.html:159
				file = "file"

//line article.html:160
			}
//line article.html:160
			qw422016.N().S(`<img src="/assets/`)
//line article.html:161
			qw422016.N().S(file)
//line article.html:161
			qw422016.N().S(`.png" width="150" height="150" loading="lazy">`)
//line article.html:162
		case img.Spoiler:
//line article.html:165
			qw422016.N().S(`<img src="/assets/spoil/default.jpg" width="150" height="150" loading="lazy">`)
//line article.html:167
		default:
//line article.html:167
			qw422016.N().S(`<img src="`)
//line article.html:168
			qw422016.N().S(assets.ThumbPath(img.ThumbType, img.SHA1))
//line article.html:168
			qw422016.N().S(`" width="`)
//line article.html:168
			qw422016.N().D(int(img.Dims[2]))
//line article.html:168
			qw422016.N().S(`" height="`)
//line article.html:168
			qw422016.N().D(int(img.Dims[3]))
//line article.html:168
			qw422016.N().S(`" loading="lazy">`)
//line article.html:169
		}
//line article.html:169
		qw422016.N().S(`</a></figure>`)
//line article.html:172
	}
//line article.html:172
	qw422016.N().S(`<blockquote>`)
//line article.html:174
	streambody(qw422016, p, c.op, c.board, c.index, c.rbText, c.pyu)
//line article.html:174
	qw422016.N().S(`</blockquote>`)
//line article.html:176
	for _, e := range p.Moderation {
//line article.html:176
		qw422016.N().S(`<b class="admin post-moderation">`)
//line article.html:178
		streampostModeration(qw422016, e)
//line article.html:178
		qw422016.N().S(`<br></b>`)
//line article.html:181
	}
//line article.html:181
	qw422016.N().S(`</div>`)
//line article.html:183
	if c.omit != 0 {
//line article.html:183
		qw422016.N().S(`<span class="omit spaced" data-omit="`)
//line article.html:184
		qw422016.N().D(c.omit)
//line article.html:184
		qw422016.N().S(`" data-image-omit="`)
//line article.html:184
		qw422016.N().D(c.imageOmit)
//line article.html:184
		qw422016.N().S(`">`)
//line article.html:185
		if c.imageOmit == 0 {
//line article.html:186
			qw422016.N().S(fmt.Sprintf(ln.Common.Format["postsOmitted"], c.omit))
//line article.html:187
		} else {
//line article.html:188
			qw422016.N().S(fmt.Sprintf(ln.Common.Format["postsAndImagesOmitted"], c.omit, c.imageOmit))
//line article.html:189
		}
//line article.html:189
		qw422016.N().S(`<span class="act"><a href="`)
//line article.html:191
		qw422016.N().S(strconv.FormatUint(c.op, 10))
//line article.html:191
		qw422016.N().S(`">`)
//line article.html:192
		qw422016.N().S(ln.Common.Posts["seeAll"])
//line article.html:192
		qw422016.N().S(`</a></span></span>`)
//line article.html:196
	}
//line article.html:197
	if bls := c.backlinks[p.ID]; len(bls) != 0 {
//line article.html:197
		qw422016.N().S(`<span class="backlinks spaced">`)
//line article.html:199
		for _, l := range bls {
//line article.html:199
			qw422016.N().S(`<em>`)
//line article.html:201
			streampostLink(qw422016, l, c.index || l.OP != c.op, c.index)
//line article.html:201
			qw422016.N().S(`</em>`)
//line article.html:203
		}
//line article.html:203
		qw422016.N().S(`</span>`)
//line article.html:205
	}
//line article.html:205
	qw422016.N().S(`</article>`)
//line article.html:207
}

//line article.html:207
func writerenderArticle(qq422016 qtio422016.Writer, p common.Post, c articleContext) {
//line article.html:207
	qw422016 := qt422016.AcquireWriter(qq422016)
//line article.html:207
	streamrenderArticle(qw422016, p, c)
//line article.html:207
	qt422016.ReleaseWriter(qw422016)
//line article.html:207
}

//line article.html:207
func renderArticle(p common.Post, c articleContext) string {
//line article.html:207
	qb422016 := qt422016.AcquireByteBuffer()
//line article.html:207
	writerenderArticle(qb422016, p, c)
//line article.html:207
	qs422016 := string(qb422016.B)
//line article.html:207
	qt422016.ReleaseByteBuffer(qb422016)
//line article.html:207
	return qs422016
//line article.html:207
}

// Render image search links according to file type

//line article.html:210
func streamimageSearch(qw422016 *qt422016.Writer, root string, img common.Image) {
//line article.html:211
	if img.ThumbType == common.NoFile || img.FileType == common.PDF {
//line article.html:212
		return
//line article.html:213
	}
//line article.html:215
	url := root + assets.ImageSearchPath(img.ImageCommon)

//line article.html:215
	qw422016.N().S(`<a class="image-search google" target="_blank" rel="nofollow" href="https://www.google.com/searchbyimage?image_url=`)
//line article.html:216
	qw422016.N().S(url)
//line article.html:216
	qw422016.N().S(`">G</a><a class="image-search yandex" target="_blank" rel="nofollow" href="https://yandex.com/images/search?source=collections&rpt=imageview&url=`)
//line article.html:219
	qw422016.N().S(url)
//line article.html:219
	qw422016.N().S(`">Yd</a><a class="image-search iqdb" target="_blank" rel="nofollow" href="http://iqdb.org/?url=`)
//line article.html:222
	qw422016.N().S(url)
//line article.html:222
	qw422016.N().S(`">Iq</a><a class="image-search saucenao" target="_blank" rel="nofollow" href="http://saucenao.com/search.php?db=999&url=`)
//line article.html:225
	qw422016.N().S(url)
//line article.html:225
	qw422016.N().S(`">Sn</a><a class="image-search tracemoe" target="_blank" rel="nofollow" href="https://trace.moe/?url=`)
//line article.html:228
	qw422016.N().S(url)
//line article.html:228
	qw422016.N().S(`">Tm</a>`)
//line article.html:231
	switch img.FileType {
//line article.html:232
	case common.JPEG, common.PNG, common.GIF, common.WEBM:
//line article.html:232
		qw422016.N().S(`<a class="image-search desuarchive" target="_blank" rel="nofollow" href="https://desuarchive.org/_/search/image/`)
//line article.html:233
		qw422016.N().S(img.MD5)
//line article.html:233
		qw422016.N().S(`">Da</a>`)
//line article.html:236
	}
//line article.html:237
	switch img.FileType {
//line article.html:238
	case common.JPEG, common.PNG:
//line article.html:238
		qw422016.N().S(`<a class="image-search exhentai" target="_blank" rel="nofollow" href="http://exhentai.org/?fs_similar=1&fs_exp=1&f_shash=`)
//line article.html:239
		qw422016.N().S(img.SHA1)
//line article.html:239
		qw422016.N().S(`">Ex</a>`)
//line article.html:242
	}
//line article.html:243
}

//line article.html:243
func writeimageSearch(qq422016 qtio422016.Writer, root string, img common.Image) {
//line article.html:243
	qw422016 := qt422016.AcquireWriter(qq422016)
//line article.html:243
	streamimageSearch(qw422016, root, img)
//line article.html:243
	qt422016.ReleaseWriter(qw422016)
//line article.html:243
}

//line article.html:243
func imageSearch(root string, img common.Image) string {
//line article.html:243
	qb422016 := qt422016.AcquireByteBuffer()
//line article.html:243
	writeimageSearch(qb422016, root, img)
//line article.html:243
	qs422016 := string(qb422016.B)
//line article.html:243
	qt422016.ReleaseByteBuffer(qb422016)
//line article.html:243
	return qs422016
//line article.html:243
}
