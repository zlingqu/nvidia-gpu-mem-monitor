package main

import (
	"fmt"
	"net/http"
	"github.com/zlingqu/nvidia-gpu-mem-monitor/dockercli"
	"github.com/zlingqu/nvidia-gpu-mem-monitor/myexec"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/metrics", func(c *gin.Context) {
		r := httpGetRespon()
		c.String(http.StatusOK, r)
	})
	r.Run(":80")

}

func httpGetRespon() string {

	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.37", nil, nil) //使用socket通信
	if err != nil {
		panic(err)
	}

	records := myexec.GetExecOutByCSV("nvidia-smi --query-compute-apps=pid,used_gpu_memory,gpu_name,gpu_uuid --format=csv,noheader,nounits")

	var respon string = `# HELP pod_used_gpu_mem_MB . Pod使用的GPU显存大小
# TYPE pod_used_gpu_mem_MB gauge
`

	for _, row := range records {
		cmd := "cat /proc/" + row[0] + "/cgroup |head -1 | awk -F'/' '{print $5}'"
		containID := myexec.GetExecOutByString(cmd)
		podName, podNamespace := "null", "null" //非pod使用gpu的进程
		if containID != "" {
			podName, podNamespace = dockercli.GetContainsPodInfo(cli, containID) //获取pod信息
		}
		respon = fmt.Sprintf("%spod_used_gpu_mem_MB{app_pid=\"%s\",gpu_name=\"%s\",gpu_uuid=\"%s\",pod_name=\"%s\",pod_namespace=\"%s\"} %s\n",
			respon, row[0], row[2], row[3], podName, podNamespace, row[1])
	}
	return respon
}
