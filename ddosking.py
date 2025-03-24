#!/usr/bin/env python3
# encoding: utf-8

from seedemu.layers import Base, Routing, Ebgp, Ibgp, Ospf, PeerRelationship
from seedemu.compiler import Docker, Platform, DockerImage
from seedemu.core import Emulator, Binding, Filter
from seedemu.services import WebService, DomainNameService, DomainNameCachingService
from seedemu.utilities import Makers
import os, sys
import shutil
from pathlib import Path

def copy_files_to_nodes():
    # 定义基础路径
    current_dir = Path(__file__).parent
    code_dir = current_dir / "code"
    container_dir = current_dir / "container"
    output_dir = current_dir / "output"

    # 节点类型到源目录的映射
    node_mappings = {
        "c2-server": {
            "code": ["c2"],
            "container": ["c2"]
        },
        "dns-auth": {
            "code": ["auth"],
            "container": ["auth"]
        },
        "unbound-server": {
            "container": ["unbound"]
        },
        "deepseek-server": {
            "container": ["deepseek"]
        },
        "bot-": {
            "code": ["bot"],
            "container": ["bot"]
        },
        "reflector-": {
            "code": ["reflector"],
            "container": ["reflector"]
        },
        "victim": {
            "container": ["victim"]
        }
    }

    # 遍历output目录中的每个节点目录
    for node_dir in output_dir.iterdir():
        if not node_dir.is_dir() or not node_dir.name.startswith("hnode_"):
            continue

        node_name = node_dir.name
        print(f"Processing {node_name}...")

        # 确定节点类型
        node_type = None
        for key in node_mappings.keys():
            if key in node_name:
                node_type = key
                break

        if not node_type:
            print(f"  Unknown node type: {node_name}")
            continue

        mapping = node_mappings[node_type]
        print(f"  Identified as {node_type}")

        # 复制代码目录
        if "code" in mapping:
            for code_folder in mapping["code"]:
                src = code_dir / code_folder
                dest = node_dir / code_folder
                if dest.exists():
                    shutil.rmtree(dest)
                print(f"  Copying code {src} to {dest}")
                shutil.copytree(src, dest)

        # 复制容器配置文件
        if "container" in mapping:
            for container_folder in mapping["container"]:
                src = container_dir / container_folder
                for item in src.iterdir():
                    if item.is_dir():
                        dest = node_dir / item.name
                        shutil.copytree(item, dest, dirs_exist_ok=True)
                    else:
                        shutil.copy2(item, node_dir)
                print(f"  Copied container files from {src} to {node_dir}")

    print("File copy completed.")

def run(dumpfile=None, hosts_per_as=2):
    ###############################################################################
    # Set the platform information
    script_name = os.path.basename(__file__)

    if len(sys.argv) == 1:
        platform = Platform.AMD64
    elif len(sys.argv) == 2:
        if sys.argv[1].lower() == 'amd':
            platform = Platform.AMD64
        elif sys.argv[1].lower() == 'arm':
            platform = Platform.ARM64
        else:
            print(f"Usage: {script_name} amd|arm")
            sys.exit(1)
    else:
        print(f"Usage: {script_name} amd|arm")
        sys.exit(1)

    # Initialize the emulator and layers
    emu = Emulator()
    base = Base()
    routing = Routing()
    ebgp = Ebgp()
    web = WebService()
    dns = DomainNameService()
    dnsCache = DomainNameCachingService()
    
    ###############################################################################
    # Create internet exchanges (smaller network than mini_internet.py)
    ix100 = base.createInternetExchange(100)
    ix101 = base.createInternetExchange(101)
    ix102 = base.createInternetExchange(102)
    
    # Customize names (for visualization purpose)
    ix100.getPeeringLan().setDisplayName('NYC-100')
    ix101.getPeeringLan().setDisplayName('Chicago-101')
    ix102.getPeeringLan().setDisplayName('Miami-102')
    
    ###############################################################################
    # Create Transit Autonomous Systems
    
    # Tier 1 ASes
    Makers.makeTransitAs(base, 2, [100, 101, 102], [(100, 101), (101, 102), (100, 102)])
    Makers.makeTransitAs(base, 3, [100, 102], [(100, 102)])
    
    ###############################################################################
    # Create stub ASes for various components
    
    # AS-150: Will host the C2 server
    as150 = base.createAutonomousSystem(150)
    as150.createNetwork('net0')
    as150.createRouter('router0').joinNetwork('net0').joinNetwork('ix100')
    
    # Create C2 server in AS-150
    c2_server = as150.createHost('c2-server')
    c2_server.joinNetwork('net0')
    
    # # Add startup script for C2 server
    # c2_server.appendStartCommand('cd /root/c2', False)
    # c2_server.appendStartCommand('go run main.go > /tmp/c2.log 2>&1', True)
    
    # AS-151: Will host the DNS authority server
    as151 = base.createAutonomousSystem(151)
    as151.createNetwork('net0')
    as151.createRouter('router0').joinNetwork('net0').joinNetwork('ix100')
    
    # Create DNS authority server in AS-151
    dns_auth = as151.createHost('dns-auth')
    dns_auth.joinNetwork('net0')
    
    # AS-152: Will host the Unbound server
    as152 = base.createAutonomousSystem(152)
    as152.createNetwork('net0')
    as152.createRouter('router0').joinNetwork('net0').joinNetwork('ix100')
    
    # Create Unbound server in AS-152
    unbound_server = as152.createHost('unbound-server')
    unbound_server.joinNetwork('net0')
    
    # AS-153: Will host the Deepseek server
    as153 = base.createAutonomousSystem(153)
    as153.createNetwork('net0')
    as153.createRouter('router0').joinNetwork('net0').joinNetwork('ix100')
    
    # Create Deepseek server in AS-153
    deepseek_server = as153.createHost('deepseek-server')
    deepseek_server.joinNetwork('net0')
    
    # AS-160-161: Will host 2 bot hosts
    for asn in range(160, 162):  # Changed to create AS160 and AS161 (2 bots)
        current_as = base.createAutonomousSystem(asn)
        current_as.createNetwork('net0')
        current_as.createRouter('router0').joinNetwork('net0').joinNetwork('ix101')
        
        # Create bot host in the current AS
        bot_host = current_as.createHost(f'bot-{asn}')
        bot_host.joinNetwork('net0')
        
    
    # AS-170-174: Will host 5 reflector servers
    for asn in range(170, 175):  # Creates AS170 through AS174 (5 reflectors)
        current_as = base.createAutonomousSystem(asn)
        current_as.createNetwork('net0')
        current_as.createRouter('router0').joinNetwork('net0').joinNetwork('ix102')
        
        # Create reflector host in the current AS
        reflector = current_as.createHost(f'reflector-{asn}')
        reflector.joinNetwork('net0')
        
        # Add startup script for reflector - need to start program and then input 1
        # Create a script that will start reflector and then send "1" after a short delay
    
    # Create a victim AS with a web server
    as180 = base.createAutonomousSystem(180)
    as180.createNetwork('net0')
    as180.createRouter('router0').joinNetwork('net0').joinNetwork('ix102')
    victim = as180.createHost('victim')
    victim.joinNetwork('net0')
    
    # Create a web service on victim
    web.install('web-victim')
    emu.addBinding(Binding('web-victim', filter=Filter(nodeName='victim', asn=180)))
    
    ###############################################################################
    # Set up the custom images for all relevant nodes
    
    docker = Docker(internetMapEnabled=True, platform=platform)
    
    ###############################################################################
    # Peering setup
    
    # RS (route server) peering - default peering mode is PeerRelationship.Peer
    ebgp.addRsPeers(100, [2, 3, 150, 151, 152, 153])
    ebgp.addRsPeers(101, [2, 160, 161])  # Changed to only include the 2 bot ASes
    ebgp.addRsPeers(102, [2, 3, 170, 171, 172, 173, 174, 180])  # Updated to include all 5 reflector ASes
    
    # Set up private peerings
    ebgp.addPrivatePeerings(100, [2, 3], [150, 151, 152, 153], PeerRelationship.Provider)
    ebgp.addPrivatePeerings(101, [2], [160, 161], PeerRelationship.Provider)  # Changed to match bot ASes
    ebgp.addPrivatePeerings(102, [2, 3], [170, 171, 172, 173, 174, 180], PeerRelationship.Provider)  # Updated to match reflector ASes
    
    ###############################################################################
    # Add layers to the emulator
    
    emu.addLayer(base)
    emu.addLayer(routing)
    emu.addLayer(ebgp)
    emu.addLayer(Ibgp())
    emu.addLayer(Ospf())
    emu.addLayer(web)
    emu.addLayer(dns)
    emu.addLayer(dnsCache)
    
    if dumpfile is not None:
        # Save it to a file, so it can be used by other emulators
        emu.dump(dumpfile)
    else:
        emu.render()
        emu.compile(docker, './output', override=True)
        
if __name__ == "__main__":
    run()
    copy_files_to_nodes()
