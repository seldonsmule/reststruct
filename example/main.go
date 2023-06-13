package main

import (
	"os"
	"fmt"
//        "github.com/seldonsmule/restapi"
        "github.com/seldonsmule/reststruct"
//	"strings"
//	"io/ioutil"
        //"time"
//        "bufio"
        //"syscall"
 //       "strconv"
  //      "strings"
//	"net/http"
//	"io/ioutil"
 //       "encoding/json"
//        "database/sql"
  //      "time"
//        _ "github.com/mattn/go-sqlite3" 
        "github.com/seldonsmule/logmsg"
//        "golang.org/x/crypto/ssh/terminal"
)

type MyExample struct {

  oEndpoints *reststruct.RestStruct

}

func (pMyExample *MyExample) SetEndpoints(oEndpoints *reststruct.RestStruct) {
  pMyExample.oEndpoints = oEndpoints
}

func (pMyExample *MyExample) SetObject(obj interface{}){

  pMyExample.oEndpoints.SetObject(obj)


}

func (pMyExample *MyExample) GetStruct(name string, bDebug bool) bool{

  r := pMyExample.oEndpoints.NewGetRestapi(name) 

  if(r == nil){
    fmt.Println("NewGetRestapi failed!!!!")
    return false
  }

  // if you need to set token use
  // r.SetBearerAccessToken()
  //
  // if you need to set a certificate file use
  // r.SetCertificateFile()
  //

  if(!pMyExample.oEndpoints.Send(r, false)){
    fmt.Println("Simpel send faile")
    r.Dump()
    return false
  }


  return true
}


func usestruct(pEndpoints *reststruct.RestStruct, bDebug bool) bool{

  fmt.Println(" --- usestruct --- ")

  var oMyExample MyExample

  oMyExample.SetEndpoints(pEndpoints)

  worked, auth := oMyExample.GetAuthentication()

  if(!worked){
    fmt.Println("GetAuthentication failed")
    return false
  }

  fmt.Println("auth:", auth)

  fmt.Println("auth:", auth.AccessToken)

  worked, vehicles := oMyExample.GetVehicles()

  if(!worked){
    fmt.Println("GetVehicles failed")
    return false
  }

  fmt.Println("vehicles:", vehicles)

  fmt.Println("vehicles:", vehicles.Response[0].DisplayName)


  return true
}


func listsave(pEndpoints *reststruct.RestStruct, bDebug bool) bool{

  fmt.Println(" --- listsave --- ")

  // first get our indexes for running through the list

  keys := pEndpoints.GetIndex()

  // everyone will use the same packge info
  pEndpoints.PackageInfo.SetSave(true)
  pEndpoints.PackageInfo.SetPackageName("example")
  pEndpoints.PackageInfo.SetClassName("MyExample")

  // now run though the list

  for _, k := range keys {
    fmt.Println("key:", k)

    r := pEndpoints.NewGetRestapi(k) 

    if(r == nil){
      fmt.Println("NewGetRestapi failed!!!!")
      return false
    }

    // if you need to set token use
    // r.SetBearerAccessToken()
    //
    // if you need to set a certificate file use
    // r.SetCertificateFile()
    //

    if(!pEndpoints.Send(r, true)){
      fmt.Println("Simpel send faile")
      r.Dump()
      return false
    }
  }

  return true
}

func rawsave(bDebug bool) bool{

  fmt.Println(" --- rawsave --- ")
  oEndpoints := reststruct.NewRestStruct()

  oEndpoints.SetDebug(true) // turn on debug
  oEndpoints.SetRawURL("http://localhost:3000/api/1/vehicles") // set the url prefix
  //oEndpoints.SetRawURL("http://localhost:3000/api/1/authentication") // set the url prefix
  oEndpoints.SetDirName("SampleDir") // set the save directory

  r := oEndpoints.NewGetRestapi("authentication") 

  //r.Dump()

  // if you need to set token use
  // r.SetBearerAccessToken()
  //
  // if you need to set a certificate file use
  // r.SetCertificateFile()
  //


  if(!oEndpoints.Send(r, true)){
    fmt.Println("Simpel send faile")
    r.Dump()
    return false
  }


  return true
}

func simplesave(pEndpoints *reststruct.RestStruct, bDebug bool) bool{

  r := pEndpoints.NewGetRestapi("authentication") 

  if(r == nil){
    fmt.Println("NewGetRestapi failed!!!!")
    return false
  }

  // if you need to set token use
  // r.SetBearerAccessToken()
  //
  // if you need to set a certificate file use
  // r.SetCertificateFile()
  //

  pEndpoints.PackageInfo.SetSave(true)
  pEndpoints.PackageInfo.SetPackageName("example")
  pEndpoints.PackageInfo.SetClassName("MyExample")



  if(!pEndpoints.Send(r, true)){
    fmt.Println("Simpel send faile")
    r.Dump()
    return false
  }


  return true
}

func help(){

  fmt.Println("usage example test_name [debug]")
  fmt.Println()
  fmt.Println("simplesave - how to save off a struct")
  fmt.Println("rawssave - how to save off a raw struct")
  fmt.Println("listsave - how to save off a list of structs")
  fmt.Println("usestruct - Example of using saved off structs")

}

func main() {

  bDebug := false

  oEndpoints := reststruct.NewRestStruct()

  oEndpoints.SetDebug(true) // turn on debug
  oEndpoints.SetURLPrefix("http://localhost:3000/api/1") // set the url prefix
  oEndpoints.SetDirName("SampleDir") // set the save directory

  oEndpoints.AddEndpoint("authentication", "Authentication", "authentication", false)
  oEndpoints.AddEndpoint("vehicles", "Vehicles", "vehicles", false)
  oEndpoints.AddEndpoint("vehicles_multi", "Vehicles_multi", "vehicles_multi", false)
  oEndpoints.AddEndpoint("drive_state", "Drive_state", "drive_state", false)


//  oEndpoints.Dump() // dump the endpoints


  fmt.Println("starting example for restapi")
  logmsg.SetLogFile("example.log");

  args := os.Args

  if(len(args) < 2){
    help()
    os.Exit(1)
  }

  if(len(args) == 3){
  
    switch args[2] {
      case "debug":
        bDebug = true
        fmt.Println("Debug on")

      default:
        help()
        os.Exit(2)
    }
  }

  switch args[1] {

    case "simplesave":
      simplesave(oEndpoints, bDebug)

    case "rawsave":
      rawsave(bDebug)

     
    case "listsave":
      listsave(oEndpoints, bDebug)
      
    case "usestruct":
      usestruct(oEndpoints, bDebug)

    default:
      help()

  }

}
