// This file is automatically generated by qtc from "auth.qtpl".
// See https://github.com/valyala/quicktemplate for details.

//line auth.qtpl:1
package templates

//line auth.qtpl:1
import "fmt"

//line auth.qtpl:2
import "time"

//line auth.qtpl:3
import "strconv"

//line auth.qtpl:4
import "meguca/auth"

//line auth.qtpl:5
import "meguca/config"

//line auth.qtpl:6
import "meguca/lang"

//line auth.qtpl:7
import "meguca/common"

//line auth.qtpl:8
import "github.com/bakape/mnemonics"

// Header of a standalone HTML page

//line auth.qtpl:11
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line auth.qtpl:11
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line auth.qtpl:11
func streamhtmlHeader(qw422016 *qt422016.Writer) {
	//line auth.qtpl:11
	qw422016.N().S(`<!DOCTYPE html><html><head><meta charset="utf-8"/></head><body>`)
//line auth.qtpl:18
}

//line auth.qtpl:18
func writehtmlHeader(qq422016 qtio422016.Writer) {
	//line auth.qtpl:18
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line auth.qtpl:18
	streamhtmlHeader(qw422016)
	//line auth.qtpl:18
	qt422016.ReleaseWriter(qw422016)
//line auth.qtpl:18
}

//line auth.qtpl:18
func htmlHeader() string {
	//line auth.qtpl:18
	qb422016 := qt422016.AcquireByteBuffer()
	//line auth.qtpl:18
	writehtmlHeader(qb422016)
	//line auth.qtpl:18
	qs422016 := string(qb422016.B)
	//line auth.qtpl:18
	qt422016.ReleaseByteBuffer(qb422016)
	//line auth.qtpl:18
	return qs422016
//line auth.qtpl:18
}

// End of a standalone HTML page

//line auth.qtpl:21
func streamhtmlEnd(qw422016 *qt422016.Writer) {
	//line auth.qtpl:21
	qw422016.N().S(`</body></html>`)
//line auth.qtpl:24
}

//line auth.qtpl:24
func writehtmlEnd(qq422016 qtio422016.Writer) {
	//line auth.qtpl:24
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line auth.qtpl:24
	streamhtmlEnd(qw422016)
	//line auth.qtpl:24
	qt422016.ReleaseWriter(qw422016)
//line auth.qtpl:24
}

//line auth.qtpl:24
func htmlEnd() string {
	//line auth.qtpl:24
	qb422016 := qt422016.AcquireByteBuffer()
	//line auth.qtpl:24
	writehtmlEnd(qb422016)
	//line auth.qtpl:24
	qs422016 := string(qb422016.B)
	//line auth.qtpl:24
	qt422016.ReleaseByteBuffer(qb422016)
	//line auth.qtpl:24
	return qs422016
//line auth.qtpl:24
}

// BanPage renders a ban page for a banned user

//line auth.qtpl:27
func StreamBanPage(qw422016 *qt422016.Writer, rec auth.BanRecord) {
	//line auth.qtpl:28
	streamhtmlHeader(qw422016)
	//line auth.qtpl:29
	ln := lang.Get().Templates["banPage"]

	//line auth.qtpl:29
	qw422016.N().S(`<div class="ban-page glass">`)
	//line auth.qtpl:31
	qw422016.N().S(fmt.Sprintf(ln[0], bold(rec.Board), bold(rec.By)))
	//line auth.qtpl:31
	qw422016.N().S(`<br><br><b>`)
	//line auth.qtpl:35
	qw422016.E().S(rec.Reason)
	//line auth.qtpl:35
	qw422016.N().S(`</b><br><br>`)
	//line auth.qtpl:39
	exp := rec.Expires.Round(time.Second)

	//line auth.qtpl:40
	date := exp.Format(time.UnixDate)

	//line auth.qtpl:41
	till := exp.Sub(time.Now().Round(time.Second)).String()

	//line auth.qtpl:42
	qw422016.N().S(fmt.Sprintf(ln[1], bold(date), bold(till)))
	//line auth.qtpl:42
	qw422016.N().S(`<br><br>`)
	//line auth.qtpl:45
	qw422016.N().S(fmt.Sprintf(ln[2], bold(rec.IP)))
	//line auth.qtpl:45
	qw422016.N().S(`<br></div>`)
	//line auth.qtpl:48
	streamhtmlEnd(qw422016)
//line auth.qtpl:49
}

//line auth.qtpl:49
func WriteBanPage(qq422016 qtio422016.Writer, rec auth.BanRecord) {
	//line auth.qtpl:49
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line auth.qtpl:49
	StreamBanPage(qw422016, rec)
	//line auth.qtpl:49
	qt422016.ReleaseWriter(qw422016)
//line auth.qtpl:49
}

//line auth.qtpl:49
func BanPage(rec auth.BanRecord) string {
	//line auth.qtpl:49
	qb422016 := qt422016.AcquireByteBuffer()
	//line auth.qtpl:49
	WriteBanPage(qb422016, rec)
	//line auth.qtpl:49
	qs422016 := string(qb422016.B)
	//line auth.qtpl:49
	qt422016.ReleaseByteBuffer(qb422016)
	//line auth.qtpl:49
	return qs422016
//line auth.qtpl:49
}

// Renders a list of bans for a specific page with optional unbanning API links

//line auth.qtpl:52
func StreamBanList(qw422016 *qt422016.Writer, bans []auth.BanRecord, board string, canUnban bool) {
	//line auth.qtpl:53
	streamhtmlHeader(qw422016)
	//line auth.qtpl:54
	streamtableStyle(qw422016)
	//line auth.qtpl:54
	qw422016.N().S(`<form method="post" action="/api/unban/`)
	//line auth.qtpl:55
	qw422016.N().S(board)
	//line auth.qtpl:55
	qw422016.N().S(`"><table>`)
	//line auth.qtpl:57
	headers := []string{
		"reason", "by", "post", "posterID", "expires",
	}

	//line auth.qtpl:60
	if canUnban {
		//line auth.qtpl:61
		headers = append(headers, "unban")

		//line auth.qtpl:62
	}
	//line auth.qtpl:63
	streamtableHeaders(qw422016, headers...)
	//line auth.qtpl:64
	salt := config.Get().Salt

	//line auth.qtpl:65
	for _, b := range bans {
		//line auth.qtpl:65
		qw422016.N().S(`<tr><td>`)
		//line auth.qtpl:67
		qw422016.E().S(b.Reason)
		//line auth.qtpl:67
		qw422016.N().S(`</td><td>`)
		//line auth.qtpl:68
		qw422016.E().S(b.By)
		//line auth.qtpl:68
		qw422016.N().S(`</td><td>`)
		//line auth.qtpl:69
		streamstaticPostLink(qw422016, b.ForPost)
		//line auth.qtpl:69
		qw422016.N().S(`</td>`)
		//line auth.qtpl:70
		buf := make([]byte, 0, len(salt)+len(b.IP))

		//line auth.qtpl:71
		buf = append(buf, salt...)

		//line auth.qtpl:72
		buf = append(buf, b.IP...)

		//line auth.qtpl:72
		qw422016.N().S(`<td>`)
		//line auth.qtpl:73
		qw422016.E().S(mnemonic.FantasyName(buf))
		//line auth.qtpl:73
		qw422016.N().S(`</td><td>`)
		//line auth.qtpl:74
		qw422016.E().S(b.Expires.Format(time.UnixDate))
		//line auth.qtpl:74
		qw422016.N().S(`</td>`)
		//line auth.qtpl:75
		if canUnban {
			//line auth.qtpl:75
			qw422016.N().S(`<td><input type="checkbox" name="`)
			//line auth.qtpl:77
			qw422016.E().S(strconv.FormatUint(b.ForPost, 10))
			//line auth.qtpl:77
			qw422016.N().S(`"></td>`)
			//line auth.qtpl:79
		}
		//line auth.qtpl:79
		qw422016.N().S(`</tr>`)
		//line auth.qtpl:81
	}
	//line auth.qtpl:81
	qw422016.N().S(`</table>`)
	//line auth.qtpl:83
	if canUnban {
		//line auth.qtpl:84
		streamsubmit(qw422016, false)
		//line auth.qtpl:85
	}
	//line auth.qtpl:85
	qw422016.N().S(`</form>`)
	//line auth.qtpl:87
	streamhtmlEnd(qw422016)
//line auth.qtpl:88
}

//line auth.qtpl:88
func WriteBanList(qq422016 qtio422016.Writer, bans []auth.BanRecord, board string, canUnban bool) {
	//line auth.qtpl:88
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line auth.qtpl:88
	StreamBanList(qw422016, bans, board, canUnban)
	//line auth.qtpl:88
	qt422016.ReleaseWriter(qw422016)
//line auth.qtpl:88
}

//line auth.qtpl:88
func BanList(bans []auth.BanRecord, board string, canUnban bool) string {
	//line auth.qtpl:88
	qb422016 := qt422016.AcquireByteBuffer()
	//line auth.qtpl:88
	WriteBanList(qb422016, bans, board, canUnban)
	//line auth.qtpl:88
	qs422016 := string(qb422016.B)
	//line auth.qtpl:88
	qt422016.ReleaseByteBuffer(qb422016)
	//line auth.qtpl:88
	return qs422016
//line auth.qtpl:88
}

// Common style for plain html tables

//line auth.qtpl:91
func streamtableStyle(qw422016 *qt422016.Writer) {
	//line auth.qtpl:91
	qw422016.N().S(`<style>table, th, td {border: 1px solid black;}.hash-link {display: none;}</style>`)
//line auth.qtpl:100
}

//line auth.qtpl:100
func writetableStyle(qq422016 qtio422016.Writer) {
	//line auth.qtpl:100
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line auth.qtpl:100
	streamtableStyle(qw422016)
	//line auth.qtpl:100
	qt422016.ReleaseWriter(qw422016)
//line auth.qtpl:100
}

//line auth.qtpl:100
func tableStyle() string {
	//line auth.qtpl:100
	qb422016 := qt422016.AcquireByteBuffer()
	//line auth.qtpl:100
	writetableStyle(qb422016)
	//line auth.qtpl:100
	qs422016 := string(qb422016.B)
	//line auth.qtpl:100
	qt422016.ReleaseByteBuffer(qb422016)
	//line auth.qtpl:100
	return qs422016
//line auth.qtpl:100
}

// Post link, that will redirect to the post from any page

//line auth.qtpl:103
func streamstaticPostLink(qw422016 *qt422016.Writer, id uint64) {
	//line auth.qtpl:104
	streampostLink(qw422016, common.Link{id, id, "all"}, true, true)
//line auth.qtpl:105
}

//line auth.qtpl:105
func writestaticPostLink(qq422016 qtio422016.Writer, id uint64) {
	//line auth.qtpl:105
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line auth.qtpl:105
	streamstaticPostLink(qw422016, id)
	//line auth.qtpl:105
	qt422016.ReleaseWriter(qw422016)
//line auth.qtpl:105
}

//line auth.qtpl:105
func staticPostLink(id uint64) string {
	//line auth.qtpl:105
	qb422016 := qt422016.AcquireByteBuffer()
	//line auth.qtpl:105
	writestaticPostLink(qb422016, id)
	//line auth.qtpl:105
	qs422016 := string(qb422016.B)
	//line auth.qtpl:105
	qt422016.ReleaseByteBuffer(qb422016)
	//line auth.qtpl:105
	return qs422016
//line auth.qtpl:105
}

// Renders a moderation log page

//line auth.qtpl:108
func StreamModLog(qw422016 *qt422016.Writer, log []auth.ModLogEntry) {
	//line auth.qtpl:109
	streamhtmlHeader(qw422016)
	//line auth.qtpl:110
	ln := lang.Get()

	//line auth.qtpl:111
	streamtableStyle(qw422016)
	//line auth.qtpl:111
	qw422016.N().S(`<table>`)
	//line auth.qtpl:113
	streamtableHeaders(qw422016, "type", "by", "post", "time", "data", "duration")
	//line auth.qtpl:114
	for _, l := range log {
		//line auth.qtpl:114
		qw422016.N().S(`<tr><td>`)
		//line auth.qtpl:117
		switch l.Type {
		//line auth.qtpl:118
		case common.BanPost:
			//line auth.qtpl:119
			qw422016.E().S(ln.UI["ban"])
		//line auth.qtpl:120
		case common.UnbanPost:
			//line auth.qtpl:121
			qw422016.E().S(ln.UI["unban"])
		//line auth.qtpl:122
		case common.DeletePost:
			//line auth.qtpl:123
			qw422016.E().S(ln.UI["deletePost"])
		//line auth.qtpl:124
		case common.DeleteImage:
			//line auth.qtpl:125
			qw422016.E().S(ln.UI["deleteImage"])
		//line auth.qtpl:126
		case common.SpoilerImage:
			//line auth.qtpl:127
			qw422016.E().S(ln.UI["spoilerImage"])
		//line auth.qtpl:128
		case common.LockThread:
			//line auth.qtpl:129
			qw422016.E().S(ln.Common.UI["lockThread"])
		//line auth.qtpl:130
		case common.DeleteBoard:
			//line auth.qtpl:131
			qw422016.E().S(ln.Common.UI["deleteBoard"])
		//line auth.qtpl:132
		case common.MeidoVision:
			//line auth.qtpl:133
			qw422016.E().S(ln.Common.UI["meidoVisionPost"])
		//line auth.qtpl:134
		case common.PurgePost:
			//line auth.qtpl:135
			qw422016.E().S(ln.UI["purgePost"])
			//line auth.qtpl:136
		}
		//line auth.qtpl:136
		qw422016.N().S(`</td><td>`)
		//line auth.qtpl:138
		qw422016.E().S(l.By)
		//line auth.qtpl:138
		qw422016.N().S(`</td><td>`)
		//line auth.qtpl:140
		if l.ID != 0 {
			//line auth.qtpl:141
			streamstaticPostLink(qw422016, l.ID)
			//line auth.qtpl:142
		}
		//line auth.qtpl:142
		qw422016.N().S(`</td><td>`)
		//line auth.qtpl:144
		qw422016.E().S(l.Created.Format(time.UnixDate))
		//line auth.qtpl:144
		qw422016.N().S(`</td><td>`)
		//line auth.qtpl:145
		qw422016.E().S(l.Data)
		//line auth.qtpl:145
		qw422016.N().S(`</td><td>`)
		//line auth.qtpl:147
		if l.Length != 0 {
			//line auth.qtpl:148
			qw422016.E().S((time.Second * time.Duration(l.Length)).String())
			//line auth.qtpl:149
		}
		//line auth.qtpl:149
		qw422016.N().S(`</td></tr>`)
		//line auth.qtpl:152
	}
	//line auth.qtpl:152
	qw422016.N().S(`</table>`)
	//line auth.qtpl:154
	streamhtmlEnd(qw422016)
//line auth.qtpl:155
}

//line auth.qtpl:155
func WriteModLog(qq422016 qtio422016.Writer, log []auth.ModLogEntry) {
	//line auth.qtpl:155
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line auth.qtpl:155
	StreamModLog(qw422016, log)
	//line auth.qtpl:155
	qt422016.ReleaseWriter(qw422016)
//line auth.qtpl:155
}

//line auth.qtpl:155
func ModLog(log []auth.ModLogEntry) string {
	//line auth.qtpl:155
	qb422016 := qt422016.AcquireByteBuffer()
	//line auth.qtpl:155
	WriteModLog(qb422016, log)
	//line auth.qtpl:155
	qs422016 := string(qb422016.B)
	//line auth.qtpl:155
	qt422016.ReleaseByteBuffer(qb422016)
	//line auth.qtpl:155
	return qs422016
//line auth.qtpl:155
}
