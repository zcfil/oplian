package oplocal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/multiformats/go-base32"
	"golang.org/x/xerrors"
	"log"
	"oplian/build"
	"oplian/config"
	"oplian/define"
	"oplian/global"
	"oplian/lotusrpc"
	"oplian/service/pb"
	"oplian/utils"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type OpLotusService struct{}

var (
	kstrPermissionMsg = "permissions of key: '%s' are too relaxed, " +
		"required: 0600, got: %#o"
	ErrClosedRepo = errors.New("repo is no longer open")
)

// GetWalletList
// @author: nathan
// @function: GetWalletList
// @description: Get a list of wallet addresses
// @param: info request.PageInfo
// @return: *pb.WalletList,error
func (deploy *OpLotusService) GetWalletList(args *pb.RequestConnect) (*pb.WalletList, error) {

	dir, err := os.Open(define.PathIpfsLotusKeystore)
	if err != nil {
		return nil, xerrors.Errorf("opening dir to list keystore: %w %s", err, define.PathIpfsLotusKeystore)
	}
	defer dir.Close() //nolint:errcheck
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, xerrors.Errorf("reading keystore dir: %w", err)
	}
	var wallets = []*pb.Wallet{}
	goRepeat := make(map[string]struct{})
	for _, f := range files {
		if f.Mode()&0077 != 0 {
			return nil, xerrors.Errorf(kstrPermissionMsg, f.Name(), f.Mode())
		}
		name, err := base32.RawStdEncoding.DecodeString(f.Name())
		if err != nil {
			return nil, xerrors.Errorf("decoding key: '%s': %w", f.Name(), err)
		}
		if !strings.Contains(string(name), define.WalletPrefix) {
			continue
		}
		addr := strings.Replace(string(name), define.WalletPrefix, "", 1)
		// Remove weight
		if _, ok := goRepeat[addr]; ok {
			continue
		}
		goRepeat[addr] = struct{}{}

		balance, err := lotusrpc.FullApi.WalletBalance(args.Token, args.Ip, addr)
		if err != nil {
			balance = 0.0
			log.Println(args.Token, args.Ip, addr)
			global.ZC_LOG.Error("Failed to obtain the wallet balance：" + err.Error())
		}
		wallets = append(wallets, &pb.Wallet{Address: addr, Balance: balance, OpId: global.OpUUID.String()})
	}
	return &pb.WalletList{Wallets: wallets}, nil
}

// RunNewLotus
// @author: nathan
// @function: RunNewLotus
// @description: Running lotus
// @param: args *pb.LotusInfo
// @return: error
func (p *OpLotusService) RunNewLotus(args *pb.LotusInfo) error {
	var outs []byte
	var err error
	var errMsg error
	if !utils.ExistFileOrDir(define.MainDisk) {
		return errors.New(define.MainDisk + "目录不存在！")
	}
	httpToken := ""

	exec.Command("bash", "-c", "supervisorctl stop lotus").CombinedOutput()

	if !utils.IsDir(define.PathIpfsData) {
		os.Remove(define.PathIpfsData)
	}
	if !utils.IsDir(define.PathIpfsLogs) {
		os.Remove(define.PathIpfsLogs)
	}
	_ = os.MkdirAll(define.PathIpfsLogs, 0664)
	_ = os.MkdirAll(define.PathIpfsLotusKeystore, 0600)
	_ = os.Rename(define.PathIpfsLotusKeystore, path.Join(define.PathIpfsData, "lotus_keystore_bak"))
	_ = os.RemoveAll(define.PathIpfsLotus)

	_ = os.MkdirAll(define.PathIpfsLotusKeystore, 0600)

	go func() {
		defer func() {

			if errMsg != nil {
				global.ZC_LOG.Error(errMsg.Error())
				_, err = global.OpToGatewayClient.UpdateLotus(context.Background(), &pb.ConnectInfo{Id: args.LotusId, DeployStatus: define.DeployFail.Int32(), SyncStatus: define.SyncFail.Int32(), RunStatus: define.RunStatusStop.Int32(), Token: httpToken, ErrMsg: errMsg.Error()})
				if err != nil {
					global.ZC_LOG.Error(err.Error())
					return
				}
			}
		}()

		if err1 := p.CheckProofsParameter(define.MainDisk); err1 != nil {

			if err = p.DownlodParameters(context.Background(), args.Ip); err != nil {
				errMsg = fmt.Errorf("Failed to download proof parameters：%s", err.Error())
				return
			}
		}

		switch define.LotusInitModel(args.ImportMode) {
		case define.LotusInitCopyModel:

			opInfo := []*pb.OpInfo{{Ip: args.Ip, Port: define.OpPort}}
			var res *pb.ResponseMsg
			if res, err = global.OpToGatewayClient.DownLoadFiles(context.Background(), &pb.DownLoadInfo{GateWayId: args.GateId, OpInfo: opInfo, FileInfo: []*pb.FileInfo{{FileName: args.FileName}}, DownloadPath: define.PathIpfsLotus}); err != nil {
				errMsg = fmt.Errorf("failed to download the height file. Procedure：%s", err.Error())
				return
			}
			if res.Code != 200 {
				errMsg = fmt.Errorf("failed to download the height file. Procedure：%s", res.Msg)
				return
			}

			if err = utils.UnzipLocal(filepath.Join(define.PathIpfsLotus, args.FileName), define.PathIpfsLotusDatastore); err != nil {
				errMsg = fmt.Errorf("解压高度文件失败：%s", err.Error())
				return
			}

			outs, err = exec.Command("bash", "-c", fmt.Sprintf(define.PathIpfsScriptRunLotus+" %s %s %d", define.LotusPort, args.Ip, args.LotusId)).CombinedOutput()
			if err != nil {
				errMsg = fmt.Errorf("初始化服务失败：%s,error：%s", string(outs), err.Error())
				return
			}

			outs, err = exec.Command("bash", "-c", "supervisorctl update").CombinedOutput()
			if err != nil {
				errMsg = fmt.Errorf("更新守护进程失败：%s,error: %s", string(outs), err.Error())
				return
			}
			exec.Command("bash", "-c", "supervisorctl start lotus").CombinedOutput()
			for {
				time.Sleep(time.Second)
				if httpToken == "" {
					token, err := os.ReadFile(define.PathIpfsLotusToken)
					if err != nil {
						global.ZC_LOG.Error(err.Error())
						continue
					}

					httpToken = string(token)
				}

				height, err := lotusrpc.FullApi.LotusHeight(httpToken, args.Ip)
				if err != nil {
					global.ZC_LOG.Error(err.Error())
					continue
				}

				if utils.BlockHeight()-height > 10 {
					log.Println("Common chain height：", utils.BlockHeight(), "，Synchronizing height：", height)
					continue
				}
				break
			}
			log.Println("Service start complete！", string(outs))

		case define.LotusInitImportModel:

			opInfo := []*pb.OpInfo{{Ip: args.Ip, Port: define.OpPort}}
			var res *pb.ResponseMsg
			if res, err = global.OpToGatewayClient.DownLoadFiles(context.Background(), &pb.DownLoadInfo{GateWayId: args.GateId, OpInfo: opInfo, FileInfo: []*pb.FileInfo{{FileName: args.FileName}}, DownloadPath: define.PathIpfsLotus}); err != nil {
				errMsg = fmt.Errorf("failed to download the snapshot file. Procedure：%s", err.Error())
				return
			}
			if res.Code != 200 {
				errMsg = fmt.Errorf("failed to download the snapshot file. Procedure：%s", res.Msg)
				return
			}
			spPath := filepath.Join(define.PathIpfsLotus, args.FileName)

			if strings.HasSuffix(spPath, ".zst") {
				outs, err = exec.Command("bash", "-c", fmt.Sprintf("zstd -d %s", spPath)).CombinedOutput()
				if err != nil || len(spPath) < 4 {
					errMsg = fmt.Errorf("zstd -d 快照文件失败：%s", res.Msg)
					return
				}
				spPath = spPath[:len(spPath)-4]

			}

			exec.Command("bash", "-c", "supervisorctl stop lotus").CombinedOutput()

			outs, err = exec.Command("bash", "-c", fmt.Sprintf(define.PathIpfsScriptImportLotus+" %s %s %s %d %s", define.LotusPort, spPath, args.Ip, args.LotusId, define.MainDisk)).CombinedOutput()
			if err != nil {
				errMsg = fmt.Errorf("初始化守护进程失败1：%s,error: %s", string(outs), err.Error())
				return
			}

			outs, err = exec.Command("bash", "-c", "supervisorctl update").CombinedOutput()
			if err != nil {
				errMsg = fmt.Errorf("failed to update the daemon：%s,error: %s", string(outs), err.Error())
				return
			}
			outs, _ = exec.Command("bash", "-c", "supervisorctl start lotus").CombinedOutput()
			//if err != nil {
			//	global.ZC_LOG.Error(err.Error())
			//	return
			//}

			for {
				time.Sleep(time.Second)
				if httpToken == "" {
					token, err := os.ReadFile(define.PathIpfsLotusToken)
					if err != nil {
						global.ZC_LOG.Error(err.Error())
						continue
					}

					httpToken = string(token)
				}

				height, err := lotusrpc.FullApi.LotusHeight(httpToken, args.Ip)
				if err != nil {
					global.ZC_LOG.Error(err.Error())
					continue
				}

				if utils.BlockHeight()-height > 10 {
					log.Println("Common chain height：", utils.BlockHeight(), "，Synchronizing height：", height)
					continue
				}

				if define.WalletNewMode(args.WalletNewMode) == define.WalletNew {
					for i := 0; i < int(args.BlsCount); i++ {
						_, err = lotusrpc.FullApi.WalletNew(httpToken, args.Ip, define.WalletTypeBls)
						if err != nil {
							errMsg = fmt.Errorf("wallet creation failure：%s,error: %s", string(outs), err.Error())
							return
						}
					}
					for i := 0; i < int(args.SecpCount); i++ {
						_, err = lotusrpc.FullApi.WalletNew(httpToken, args.Ip, define.WalletTypeSecp256k1)
						if err != nil {
							errMsg = fmt.Errorf("wallet creation failure：%s,error: %s", string(outs), err.Error())
							return
						}
					}
				}
				break
			}

			exec.Command("bash", "-c", "supervisorctl stop lotus").CombinedOutput()

			outs, err = exec.Command("bash", "-c", fmt.Sprintf(define.PathIpfsScriptRunLotus+" %s %s %s", define.LotusPort, args.Ip, define.MainDisk)).CombinedOutput()
			if err != nil {
				errMsg = fmt.Errorf("failed to initialize the daemon 2：%s,error: %s", string(outs), err.Error())
				return
			}

			outs, _ = exec.Command("bash", "-c", "supervisorctl start lotus").CombinedOutput()

		}

		_, err = global.OpToGatewayClient.UpdateLotus(context.Background(), &pb.ConnectInfo{Id: args.LotusId, DeployStatus: define.DeployFinish.Int32(), SyncStatus: define.SyncFinish.Int32(), RunStatus: define.RunStatusRunning.Int32(), Token: httpToken})
		if err != nil {
			global.ZC_LOG.Error(err.Error())
			return
		}
	}()

	time.Sleep(time.Millisecond * 100)
	if err != nil {
		return err
	}
	return nil
}

var once sync.Once

// DownlodParameters Download proof parameter
func (p *OpLotusService) DownlodParameters(ctx context.Context, toIp string) (reserr error) {

	if err := p.CheckProofsParameter(define.MainDisk); err == nil {
		global.ZC_LOG.Error("Prove that the parameter already exists")
		return nil
	}

	once.Do(func() {

		_ = os.RemoveAll(define.PathIpfsPARAMETER)
		_ = os.MkdirAll(define.PathIpfsPARAMETER, 0644)

		fileName, err := global.OpToGatewayClient.GetFileName(ctx, &pb.FileNameInfo{GateWayId: global.GateWayID.String(), FileType: define.ProveFile.Int64()})
		if err != nil {
			reserr = fmt.Errorf("获取证明参数失败：%s,%s", err.Error(), global.GateWayID.String())
			return
		}
		log.Println("Download the proof parameter file:", fileName, ",gatewayId:", global.GateWayID.String(), ",opInfo:", toIp)

		opInfo := []*pb.OpInfo{{Ip: toIp, Port: define.OpPort}}
		res, err := global.OpToGatewayClient.DownLoadFiles(ctx, &pb.DownLoadInfo{GateWayId: global.GateWayID.String(), OpInfo: opInfo, FileInfo: []*pb.FileInfo{{FileName: fileName.Value}}, DownloadPath: define.PathIpfsData})
		if err != nil {
			reserr = fmt.Errorf("%s", err.Error())
			return
		}
		if res.Code != 200 {
			reserr = fmt.Errorf("Msg：%s,Code：%d", res.Msg, res.Code)
			return
		}

		log.Println("Unzip the proof parameter file:", filepath.Join(define.PathIpfsData, fileName.Value), "，to：", define.PathMntMd0)
		err = utils.ExtractTarGz(filepath.Join(define.PathIpfsData, fileName.Value), define.PathMntMd0)
		if err != nil {
			reserr = fmt.Errorf("failed to decompress the proof parameter file：%s", err.Error())
			return
		}

		if err = p.CheckProofsParameter(define.MainDisk); err != nil {
			reserr = fmt.Errorf("Check the gateway proof parameters: %w", err)
			return
		}
	})
	return reserr
}

// CheckProofsParameter Check proof parameter
func (p *OpLotusService) CheckProofsParameter(mainDisk string) error {
	params := build.ParametersMap()
	for paranName, _ := range params {
		paramPath := filepath.Join(mainDisk, "/filecoin-proof-parameters", paranName)
		if !utils.FileExist(paramPath) {
			log.Println(fmt.Errorf("missing file：%s", paramPath).Error())
			return fmt.Errorf("missing file：%s", paramPath)
		}
	}
	return nil
}
func (p *OpLotusService) CheckProofsParameterPath(paramPath string) error {
	params := build.ParametersMap()
	for paramName, _ := range params {
		Path := filepath.Join(paramPath, paramName)
		if !utils.FileExist(Path) {
			log.Println(fmt.Errorf("missing file：%s", Path).Error())
			return fmt.Errorf("missing file：%s", Path)
		}
	}
	return nil
}

// RunNewMiner
// @author: nathan
// @function: RunNewMiner
// @description: Run New miner
// @param: args *pb.MinerInfo
// @return: error
func (p *OpLotusService) RunNewMiner(args *pb.MinerInfo) error {
	var ctx = context.Background()
	var err error

	if !utils.ExistFileOrDir(define.MainDisk) {
		return errors.New(define.MainDisk + "define.MainDisk Directory does not exist！")
	}

	//_ = os.RemoveAll(define.PathIpfsMiner)
	_ = os.Rename(define.PathIpfsMiner, define.PathIpfsMiner+strconv.FormatInt(time.Now().Unix(), 10))

	if !utils.IsDir(define.PathIpfsData) {
		os.Remove(define.PathIpfsData)
	}
	if !utils.IsDir(define.PathIpfsLogs) {
		os.Remove(define.PathIpfsLogs)
	}
	_ = os.MkdirAll(define.PathIpfsLogs, 0664)
	_ = os.MkdirAll(define.PathIpfsMiner, 0664)

	go func() {
		var errMsg error
		var minerToken string
		defer func() {

			if errMsg != nil {
				global.ZC_LOG.Error(errMsg.Error())
				if _, err = global.OpToGatewayClient.UpdateMiner(ctx, &pb.ConnectInfo{Id: args.MinerId, Token: minerToken, DeployStatus: define.DeployFail.Int32(), RunStatus: define.RunStatusStop.Int32(), Actor: args.Actor, ErrMsg: errMsg.Error()}); err != nil {
					global.ZC_LOG.Error(err.Error())
				}
			}

			if err = p.CheckProofsParameter(define.MainDisk); err != nil {
				if errMsg = p.DownlodParameters(ctx, args.Ip); errMsg != nil {
					return
				}
			}
		}()
		switch args.AddType {
		case define.MinerDepolyFile:
			errMsg = p.FileDeployMiner(ctx, args)
		case define.MinerDepolyNew, define.MinerDepolyWorker:
			errMsg = p.InitMiner(ctx, args)
		}
	}()
	time.Sleep(time.Millisecond * 100)
	if err != nil {
		return err
	}

	return nil
}

// FileDeployMiner Copy miner files
func (p *OpLotusService) FileDeployMiner(ctx context.Context, args *pb.MinerInfo) error {
	partition := "all"
	FULLNODE_API := utils.LotusApiInfoMerge(args.LotusToken, args.LotusIp, define.LotusPort)

	if args.IsWdpost && args.Partitions != "" {
		partition = args.Partitions
	}

	opInfo := []*pb.OpInfo{{Ip: args.Ip, Port: define.OpPort}}
	miner := define.MinerName + "_" + args.Actor + ".zip"

	res, err := global.OpToGatewayClient.DownLoadFiles(ctx, &pb.DownLoadInfo{OpInfo: opInfo, FileInfo: []*pb.FileInfo{{FileName: miner}}, DownloadPath: define.PathIpfsMiner})
	if err != nil {
		return fmt.Errorf("download failed0：%s %s", err.Error(), miner)
	}
	if res.Code != 200 {
		return fmt.Errorf("download failed1：%s %s", res.Msg, miner)
	}

	if err := utils.UnzipLocal(filepath.Join(define.PathIpfsMiner, miner), define.PathIpfsMiner); err != nil {
		return fmt.Errorf("decompression failure：%s", err.Error())
	}
	if err := p.CreateMinerTables(ctx, args.Actor, args.StorageType); err != nil {
		return err
	}

	os.Chmod(filepath.Join(define.PathIpfsMinerKeystore, define.JwtHmacSecret), 0600)
	os.Chmod(filepath.Join(define.PathIpfsMinerKeystore, define.Libp2pHost), 0600)

	os.Remove(filepath.Join(define.PathIpfsMiner, define.ConfigName))

	var b []byte

	cmd := exec.Command("bash", "-c", fmt.Sprintf("%s %v %v %v %s %s %s %s %s %s", define.PathIpfsScriptRunMiner, args.IsManage, args.IsWnpost, args.IsWdpost, partition, FULLNODE_API, args.Ip, define.MinerPort, args.Actor, define.MainDisk))
	b, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("初始化脚本失败：%s,error: %s", string(b), err.Error())
	}

	b, err = exec.Command("bash", "-c", "supervisorctl update").CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update the daemon：%s,error: %s", string(b), err.Error())
	}

	b, _ = exec.Command("bash", "-c", "supervisorctl start lotus-miner").CombinedOutput()
	log.Println("Start up complete；", string(b))
	return nil
}

// InitMiner initialize miner
func (p *OpLotusService) InitMiner(ctx context.Context, args *pb.MinerInfo) error {
	partition := "all"
	FULLNODE_API := utils.LotusApiInfoMerge(args.LotusToken, args.LotusIp, define.LotusPort)

	if args.IsWdpost && args.Partitions != "" {
		partition = args.Partitions
	}
	var cmd *exec.Cmd
	if args.AddType == define.MinerDepolyWorker {

		fil := &pb.FilParam{Param: args.Actor, Token: args.LotusToken, Ip: args.LotusIp}
		info, err := global.OpToGatewayClient.StateMinerInfo(ctx, fil)
		if err != nil {
			return fmt.Errorf("StateMinerInfo error:%s", err.Error())
		}
		fil.Param = info.Worker
		worker, err := global.OpToGatewayClient.StateAccountKey(ctx, fil)
		if err != nil {
			return fmt.Errorf("StateMinerInfo error:%s", err.Error())
		}
		cmd = exec.Command("bash", "-c", fmt.Sprintf("%s %s %s %s %s %s %s", define.PathIpfsScriptInitMinerWorker, worker.Value, args.Actor, define.MainDisk, FULLNODE_API, args.Ip, define.MinerPort))
	} else {

		cmd = exec.Command("bash", "-c", fmt.Sprintf("%s %s %d %s %s %s %s", define.PathIpfsScriptInitMiner, args.Owner, args.SectorSize, define.MainDisk, FULLNODE_API, args.Ip, define.MinerPort))
	}
	outs, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("initialization failure：%s,error: %s", string(outs), err.Error())
	}
	reg := regexp.MustCompile(`0\w*`)
	actor := reg.FindString(string(outs))
	if len(actor) == 0 {
		return fmt.Errorf("initialization failure：%s", string(outs))
	}
	actor = "f" + actor

	fmt.Println(actor)
	args.Actor = actor

	lotusminer := fmt.Sprintf(define.MinerName+"_%s", args.Actor)
	if _, err = global.OpToGatewayClient.FileOpSynGateWay(ctx, &pb.AddFileInfo{FileType: define.MinerFile.Int64(), ZipFileName: lotusminer, AddType: 3, OpId: global.OpUUID.String(), GateWayId: global.GateWayID.String()}); err != nil {
		return fmt.Errorf("failed to upload miner file：%s", err.Error())
	}
	if err = p.CreateMinerTables(ctx, args.Actor, args.StorageType); err != nil {
		return err
	}

	localProve := args.StorageType
	if args.StorageType != define.StorageTypeWorker {
		localProve = 1
	}

	outs, err = exec.Command("bash", "-c", fmt.Sprintf("%s %v %v %v %s %s %s %s %s %s %d", define.PathIpfsScriptRunMiner, args.IsManage, args.IsWnpost, args.IsWdpost, partition, FULLNODE_API, args.Ip, define.MinerPort, args.Actor, define.MainDisk, localProve)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to start miner：%s,error: %s", err.Error(), string(outs))
	}

	outs, err = exec.Command("bash", "-c", "supervisorctl update").CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to update the daemon：%s,error: %s", string(outs), err.Error())
	}
	outs, _ = exec.Command("bash", "-c", "supervisorctl start lotus-miner").CombinedOutput()
	log.Println(string(outs))
	return nil
}

func (p *OpLotusService) CreateMinerTables(ctx context.Context, actor string, ColonyType int32) (err error) {

	if _, err = global.OpToGatewayClient.CreateSectorTable(ctx, &pb.Actor{MinerId: actor}); err != nil {
		return fmt.Errorf("failed to create sector table：%s", err.Error())
	}

	if _, err = global.OpToGatewayClient.CreateSectorPieceTable(ctx, &pb.Actor{MinerId: actor}); err != nil {
		return fmt.Errorf("failed to create the sector order table：%s", err.Error())
	}

	if _, err = global.OpToGatewayClient.CreateSectorLogTable(ctx, &pb.Actor{MinerId: actor}); err != nil {
		return fmt.Errorf("failed to create the sector log table. Procedure：%s", err.Error())
	}

	if _, err = global.OpToGatewayClient.CreateSectorQueueDetailTable(ctx, &pb.Actor{MinerId: actor}); err != nil {
		return fmt.Errorf("failed to create the sector task queue detail table：%s", err.Error())
	}

	if _, err = global.OpToGatewayClient.AddColony(ctx, &pb.Colony{ColonyName: actor, ColonyType: ColonyType}); err != nil {
		return fmt.Errorf("failed to create the sector task queue detail table：%s", err.Error())
	}
	return nil
}

// RunMiner
// @author: nathan
// @function: RunMiner
// @description: Run miner
// @param: args *pb.MinerRun
// @return: error
func (p *OpLotusService) RunMiner(args *pb.MinerRun) error {
	partition := "all"
	full := utils.LotusApiInfoMerge(args.LotusToken, args.LotusIp, define.LotusPort)
	if args.IsWdpost && args.Partitions != "" {
		partition = args.Partitions
	}

	exec.Command("bash", "-c", fmt.Sprintf("supervisorctl stop %s", define.ProgramMiner)).CombinedOutput()

	outs, err := exec.Command("bash", "-c", fmt.Sprintf("%s %v %v %v %s %s %s %s %s %s", define.PathIpfsScriptRunMiner, args.IsManage, args.IsWnpost, args.IsWdpost, partition, full, args.Ip, define.MinerPort, args.Actor, define.MainDisk)).CombinedOutput()
	if err != nil {
		return err
	}
	log.Println(string(outs))

	outs, _ = exec.Command("bash", "-c", fmt.Sprintf("supervisorctl start %s", define.ProgramMiner)).CombinedOutput()
	log.Println(string(outs))
	return nil
}

// RunNewWorker
// @author: nathan
// @function: RunNewWorker
// @description: Run a new worker
// @param: args *pb.WorkerInfo
// @return:  error
func (p *OpLotusService) RunNewWorker(args *pb.WorkerInfo) error {
	if !utils.ExistFileOrDir(define.MainDisk) {
		return errors.New(define.MainDisk + "目录不存在！")
	}

	var errMsg error
	httpToken := ""
	go func() {
		defer func() {
			if errMsg != nil {
				for i := 0; i < 3600; i++ {

					if _, err := global.OpToGatewayClient.UpdateWorker(context.Background(), &pb.ConnectInfo{Id: args.Id, Token: httpToken, DeployStatus: define.DeployFail.Int32(), RunStatus: define.RunStatusStop.Int32(), ErrMsg: errMsg.Error()}); err != nil {
						global.ZC_LOG.Error(err.Error())
						if strings.Contains(err.Error(), "connection closed") {
							time.Sleep(time.Second)
							continue
						}
					}
					break
				}
			}
		}()


		StorageService{}.ClearWorker(define.PathIpfsWorker)
		//_ = os.Rename(dataPath.String(), dataPath.String()+strconv.FormatInt(time.Now().Unix(), 10))

		if !utils.IsDir(define.PathIpfsData) {
			os.Remove(define.PathIpfsData)
		}
		if !utils.IsDir(define.PathIpfsLogs) {
			os.Remove(define.PathIpfsLogs)
		}
		_ = os.MkdirAll(define.PathIpfsWorker, 0664)
		_ = os.MkdirAll(define.PathIpfsLogs, 0664)

		exec.Command("bash", "-c", fmt.Sprintf("supervisorctl stop %s", define.ProgramWorkerTask)).CombinedOutput()
		minerToken := utils.LotusApiInfoMerge(args.MinerToken, args.MinerIp, define.MinerPort)

		startNo, p2StartNo, p2StartNo1 := 0, 0, 0
		p2EndNo2, p2EndNo := 0, 0
		cupNo := runtime.NumCPU()
		endNo := cupNo - 1
		if cupNo <= 32 {
			endNo -= 4
			p2StartNo = endNo + 1
			p2EndNo = cupNo - 1
		} else {
			endNo -= 16
			p2StartNo = endNo + 1
			p2EndNo = p2StartNo + 7
			p2StartNo1 = p2EndNo + 1
			p2EndNo2 = p2EndNo + 8
		}

		outs, err := exec.Command("bash", "-c", fmt.Sprintf("%s %s %s %s %s %d %d", define.PathIpfsScriptRunWorker, define.WorkerPort, minerToken, define.MainDisk, global.ROOM_CONFIG.Gateway.IP+":"+define.SlotUnsealedPort, startNo, endNo)).CombinedOutput()
		if err != nil {
			errMsg = fmt.Errorf("Initialization script error：%s %v", string(outs), err)
			return
		}

		outs, err = exec.Command("bash", "-c", "supervisorctl update").CombinedOutput()
		if err != nil {
			errMsg = fmt.Errorf("更新守护进程失败：%s,error: %s", string(outs), err.Error())
			return
		}

		outs, _ = exec.Command("bash", "-c", fmt.Sprintf("supervisorctl start %s", define.ProgramWorkerTask)).CombinedOutput()

		portMap := make(map[string]string)
		if p2StartNo > 0 && p2EndNo > 0 {
			portMap[define.OpP2Port1] = fmt.Sprintf("%d,%d", p2StartNo, p2EndNo)
		}
		if p2StartNo1 > 0 && p2EndNo2 > 0 {
			portMap[define.OpP2Port2] = fmt.Sprintf("%d,%d", p2StartNo1, p2EndNo2)
		}

		for k, v := range portMap {

			vAr := strings.Split(v, ",")
			start, _ := strconv.Atoi(vAr[0])
			end, _ := strconv.Atoi(vAr[1])
			log.Println(fmt.Sprintf("CPUs:%d,P2 bind %d-%d", cupNo, start, end))
			outs, err = exec.Command("bash", "-c", fmt.Sprintf("%s %s %d %d %s", define.PathIpfsScriptRunWorkerP2, define.MainDisk, start, end, k)).CombinedOutput()
			if err != nil {
				errMsg = fmt.Errorf("worker-p2 init error：：%s %v", string(outs), err)
				return
			}

			outs, err = exec.Command("bash", "-c", "supervisorctl update").CombinedOutput()
			if err != nil {
				errMsg = fmt.Errorf("worker-p2 update：%s,error: %s", string(outs), err.Error())
				return
			}

			outs, _ = exec.Command("bash", "-c", fmt.Sprintf("supervisorctl start %s", fmt.Sprintf("p2-%s", v))).CombinedOutput()
		}



		//获取token、等待服务起来
		for i := 0; i < 600; i++ {
			if i == 599 {
				errMsg = fmt.Errorf("failed to start worker！")
				return
			}
			time.Sleep(time.Second)
			if httpToken == "" {
				token, err := os.ReadFile(define.PathIpfsWorkerToken)
				if err != nil {
					log.Println(define.PathIpfsWorkerToken, err.Error())
					global.ZC_LOG.Error(err.Error())
					continue
				}
				httpToken = string(token)
			}

			if _, err = global.OpToGatewayClient.UpdateWorker(context.Background(), &pb.ConnectInfo{Id: args.Id, Token: httpToken, DeployStatus: define.DeployFinish.Int32(), RunStatus: define.RunStatusRunning.Int32()}); err != nil {
				errMsg = fmt.Errorf("failed to update status！%s", err.Error())
				log.Println("The service is started. 3", i, err)
				return
			}
			break
		}
		disks := utils.GetOpDiskInfo()

		for i, disk := range disks {
			if utils.FilterIP(disk.Mounted) != "" {
				continue
			}
			path := filepath.Join(disk.Mounted, define.WorkerSeal)
			err = StorageService{}.StorageAttachInit(path, true, false)
			if err != nil {
				global.ZC_LOG.Error(err.Error())
				continue
			}

			log.Printf("Add encapsulated hard disk%d：%s\n", i, path)
			for {
				if err = lotusrpc.FullApi.StorageAddLocal(httpToken, args.Ip, path, define.TaskWorker); err != nil {
					if strings.Contains(err.Error(), "connect") {
						time.Sleep(time.Second)
						continue
					}
					global.ZC_LOG.Error(err.Error())
				}
				break
			}
		}
	}()
	time.Sleep(time.Millisecond * 100)
	return errMsg
}

// RunNewStorage
// @author: nathan
// @function: RunNewStorage
// @description: Adding a storage worker
// @param: args *pb.WorkerInfo
// @return:  error
func (p *OpLotusService) RunNewStorage(args *pb.WorkerInfo) error {
	if !utils.ExistFileOrDir(define.MainDisk) {
		return errors.New(define.MainDisk + "目录不存在！")
	}
	var errMsg error
	httpToken := ""
	go func() {
		defer func() {
			if errMsg != nil {
				for i := 0; i < 600; i++ {

					time.Sleep(time.Second)
					if _, err := global.OpToGatewayClient.UpdateStorage(context.Background(), &pb.ConnectInfo{Id: args.Id, Token: httpToken, DeployStatus: define.DeployFail.Int32(), RunStatus: define.RunStatusStop.Int32(), ErrMsg: errMsg.Error()}); err != nil {
						global.ZC_LOG.Error(err.Error())
						if strings.Contains(err.Error(), "connection closed") {
							continue
						}
					}
					break
				}
			}
		}()

		if err := p.CheckProofsParameter(define.MainDisk); err != nil {

			if err := p.DownlodParameters(context.Background(), args.Ip); err != nil {
				errMsg = fmt.Errorf("failed to download proof parameters：%s", err.Error())
				return
			}
		}

		_ = os.RemoveAll(define.PathIpfsStorage)

		if !utils.IsDir(define.PathIpfsData) {
			os.Remove(define.PathIpfsData)
		}
		if !utils.IsDir(define.PathIpfsLogs) {
			os.Remove(define.PathIpfsLogs)
		}
		_ = os.MkdirAll(define.PathIpfsLogs, 0664)
		_ = os.MkdirAll(define.PathIpfsStorage, 0664)

		exec.Command("bash", "-c", fmt.Sprintf("supervisorctl stop %s", define.ProgramWorkerStorage)).CombinedOutput()
		minerToken := utils.LotusApiInfoMerge(args.MinerToken, args.MinerIp, define.MinerPort)

		outs, err := exec.Command("bash", "-c", fmt.Sprintf("%s %s %s %s", define.PathIpfsScriptRunStorage, define.StoragePort, minerToken, define.MainDisk)).CombinedOutput()
		if err != nil {
			errMsg = fmt.Errorf("Initialization script error：%s %v", string(outs), err)
			return
		}

		outs, err = exec.Command("bash", "-c", "supervisorctl update").CombinedOutput()
		if err != nil {
			errMsg = fmt.Errorf("failed to update the daemon：%s,error: %s", string(outs), err.Error())
			return
		}

		outs, _ = exec.Command("bash", "-c", fmt.Sprintf("supervisorctl start %s", define.ProgramWorkerStorage)).CombinedOutput()
		log.Println("Service start complete：" + string(outs))
		//获取token、等待服务起来
		for i := 0; i < 600; i++ {
			if i == 599 {
				errMsg = fmt.Errorf("failed to start storage！")
				return
			}
			time.Sleep(time.Second)
			if httpToken == "" {
				token, err := os.ReadFile(define.PathIpfsStorageToken)
				if err != nil {
					log.Println(err.Error(), define.PathIpfsStorageToken)
					global.ZC_LOG.Error(err.Error())
					continue
				}
				httpToken = string(token)
			}
			//修改状态
			if _, err = global.OpToGatewayClient.UpdateStorage(context.Background(), &pb.ConnectInfo{Id: args.Id, Token: httpToken, DeployStatus: define.DeployFinish.Int32(), RunStatus: define.RunStatusRunning.Int32()}); err != nil {

				errMsg = fmt.Errorf("failed to update status！%s", err.Error())
				return
			}
			break

		}

		disks := utils.GetOpDiskInfo()
		for _, disk := range disks {
			path := filepath.Join(disk.Mounted, define.StoragePath)
			err = StorageService{}.StorageAttachInit(path, false, true)
			if err != nil {
				global.ZC_LOG.Error(err.Error())
				continue
			}

			for {
				if err = lotusrpc.FullApi.StorageAddLocal(httpToken, args.Ip, path, define.StorageWorker); err != nil {
					if strings.Contains(err.Error(), "connect") {
						time.Sleep(time.Second)
						continue
					}
					global.ZC_LOG.Error(err.Error())
				}
				break
			}
		}
	}()
	time.Sleep(time.Millisecond * 100)
	return errMsg
}

// AddNodeStorage
// @author: lex
// @function: AddNodeStorage
// @description: Node storage configuration file Add the mounted nfs directory
// @param: args *pb.WorkerInfo
// @return:  error
func (p *OpLotusService) AddNodeStorage(args *pb.StorageInfo) error {

	if !utils.ExistFileOrDir(define.MainDisk) {
		return errors.New(define.MainDisk + "Directory does not exist！")
	}

	token, err := os.ReadFile(define.PathIpfsMinerToken)
	if err != nil {
		log.Println(err.Error(), define.PathIpfsMinerToken)
		global.ZC_LOG.Error(err.Error())
		return err
	}
	httpToken := string(token)
	var errMsg error

	args.MountDir = utils.DealMountDir(args.MountDir, args.StorageIp)

	mountDirs := strings.Split(args.MountDir, ",")

	for _, disk := range mountDirs {
		path := filepath.Join(disk, define.StoragePath)
		err := StorageService{}.StorageAttachInit(path, false, true)
		if err != nil {
			global.ZC_LOG.Error(err.Error())
			continue
		}

		for {
			if err = lotusrpc.FullApi.StorageAddLocal(httpToken, args.NodeIp, path, define.MinerWorker); err != nil {
				if strings.Contains(err.Error(), "connect") {
					time.Sleep(time.Second)
					continue
				}
				global.ZC_LOG.Error(err.Error())
			}
			break
		}
	}

	log.Println("AddNodeStorage success")

	time.Sleep(time.Millisecond * 100)
	return errMsg
}

// RunWorker
// @author: nathan
// @function: RunWorker
// @description: Run worker
// @param: args *pb.FilParam
// @return: error
func (p *OpLotusService) RunWorker(args *pb.FilParam) error {
	var runPath string
	var pg define.ProgramName
	var port string
	workerType, _ := strconv.Atoi(args.Param)
	switch define.WorkerType(workerType) {
	case define.TaskWorker:
		runPath = define.PathIpfsScriptRunWorker
		pg = define.ProgramWorkerTask
		port = define.WorkerPort
	case define.StorageWorker:
		runPath = define.PathIpfsScriptRunStorage
		pg = define.ProgramWorkerStorage
		port = define.StoragePort
	default:
		msg := "未知worker类型"
		return errors.New(msg)
	}

	exec.Command("bash", "-c", fmt.Sprintf("supervisorctl stop %s", pg)).CombinedOutput()
	minerToken := utils.LotusApiInfoMerge(args.Token, args.Ip, define.MinerPort)

	cmd := exec.Command("bash", "-c", fmt.Sprintf("%s %s %s %s %s", runPath, port, minerToken, define.MainDisk, global.ROOM_CONFIG.Gateway.IP+":"+define.SlotUnsealedPort))
	b, err := cmd.CombinedOutput()
	if err != nil {
		global.ZC_LOG.Error(fmt.Sprintf("初始化脚本错误：%s %v", string(b), err))
		return err
	}

	b, err = exec.Command("bash", "-c", fmt.Sprintf("supervisorctl start %s", pg)).CombinedOutput()
	if err != nil {
		global.ZC_LOG.Error(string(b) + err.Error())
		return err
	}
	log.Println("Service start complete：" + string(b))

	return nil
}

// RunBoost
// @author: nathan
// @function: RunBoost
// @description: Run boost
// @param: args *pb.BoostInfo
// @return: error
func (p *OpLotusService) RunBoost(args *pb.BoostInfo) error {
	listenIP := args.LanIp
	listenPort := args.LanPort
	if args.NetworkType == define.NetworkPub.Int32() {
		listenIP = args.InternetIp
		listenPort = args.InternetPort
	}
	go func() {
		var errMsg error
		boostToken := ""
		defer func() {
			if errMsg != nil {
				for i := 0; i < 3600; i++ {
					if _, err := global.OpToGatewayClient.UpdateBoost(context.Background(), &pb.ConnectInfo{Id: args.Id, Token: boostToken, DeployStatus: define.DeployFail.Int32(), RunStatus: define.RunStatusStop.Int32(), ErrMsg: errMsg.Error()}); err != nil {
						log.Println(err.Error())
						if strings.Contains(err.Error(), "connection closed") {
							continue
						}
					}
					break
				}
			}
		}()
		if !utils.IsDir(define.PathIpfsData) {
			os.Remove(define.PathIpfsData)
			_ = os.MkdirAll(define.PathIpfsData, 0664)
		}
		if !utils.IsDir(define.PathIpfsLogs) {
			os.Remove(define.PathIpfsLogs)
			_ = os.MkdirAll(define.PathIpfsLogs, 0664)
		}
		if utils.FileExist(define.PathIpfsBoostToken) {

			var boostConfig config.Boost
			if _, err := toml.DecodeFile(define.PathIpfsBoostConfig, &boostConfig); err != nil {
				errMsg = fmt.Errorf("failed to read the configuration file：%s", err.Error())
				log.Println(fmt.Sprintf("Failed to read the configuration file：%s", err.Error()))
				return
			}
			boostConfig.Libp2p.ListenAddresses = []string{
				fmt.Sprintf("/ip4/%s/tcp/%s", listenIP, listenPort),
				fmt.Sprintf("/ip6/::/tcp/%s", listenPort),
			}
			boostConfig.Libp2p.AnnounceAddresses = []string{fmt.Sprintf("/ip4/%s/tcp/%s", listenIP, listenPort)}
			buf := new(bytes.Buffer)
			e := toml.NewEncoder(buf)
			if err := e.Encode(boostConfig); err != nil {
				errMsg = fmt.Errorf("description Failed to modify the configuration file：%s", err.Error())
				return
			}
			os.Remove(define.PathIpfsBoostConfig)
			err := os.WriteFile(define.PathIpfsBoostConfig, []byte(buf.String()), 0644)
			if err != nil {
				fmt.Printf("writing config file: %v", err)
			}

			outs, err := exec.Command("bash", "-c", fmt.Sprintf("supervisorctl restart %s", define.ProgramBoost)).CombinedOutput()
			if err != nil {
				log.Println(fmt.Sprintf("Process restart failed：%s,error: %s", string(outs), err.Error()))
				return
			}

			if _, err = global.OpToGatewayClient.UpdateBoost(context.Background(), &pb.ConnectInfo{Id: args.Id, Token: boostToken, DeployStatus: define.DeployFinish.Int32(), RunStatus: define.RunStatusRunning.Int32()}); err != nil {
				errMsg = fmt.Errorf("更新状态失败！%s", err.Error())
				return
			}
			log.Println("Restart service complete：" + string(outs))
		} else {

			exec.Command("bash", "-c", fmt.Sprintf("supervisorctl stop %s", define.ProgramBoost)).CombinedOutput()
			minerToken := utils.LotusApiInfoMerge(args.MinerToken, args.MinerIp, define.MinerPort)
			lotusToken := utils.LotusApiInfoMerge(args.LotusToken, args.LotusIp, define.LotusPort)
			// 获取存储信息
			outs, err := exec.Command("bash", "-c", fmt.Sprintf("%s %s %s %s %s %s %s", define.PathIpfsScriptRunBoost, args.WorkerWallet, lotusToken, minerToken, listenIP, listenPort, define.MainDisk)).CombinedOutput()
			if err != nil {
				errMsg = fmt.Errorf("Initialization script error：%s %v", string(outs), err)
				return
			}

			outs, err = exec.Command("bash", "-c", "supervisorctl update").CombinedOutput()
			if err != nil {
				errMsg = fmt.Errorf("failed to update the daemon：%s,error: %s", string(outs), err.Error())
				return
			}

			outs, err = exec.Command("bash", "-c", fmt.Sprintf("supervisorctl start %s", define.ProgramBoost)).CombinedOutput()
			if err != nil {
				log.Println(err.Error())
			}
			ID, err := lotusrpc.FullApi.PreeID(args.MinerToken, args.MinerIp, define.MinerPort)
			if err != nil {
				errMsg = fmt.Errorf("failed to obtain the PreeID. Procedure：%s", err.Error())
				log.Println(err.Error())
				return
			}
			//获取token、等待服务起来
			for i := 0; i < 600; i++ {
				if i == 599 {
					errMsg = fmt.Errorf("failed to start boost！")
					return
				}
				time.Sleep(time.Second)
				if boostToken == "" {
					token, err := os.ReadFile(define.PathIpfsBoostToken)
					if err != nil {
						log.Println(err.Error())
						continue
					}
					boostToken = string(token)
				}
				if _, err = lotusrpc.FullApi.ClientQueryAsk(args.LotusToken, args.LotusIp, ID, args.Actor); err != nil {

					//宣告公网IP
					minerPath := filepath.Join(define.PathIpfsProgram, "lotus-miner")
					outs, err = exec.Command("bash", "-c", fmt.Sprintf("%s actor set-addrs /ip4/%s/tcp/%s", minerPath, args.InternetIp, args.InternetPort)).CombinedOutput()
					if err != nil {
						errMsg = fmt.Errorf("declare the public IP failed：%s,error: %s", string(outs), err.Error())
						log.Println(err.Error())
						return
					}
					if !strings.Contains(string(outs), define.MsgPrefix) {
						log.Println(string(outs))
						errMsg = fmt.Errorf("declare the public IP failed：%s", string(outs))
						return
					}

					outs, err = exec.Command("bash", "-c", fmt.Sprintf("%s actor set-peer-id %s", minerPath, ID)).CombinedOutput()
					if err != nil {
						errMsg = fmt.Errorf("failed to declare peer address：%s,error: %s", string(outs), err.Error())
						log.Println(err.Error())
						return
					}
					if !strings.Contains(string(outs), define.MsgPrefix) {
						log.Println(string(outs))
						errMsg = fmt.Errorf("failed to declare peer address：%s", string(outs))
						return
					}
				}
				if errMsg == nil {

					if _, err = global.OpToGatewayClient.UpdateBoost(context.Background(), &pb.ConnectInfo{Id: args.Id, Token: boostToken, DeployStatus: define.DeployFinish.Int32(), RunStatus: define.RunStatusRunning.Int32()}); err != nil {
						errMsg = fmt.Errorf("failed to update status！%s", err.Error())
						return
					}
				}
				break
			}
			log.Println("Service start complete：" + string(outs))
		}
	}()

	return nil
}
