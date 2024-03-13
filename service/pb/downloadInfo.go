package pb

func (x *DownloadInfo) Write(p []byte) (int, error) {
	n := len(p)
	x.Load += uint64(n)
	//fmt.Printf("\rDownloading... %.6f"+`%% complete`, float32(wc.Total)/float32(total))
	//fmt.Print(float32(wc.Total) / float32(total))
	return n, nil
}
