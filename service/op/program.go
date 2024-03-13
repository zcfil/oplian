package op

import (
	"log"
	"oplian/define"
	"oplian/service/pb"
	"os/exec"
)

type ProgramService struct{}

//StopService
//@author: nathan
//@function: StopService
//@description: Start-stop service
//@param: args *pb.RunStop
//@return: *pb.ResponseMsg, error

func (p *ProgramService) StopService(args *pb.RunStopType) (*pb.ResponseMsg, error) {
	var b []byte
	var err error
	switch define.ServiceType(args.ServiceType) {
	case define.ServiceLotus:
		//cmd := exec.Command("bash", "-c", "ps aux |grep "+define.ProgramPsLotus.String()+"|grep -v grep|awk '{print $2}'|xargs kill -9")
		if !args.IsRun {
			b, err = exec.Command("bash", "-c", "supervisorctl stop "+define.ProgramLotus.String()).CombinedOutput()
		} else {
			b, err = exec.Command("bash", "-c", "supervisorctl start "+define.ProgramLotus.String()).CombinedOutput()
		}
	case define.ServiceMiner:
		//cmd := exec.Command("bash", "-c", "ps aux |grep "+define.ProgramPsMiner.String()+"|grep -v grep|awk '{print $2}'|xargs kill -9")
		if !args.IsRun {
			b, err = exec.Command("bash", "-c", "supervisorctl stop "+define.ProgramMiner.String()).CombinedOutput()
		} else {
			b, err = exec.Command("bash", "-c", "supervisorctl start "+define.ProgramMiner.String()).CombinedOutput()
		}
	case define.ServiceWorkerTask:
		//cmd := exec.Command("bash", "-c", "ps aux |grep "+define.ProgramPsWorkerTask.String()+"|grep -v grep|awk '{print $2}'|xargs kill -9")
		if !args.IsRun {
			b, err = exec.Command("bash", "-c", "supervisorctl stop "+define.ProgramWorkerTask.String()).CombinedOutput()
		} else {
			b, err = exec.Command("bash", "-c", "supervisorctl start "+define.ProgramWorkerTask.String()).CombinedOutput()
		}
	case define.ServiceWorkerStorage:
		//cmd := exec.Command("bash", "-c", "ps aux |grep "+define.ProgramPsWorkerStorage.String()+"|grep -v grep|awk '{print $2}'|xargs kill -9")
		if !args.IsRun {
			b, err = exec.Command("bash", "-c", "supervisorctl stop "+define.ProgramWorkerStorage.String()).CombinedOutput()
		} else {
			b, err = exec.Command("bash", "-c", "supervisorctl start "+define.ProgramWorkerStorage.String()).CombinedOutput()
		}
	}
	log.Println(string(b))
	if err != nil {
		return &pb.ResponseMsg{Code: 500, Msg: err.Error()}, err
	}
	return &pb.ResponseMsg{Code: 200, Msg: "Successful ！"}, nil
}
