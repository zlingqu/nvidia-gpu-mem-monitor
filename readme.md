## gpu节点加入k8s集群，做AI分析是一种比较常见的情景，当pod使用GPU卡时，通常有以下两种方案
## 一、安装nvidia官方插件
此时资源分配时，是按照==卡的数量==进行资源分配，k8s的yml文件类似于如下
```yaml

          resources:
            limits:
              nvidia.com/gpu: '1'
            requests:
              nvidia.com/gpu: '1'

```
其中1表示使用1张GPU卡。


使用prometheus从cAdvisor总获取的有效监控项有：

```bash
container_accelerator_memory_total_bytes #pod所在卡的显存总大小
container_accelerator_memory_used_bytes  #pod所在卡的显存使用
```



## 二、安装第三方插件

比如阿里云的插件，git地址：https://github.com/AliyunContainerService/gpushare-device-plugin/

此时资源分配时，是可以按照==卡的显存大小==进行资源分配，k8s的yml文件类似于如下

```yaml
          resources:
            limits:
              aliyun.com/gpu-mem: '4'
            requests:
              aliyun.com/gpu-mem: '2'
```
其中2、4表示使用的显存大小
查看现场资源分配
```
# kubectl-inspect-gpushare 
NAME          IPADDRESS     GPU0(Allocated/Total)  GPU1(Allocated/Total)  GPU2(Allocated/Total)  GPU3(Allocated/Total)  GPU4(Allocated/Total)  GPU5(Allocated/Total)  GPU6(Allocated/Total)  GPU7(Allocated/Total)  GPU Memory(GiB)
192.168.3.4   192.168.3.4   4/11                   6/11                   10/11                  0/11                   0/11                   0/11                   0/11                   0/11                   20/88
192.168.68.4  192.168.68.4  4/10                   0/10                   10/10                  9/10                   0/10                   10/10                  6/10                   6/10                   45/80
-------------------------------------------------------------------------------------------
Allocated/Total GPU Memory In Cluster:
65/168 (38%)  
-------------------------------------------------------------------------------------------
Allocated/Total GPU Memory In Cluster:
65/168 (38%)  
```



使用prometheus从cAdvisor总获取的有效监控项有：

```bash
container_accelerator_memory_total_bytes #pod所在卡的显存总大小
container_accelerator_memory_used_bytes  #pod所在卡的显存使用
```

## 注意
使用第二种方式时container_accelerator_memory_used_bytes获取的仍然pod所在卡的显存使用情况，而不是这个pod使用的显存情况。

当一张卡上跑多个进程时，此时获取到的数据是失真的，比如下面，第2张第6张卡都跑了2个进程

```
# nvidia-smi 
Thu Dec 17 09:36:05 2020       
+-----------------------------------------------------------------------------+
| NVIDIA-SMI 450.57       Driver Version: 450.57       CUDA Version: 11.0     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|                               |                      |               MIG M. |
|===============================+======================+======================|
|   0  TITAN V             Off  | 00000000:1A:00.0 Off |                  N/A |
| 29%   42C    P8    26W / 250W |   2639MiB / 12066MiB |      0%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+
|   1  TITAN V             Off  | 00000000:1B:00.0 Off |                  N/A |
| 28%   42C    P8    26W / 250W |      4MiB / 12066MiB |      0%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+
|   2  TITAN V             Off  | 00000000:3D:00.0 Off |                  N/A |
| 28%   39C    P2    36W / 250W |   8702MiB / 12066MiB |      1%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+
|   3  TITAN V             Off  | 00000000:3E:00.0 Off |                  N/A |
| 28%   40C    P8    27W / 250W |      4MiB / 12066MiB |      0%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+
|   4  TITAN V             Off  | 00000000:88:00.0 Off |                  N/A |
| 28%   37C    P8    27W / 250W |      4MiB / 12066MiB |      0%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+
|   5  TITAN V             Off  | 00000000:89:00.0 Off |                  N/A |
| 28%   39C    P8    26W / 250W |      4MiB / 12066MiB |      0%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+
|   6  TITAN V             Off  | 00000000:B1:00.0 Off |                  N/A |
| 33%   48C    P2    43W / 250W |   7778MiB / 12066MiB |     16%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+
|   7  TITAN V             Off  | 00000000:B2:00.0 Off |                  N/A |
| 28%   38C    P8    25W / 250W |      4MiB / 12066MiB |      0%      Default |
|                               |                      |                  N/A |
+-------------------------------+----------------------+----------------------+
                                                                               
+-----------------------------------------------------------------------------+
| Processes:                                                                  |
|  GPU   GI   CI        PID   Type   Process name                  GPU Memory |
|        ID   ID                                                   Usage      |
|=============================================================================|
|    0   N/A  N/A    159584      C   ./bin/main                       2635MiB |
|    2   N/A  N/A     41728      C   python                           5133MiB |
|    2   N/A  N/A    152095      C   ./bin/main                       3565MiB |
|    6   N/A  N/A      3148      C   python3.6                        3887MiB |
|    6   N/A  N/A     11020      C   python3.6                        3887MiB |
+-----------------------------------------------------------------------------+
```


该项目就是为了解决这个文件，通过以下监控线即可获取到pod本身使用的显存大小。
```
pod_used_gpu_mem_MB
```

获取到的监控项类似如下
```
pod_used_gpu_mem_MB{app="nvidia-gpu-mem-monitor",app_pid="31563",gpu_name="GeForce GTX 1080 Ti",gpu_uuid="GPU-78d64296-8254-ef39-35ec-cb35bd6e6192",instance="10.244.19.248:80",job="nvidia-gpu-mem-monitor",kubernetes_name="nvidia-gpu-mem-monitor",kubernetes_namespace="devops",pod_name="xmcvt-speech-speed-detect-77b7bbdb96-sjjdq",pod_namespace="xmcvt"}
```