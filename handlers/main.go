package handlers

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/docker/docker/client"
	svc "github.com/zlingqu/nvidia-gpu-mem-monitor/service"
)

// Metrics 提供metrics接口
func Metrics() string {

	cli, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"), client.WithVersion("v1.38")) //使用socket通信

	if err != nil {
		log.Print("docker client 初始化错误" + err.Error())
		return ""
	}

	if cli == nil {
		log.Print("docker client 初始化错误")
		return ""
	}
	defer cli.Close() //记得释放

	records := svc.GetExecOutByCSV("nvidia-smi --query-compute-apps=pid,used_gpu_memory,gpu_name,gpu_uuid --format=csv,noheader,nounits")
	/*
		31756, 1267, GeForce GTX 1080 Ti, GPU-78d64296-8254-ef39-35ec-cb35bd6e6192
		25580, 753, GeForce GTX 1080 Ti, GPU-78d64296-8254-ef39-35ec-cb35bd6e6192
	*/
	gpuLists := svc.GetExecOutByCSV("nvidia-smi -L|awk '{print $NF,$2}'|sed 's/)/,/g'|sed 's/://g'")
	/*
		GPU-78d64296-8254-ef39-35ec-cb35bd6e6192, 0
		GPU-2b8215f8-eb7c-0ae4-328f-a678a84f8d08, 1
		GPU-9d5d5439-4397-7189-1a46-801b59248301, 2
		GPU-55da7249-18e7-c3e7-beb8-4e1f661f5461, 3
	*/
	response := `# HELP pod_used_gpu_mem_MB . Pod使用的GPU显存大小
# TYPE pod_used_gpu_mem_MB gauge
`
	gpu := ""
	for _, row := range records {
		cmd := "cat /proc/" + row[0] + "/cgroup |head -1 | awk -F'/' '{print $NF}'"
		containID := svc.GetExecOutByString(cmd)
		podName, podNamespace := "服务器直接运行的程序", "null" //非pod使用gpu的进程
		if containID != "" {
			podName, podNamespace = svc.GetContainsPodInfo(cli, containID) //获取pod信息
			if podName == "" || podNamespace == "" {                       //排除docker run起来的进程
				podName = "docker run运行的程序"
				podNamespace = "null"
			}
		}

		for _, gpuOne := range gpuLists {
			if gpuOne[0] == row[3] {
				gpu = gpuOne[1]
				break
			}
		}
		response = fmt.Sprintf("%spod_used_gpu_mem_MB{hostIP=\"%s\",app_pid=\"%s\",gpu_name=\"%s\",UUID=\"%s\",gpu=\"%s\",pod=\"%s\",namespace=\"%s\"} %s\n",
			response, getIP(), row[0], row[2], row[3], gpu, podName, podNamespace, row[1])
	}
	return response
}

func getIP() string {
	if hostIP := os.Getenv("hostIP"); hostIP != "" { //如果部署到k8s中会注入hostIP变量
		return hostIP
	}
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("net.Interfaces failed, err:", err.Error())
		return ""
	}
	for i := 0; i < len(netInterfaces); i++ {
		//fmt.Println(netInterfaces[i],net.FlagUp)
		if (netInterfaces[i].Flags&net.FlagUp) != 0 && interFaceFields(netInterfaces[i]) {
			adds, _ := netInterfaces[i].Addrs()

			for _, address := range adds {
				//fmt.Println(address)
				if inet, ok := address.(*net.IPNet); ok && !inet.IP.IsLoopback() {
					if inet.Contains(inet.IP) && inet.IP.To4() != nil {
						return inet.IP.String()
					}
				}
			}
		}
	}
	return ""
}

func interFaceFields(myInterFace net.Interface) bool {
	if myInterFace.MTU != 1500 {
		return false
	}
	if len(myInterFace.HardwareAddr) > 17 { //排除ib网络的网卡
		return false
	}
	for _, v := range []string{"cni0", "flannel.1", "docker0", "virbr0"} { //排除特殊的网卡设备
		if myInterFace.Name == v {
			return false
		}
	}
	return true
}
