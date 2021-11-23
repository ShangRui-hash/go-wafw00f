# GO-WAFW00F

## 介绍

- WAFW00F是一款优秀的web应用防火墙识别开源工具：https://github.com/EnableSecurity/wafw00f
- GO-WAFW00F forked from EmYiQing/go-wafw00f,由于作者不在维护，所以我接盘了
- 使用Golang重写的原因：Python环境配置不便利，Golang打包生成可执行文件直接运行

## 安装

```shell
go install github.com/ShangRui-hash/go-wafw00f@latest
```

## 快速开始

直接在github的release页面下载可执行文件：[下载地址](https://github.com/EmYiQing/go-wafw00f/releases/)

```shell
# 探测waf
go-wafw00f run  -u http://www.xxx.com -r rule.json

# 解析wafw00f的插件库为json文件
go-wafw00f parse --lib ./lib -dst rule.json
```


## 简单原理

- parse 命令使用正则表达式来解析wafw00f的插件库，生成json文件 (todo:这是原作者的思路，但是这个思路问题大大的：对于any,all 没法完全解析，需要重新构思一个解决方案)
- run 命令读取json格式的规则文件，探测waf

## 支持WAF列表
```
  WAF Name                        Manufacturer
  --------                        ------------

  ACE XML Gateway                  Cisco                            
  aeSecure                         aeSecure                         
  AireeCDN                         Airee                            
  Airlock                          Phion/Ergon                      
  Alert Logic                      Alert Logic                      
  AliYunDun                        Alibaba Cloud Computing          
  Anquanbao                        Anquanbao                        
  AnYu                             AnYu Technologies                
  Approach                         Approach                         
  AppWall                          Radware                          
  Armor Defense                    Armor                            
  ArvanCloud                       ArvanCloud                       
  ASP.NET Generic                  Microsoft                        
  ASPA Firewall                    ASPA Engineering Co.             
  Astra                            Czar Securities                  
  AWS Elastic Load Balancer        Amazon                           
  AzionCDN                         AzionCDN                         
  Azure Front Door                 Microsoft                        
  Barikode                         Ethic Ninja                      
  Barracuda                        Barracuda Networks               
  Bekchy                           Faydata Technologies Inc.        
  Beluga CDN                       Beluga                           
  BIG-IP Local Traffic Manager     F5 Networks                      
  BinarySec                        BinarySec                        
  BitNinja                         BitNinja                         
  BlockDoS                         BlockDoS                         
  Bluedon                          Bluedon IST                      
  BulletProof Security Pro         AITpro Security                  
  CacheWall                        Varnish                          
  CacheFly CDN                     CacheFly                         
  Comodo cWatch                    Comodo CyberSecurity             
  CdnNS Application Gateway        CdnNs/WdidcNet                   
  ChinaCache Load Balancer         ChinaCache                       
  Chuang Yu Shield                 Yunaq                            
  Cloudbric                        Penta Security                   
  Cloudflare                       Cloudflare Inc.                  
  Cloudfloor                       Cloudfloor DNS                   
  Cloudfront                       Amazon                           
  CrawlProtect                     Jean-Denis Brun                  
  DataPower                        IBM                              
  DenyALL                          Rohde & Schwarz CyberSecurity    
  Distil                           Distil Networks                  
  DOSarrest                        DOSarrest Internet Security      
  DotDefender                      Applicure Technologies           
  DynamicWeb Injection Check       DynamicWeb                       
  Edgecast                         Verizon Digital Media            
  Eisoo Cloud Firewall             Eisoo                            
  Expression Engine                EllisLab                         
  BIG-IP AppSec Manager            F5 Networks                      
  BIG-IP AP Manager                F5 Networks                      
  Fastly                           Fastly CDN                       
  FirePass                         F5 Networks                      
  FortiWeb                         Fortinet                         
  GoDaddy Website Protection       GoDaddy                          
  Greywizard                       Grey Wizard                      
  Huawei Cloud Firewall            Huawei                           
  HyperGuard                       Art of Defense                   
  Imunify360                       CloudLinux                       
  Incapsula                        Imperva Inc.                     
  IndusGuard                       Indusface                        
  Instart DX                       Instart Logic                    
  ISA Server                       Microsoft                        
  Janusec Application Gateway      Janusec                          
  Jiasule                          Jiasule                          
  Kona SiteDefender                Akamai                           
  KS-WAF                           KnownSec                         
  KeyCDN                           KeyCDN                           
  LimeLight CDN                    LimeLight                        
  LiteSpeed                        LiteSpeed Technologies           
  Open-Resty Lua Nginx             FLOSS                            
  Oracle Cloud                     Oracle                           
  Malcare                          Inactiv                          
  MaxCDN                           MaxCDN                           
  Mission Control Shield           Mission Control                  
  ModSecurity                      SpiderLabs                       
  NAXSI                            NBS Systems                      
  Nemesida                         PentestIt                        
  NevisProxy                       AdNovum                          
  NetContinuum                     Barracuda Networks               
  NetScaler AppFirewall            Citrix Systems                   
  Newdefend                        NewDefend                        
  NexusGuard Firewall              NexusGuard                       
  NinjaFirewall                    NinTechNet                       
  NullDDoS Protection              NullDDoS                         
  NSFocus                          NSFocus Global Inc.              
  OnMessage Shield                 BlackBaud                        
  Palo Alto Next Gen Firewall      Palo Alto Networks               
  PerimeterX                       PerimeterX                       
  PentaWAF                         Global Network Services          
  pkSecurity IDS                   pkSec                            
  PT Application Firewall          Positive Technologies            
  PowerCDN                         PowerCDN                         
  Profense                         ArmorLogic                       
  Puhui                            Puhui                            
  Qcloud                           Tencent Cloud                    
  Qiniu                            Qiniu CDN                        
  Reblaze                          Reblaze                          
  RSFirewall                       RSJoomla!                        
  RequestValidationMode            Microsoft                        
  Sabre Firewall                   Sabre                            
  Safe3 Web Firewall               Safe3                            
  Safedog                          SafeDog                          
  Safeline                         Chaitin Tech.                    
  SecKing                          SecKing                          
  eEye SecureIIS                   BeyondTrust                      
  SecuPress WP Security            SecuPress                        
  SecureSphere                     Imperva Inc.                     
  Secure Entry                     United Security Providers        
  SEnginx                          Neusoft                          
  ServerDefender VP                Port80 Software                  
  Shield Security                  One Dollar Plugin                
  Shadow Daemon                    Zecure                           
  SiteGround                       SiteGround                       
  SiteGuard                        Sakura Inc.                      
  Sitelock                         TrueShield                       
  SonicWall                        Dell                             
  UTM Web Protection               Sophos                           
  Squarespace                      Squarespace                      
  SquidProxy IDS                   SquidProxy                       
  StackPath                        StackPath                        
  Sucuri CloudProxy                Sucuri Inc.                      
  Tencent Cloud Firewall           Tencent Technologies             
  Teros                            Citrix Systems                   
  Trafficshield                    F5 Networks                      
  TransIP Web Firewall             TransIP                          
  URLMaster SecurityCheck          iFinity/DotNetNuke               
  URLScan                          Microsoft                        
  UEWaf                            UCloud                           
  Varnish                          OWASP                            
  Viettel                          Cloudrity                        
  VirusDie                         VirusDie LLC                     
  Wallarm                          Wallarm Inc.                     
  WatchGuard                       WatchGuard Technologies          
  WebARX                           WebARX Security Solutions        
  WebKnight                        AQTRONIX                         
  WebLand                          WebLand                          
  RayWAF                           WebRay Solutions                 
  WebSEAL                          IBM                              
  WebTotem                         WebTotem                         
  West263 CDN                      West263CDN                       
  Wordfence                        Defiant                          
  WP Cerber Security               Cerber Tech                      
  WTS-WAF                          WTS                              
  360WangZhanBao                   360 Technologies                 
  XLabs Security WAF               XLabs                            
  Xuanwudun                        Xuanwudun                        
  Yundun                           Yundun                           
  Yunsuo                           Yunsuo                           
  Yunjiasu                         Baidu Cloud Computing            
  YXLink                           YxLink Technologies              
  Zenedge                          Zenedge                          
  ZScaler                          Accenture
```
