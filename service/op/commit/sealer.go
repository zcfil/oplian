package commit

import (
	"context"
	"errors"
	"fmt"
	uuidGo "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"oplian/define"
	"oplian/global"
	"oplian/service/pb"
	"oplian/utils"
	"os"
	"strconv"
	"sync"
	"time"
)

var SealerService = new(Sealer)

var Sl = Sealer{
	SealPort: make(map[string]string),
}

type Sealer struct {
	SealPort   map[string]string
	SealPortRL sync.RWMutex
}

func (s *Sealer) SetSealPort(key, val string) {
	Sl.SealPortRL.Lock()
	defer Sl.SealPortRL.Unlock()
	Sl.SealPort[key] = val
}

func (s *Sealer) GetSealPort(key string) string {
	Sl.SealPortRL.RLock()
	defer Sl.SealPortRL.RUnlock()
	return Sl.SealPort[key]
}

// RunOpC2Client Start the opC2 client
func (s *Sealer) RunOpC2Client() error {

	time.Sleep(time.Minute)
	total, err := global.OpToGatewayClient.GetC2WorkerInfo(context.TODO(), &pb.String{Value: global.OpUUID.String()})
	if err != nil {
		log.Println("RunOpC2Client GetC2WorkerInfo err:", err)
		return err
	}

	if total.Value != "1" {
		msg := fmt.Sprintf("%s,This machine is not a C2_worker machine and op_c2 does not need to be started", global.OpUUID.String())
		log.Println(msg)
		return errors.New(msg)
	} else {

		_, err = global.OpToGatewayClient.RedoC2Task(context.TODO(), &pb.String{Value: global.OpUUID.String()})
		if err != nil {
			log.Println("RunOpC2Client RedoC2Task err:", err)
			return err
		}
	}

	script := "killall -9 oplian-op-c2 oplian-sectors-seal"
	_, err = utils.ExecuteScript(script)
	if err != nil {
		log.Println("kill opc2 ExecuteScript err:", err)
	}

	for {
		checkOpC2 := utils.ProgramRun("oplian-op-c2")
		if checkOpC2 {
			time.Sleep(time.Minute)
			continue
		}
		break
	}

	for {
		checkOpC2 := utils.ProgramRun("oplian-sectors-seal")
		if checkOpC2 {
			time.Sleep(time.Minute)
			continue
		}
		break
	}

	var opServer utils.OpServer
	opServer.GetGPUInfo()
	num := len(opServer.GPU.Gpus)

	if num == utils.ZERO {
		return errors.New("the GPU is not enabled on the machine")
	}

	port, _ := strconv.Atoi(define.OpC2Port)
	port1, _ := strconv.Atoi(define.OpSectorC2Port)

	Ram, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal("Memory acquisition failure！", err)
	}

	Mem := 0
	runCount := 0
	nvCount := 1
	runGpuCount := num * nvCount
	runRamCount := Ram.Available / define.Ss234GiB
	if int(runRamCount) >= runGpuCount {
		Mem = runGpuCount
	} else {
		Mem = int(runRamCount)
	}

	if Mem%2 == 0 {
		nvCount = Mem / num
	} else if Mem < 4 {
		nvCount = 2
	}

	log.Println(fmt.Sprintf("Number of Gpus :%d, memory size :%d, number of opC2 starts:%d,nvCount:%d", num, Ram.Available, Mem, nvCount))
	for j := 0; j < num; j++ {

		os.MkdirAll("/tmp/gpu"+strconv.Itoa(j), 0664)
		os.MkdirAll("/tmp/"+strconv.Itoa(j), 0664)
		runNv := 0
		for i := 0; i < Mem; i++ {

			script := fmt.Sprintf(define.PathIpfsScriptRunOpC2+" %d %s", port, define.MainDisk)
			_, err := utils.ExecuteScript(script)
			if err != nil {
				log.Println("run opc2 ExecuteScript err:", err)
			}
			log.Println("Start OpC2:", script)
			script = fmt.Sprintf(define.PathIpfsScriptRunSectorC2+" %d %d %s", j, port1, define.MainDisk)
			_, err = utils.ExecuteScript(script)
			if err != nil {
				log.Println("run sectorSeal ExecuteScript err:", err)
			}
			log.Println("Start the sector seal task:", script)
			Sl.SetSealPort(fmt.Sprintf("%d", port), fmt.Sprintf("%d", port1))
			port++
			port1++
			runCount++
			runNv++

			if Mem == runCount || runNv == nvCount {
				break
			}
		}

		if Mem == runCount {
			break
		}
	}

	log.Println("Starting OpC2 is complete")
	return nil
}

// StopOpC2Client Stop the opC2 client
func (s *Sealer) StopOpC2Client() error {

	script := "killall -9 oplian-op-c2 oplian-sectors-seal"
	_, err := utils.ExecuteScript(script)
	if err != nil {
		log.Println("killall opc2 ExecuteScript err:", err)
		return err
	}

	return nil
}

// InitOpC2Uid Initialize the opc2 uid
func (s *Sealer) InitOpC2Uid() {
	global.OpC2UUID = uuidGo.NewV4()
}
