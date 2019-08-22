// Code generated by qtc from "auth.html". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line auth.html:1
package templates

//line auth.html:1
import "fmt"

//line auth.html:2
import "time"

//line auth.html:3
import "strconv"

//line auth.html:4
import "github.com/bakape/meguca/auth"

//line auth.html:5
import "github.com/bakape/meguca/config"

//line auth.html:6
import "github.com/bakape/meguca/lang"

//line auth.html:7
import "github.com/bakape/meguca/common"

//line auth.html:8
import "github.com/bakape/mnemonics"

// Header of a standalone HTML page

//line auth.html:11
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line auth.html:11
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line auth.html:11
func streamhtmlHeader(qw422016 *qt422016.Writer) {
//line auth.html:11
	qw422016.N().S(`<!DOCTYPE html><html><head><meta charset="utf-8"/></head><body>`)
//line auth.html:18
}

//line auth.html:18
func writehtmlHeader(qq422016 qtio422016.Writer) {
//line auth.html:18
	qw422016 := qt422016.AcquireWriter(qq422016)
//line auth.html:18
	streamhtmlHeader(qw422016)
//line auth.html:18
	qt422016.ReleaseWriter(qw422016)
//line auth.html:18
}

//line auth.html:18
func htmlHeader() string {
//line auth.html:18
	qb422016 := qt422016.AcquireByteBuffer()
//line auth.html:18
	writehtmlHeader(qb422016)
//line auth.html:18
	qs422016 := string(qb422016.B)
//line auth.html:18
	qt422016.ReleaseByteBuffer(qb422016)
//line auth.html:18
	return qs422016
//line auth.html:18
}

// End of a standalone HTML page

//line auth.html:21
func streamhtmlEnd(qw422016 *qt422016.Writer) {
//line auth.html:21
	qw422016.N().S(`</body></html>`)
//line auth.html:24
}

//line auth.html:24
func writehtmlEnd(qq422016 qtio422016.Writer) {
//line auth.html:24
	qw422016 := qt422016.AcquireWriter(qq422016)
//line auth.html:24
	streamhtmlEnd(qw422016)
//line auth.html:24
	qt422016.ReleaseWriter(qw422016)
//line auth.html:24
}

//line auth.html:24
func htmlEnd() string {
//line auth.html:24
	qb422016 := qt422016.AcquireByteBuffer()
//line auth.html:24
	writehtmlEnd(qb422016)
//line auth.html:24
	qs422016 := string(qb422016.B)
//line auth.html:24
	qt422016.ReleaseByteBuffer(qb422016)
//line auth.html:24
	return qs422016
//line auth.html:24
}

// BanPage renders a ban page for a banned user

//line auth.html:27
func StreamBanPage(qw422016 *qt422016.Writer, rec auth.BanRecord) {
//line auth.html:28
	streamhtmlHeader(qw422016)
//line auth.html:29
	ln := lang.Get().Templates["banPage"]

//line auth.html:30
	if len(ln) < 3 {
//line auth.html:31
		panic(fmt.Errorf("invalid ban format strings: %v", ln))

//line auth.html:32
	}
//line auth.html:32
	qw422016.N().S(`<div class="ban-page glass">`)
//line auth.html:34
	qw422016.N().S(fmt.Sprintf(ln[0], bold(rec.Board), bold(rec.By)))
//line auth.html:34
	qw422016.N().S(`<br><br><b>`)
//line auth.html:38
	qw422016.E().S(rec.Reason)
//line auth.html:38
	qw422016.N().S(`</b><br><br>`)
//line auth.html:42
	exp := rec.Expires.Round(time.Second)

//line auth.html:43
	date := exp.Format(time.UnixDate)

//line auth.html:44
	till := exp.Sub(time.Now().Round(time.Second)).String()

//line auth.html:45
	qw422016.N().S(fmt.Sprintf(ln[1], bold(date), bold(till)))
//line auth.html:45
	qw422016.N().S(`<br><br>`)
//line auth.html:48
	qw422016.N().S(fmt.Sprintf(ln[2], bold(rec.IP)))
//line auth.html:48
	qw422016.N().S(`<br></div>`)
//line auth.html:51
	streamhtmlEnd(qw422016)
//line auth.html:52
}

//line auth.html:52
func WriteBanPage(qq422016 qtio422016.Writer, rec auth.BanRecord) {
//line auth.html:52
	qw422016 := qt422016.AcquireWriter(qq422016)
//line auth.html:52
	StreamBanPage(qw422016, rec)
//line auth.html:52
	qt422016.ReleaseWriter(qw422016)
//line auth.html:52
}

//line auth.html:52
func BanPage(rec auth.BanRecord) string {
//line auth.html:52
	qb422016 := qt422016.AcquireByteBuffer()
//line auth.html:52
	WriteBanPage(qb422016, rec)
//line auth.html:52
	qs422016 := string(qb422016.B)
//line auth.html:52
	qt422016.ReleaseByteBuffer(qb422016)
//line auth.html:52
	return qs422016
//line auth.html:52
}

// Renders a list of bans for a specific page with optional unbanning API links

//line auth.html:55
func StreamBanList(qw422016 *qt422016.Writer, bans []auth.BanRecord, board string, canUnban bool) {
//line auth.html:56
	streamhtmlHeader(qw422016)
//line auth.html:57
	streamtableStyle(qw422016)
//line auth.html:58
	ln := lang.Get()

//line auth.html:58
	qw422016.N().S(`<form method="post" action="/api/unban/`)
//line auth.html:59
	qw422016.N().S(board)
//line auth.html:59
	qw422016.N().S(`"><table>`)
//line auth.html:61
	headers := []string{
		"reason", "by", "post", "posterID", "expires", "type",
	}

//line auth.html:64
	if canUnban {
//line auth.html:65
		headers = append(headers, "unban")

//line auth.html:66
	}
//line auth.html:67
	streamtableHeaders(qw422016, headers...)
//line auth.html:68
	salt := config.Get().Salt

//line auth.html:69
	for _, b := range bans {
//line auth.html:69
		qw422016.N().S(`<tr><td>`)
//line auth.html:71
		qw422016.E().S(b.Reason)
//line auth.html:71
		qw422016.N().S(`</td><td>`)
//line auth.html:72
		qw422016.E().S(b.By)
//line auth.html:72
		qw422016.N().S(`</td><td>`)
//line auth.html:73
		streamstaticPostLink(qw422016, b.ForPost)
//line auth.html:73
		qw422016.N().S(`</td>`)
//line auth.html:74
		buf := make([]byte, 0, len(salt)+len(b.IP))

//line auth.html:75
		buf = append(buf, salt...)

//line auth.html:76
		buf = append(buf, b.IP...)

//line auth.html:76
		qw422016.N().S(`<td>`)
//line auth.html:77
		qw422016.E().S(mnemonic.FantasyName(buf))
//line auth.html:77
		qw422016.N().S(`</td><td>`)
//line auth.html:78
		qw422016.E().S(b.Expires.Format(time.UnixDate))
//line auth.html:78
		qw422016.N().S(`</td><td>`)
//line auth.html:79
		qw422016.E().S(ln.UI[b.Type])
//line auth.html:79
		qw422016.N().S(`</td>`)
//line auth.html:80
		if canUnban {
//line auth.html:80
			qw422016.N().S(`<td><input type="checkbox" name="`)
//line auth.html:82
			qw422016.E().S(strconv.FormatUint(b.ForPost, 10))
//line auth.html:82
			qw422016.N().S(`"></td>`)
//line auth.html:84
		}
//line auth.html:84
		qw422016.N().S(`</tr>`)
//line auth.html:86
	}
//line auth.html:86
	qw422016.N().S(`</table>`)
//line auth.html:88
	if canUnban {
//line auth.html:89
		streamsubmit(qw422016, false)
//line auth.html:90
	}
//line auth.html:90
	qw422016.N().S(`</form>`)
//line auth.html:92
	streamhtmlEnd(qw422016)
//line auth.html:93
}

//line auth.html:93
func WriteBanList(qq422016 qtio422016.Writer, bans []auth.BanRecord, board string, canUnban bool) {
//line auth.html:93
	qw422016 := qt422016.AcquireWriter(qq422016)
//line auth.html:93
	StreamBanList(qw422016, bans, board, canUnban)
//line auth.html:93
	qt422016.ReleaseWriter(qw422016)
//line auth.html:93
}

//line auth.html:93
func BanList(bans []auth.BanRecord, board string, canUnban bool) string {
//line auth.html:93
	qb422016 := qt422016.AcquireByteBuffer()
//line auth.html:93
	WriteBanList(qb422016, bans, board, canUnban)
//line auth.html:93
	qs422016 := string(qb422016.B)
//line auth.html:93
	qt422016.ReleaseByteBuffer(qb422016)
//line auth.html:93
	return qs422016
//line auth.html:93
}

// Common style for plain html tables

//line auth.html:96
func streamtableStyle(qw422016 *qt422016.Writer) {
//line auth.html:96
	qw422016.N().S(`<style>table, th, td {border: 1px solid black;}.hash-link {display: none;}</style>`)
//line auth.html:105
}

//line auth.html:105
func writetableStyle(qq422016 qtio422016.Writer) {
//line auth.html:105
	qw422016 := qt422016.AcquireWriter(qq422016)
//line auth.html:105
	streamtableStyle(qw422016)
//line auth.html:105
	qt422016.ReleaseWriter(qw422016)
//line auth.html:105
}

//line auth.html:105
func tableStyle() string {
//line auth.html:105
	qb422016 := qt422016.AcquireByteBuffer()
//line auth.html:105
	writetableStyle(qb422016)
//line auth.html:105
	qs422016 := string(qb422016.B)
//line auth.html:105
	qt422016.ReleaseByteBuffer(qb422016)
//line auth.html:105
	return qs422016
//line auth.html:105
}

// Post link, that will redirect to the post from any page

//line auth.html:108
func streamstaticPostLink(qw422016 *qt422016.Writer, id uint64) {
//line auth.html:109
	streampostLink(qw422016, common.Link{id, id, "all"}, true, true)
//line auth.html:110
}

//line auth.html:110
func writestaticPostLink(qq422016 qtio422016.Writer, id uint64) {
//line auth.html:110
	qw422016 := qt422016.AcquireWriter(qq422016)
//line auth.html:110
	streamstaticPostLink(qw422016, id)
//line auth.html:110
	qt422016.ReleaseWriter(qw422016)
//line auth.html:110
}

//line auth.html:110
func staticPostLink(id uint64) string {
//line auth.html:110
	qb422016 := qt422016.AcquireByteBuffer()
//line auth.html:110
	writestaticPostLink(qb422016, id)
//line auth.html:110
	qs422016 := string(qb422016.B)
//line auth.html:110
	qt422016.ReleaseByteBuffer(qb422016)
//line auth.html:110
	return qs422016
//line auth.html:110
}

// Renders a moderation log page

//line auth.html:113
func StreamModLog(qw422016 *qt422016.Writer, log []auth.ModLogEntry) {
//line auth.html:114
	streamhtmlHeader(qw422016)
//line auth.html:115
	ln := lang.Get()

//line auth.html:116
	streamtableStyle(qw422016)
//line auth.html:116
	qw422016.N().S(`<table>`)
//line auth.html:118
	streamtableHeaders(qw422016, "type", "by", "post", "time", "data", "duration")
//line auth.html:119
	for _, l := range log {
//line auth.html:119
		qw422016.N().S(`<tr><td>`)
//line auth.html:122
		switch l.Type {
//line auth.html:123
		case common.BanPost:
//line auth.html:124
			qw422016.E().S(ln.UI["ban"])
//line auth.html:125
		case common.ShadowBinPost:
//line auth.html:126
			qw422016.E().S(ln.UI["shadowBin"])
//line auth.html:127
		case common.UnbanPost:
//line auth.html:128
			qw422016.E().S(ln.UI["unban"])
//line auth.html:129
		case common.DeletePost:
//line auth.html:130
			qw422016.E().S(ln.UI["deletePost"])
//line auth.html:131
		case common.DeleteImage:
//line auth.html:132
			qw422016.E().S(ln.UI["deleteImage"])
//line auth.html:133
		case common.SpoilerImage:
//line auth.html:134
			qw422016.E().S(ln.UI["spoilerImage"])
//line auth.html:135
		case common.LockThread:
//line auth.html:136
			qw422016.E().S(ln.Common.UI["lockThread"])
//line auth.html:137
		case common.DeleteBoard:
//line auth.html:138
			qw422016.E().S(ln.Common.UI["deleteBoard"])
//line auth.html:139
		case common.MeidoVision:
//line auth.html:140
			qw422016.E().S(ln.Common.UI["meidoVisionPost"])
//line auth.html:141
		case common.PurgePost:
//line auth.html:142
			qw422016.E().S(ln.UI["purgePost"])
//line auth.html:143
		}
//line auth.html:143
		qw422016.N().S(`</td><td>`)
//line auth.html:145
		qw422016.E().S(l.By)
//line auth.html:145
		qw422016.N().S(`</td><td>`)
//line auth.html:147
		if l.ID != 0 {
//line auth.html:148
			streamstaticPostLink(qw422016, l.ID)
//line auth.html:149
		}
//line auth.html:149
		qw422016.N().S(`</td><td>`)
//line auth.html:151
		qw422016.E().S(l.Created.Format(time.UnixDate))
//line auth.html:151
		qw422016.N().S(`</td><td>`)
//line auth.html:152
		qw422016.E().S(l.Data)
//line auth.html:152
		qw422016.N().S(`</td><td>`)
//line auth.html:154
		if l.Length != 0 {
//line auth.html:155
			qw422016.E().S((time.Second * time.Duration(l.Length)).String())
//line auth.html:156
		}
//line auth.html:156
		qw422016.N().S(`</td></tr>`)
//line auth.html:159
	}
//line auth.html:159
	qw422016.N().S(`</table>`)
//line auth.html:161
	streamhtmlEnd(qw422016)
//line auth.html:162
}

//line auth.html:162
func WriteModLog(qq422016 qtio422016.Writer, log []auth.ModLogEntry) {
//line auth.html:162
	qw422016 := qt422016.AcquireWriter(qq422016)
//line auth.html:162
	StreamModLog(qw422016, log)
//line auth.html:162
	qt422016.ReleaseWriter(qw422016)
//line auth.html:162
}

//line auth.html:162
func ModLog(log []auth.ModLogEntry) string {
//line auth.html:162
	qb422016 := qt422016.AcquireByteBuffer()
//line auth.html:162
	WriteModLog(qb422016, log)
//line auth.html:162
	qs422016 := string(qb422016.B)
//line auth.html:162
	qt422016.ReleaseByteBuffer(qb422016)
//line auth.html:162
	return qs422016
//line auth.html:162
}
