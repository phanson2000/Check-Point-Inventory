Easily inventory all firewalls connectioned to a Check Point Management Server. This tool utilizes Check Point WEBAPI and should work for R80.x and R81.X. This will go through all firewalls connected through all domains on the management.  This will export the details intoa summary file per gateway as well as CSV files for reporting.  This can be run on a Windows client, Linux such as Ubuntu/Redhat/Centos or Gaia management directly.

You can download the source and compile with go.  If compiling to run on the mangement platform be sure CGO_ENABLED is set to 0.  i.e. export CGO_ENABLED=0

For compiling you will alsot need Check Point GO at https://github.com/CheckPointSW/cp-mgmt-api-go-sdk.  Biniaries are uploaded for both windows and linux including Gaia.

To run either execute CheckPointInventory and fill in the required fields or utilize the command line arguments.

Comamand Line Arguments Example:

![linuxwithargs](https://user-images.githubusercontent.com/2261078/193640912-63edd5e4-5961-4f3e-aa42-89a966993389.png)

No Arguments Example:

![WindowsNoArgs](https://user-images.githubusercontent.com/2261078/193640914-974862b4-1438-4999-b114-8e1531e4e39f.PNG)



This is an early verison.  If you find a problem, please do not hesitate to reach out to me here, or at Check Point at phanson@checkpoint.com.  Some additional features requested and coming include:

  Option to use APIKEY instead for password
  argument for option detail firewall output
  Remove timeout and use wait for completion to 100% and success TRUE when getting task output
  Use golang terminal library for password collection instead of gopass
  Store configuration as well as password or api key in YAML file
  Generic report and summary creation as well as CSV file.

