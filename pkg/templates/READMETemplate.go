package templates

const READMETemplate = `
# Egg Framework

	           ████████████████        
	         ██                ██      
	     ████    ░░░░░░░░        ██    
	   ██      ░░      ░░░░        ██  
	 ██      ░░          ░░░░        ██
	 ██      ░░          ░░░░        ██
	██        ░░▒▒░░  ░░░░░░░░        ██
	██░░        ░░░░░░░░░░░░        ░░██
	  ██░░        ░░░░░░░░        ░░██  
	  ██░░░░                    ░░██    
	    ████░░░░            ░░░░██      
 	       ████░░░░░░░░░░░░████        
	           ████████████            

The Egg framework is based on getting things done. In fact it is a framework made to be perfect for the solo developer. You are given all you need to get started extremely quickly. In this generated repository are also ci/cd to have this automatically test, build, and deploy your website to a vps using coolify with docker. 

If you do not want to use docker or coolify, great! That is not the point of this framework. You just want to get up and running without having to rebuild the same starter over and over again. This is that. Ok Cheers! 

## Getting started 

make sure that you have a postgres database running the connection string in the ` + "`config/development.yaml`" +` is correct. and make sure that the configured s3 storage configuration is correct as well. 

Run the command ` + "`air`" + ` to start the server,
or Run the command ` + "`go run main.go`" + ` to start the server

You can change the default port by going into the config/ directory and changing the ` + "`server.port`" + ` value in the development.yaml to what you want. 
`

