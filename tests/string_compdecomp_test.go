package g_test

import (
	"testing"

	. "github.com/enetx/g"
)

func TestStringCompressionAndDecompression(t *testing.T) {
	// Test data
	inputData := String("hello world")

	// Test Zstd compression and decompression
	zstdCompressed := inputData.Compress().Zstd()
	zstdDecompressed := zstdCompressed.Decompress().Zstd()
	if zstdDecompressed.IsErr() || zstdDecompressed.Unwrap().Ne(inputData) {
		t.Errorf(
			"Zstd compression and decompression failed. Input: %s, Decompressed: %s",
			inputData,
			zstdDecompressed.Unwrap(),
		)
	}

	// Test Brotli compression and decompression
	brotliCompressed := inputData.Compress().Brotli()
	brotliDecompressed := brotliCompressed.Decompress().Brotli()
	if brotliDecompressed.IsErr() || brotliDecompressed.Unwrap().Ne(inputData) {
		t.Errorf(
			"Brotli compression and decompression failed. Input: %s, Decompressed: %s",
			inputData,
			brotliDecompressed.Unwrap(),
		)
	}

	// Test Zlib compression and decompression
	zlibCompressed := inputData.Compress().Zlib()
	zlibDecompressed := zlibCompressed.Decompress().Zlib()
	if zlibDecompressed.IsErr() || zlibDecompressed.Unwrap().Ne(inputData) {
		t.Errorf(
			"Zlib compression and decompression failed. Input: %s, Decompressed: %s",
			inputData,
			zlibDecompressed.Unwrap(),
		)
	}

	// Test Gzip compression and decompression
	gzipCompressed := inputData.Compress().Gzip()
	gzipDecompressed := gzipCompressed.Decompress().Gzip()
	if gzipDecompressed.IsErr() || gzipDecompressed.Unwrap().Ne(inputData) {
		t.Errorf(
			"Gzip compression and decompression failed. Input: %s, Decompressed: %s",
			inputData,
			gzipDecompressed.Unwrap(),
		)
	}

	// Test Flate compression and decompression
	flateCompressed := inputData.Compress().Flate()
	flateDecompressed := flateCompressed.Decompress().Flate()
	if flateDecompressed.IsErr() || flateDecompressed.Unwrap().Ne(inputData) {
		t.Errorf(
			"Flate compression and decompression failed. Input: %s, Decompressed: %s",
			inputData,
			flateDecompressed.Unwrap(),
		)
	}
}

// go test -bench=. -benchmem -count=4

var alice = String(
	"IEFsaWNlIHdhcyBiZWdpbm5pbmcgdG8gZ2V0IHZlcnkgdGlyZWQgb2Ygc2l0dGluZyBieSBoZXIgc2lzdGVyIG9uIHRoZSBiYW5rLCBhbmQgb2YgaGF2aW5nIG5vdGhpbmcgdG8gZG86IG9uY2Ugb3IgdHdpY2Ugc2hlIGhhZCBwZWVwZWQgaW50byB0aGUgYm9vayBoZXIgc2lzdGVyIHdhcyByZWFkaW5nLCBidXQgaXQgaGFkIG5vIHBpY3R1cmVzIG9yIGNvbnZlcnNhdGlvbnMgaW4gaXQsIOKAnGFuZCB3aGF0IGlzIHRoZSB1c2Ugb2YgYSBib29rLOKAnSB0aG91Z2h0IEFsaWNlIOKAnHdpdGhvdXQgcGljdHVyZXMgb3IgY29udmVyc2F0aW9ucz/igJ0KClNvIHNoZSB3YXMgY29uc2lkZXJpbmcgaW4gaGVyIG93biBtaW5kIChhcyB3ZWxsIGFzIHNoZSBjb3VsZCwgZm9yIHRoZSBob3QgZGF5IG1hZGUgaGVyIGZlZWwgdmVyeSBzbGVlcHkgYW5kIHN0dXBpZCksIHdoZXRoZXIgdGhlIHBsZWFzdXJlIG9mIG1ha2luZyBhIGRhaXN5LWNoYWluIHdvdWxkIGJlIHdvcnRoIHRoZSB0cm91YmxlIG9mIGdldHRpbmcgdXAgYW5kIHBpY2tpbmcgdGhlIGRhaXNpZXMsIHdoZW4gc3VkZGVubHkgYSBXaGl0ZSBSYWJiaXQgd2l0aCBwaW5rIGV5ZXMgcmFuIGNsb3NlIGJ5IGhlci4KClRoZXJlIHdhcyBub3RoaW5nIHNvIHZlcnkgcmVtYXJrYWJsZSBpbiB0aGF0OyBub3IgZGlkIEFsaWNlIHRoaW5rIGl0IHNvIHZlcnkgbXVjaCBvdXQgb2YgdGhlIHdheSB0byBoZWFyIHRoZSBSYWJiaXQgc2F5IHRvIGl0c2VsZiwg4oCcT2ggZGVhciEgT2ggZGVhciEgSSBzaGFsbCBiZSBsYXRlIeKAnSAod2hlbiBzaGUgdGhvdWdodCBpdCBvdmVyIGFmdGVyd2FyZHMsIGl0IG9jY3VycmVkIHRvIGhlciB0aGF0IHNoZSBvdWdodCB0byBoYXZlIHdvbmRlcmVkIGF0IHRoaXMsIGJ1dCBhdCB0aGUgdGltZSBpdCBhbGwgc2VlbWVkIHF1aXRlIG5hdHVyYWwpOyBidXQgd2hlbiB0aGUgUmFiYml0IGFjdHVhbGx5IHRvb2sgYSB3YXRjaCBvdXQgb2YgaXRzIHdhaXN0Y29hdC1wb2NrZXQsIGFuZCBsb29rZWQgYXQgaXQsIGFuZCB0aGVuIGh1cnJpZWQgb24sIEFsaWNlIHN0YXJ0ZWQgdG8gaGVyIGZlZXQsIGZvciBpdCBmbGFzaGVkIGFjcm9zcyBoZXIgbWluZCB0aGF0IHNoZSBoYWQgbmV2ZXIgYmVmb3JlIHNlZW4gYSByYWJiaXQgd2l0aCBlaXRoZXIgYSB3YWlzdGNvYXQtcG9ja2V0LCBvciBhIHdhdGNoIHRvIHRha2Ugb3V0IG9mIGl0LCBhbmQgYnVybmluZyB3aXRoIGN1cmlvc2l0eSwgc2hlIHJhbiBhY3Jvc3MgdGhlIGZpZWxkIGFmdGVyIGl0LCBhbmQgZm9ydHVuYXRlbHkgd2FzIGp1c3QgaW4gdGltZSB0byBzZWUgaXQgcG9wIGRvd24gYSBsYXJnZSByYWJiaXQtaG9sZSB1bmRlciB0aGUgaGVkZ2UuCgpJbiBhbm90aGVyIG1vbWVudCBkb3duIHdlbnQgQWxpY2UgYWZ0ZXIgaXQsIG5ldmVyIG9uY2UgY29uc2lkZXJpbmcgaG93IGluIHRoZSB3b3JsZCBzaGUgd2FzIHRvIGdldCBvdXQgYWdhaW4uCgpUaGUgcmFiYml0LWhvbGUgd2VudCBzdHJhaWdodCBvbiBsaWtlIGEgdHVubmVsIGZvciBzb21lIHdheSwgYW5kIHRoZW4gZGlwcGVkIHN1ZGRlbmx5IGRvd24sIHNvIHN1ZGRlbmx5IHRoYXQgQWxpY2UgaGFkIG5vdCBhIG1vbWVudCB0byB0aGluayBhYm91dCBzdG9wcGluZyBoZXJzZWxmIGJlZm9yZSBzaGUgZm91bmQgaGVyc2VsZiBmYWxsaW5nIGRvd24gYSB2ZXJ5IGRlZXAgd2VsbC4KCkVpdGhlciB0aGUgd2VsbCB3YXMgdmVyeSBkZWVwLCBvciBzaGUgZmVsbCB2ZXJ5IHNsb3dseSwgZm9yIHNoZSBoYWQgcGxlbnR5IG9mIHRpbWUgYXMgc2hlIHdlbnQgZG93biB0byBsb29rIGFib3V0IGhlciBhbmQgdG8gd29uZGVyIHdoYXQgd2FzIGdvaW5nIHRvIGhhcHBlbiBuZXh0LiBGaXJzdCwgc2hlIHRyaWVkIHRvIGxvb2sgZG93biBhbmQgbWFrZSBvdXQgd2hhdCBzaGUgd2FzIGNvbWluZyB0bywgYnV0IGl0IHdhcyB0b28gZGFyayB0byBzZWUgYW55dGhpbmc7IHRoZW4gc2hlIGxvb2tlZCBhdCB0aGUgc2lkZXMgb2YgdGhlIHdlbGwsIGFuZCBub3RpY2VkIHRoYXQgdGhleSB3ZXJlIGZpbGxlZCB3aXRoIGN1cGJvYXJkcyBhbmQgYm9vay1zaGVsdmVzOyBoZXJlIGFuZCB0aGVyZSBzaGUgc2F3IG1hcHMgYW5kIHBpY3R1cmVzIGh1bmcgdXBvbiBwZWdzLiBTaGUgdG9vayBkb3duIGEgamFyIGZyb20gb25lIG9mIHRoZSBzaGVsdmVzIGFzIHNoZSBwYXNzZWQ7IGl0IHdhcyBsYWJlbGxlZCDigJxPUkFOR0UgTUFSTUFMQURF4oCdLCBidXQgdG8gaGVyIGdyZWF0IGRpc2FwcG9pbnRtZW50IGl0IHdhcyBlbXB0eTogc2hlIGRpZCBub3QgbGlrZSB0byBkcm9wIHRoZSBqYXIgZm9yIGZlYXIgb2Yga2lsbGluZyBzb21lYm9keSB1bmRlcm5lYXRoLCBzbyBtYW5hZ2VkIHRvIHB1dCBpdCBpbnRvIG9uZSBvZiB0aGUgY3VwYm9hcmRzIGFzIHNoZSBmZWxsIHBhc3QgaXQuCgrigJxXZWxsIeKAnSB0aG91Z2h0IEFsaWNlIHRvIGhlcnNlbGYsIOKAnGFmdGVyIHN1Y2ggYSBmYWxsIGFzIHRoaXMsIEkgc2hhbGwgdGhpbmsgbm90aGluZyBvZiB0dW1ibGluZyBkb3duIHN0YWlycyEgSG93IGJyYXZlIHRoZXnigJlsbCBhbGwgdGhpbmsgbWUgYXQgaG9tZSEgV2h5LCBJIHdvdWxkbuKAmXQgc2F5IGFueXRoaW5nIGFib3V0IGl0LCBldmVuIGlmIEkgZmVsbCBvZmYgdGhlIHRvcCBvZiB0aGUgaG91c2Uh4oCdIChXaGljaCB3YXMgdmVyeSBsaWtlbHkgdHJ1ZS4pCgpEb3duLCBkb3duLCBkb3duLiBXb3VsZCB0aGUgZmFsbCBuZXZlciBjb21lIHRvIGFuIGVuZD8g4oCcSSB3b25kZXIgaG93IG1hbnkgbWlsZXMgSeKAmXZlIGZhbGxlbiBieSB0aGlzIHRpbWU/4oCdIHNoZSBzYWlkIGFsb3VkLiDigJxJIG11c3QgYmUgZ2V0dGluZyBzb21ld2hlcmUgbmVhciB0aGUgY2VudHJlIG9mIHRoZSBlYXJ0aC4gTGV0IG1lIHNlZTogdGhhdCB3b3VsZCBiZSBmb3VyIHRob3VzYW5kIG1pbGVzIGRvd24sIEkgdGhpbmvigJTigJ0gKGZvciwgeW91IHNlZSwgQWxpY2UgaGFkIGxlYXJudCBzZXZlcmFsIHRoaW5ncyBvZiB0aGlzIHNvcnQgaW4gaGVyIGxlc3NvbnMgaW4gdGhlIHNjaG9vbHJvb20sIGFuZCB0aG91Z2ggdGhpcyB3YXMgbm90IGEgdmVyeSBnb29kIG9wcG9ydHVuaXR5IGZvciBzaG93aW5nIG9mZiBoZXIga25vd2xlZGdlLCBhcyB0aGVyZSB3YXMgbm8gb25lIHRvIGxpc3RlbiB0byBoZXIsIHN0aWxsIGl0IHdhcyBnb29kIHByYWN0aWNlIHRvIHNheSBpdCBvdmVyKSDigJzigJR5ZXMsIHRoYXTigJlzIGFib3V0IHRoZSByaWdodCBkaXN0YW5jZeKAlGJ1dCB0aGVuIEkgd29uZGVyIHdoYXQgTGF0aXR1ZGUgb3IgTG9uZ2l0dWRlIEnigJl2ZSBnb3QgdG8/4oCdIChBbGljZSBoYWQgbm8gaWRlYSB3aGF0IExhdGl0dWRlIHdhcywgb3IgTG9uZ2l0dWRlIGVpdGhlciwgYnV0IHRob3VnaHQgdGhleSB3ZXJlIG5pY2UgZ3JhbmQgd29yZHMgdG8gc2F5LikKClByZXNlbnRseSBzaGUgYmVnYW4gYWdhaW4uIOKAnEkgd29uZGVyIGlmIEkgc2hhbGwgZmFsbCByaWdodCB0aHJvdWdoIHRoZSBlYXJ0aCEgSG93IGZ1bm55IGl04oCZbGwgc2VlbSB0byBjb21lIG91dCBhbW9uZyB0aGUgcGVvcGxlIHRoYXQgd2FsayB3aXRoIHRoZWlyIGhlYWRzIGRvd253YXJkISBUaGUgQW50aXBhdGhpZXMsIEkgdGhpbmvigJTigJ0gKHNoZSB3YXMgcmF0aGVyIGdsYWQgdGhlcmUgd2FzIG5vIG9uZSBsaXN0ZW5pbmcsIHRoaXMgdGltZSwgYXMgaXQgZGlkbuKAmXQgc291bmQgYXQgYWxsIHRoZSByaWdodCB3b3JkKSDigJzigJRidXQgSSBzaGFsbCBoYXZlIHRvIGFzayB0aGVtIHdoYXQgdGhlIG5hbWUgb2YgdGhlIGNvdW50cnkgaXMsIHlvdSBrbm93LiBQbGVhc2UsIE1h4oCZYW0sIGlzIHRoaXMgTmV3IFplYWxhbmQgb3IgQXVzdHJhbGlhP+KAnSAoYW5kIHNoZSB0cmllZCB0byBjdXJ0c2V5IGFzIHNoZSBzcG9rZeKAlGZhbmN5IGN1cnRzZXlpbmcgYXMgeW914oCZcmUgZmFsbGluZyB0aHJvdWdoIHRoZSBhaXIhIERvIHlvdSB0aGluayB5b3UgY291bGQgbWFuYWdlIGl0Pykg4oCcQW5kIHdoYXQgYW4gaWdub3JhbnQgbGl0dGxlIGdpcmwgc2hl4oCZbGwgdGhpbmsgbWUgZm9yIGFza2luZyEgTm8sIGl04oCZbGwgbmV2ZXIgZG8gdG8gYXNrOiBwZXJoYXBzIEkgc2hhbGwgc2VlIGl0IHdyaXR0ZW4gdXAgc29tZXdoZXJlLuKAnQoKRG93biwgZG93biwgZG93bi4gVGhlcmUgd2FzIG5vdGhpbmcgZWxzZSB0byBkbywgc28gQWxpY2Ugc29vbiBiZWdhbiB0YWxraW5nIGFnYWluLiDigJxEaW5haOKAmWxsIG1pc3MgbWUgdmVyeSBtdWNoIHRvLW5pZ2h0LCBJIHNob3VsZCB0aGluayHigJ0gKERpbmFoIHdhcyB0aGUgY2F0Likg4oCcSSBob3BlIHRoZXnigJlsbCByZW1lbWJlciBoZXIgc2F1Y2VyIG9mIG1pbGsgYXQgdGVhLXRpbWUuIERpbmFoIG15IGRlYXIhIEkgd2lzaCB5b3Ugd2VyZSBkb3duIGhlcmUgd2l0aCBtZSEgVGhlcmUgYXJlIG5vIG1pY2UgaW4gdGhlIGFpciwgSeKAmW0gYWZyYWlkLCBidXQgeW91IG1pZ2h0IGNhdGNoIGEgYmF0LCBhbmQgdGhhdOKAmXMgdmVyeSBsaWtlIGEgbW91c2UsIHlvdSBrbm93LiBCdXQgZG8gY2F0cyBlYXQgYmF0cywgSSB3b25kZXI/4oCdIEFuZCBoZXJlIEFsaWNlIGJlZ2FuIHRvIGdldCByYXRoZXIgc2xlZXB5LCBhbmQgd2VudCBvbiBzYXlpbmcgdG8gaGVyc2VsZiwgaW4gYSBkcmVhbXkgc29ydCBvZiB3YXksIOKAnERvIGNhdHMgZWF0IGJhdHM/IERvIGNhdHMgZWF0IGJhdHM/4oCdIGFuZCBzb21ldGltZXMsIOKAnERvIGJhdHMgZWF0IGNhdHM/4oCdIGZvciwgeW91IHNlZSwgYXMgc2hlIGNvdWxkbuKAmXQgYW5zd2VyIGVpdGhlciBxdWVzdGlvbiwgaXQgZGlkbuKAmXQgbXVjaCBtYXR0ZXIgd2hpY2ggd2F5IHNoZSBwdXQgaXQuIFNoZSBmZWx0IHRoYXQgc2hlIHdhcyBkb3ppbmcgb2ZmLCBhbmQgaGFkIGp1c3QgYmVndW4gdG8gZHJlYW0gdGhhdCBzaGUgd2FzIHdhbGtpbmcgaGFuZCBpbiBoYW5kIHdpdGggRGluYWgsIGFuZCBzYXlpbmcgdG8gaGVyIHZlcnkgZWFybmVzdGx5LCDigJxOb3csIERpbmFoLCB0ZWxsIG1lIHRoZSB0cnV0aDogZGlkIHlvdSBldmVyIGVhdCBhIGJhdD/igJ0gd2hlbiBzdWRkZW5seSwgdGh1bXAhIHRodW1wISBkb3duIHNoZSBjYW1lIHVwb24gYSBoZWFwIG9mIHN0aWNrcyBhbmQgZHJ5IGxlYXZlcywgYW5kIHRoZSBmYWxsIHdhcyBvdmVyLgoKQWxpY2Ugd2FzIG5vdCBhIGJpdCBodXJ0LCBhbmQgc2hlIGp1bXBlZCB1cCBvbiB0byBoZXIgZmVldCBpbiBhIG1vbWVudDogc2hlIGxvb2tlZCB1cCwgYnV0IGl0IHdhcyBhbGwgZGFyayBvdmVyaGVhZDsgYmVmb3JlIGhlciB3YXMgYW5vdGhlciBsb25nIHBhc3NhZ2UsIGFuZCB0aGUgV2hpdGUgUmFiYml0IHdhcyBzdGlsbCBpbiBzaWdodCwgaHVycnlpbmcgZG93biBpdC4gVGhlcmUgd2FzIG5vdCBhIG1vbWVudCB0byBiZSBsb3N0OiBhd2F5IHdlbnQgQWxpY2UgbGlrZSB0aGUgd2luZCwgYW5kIHdhcyBqdXN0IGluIHRpbWUgdG8gaGVhciBpdCBzYXksIGFzIGl0IHR1cm5lZCBhIGNvcm5lciwg4oCcT2ggbXkgZWFycyBhbmQgd2hpc2tlcnMsIGhvdyBsYXRlIGl04oCZcyBnZXR0aW5nIeKAnSBTaGUgd2FzIGNsb3NlIGJlaGluZCBpdCB3aGVuIHNoZSB0dXJuZWQgdGhlIGNvcm5lciwgYnV0IHRoZSBSYWJiaXQgd2FzIG5vIGxvbmdlciB0byBiZSBzZWVuOiBzaGUgZm91bmQgaGVyc2VsZiBpbiBhIGxvbmcsIGxvdyBoYWxsLCB3aGljaCB3YXMgbGl0IHVwIGJ5IGEgcm93IG9mIGxhbXBzIGhhbmdpbmcgZnJvbSB0aGUgcm9vZi4KClRoZXJlIHdlcmUgZG9vcnMgYWxsIHJvdW5kIHRoZSBoYWxsLCBidXQgdGhleSB3ZXJlIGFsbCBsb2NrZWQ7IGFuZCB3aGVuIEFsaWNlIGhhZCBiZWVuIGFsbCB0aGUgd2F5IGRvd24gb25lIHNpZGUgYW5kIHVwIHRoZSBvdGhlciwgdHJ5aW5nIGV2ZXJ5IGRvb3IsIHNoZSB3YWxrZWQgc2FkbHkgZG93biB0aGUgbWlkZGxlLCB3b25kZXJpbmcgaG93IHNoZSB3YXMgZXZlciB0byBnZXQgb3V0IGFnYWluLgoKU3VkZGVubHkgc2hlIGNhbWUgdXBvbiBhIGxpdHRsZSB0aHJlZS1sZWdnZWQgdGFibGUsIGFsbCBtYWRlIG9mIHNvbGlkIGdsYXNzOyB0aGVyZSB3YXMgbm90aGluZyBvbiBpdCBleGNlcHQgYSB0aW55IGdvbGRlbiBrZXksIGFuZCBBbGljZeKAmXMgZmlyc3QgdGhvdWdodCB3YXMgdGhhdCBpdCBtaWdodCBiZWxvbmcgdG8gb25lIG9mIHRoZSBkb29ycyBvZiB0aGUgaGFsbDsgYnV0LCBhbGFzISBlaXRoZXIgdGhlIGxvY2tzIHdlcmUgdG9vIGxhcmdlLCBvciB0aGUga2V5IHdhcyB0b28gc21hbGwsIGJ1dCBhdCBhbnkgcmF0ZSBpdCB3b3VsZCBub3Qgb3BlbiBhbnkgb2YgdGhlbS4gSG93ZXZlciwgb24gdGhlIHNlY29uZCB0aW1lIHJvdW5kLCBzaGUgY2FtZSB1cG9uIGEgbG93IGN1cnRhaW4gc2hlIGhhZCBub3Qgbm90aWNlZCBiZWZvcmUsIGFuZCBiZWhpbmQgaXQgd2FzIGEgbGl0dGxlIGRvb3IgYWJvdXQgZmlmdGVlbiBpbmNoZXMgaGlnaDogc2hlIHRyaWVkIHRoZSBsaXR0bGUgZ29sZGVuIGtleSBpbiB0aGUgbG9jaywgYW5kIHRvIGhlciBncmVhdCBkZWxpZ2h0IGl0IGZpdHRlZCEg",
)

func BenchmarkGzip(b *testing.B) {
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		alice.Compress().Gzip()
	}
}

func BenchmarkFlate(b *testing.B) {
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		alice.Compress().Flate()
	}
}

func BenchmarkZlib(b *testing.B) {
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		alice.Compress().Zlib()
	}
}

func TestZstdDecompressionError(t *testing.T) {
	// Test Zstd decompression with invalid data
	invalidData := String("invalid zstd data")
	result := invalidData.Decompress().Zstd()

	if result.IsOk() {
		t.Error("Zstd decompression of invalid data should fail")
	}
}

func TestBrotliDecompressionError(t *testing.T) {
	// Test Brotli decompression with invalid data
	invalidData := String("invalid brotli data")
	result := invalidData.Decompress().Brotli()

	if result.IsOk() {
		t.Error("Brotli decompression of invalid data should fail")
	}
}

func TestZlibDecompressionError(t *testing.T) {
	// Test Zlib decompression with invalid data
	invalidData := String("invalid zlib data")
	result := invalidData.Decompress().Zlib()

	if result.IsOk() {
		t.Error("Zlib decompression of invalid data should fail")
	}
}

func TestGzipDecompressionError(t *testing.T) {
	// Test Gzip decompression with invalid data
	invalidData := String("invalid gzip data")
	result := invalidData.Decompress().Gzip()

	if result.IsOk() {
		t.Error("Gzip decompression of invalid data should fail")
	}
}

func TestFlateCompressionDecompression(t *testing.T) {
	// Test Flate compression with larger data to hit more code paths
	inputData := String("hello world, this is a longer string for testing flate compression and decompression")

	// Test compression
	compressed := inputData.Compress().Flate()

	// Test decompression
	decompressed := compressed.Decompress().Flate()
	if decompressed.IsErr() {
		t.Fatalf("Flate decompression failed: %v", decompressed.Err())
	}

	if decompressed.Ok().Ne(inputData) {
		t.Errorf("Flate round trip failed. Expected: %s, Got: %s", inputData, decompressed.Ok())
	}
}

func TestFlateDecompressionError(t *testing.T) {
	// Test Flate decompression with invalid data
	invalidData := String("invalid flate data")
	result := invalidData.Decompress().Flate()

	if result.IsOk() {
		t.Error("Flate decompression of invalid data should fail")
	}
}

func TestEmptyStringCompression(t *testing.T) {
	// Test compression of empty string
	emptyStr := String("")

	// Test all compression methods with empty string
	zstdCompressed := emptyStr.Compress().Zstd()
	brotliCompressed := emptyStr.Compress().Brotli()
	zlibCompressed := emptyStr.Compress().Zlib()
	gzipCompressed := emptyStr.Compress().Gzip()
	flateCompressed := emptyStr.Compress().Flate()

	// Test that empty strings can be compressed and decompressed properly
	// Some compression algorithms may produce empty compressed data for empty input

	// Test decompression of empty compressed data
	zstdDecompressed := zstdCompressed.Decompress().Zstd()
	if zstdDecompressed.IsErr() {
		t.Errorf("Zstd decompression of empty string failed: %v", zstdDecompressed.Err())
	} else if zstdDecompressed.Ok().Ne(emptyStr) {
		t.Errorf("Zstd round trip failed for empty string")
	}

	brotliDecompressed := brotliCompressed.Decompress().Brotli()
	if brotliDecompressed.IsErr() {
		t.Errorf("Brotli decompression of empty string failed: %v", brotliDecompressed.Err())
	} else if brotliDecompressed.Ok().Ne(emptyStr) {
		t.Errorf("Brotli round trip failed for empty string")
	}

	zlibDecompressed := zlibCompressed.Decompress().Zlib()
	if zlibDecompressed.IsErr() {
		t.Errorf("Zlib decompression of empty string failed: %v", zlibDecompressed.Err())
	} else if zlibDecompressed.Ok().Ne(emptyStr) {
		t.Errorf("Zlib round trip failed for empty string")
	}

	gzipDecompressed := gzipCompressed.Decompress().Gzip()
	if gzipDecompressed.IsErr() {
		t.Errorf("Gzip decompression of empty string failed: %v", gzipDecompressed.Err())
	} else if gzipDecompressed.Ok().Ne(emptyStr) {
		t.Errorf("Gzip round trip failed for empty string")
	}

	flateDecompressed := flateCompressed.Decompress().Flate()
	if flateDecompressed.IsErr() {
		t.Errorf("Flate decompression of empty string failed: %v", flateDecompressed.Err())
	} else if flateDecompressed.Ok().Ne(emptyStr) {
		t.Errorf("Flate round trip failed for empty string")
	}
}
