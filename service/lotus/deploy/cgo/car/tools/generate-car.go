package tools

import (
	"bufio"
	"context"
	"io"
	"os"
	"path"

	"oplian/service/lotus/deploy/cgo/car/util"

	commcid "github.com/filecoin-project/go-fil-commcid"
	commp "github.com/filecoin-project/go-fil-commp-hashhash"
	"github.com/google/uuid"
	cid "github.com/ipfs/go-cid/_rsrch/cidiface"
	cbor "github.com/ipfs/go-ipld-cbor"
)

type CommpResult struct {
	commp     string
	pieceSize uint64
}

type Result struct {
	Ipld      *util.FsNode
	DataCid   string
	PieceCid  string
	PieceSize uint64
	CarSize   int64
	CidMap    map[string]util.CidMapValue
}

type Input []util.Finfo

type CarHeader struct {
	Roots   []cid.Cid
	Version uint64
}

func init() {
	cbor.RegisterCborType(CarHeader{})
}

const BufSize = (4 << 20) / 128 * 127

func GenerateCar(ctx context.Context, outDir, tmpDir, parent string, input Input) (*Result, error) {
	outFilename := uuid.New().String() + ".car"
	outPath := path.Join(outDir, outFilename)
	carF, err := os.Create(outPath)
	if err != nil {
		return nil, err
	}
	cp := new(commp.Calc)
	writer := bufio.NewWriterSize(io.MultiWriter(carF, cp), BufSize)
	ipld, cid, cidMap, err := util.GenerateCar(ctx, input, parent, tmpDir, writer)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	err = carF.Close()
	if err != nil {
		return nil, err
	}
	rawCommP, pieceSize, err := cp.Digest()
	if err != nil {
		return nil, err
	}
	commCid, err := commcid.DataCommitmentV1ToCID(rawCommP)
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(outPath)
	if err != nil {
		return nil, err
	}
	carSize := stat.Size()
	err = os.Rename(outPath, path.Join(outDir, commCid.String()+".car"))
	if err != nil {
		return nil, err
	}

	return &Result{
		Ipld:      ipld,
		DataCid:   cid,
		PieceCid:  commCid.String(),
		PieceSize: pieceSize,
		CarSize:   carSize,
		CidMap:    cidMap,
	}, nil
}

func GenerateCarNew(outDir, tmpDir, parent string, input Input) (*Result, error) {
	outFilename := uuid.New().String() + ".car"
	outPath := path.Join(outDir, outFilename)
	carF, err := os.Create(outPath)
	if err != nil {
		return nil, err
	}
	cp := new(commp.Calc)
	writer := bufio.NewWriterSize(io.MultiWriter(carF, cp), BufSize)
	ipld, cid, cidMap, err := util.GenerateCarNew(input, parent, tmpDir, writer)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	err = carF.Close()
	if err != nil {
		return nil, err
	}
	rawCommP, pieceSize, err := cp.Digest()
	if err != nil {
		return nil, err
	}
	commCid, err := commcid.DataCommitmentV1ToCID(rawCommP)
	if err != nil {
		return nil, err
	}
	stat, err := os.Stat(outPath)
	if err != nil {
		return nil, err
	}
	carSize := stat.Size()
	err = os.Rename(outPath, path.Join(outDir, commCid.String()+".car"))
	if err != nil {
		return nil, err
	}

	return &Result{
		Ipld:      ipld,
		DataCid:   cid,
		PieceCid:  commCid.String(),
		PieceSize: pieceSize,
		CarSize:   carSize,
		CidMap:    cidMap,
	}, nil
}
