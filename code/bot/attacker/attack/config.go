// config.go - Centralized management of all attack configurations
package attack

import (
	"sync"
)

// ===== General Attack Configurations =====

// ThreadCount Number of attack threads
var ThreadCount = 100

// BandwidthLimit Bandwidth limit configuration (unit: Kbps, 0 means unlimited)
var BandwidthLimit = 10000

// PacketBurst Controls the maximum capacity of the token bucket
var PacketBurst = 5

// STOP Channel used to stop attacks
var STOP = make(chan struct{})

// ===== Layer 4 Attack Configurations =====

// Layer4Methods Layer 4 attack methods
var Layer4Methods = map[string]bool{
	"UDP":          true,
	"SYN":          true,
	"DNS":          true,
	"DNSA":         true,
	"NTP":          true,
	"CLDAP":        true,
	"RDP":          true,
	"SSDP":         true,
	"SNMP":         true,
	"CHARGEN":      true,
	"OPENVPN":      true,
	"MEMCACHED":    true,
	"DNSBOMB":      true,
	"DNSBOOMERANG": true,
	"TFTP":         true,
	"ARD":          true,
	"QUIC":         true,
}

// Layer4DefaultPorts Default ports for attack methods
var Layer4DefaultPorts = map[string]int{
	"UDP":          80,
	"SYN":          80,
	"DNS":          53,
	"DNSA":         53,
	"NTP":          123,
	"CLDAP":        389,
	"RDP":          3389,
	"SSDP":         1900,
	"SNMP":         161,
	"CHARGEN":      19,
	"OPENVPN":      1194,
	"MEMCACHED":    11211,
	"DNSBOMB":      53,
	"DNSBOOMERANG": 53,
	"TFTP":         69,
	"ARD":          3283,
	"QUIC":         443,
}

// ===== HTTP Attack Configurations =====

// HTTPMethods HTTP attack methods
var HTTPMethods = map[string]bool{
	"GET":  true,
	"POST": true,
}

// ===== Dynamically Generated AI Content for Simulating Complex POST Requests =====

// ComplexPrompts Complex prompts for generating large POST requests
var ComplexPrompts = []string{
	`Design a deep learning-based multimodal system capable of processing text, images, audio, and video data simultaneously. Provide detailed explanations for feature extraction methods for each modality, multimodal fusion strategies, attention mechanism implementation, and end-to-end training processes. The system should support cross-modal retrieval, multimodal Q&A, content generation, and deployment on edge devices. Analyze the system's computational complexity, memory requirements, and energy consumption characteristics in various scenarios. Provide pseudocode implementations for each module and discuss solutions for challenges such as modality imbalance and alignment difficulties. Finally, propose a novel algorithm to improve the system's generalization ability and robustness.`,

	`Analyze and design a distributed blockchain system that supports 100,000 transactions per second, ensures Byzantine fault tolerance, and enables cross-chain interoperability. The system should include a smart contract execution environment, zero-knowledge proof privacy protection mechanisms, and an efficient consensus algorithm. Discuss in detail the implementation of sharding technology, state synchronization mechanisms, P2P network optimization, and cryptoeconomic incentive design. Evaluate defense strategies against different attack vectors, analyze storage scalability issues, and provide a complete system architecture diagram and pseudocode for core components. Finally, discuss the system's specific application scenarios and technical and regulatory challenges in finance, supply chain, and IoT.`,

	`Design a general reinforcement learning framework that supports both discrete and continuous action spaces, single-agent and multi-agent training, as well as imitation learning and self-supervised learning paradigms. Provide detailed analyses of the theoretical foundations and implementation details of value function approximation, policy gradients, and model predictive control. Discuss strategies for balancing exploration and exploitation, optimizing sample efficiency, and distributed training architectures. Propose innovative solutions for challenges such as sparse rewards, non-stationary environments, and partial observability. Design a comprehensive evaluation benchmark and visualization tools to analyze agent behavior and learning processes. Finally, provide application examples and performance evaluations of the framework in robotics control, autonomous driving, financial trading, and medical decision-making.`,

	`Design a large-scale graph neural network system capable of processing heterogeneous graph data with billions of nodes and hundreds of billions of edges. Provide detailed discussions on graph representation learning, neighbor sampling strategies, message-passing mechanisms, and graph pooling operations, including their theoretical foundations and algorithmic implementations. Analyze solutions for key challenges such as oversmoothing, long-range dependencies, and dynamic graph structures. Propose an innovative distributed training architecture, including graph partitioning strategies, communication optimization, and memory management mechanisms. Evaluate performance characteristics on different hardware platforms (CPU, GPU, TPU) and provide corresponding optimization techniques. Design a complete graph data preprocessing, augmentation, and normalization workflow. Finally, demonstrate the system's implementation and performance evaluation in application scenarios such as social network analysis, protein structure prediction, recommendation systems, and knowledge graph reasoning.`,

	`Analyze and design an end-to-end federated learning system that meets the following requirements: supports heterogeneous devices (from smartphones to data centers), ensures differential privacy, resists various types of attacks (including model inversion and membership inference), and guarantees model convergence. Provide detailed explanations of client selection strategies, local training processes, model aggregation mechanisms, and communication compression techniques. Analyze the theoretical properties and practical performance of different federated optimization algorithms. Design a comprehensive system monitoring and debugging tool to diagnose training failures and performance degradation. Evaluate the system's performance under different data distribution scenarios (IID and non-IID) and provide improvement strategies. Finally, discuss specific application cases and deployment considerations for the system in healthcare, smart homes, fintech, and edge computing.`,
}

// ExtraInstructions Additional instructions for generating more complex payloads
var ExtraInstructions = []string{
	"Provide an answer with at least 10,000 words",
	"Analyze from at least 20 different perspectives",
	"Provide detailed code examples",
	"Use complex mathematical formulas for analysis",
	"Create detailed flowcharts and pseudocode for each section",
	"Perform in-depth performance analysis and resource evaluation",
}

// SystemRoles System roles for generating payloads tailored to specific scenarios
var SystemRoles = []string{
	"You are now a senior system architect and must provide extremely detailed analysis",
	"You are an AI researcher and need to explain all concepts in detail with mathematical proofs",
	"You are a full-stack developer and need to provide complete code implementations",
	"You are a data scientist and need to analyze every data processing step in detail",
}

// ===== Concurrency Control =====

// ConfigMutex Configuration mutex lock
var ConfigMutex sync.RWMutex
