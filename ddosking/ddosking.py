#!/usr/bin/env python3
# encoding: utf-8

from seedemu.layers import Base, Routing, Ebgp, Ibgp, Ospf, PeerRelationship
from seedemu.compiler import Docker, Platform, DockerImage
from seedemu.core import Emulator, Binding, Filter
from seedemu.services import WebService, DomainNameService, DomainNameCachingService
from seedemu.utilities import Makers
import os, sys

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
    
    # AS-160-164: Will host the bot hosts (5 ASes with 1 bot each)
    for asn in range(160, 165):
        current_as = base.createAutonomousSystem(asn)
        current_as.createNetwork('net0')
        current_as.createRouter('router0').joinNetwork('net0').joinNetwork('ix101')
        
        # Create bot host in the current AS
        bot_host = current_as.createHost(f'bot-{asn}')
        bot_host.joinNetwork('net0')
        
        # Add startup script for bot host
        bot_host.appendStartCommand('cd /root/bot', False)
        bot_host.appendStartCommand('go run main.go > /tmp/bot.log 2>&1', True)
    
    # AS-170, 171: Will host the reflector servers (2 ASes with 1 reflector each)
    for asn in range(170, 172):
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
    
    # Create DockerImage objects for all three required images
    docker.addImage(DockerImage('ddosking', [], local=True))
    docker.addImage(DockerImage('unbound', [], local=True))
    docker.addImage(DockerImage('deepseek', [], local=True))
    
    # Apply the ddosking image to the C2 server
    docker.setImageOverride(c2_server, 'ddosking')
    
    # Apply the ddosking image to the DNS authority server
    docker.setImageOverride(dns_auth, 'ddosking')
    
    # Apply the unbound image to the unbound server
    docker.setImageOverride(unbound_server, 'unbound')
    
    # Apply the deepseek image to the deepseek server
    docker.setImageOverride(deepseek_server, 'deepseek')
    
    # Apply the ddosking image to all bot hosts
    for asn in range(160, 165):
        current_as = base.getAutonomousSystem(asn)
        bot_host = current_as.getHost(f'bot-{asn}')
        docker.setImageOverride(bot_host, 'ddosking')
    
    # Apply the ddosking image to all reflectors
    for asn in range(170, 172):
        current_as = base.getAutonomousSystem(asn)
        reflector = current_as.getHost(f'reflector-{asn}')
        docker.setImageOverride(reflector, 'ddosking')
    
    ###############################################################################
    # Peering setup
    
    # RS (route server) peering - default peering mode is PeerRelationship.Peer
    ebgp.addRsPeers(100, [2, 3, 150, 151, 152, 153])
    ebgp.addRsPeers(101, [2, 160, 161, 162, 163, 164])
    ebgp.addRsPeers(102, [2, 3, 170, 171, 180])
    
    # Set up private peerings
    ebgp.addPrivatePeerings(100, [2, 3], [150, 151, 152, 153], PeerRelationship.Provider)
    ebgp.addPrivatePeerings(101, [2], [160, 161, 162, 163, 164], PeerRelationship.Provider)
    ebgp.addPrivatePeerings(102, [2, 3], [170, 171, 180], PeerRelationship.Provider)
    
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
        
        # Copy all container files to the output directory
        os.system('cp -r container/ddosking ./output')
        os.system('cp -r container/unbound ./output')
        os.system('cp -r container/deepseek ./output')

if __name__ == "__main__":
    run()