package commit2

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/xerrors"
	"io/ioutil"
	"log"
	"oplian/define"
	"oplian/model/lotus/request"
	"oplian/service/pb"
	"oplian/utils"
)

var Commit2ServiceApi = new(Commit2Service)

type Commit2Service struct {
}

// RunCommit2Execute C2 tasks
func (c *Commit2Service) RunCommit2(sp *pb.SealerParam) error {

	filePathC1 := sp.OpMainDisk + "/ipfs/data/c2task" + define.OpCsPathC1
	log.Println("RunC2Task sp:", sp.Sector.Id.Miner, sp.Sector.Id.Number)
	fileName := fmt.Sprintf("%s-%d.json", sp.Sector.Id.Miner, sp.Sector.Id.Number)
	c2FilePath := filePathC1 + "/" + fileName
	log.Println("c2FilePath 1:", c2FilePath)
	inb, err := ioutil.ReadFile(c2FilePath)
	if err != nil {

		fileName = fmt.Sprintf("s-%s-%d.json", sp.Sector.Id.Miner, sp.Sector.Id.Number)
		c2FilePath = filePathC1 + "/" + fileName
		log.Println("c2FilePath 2:", c2FilePath)
		inb, err = ioutil.ReadFile(c2FilePath)
		if err != nil {
			return xerrors.Errorf("RunC2Task reading input file: %w", err)
		}
	}

	var c2in request.Commit2In
	if err := json.Unmarshal(inb, &c2in); err != nil {
		return xerrors.Errorf("RunC2Task unmarshalling input file: %w", err)
	}

	sp.Phase1Out = c2in.Phase1Out
	number := sp.Sector.Id.Number
	miner, err := utils.FileCoinStrToUint64(sp.Sector.Id.Miner)
	if err != nil {
		return err
	}

	err = ImplRunCommit2(sp)
	if err != nil {
		return errors.New("ImplRunCommit2 err:" + err.Error())
	}
	log.Println("RunC2Task Start task：", miner, "|", number)

	return nil
}
