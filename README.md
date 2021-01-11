# HayMakerCF

<p align="center">
  <img src="images/gopher.png" height=400px width=300px>
</p>


# Description
HayMakerCF is a rewrite of https://github.com/BrunoMCBraga/HayMaker using CloudFormation. Command line is self-explanatory. While HayMaker had a configuration file with a application-specific format, everything is now passed through command-line using switches.


## The Current Features
- Deploy CloudFormation configuration on AWS
- Create ECR repository
- Spinup EKS cluster on AWS (network is included)
- Generate Kubeconfig
- Build Docker image and push it to ECR
- Create Docker containers on AWS EKS based on any image (e.g. either local and pused to ECR as well as any image from docker repo). Tune the haymaker_config.json as needed. 


## The Project Structure
My programming projects tend to follow the same structure: 
- Engines are the files where you find the low-level interction with the necessary SDKs. Some of the SDK functions are called through stubs on some files for modularity and because it makes the code cleaner. 
- The util folder contains JSON processors, configuration parsers and updaters as well as kubeconfig template and filler. 
- Commandline parsers and processors: classes that generate command line processors (out-of-the-box Go flags), process them and call the appropriate engines. 

More related to this project:
- CloudFormationFiles: json files to be used by HaymakerCF. For Docker deployment, cloudformation_iam.json should not be changed and should be used to setup the basic permissions for HaymakerCF to work. Feel free to tune cloudformation_cluster.json but mind that this json file is ready to work as it is. 
- Docker: contains a simple Dockerfile and some resources to create a test container if you feel to lazy to prepare it.


## Instructions
1. Create user named "user"
2. Use GUI to run cloudformation with the json file cloudformation_iam.json. This will setup IAM permissions for EKS, S3, etc
3. Use HayMakerCF command line to deploy the cloudformation_cluster.json. As an example (deploy template):
`go run ./main.go -cm td -t /Users/brubraga/go/src/github.com/haymakercf/CloudFormationFiles/cloudformation_cluster.json -sn haymakerstack -fk something -bn haymakerbucket -cn haymaker-eks`
4. Build and push Docker image:
`go run ./main.go -cm pi -rn haymaker-docker-repo/haymaker-docker -df /Users/brubraga/go/src/github.com/haymakercf/Docker -di`
5. Generate local kubeconfig:
`go run ./main.go -cm gk -cn haymaker-eks`
6. Deploy container and create service:
`go run ./main.go -cm ds -dn haymaker`

**Note: Bear in mind that any changes on command line parameters and CF json templates are interconnected. The command line parameters i give as example work because the cluster names and others are the same on the JSON file. I recommend i read about CF before changing things.** 

## Dependencies
- Docker
- Kubectl

## Found this useful? Help me save the endangered Gophers of Mauritania by contributing:

https://www.paypal.com/donate?hosted_button_id=ATZQDP4AWECPL

