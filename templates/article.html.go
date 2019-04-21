// Code generated by qtc from "article.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line article.html:1
package templates

//line article.html:1
import "fmt"

//line article.html:2
import "strconv"

//line article.html:3
import "github.com/bakape/meguca/common"

//line article.html:4
import "github.com/bakape/meguca/lang"

//line article.html:5
import "github.com/bakape/meguca/imager/assets"

//line article.html:6
import "github.com/bakape/meguca/util"

//line article.html:8
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line article.html:8
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line article.html:8
func streamrenderArticle(qw422016 *qt422016.Writer, p common.Post, c articleContext) {
	//line article.html:9
	id := strconv.FormatUint(p.ID, 10)

	//line article.html:10
	ln := lang.Get()

	//line article.html:10
	qw422016.N().S(`<article id="p`)
	//line article.html:11
	qw422016.N().S(id)
	//line article.html:11
	qw422016.N().S(`"`)
	//line article.html:11
	qw422016.N().S(` `)
	//line article.html:11
	streampostClass(qw422016, p, c.op)
	//line article.html:11
	qw422016.N().S(`>`)
	//line article.html:12
	streamdeletedToggle(qw422016)
	//line article.html:12
	qw422016.N().S(`<header class="spaced"><input type="checkbox" class="mod-checkbox hidden">`)
	//line article.html:15
	streamrenderSticky(qw422016, c.sticky)
	//line article.html:16
	streamrenderLocked(qw422016, c.locked)
	//line article.html:17
	if c.subject != "" {
		//line article.html:18
		if c.board != "" {
			//line article.html:18
			qw422016.N().S(`<b class="board">/`)
			//line article.html:20
			qw422016.N().S(c.board)
			//line article.html:20
			qw422016.N().S(`/</b>`)
			//line article.html:22
		}
		//line article.html:22
		qw422016.N().S(`<h3>「`)
		//line article.html:24
		qw422016.E().S(c.subject)
		//line article.html:24
		qw422016.N().S(`」</h3>`)
		//line article.html:26
	}
	//line article.html:26
	qw422016.N().S(`<b class="name spaced`)
	//line article.html:27
	if p.Auth != common.NotStaff {
		//line article.html:27
		qw422016.N().S(` `)
		//line article.html:27
		qw422016.N().S(`admin`)
		//line article.html:27
	}
	//line article.html:27
	if p.Sage {
		//line article.html:27
		qw422016.N().S(` `)
		//line article.html:27
		qw422016.N().S(`sage`)
		//line article.html:27
	}
	//line article.html:27
	qw422016.N().S(`">`)
	//line article.html:28
	if p.Name != "" || p.Trip == "" {
		//line article.html:28
		qw422016.N().S(`<span>`)
		//line article.html:30
		if p.Name != "" {
			//line article.html:31
			qw422016.E().S(p.Name)
			//line article.html:32
		} else {
			//line article.html:33
			qw422016.N().S(ln.Common.Posts["anon"])
			//line article.html:34
		}
		//line article.html:34
		qw422016.N().S(`</span>`)
		//line article.html:36
	}
	//line article.html:37
	if p.Trip != "" {
		//line article.html:37
		qw422016.N().S(`<code>!`)
		//line article.html:39
		qw422016.E().S(p.Trip)
		//line article.html:39
		qw422016.N().S(`</code>`)
		//line article.html:41
	}
	//line article.html:42
	if p.Auth != common.NotStaff {
		//line article.html:42
		qw422016.N().S(`<span>##`)
		//line article.html:44
		qw422016.N().S(` `)
		//line article.html:44
		qw422016.N().S(ln.Common.Posts[p.Auth.String()])
		//line article.html:44
		qw422016.N().S(`</span>`)
		//line article.html:46
	}
	//line article.html:46
	qw422016.N().S(`</b>`)
	//line article.html:48
	if p.Flag != "" {
		//line article.html:49
		title, ok := countryMap[p.Flag]

		//line article.html:50
		if !ok {
			//line article.html:51
			title = p.Flag

			//line article.html:52
		}
		//line article.html:52
		qw422016.N().S(`<img class="flag" src="/assets/flags/`)
		//line article.html:53
		qw422016.N().S(p.Flag)
		//line article.html:53
		qw422016.N().S(`.svg" title="`)
		//line article.html:53
		qw422016.N().S(title)
		//line article.html:53
		qw422016.N().S(`">`)
		//line article.html:54
	}
	//line article.html:54
	qw422016.N().S(`<time>`)
	//line article.html:56
	qw422016.N().S(formatTime(p.Time))
	//line article.html:56
	qw422016.N().S(`</time><nav>`)
	//line article.html:59
	url := "#p" + id

	//line article.html:60
	if c.index {
		//line article.html:61
		url = util.ConcatStrings("/all/", id, "?last=100", url)

		//line article.html:62
	}
	//line article.html:62
	qw422016.N().S(`<a href="`)
	//line article.html:63
	qw422016.N().S(url)
	//line article.html:63
	qw422016.N().S(`">No.</a><a class="quote" href="`)
	//line article.html:66
	qw422016.N().S(url)
	//line article.html:66
	qw422016.N().S(`">`)
	//line article.html:67
	qw422016.N().S(id)
	//line article.html:67
	qw422016.N().S(`</a></nav>`)
	//line article.html:70
	if c.index && c.subject != "" {
		//line article.html:70
		qw422016.N().S(`<span>`)
		//line article.html:72
		streamexpandLink(qw422016, "all", id)
		//line article.html:73
		streamlast100Link(qw422016, "all", id)
		//line article.html:73
		qw422016.N().S(`</span>`)
		//line article.html:75
	}
	//line article.html:76
	streamcontrolLink(qw422016)
	//line article.html:77
	if c.op == p.ID {
		//line article.html:78
		streamthreadWatcherToggle(qw422016, p.ID)
		//line article.html:79
	}
	//line article.html:79
	qw422016.N().S(`</header>`)
	//line article.html:81
	var src string

	//line article.html:82
	if p.Image != nil {
		//line article.html:83
		img := *p.Image

		//line article.html:84
		src = assets.SourcePath(img.FileType, img.SHA1)

		//line article.html:84
		qw422016.N().S(`<figcaption class="spaced"><a class="image-toggle act" hidden></a><span class="spaced image-search-container">`)
		//line article.html:88
		streamimageSearch(qw422016, c.root, img)
		//line article.html:88
		qw422016.N().S(`</span><span class="fileinfo">`)
		//line article.html:91
		if img.Artist != "" {
			//line article.html:91
			qw422016.N().S(`<span class="media-artist">`)
			//line article.html:93
			qw422016.E().S(img.Artist)
			//line article.html:93
			qw422016.N().S(`</span>`)
			//line article.html:95
		}
		//line article.html:96
		if img.Title != "" {
			//line article.html:96
			qw422016.N().S(`<span class="media-title">`)
			//line article.html:98
			qw422016.E().S(img.Title)
			//line article.html:98
			qw422016.N().S(`</span>`)
			//line article.html:100
		}
		//line article.html:101
		if img.Audio {
			//line article.html:101
			qw422016.N().S(`<span class="has-audio">♫</span>`)
			//line article.html:105
		}
		//line article.html:106
		if img.Length != 0 {
			//line article.html:106
			qw422016.N().S(`<span class="media-length">`)
			//line article.html:108
			l := img.Length

			//line article.html:109
			if l < 60 {
				//line article.html:110
				qw422016.N().S(fmt.Sprintf("0:%02d", l))
				//line article.html:111
			} else {
				//line article.html:112
				min := l / 60

				//line article.html:113
				qw422016.N().S(fmt.Sprintf("%02d:%02d", min, l-min*60))
				//line article.html:114
			}
			//line article.html:114
			qw422016.N().S(`</span>`)
			//line article.html:116
		}
		//line article.html:116
		qw422016.N().S(`<span class="filesize">`)
		//line article.html:118
		qw422016.N().S(readableFileSize(img.Size))
		//line article.html:118
		qw422016.N().S(`</span>`)
		//line article.html:120
		if img.Dims != [4]uint16{} {
			//line article.html:120
			qw422016.N().S(`<span class="dims">`)
			//line article.html:122
			qw422016.N().S(strconv.FormatUint(uint64(img.Dims[0]), 10))
			//line article.html:122
			qw422016.N().S(`x`)
			//line article.html:124
			qw422016.N().S(strconv.FormatUint(uint64(img.Dims[1]), 10))
			//line article.html:124
			qw422016.N().S(`</span>`)
			//line article.html:126
		}
		//line article.html:126
		qw422016.N().S(`</span>`)
		//line article.html:128
		name := imageName(img.FileType, img.Name)

		//line article.html:128
		qw422016.N().S(`<a href="`)
		//line article.html:129
		qw422016.N().S(assets.RelativeSourcePath(img.FileType, img.SHA1))
		//line article.html:129
		qw422016.N().S(`" download="`)
		//line article.html:129
		qw422016.N().S(name)
		//line article.html:129
		qw422016.N().S(`">`)
		//line article.html:130
		qw422016.N().S(name)
		//line article.html:130
		qw422016.N().S(`</a></figcaption>`)
		//line article.html:133
	}
	//line article.html:133
	qw422016.N().S(`<div class="post-container">`)
	//line article.html:135
	if p.Image != nil {
		//line article.html:136
		img := *p.Image

		//line article.html:136
		qw422016.N().S(`<figure><a target="_blank" href="`)
		//line article.html:138
		qw422016.N().S(src)
		//line article.html:138
		qw422016.N().S(`">`)
		//line article.html:139
		switch {
		//line article.html:140
		case img.ThumbType == common.NoFile:
			//line article.html:141
			var file string

			//line article.html:142
			switch img.FileType {
			//line article.html:143
			case common.MP4, common.MP3, common.OGG, common.FLAC:
				//line article.html:144
				file = "audio"

			//line article.html:145
			default:
				//line article.html:146
				file = "file"

				//line article.html:147
			}
			//line article.html:147
			qw422016.N().S(`<img src="/assets/`)
			//line article.html:148
			qw422016.N().S(file)
			//line article.html:148
			qw422016.N().S(`.png" width="150" height="150">`)
		//line article.html:149
		case img.Spoiler:
			//line article.html:152
			qw422016.N().S(`<img src="/assets/spoil/default.jpg" width="150" height="150">`)
		//line article.html:154
		default:
			//line article.html:154
			qw422016.N().S(`<img src="`)
			//line article.html:155
			qw422016.N().S(assets.ThumbPath(img.ThumbType, img.SHA1))
			//line article.html:155
			qw422016.N().S(`" width="`)
			//line article.html:155
			qw422016.N().D(int(img.Dims[2]))
			//line article.html:155
			qw422016.N().S(`" height="`)
			//line article.html:155
			qw422016.N().D(int(img.Dims[3]))
			//line article.html:155
			qw422016.N().S(`">`)
			//line article.html:156
		}
		//line article.html:156
		qw422016.N().S(`</a></figure>`)
		//line article.html:159
	}
	//line article.html:159
	qw422016.N().S(`<blockquote>`)
	//line article.html:161
	streambody(qw422016, p, c.op, c.board, c.index, c.rbText, c.pyu)
	//line article.html:161
	qw422016.N().S(`</blockquote>`)
	//line article.html:163
	for _, e := range p.Moderation {
		//line article.html:163
		qw422016.N().S(`<b class="admin post-moderation">`)
		//line article.html:165
		streampostModeration(qw422016, e)
		//line article.html:165
		qw422016.N().S(`<br></b>`)
		//line article.html:168
	}
	//line article.html:168
	qw422016.N().S(`</div>`)
	//line article.html:170
	if c.omit != 0 {
		//line article.html:170
		qw422016.N().S(`<span class="omit spaced" data-omit="`)
		//line article.html:171
		qw422016.N().D(c.omit)
		//line article.html:171
		qw422016.N().S(`" data-image-omit="`)
		//line article.html:171
		qw422016.N().D(c.imageOmit)
		//line article.html:171
		qw422016.N().S(`">`)
		//line article.html:172
		if c.imageOmit == 0 {
			//line article.html:173
			qw422016.N().S(fmt.Sprintf(ln.Common.Format["postsOmitted"], c.omit))
			//line article.html:174
		} else {
			//line article.html:175
			qw422016.N().S(fmt.Sprintf(ln.Common.Format["postsAndImagesOmitted"], c.omit, c.imageOmit))
			//line article.html:176
		}
		//line article.html:176
		qw422016.N().S(`<span class="act"><a href="`)
		//line article.html:178
		qw422016.N().S(strconv.FormatUint(c.op, 10))
		//line article.html:178
		qw422016.N().S(`">`)
		//line article.html:179
		qw422016.N().S(ln.Common.Posts["seeAll"])
		//line article.html:179
		qw422016.N().S(`</a></span></span>`)
		//line article.html:183
	}
	//line article.html:184
	if bls := c.backlinks[p.ID]; len(bls) != 0 {
		//line article.html:184
		qw422016.N().S(`<span class="backlinks spaced">`)
		//line article.html:186
		for _, l := range bls {
			//line article.html:186
			qw422016.N().S(`<em>`)
			//line article.html:188
			streampostLink(qw422016, l, c.index || l.OP != c.op, c.index)
			//line article.html:188
			qw422016.N().S(`</em>`)
			//line article.html:190
		}
		//line article.html:190
		qw422016.N().S(`</span>`)
		//line article.html:192
	}
	//line article.html:192
	qw422016.N().S(`</article>`)
//line article.html:194
}

//line article.html:194
func writerenderArticle(qq422016 qtio422016.Writer, p common.Post, c articleContext) {
	//line article.html:194
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line article.html:194
	streamrenderArticle(qw422016, p, c)
	//line article.html:194
	qt422016.ReleaseWriter(qw422016)
//line article.html:194
}

//line article.html:194
func renderArticle(p common.Post, c articleContext) string {
	//line article.html:194
	qb422016 := qt422016.AcquireByteBuffer()
	//line article.html:194
	writerenderArticle(qb422016, p, c)
	//line article.html:194
	qs422016 := string(qb422016.B)
	//line article.html:194
	qt422016.ReleaseByteBuffer(qb422016)
	//line article.html:194
	return qs422016
//line article.html:194
}

// Render image search links according to file type

//line article.html:197
func streamimageSearch(qw422016 *qt422016.Writer, root string, img common.Image) {
	//line article.html:198
	if img.ThumbType == common.NoFile || img.FileType == common.PDF {
		//line article.html:199
		return
		//line article.html:200
	}
	//line article.html:202
	url := root + assets.ImageSearchPath(img.ImageCommon)

	//line article.html:202
	qw422016.N().S(`<a class="image-search google" target="_blank" rel="nofollow" href="https://www.google.com/searchbyimage?image_url=`)
	//line article.html:203
	qw422016.N().S(url)
	//line article.html:203
	qw422016.N().S(`">G</a><a class="image-search yandex" target="_blank" rel="nofollow" href="https://yandex.com/images/search?source=collections&rpt=imageview&url=`)
	//line article.html:206
	qw422016.N().S(url)
	//line article.html:206
	qw422016.N().S(`">Yd</a><a class="image-search iqdb" target="_blank" rel="nofollow" href="http://iqdb.org/?url=`)
	//line article.html:209
	qw422016.N().S(url)
	//line article.html:209
	qw422016.N().S(`">Iq</a><a class="image-search saucenao" target="_blank" rel="nofollow" href="http://saucenao.com/search.php?db=999&url=`)
	//line article.html:212
	qw422016.N().S(url)
	//line article.html:212
	qw422016.N().S(`">Sn</a><a class="image-search whatAnime" target="_blank" rel="nofollow" href="https://trace.moe/?url=`)
	//line article.html:215
	qw422016.N().S(url)
	//line article.html:215
	qw422016.N().S(`">Wa</a>`)
	//line article.html:218
	switch img.FileType {
	//line article.html:219
	case common.JPEG, common.PNG, common.GIF, common.WEBM:
		//line article.html:219
		qw422016.N().S(`<a class="image-search desustorage" target="_blank" rel="nofollow" href="https://desuarchive.org/_/search/image/`)
		//line article.html:220
		qw422016.N().S(img.MD5)
		//line article.html:220
		qw422016.N().S(`">Ds</a>`)
		//line article.html:223
	}
	//line article.html:224
	switch img.FileType {
	//line article.html:225
	case common.JPEG, common.PNG:
		//line article.html:225
		qw422016.N().S(`<a class="image-search exhentai" target="_blank" rel="nofollow" href="http://exhentai.org/?fs_similar=1&fs_exp=1&f_shash=`)
		//line article.html:226
		qw422016.N().S(img.SHA1)
		//line article.html:226
		qw422016.N().S(`">Ex</a>`)
		//line article.html:229
	}
//line article.html:230
}

//line article.html:230
func writeimageSearch(qq422016 qtio422016.Writer, root string, img common.Image) {
	//line article.html:230
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line article.html:230
	streamimageSearch(qw422016, root, img)
	//line article.html:230
	qt422016.ReleaseWriter(qw422016)
//line article.html:230
}

//line article.html:230
func imageSearch(root string, img common.Image) string {
	//line article.html:230
	qb422016 := qt422016.AcquireByteBuffer()
	//line article.html:230
	writeimageSearch(qb422016, root, img)
	//line article.html:230
	qs422016 := string(qb422016.B)
	//line article.html:230
	qt422016.ReleaseByteBuffer(qb422016)
	//line article.html:230
	return qs422016
//line article.html:230
}
