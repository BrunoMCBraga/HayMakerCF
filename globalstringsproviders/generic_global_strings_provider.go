package globalstringsproviders

import "fmt"

var haymakerAsciiArt string = `                                                                                                                                                                                                                                                                                                   
_    _             __  __       _              _____ ______ 
| |  | |           |  \/  |     | |            / ____|  ____|
| |__| | __ _ _   _| \  / | __ _| | _____ _ __| |    | |__   
|  __  |/ _  | | | | |\/| |/ _  | |/ / _ \ '__| |    |  __|  
| |  | | (_| | |_| | |  | | (_| |   <  __/ |  | |____| |     
|_|  |_|\__,_|\__, |_|  |_|\__,_|_|\_\___|_|   \_____|_|     
               __/ |                                         
              |___/                                          
                   
`

var dropMicGopher string = `                           
                                                              MDND                                          
                                                             N+:::D:                                        
                                                             N+::::N                                        
                                                             N+I?::N                                        
                                                              N+N::N                                        
                                                  8NNNN$~::::::~NNNZ                                        
                                             INDI::::::::::::::::::::NN,                                    
                                          =DD+:::::::::::::::::::::::::8N                                   
                               ONN=     NN++:::::::::::::::::::::::::::::DM                                 
                             N$:::,?N 8N++::::::::::::::::::::::::::::::78:N                                
                            ,D::::N~:N+++:::::::::::::::::::::::::::::::::=:N                               
                             N+++++$N+++:::::::::::::::::::::::::::::::::::::N                              
                               ODNNN+++:::::::::::::::::::::::::NNNNNNDN::::::N                             
                                  N+++::::::::::::,,,,,,:::::NN~~~~~~~~~=N~::::N                            
                                 87++::::::::::,,,,,,,,,::::D?~~~~~~~~~~~~8N:::N                            
                                 N+++::::::::,,,,,,,,,,,::~N?~~~~~~~~~~~~~~+N~::N                           
                                ~D+++=::::::::,,,,,,,,,:::D?+~~~~~~~~~~~~=~~DI::D:                          
                                N=+NNNO~ONND:::::,,,,,,::=N?~~~~~~~~~~=?N~~D?Z:::N                          
                                NNI+~~~~~:,,$N+:::,,,,,::+N++~~~~~~~~~~~NN  ,D:::N                          
                               +N+?~~~~~~~~~,,N:::,,,,,,:?N+???++++?NNNNNN  N~:::?O                         
                               N??=~~~~~~~~~~~,NN:,,,::::=+NIDND8           N:::::N                         
                              D+??~~~~~~~~~~~~~+N:::+NNND++8D~~            D::::::N                         
                              N+??~~~~~~~+ON~:N7N::+N?+:,,N+IN~:         +N:::::::O~                        
                             ,N+??~~~~~~~~~NN   N:~N?~?NNNN:~+$N8::    NN::::::::::N                        
                             ,N8+???=INNNNNND   N:~NNNND N:::::+++IO$::::::::::::::N?8DNDO                  
                               N~~~~~~    ,    N::::~+NNN~:::::::::::::::::::::NNZ,,::~~~~:$N               
                               +N~~~          N:::::::::::::::::::::::::::::::::~~~~~~~~~~~~N               
                                 NI~        8D:::::::::::::::::::::::::::::::::::~~~~?+N+$NN:               
                                  DNNDDZNDNZ:::::::::::::::::::::::::::::::::::::~~+?NN7       N            
                                  N++++++::::::::::::::::::::::::::::::::::::::::?N8IZ       D              
                                  ,N++++~::::::::::::::::::::::::::::::::::::::::::++N     D: N ,   D       
                                   N++++::::::::::::::::::::::::::::::::::::::::::::+N  :NNZ~~MNN:  ,   :   
                                   N++++:::::::::::::::::::::::::::::::::::::::::::::D   ,NZOZ$7?NZNNN7NN~  
                                 DNN++++:::::::::::::::::::::::::::::::::::::::::::::D,   :DNOOOOOON?~8D?N, 
                                N7+ZZ+++:::::::::::::::::::::::::::::::::::::::::::::=8       +NNN8O7=D++7Z 
                               N++++N+++::::::::::::::::::::::::::::::::::::::::::::::N     ,N N   NZN7+ID  
                              N+++++N+++::::::::::::::::::::::::::::::::::::::::::::::N            ,NNNDN   
                             NI++++?D++++::::::::::::::::::::::::::::::::::::::::::D::N                     
                             N?+++N~N++++:::::::::::::::::::::::::::::::::::::::::::::N                     
                             N?+NN  N++++::::::::::::::::::::::::::::::::::::::::::::Z=                     
                                    $7+++=::::::::::::::::::::::::::::::::::::::::D::N                      
                                     N+N++:::::::::::::::::::::::::::::::::::::::::N+8                      
                                     N+++8:::::::::::::::::::::::::::::::::::::::8::N                       
                                     ~D+D++::::::::::::::::::::::::::::::::::::::::N                        
                                      N+++N+::::::::::::::::::::::::::::::::::::::N                         
                                       D++++=::::::::::::::::::::::::::::::::::::N                          
                                        N+++++:::::::::::::::::::::::::::::::::ONN                          
                                         NN++++::::::::::::::::::::::::::::::$N~~D7                         
                                        N?+DN+++=::::::::::::::::::::::::::NN++~~~N                         
                                        N++++ON8+++::::::::::::::::::::NDN =N++:~~~N                        
                                       N?+++++N  DNND++:::::::~DNNNDD        N++~~~N                        
                                       N+++++N                                N++~~N                        
                                       N+++ON                                  DN+ON                        
                                       NNZN=                                                                
                                                                                
                                                                                
							
    
	`

var optionsMenu string = `	

Commands List: 
HaymakerCF Cluster:
td: HaymakerCF CloudFormation Template Deployment (e.g. go run ./main.go -cm td -t /Users/brubraga/go/src/github.com/haymakercf/CloudFormationFiles/cloudformation_cluster.json -sn haymakerstack -fk something -bn haymakerbucket -cn haymaker-eks)
tt: HaymakerCF CloudFormation Teardown (e.g. go run ./main.go -cm tt -sn haymakerstack -bn haymakercfbucket -rn haymaker-docker-repo/haymaker-docker)

Docker Build And Push:
pi: HayMakerCF Docker Image Build And Push To ECR (e.g. go run ./main.go -cm pi -rn haymaker-docker-repo/haymaker-docker -df /Users/brubraga/go/src/github.com/haymakercf/Docker -di)

Kubernetes Orchestration:
gk: HayMakerCF Generate Kubeconfig File (e.g. go run ./main.go -cm gk -cn haymaker-eks)
sc: HayMakerCF Deploy Container And Create Service (e.g. go run ./main.go -cm sc -kp 80 -dn haymaker -in 965440066241.dkr.ecr.us-east-1.amazonaws.com/haymaker-docker-repo/haymaker-docker:latest -pr TCP -kr 2)
sd: HayMaker Delete Service And Associated Deployment (e.g. go run ./main.go -cm ds -dn haymaker)
`

func GetMenuPictureString() string {

	return fmt.Sprintf("%s\n%s\n", haymakerAsciiArt, dropMicGopher)

}

func GetMenuPictureStringWithOptions() string {

	return fmt.Sprintf("%s\n%s\n%s", haymakerAsciiArt, dropMicGopher, optionsMenu)

}

func GetOptionsMenu() string {
	return optionsMenu
}
