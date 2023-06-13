package reststruct

import (
  "fmt"
  "sort"
  "os"
  "io/ioutil"
  "strings"
  "reflect"
  "encoding/json"
  "github.com/seldonsmule/logmsg"
  "github.com/seldonsmule/restapi"
)



type Endpoints struct{

  Endpoint string    // the final piece of the URL
  Structname string  // Name to save in a created go file for the struct
  Filename string    // Name of the file to save it in

  bHidden bool // if true, will be visiable in list and saving actions

}

type PackageInfo struct{

  sPackageName string
  sClassName string
  bSave bool

}

func (pi *PackageInfo) GetPackageName() string{
	return pi.sPackageName
}

func (pi *PackageInfo) GetClassName() string{
	return pi.sClassName
}

func (pi *PackageInfo) GetSave() bool{
	return pi.bSave
}

func (pi *PackageInfo) SetPackageName(sPackageName string){
	pi.sPackageName = sPackageName
}

func (pi *PackageInfo) SetClassName(sClassName string){
	pi.sClassName = sClassName
}

func (pi *PackageInfo) SetSave(bSave bool){
	pi.bSave = bSave
}





type RestStruct struct{

  oObject interface{} // the object to be saved

  bDebug bool // if true, will print debug messages

  bAddHelperFunc bool // if true, the system will not put the helper functions in the saved .go file


  sURL string // the prefix to add to the URL

  mEndpoints map[string]Endpoints

  sDirName string // the name of the directory to save the file in
  bRawURL bool // true if not using the array of endpoints

  bStdout bool // if true, will print the json to stdout

  PackageInfo PackageInfo 


}

func NewRestStruct() *RestStruct{

  rs := new(RestStruct)

  rs.mEndpoints = make(map[string]Endpoints)

  rs.oObject = nil // default is not set

  rs.bStdout = false // default is not set

  rs.PackageInfo = PackageInfo{"", "", false}

  return rs

}

func (rs *RestStruct) SetDirName(sDirName string){
	rs.sDirName = sDirName
}

func (rs *RestStruct) SetURLPrefix(sURL string){
  rs.SetURL(sURL)
}

func (rs *RestStruct) SetRawURL(sURL string){
  rs.SetURL(sURL)
  rs.bRawURL = true
}

func (rs *RestStruct) SetURL(sURL string){
  rs.sURL = sURL
}


func (rs *RestStruct) SetStdout(bStdout bool){
	rs.bStdout = bStdout
}

func (rs *RestStruct) SetObject(oObject interface{}){

  rs.oObject = oObject
}


func (rs *RestStruct) Dump(){

  fmt.Println("Dumping RestStruct")
  fmt.Println("  DirName:", rs.sDirName)
  fmt.Println("  URL:", rs.sURL)
  fmt.Println("  Debug:", rs.bDebug)

  if(rs.PackageInfo.GetSave()){
    fmt.Println("  PackageInfo:")
    fmt.Println("    PackageName:", rs.PackageInfo.GetPackageName())
    fmt.Println("    ClassName:", rs.PackageInfo.GetClassName())
    fmt.Println("    Save:", rs.PackageInfo.GetSave())
  }else{
    fmt.Println("  PackageInfo: Not Set")
  }


  fmt.Println("   ---  Endpoints ---")

  keys := make([]string, 0, len(rs.mEndpoints))

  for k := range rs.mEndpoints {

    keys = append(keys, k)

  } 

  sort.Strings(keys)

  fmt.Println("Keys:", keys, "Type:", reflect.TypeOf(keys))

  for _, k := range keys {

    if(rs.mEndpoints[k].bHidden){
      continue
    }

    fmt.Printf("   Name:[%s] Endpoint[%s] StructName[%s] GoFileName[%s]\n", k, 
                           rs.mEndpoints[k].Endpoint,
                           rs.mEndpoints[k].Structname,
                           rs.mEndpoints[k].Filename)

  } // end for
}

func (rs *RestStruct) GetUrl(sEndpointName string) string{

	if(rs.bRawURL){
	  return rs.sURL
	}

	return rs.sURL + "/" + rs.mEndpoints[sEndpointName].Endpoint

}

func (rs *RestStruct) SetDebug(bDebug bool){
	rs.bDebug = bDebug
}

func (rs *RestStruct) AddEndpoint(sEndpointName string, sStructName string, 
                                  sFilename string, bHidden bool){

  rs.mEndpoints[sEndpointName] = Endpoints{sEndpointName, sStructName, sFilename, bHidden}


}

func (rs *RestStruct) NewGetRestapi(sEndpointName string) *restapi.Restapi{

  var sStructname string
  var sUrl string
  var sFilename string

  if(rs.bRawURL){
    sUrl = rs.GetUrl("")
    logmsg.Print(logmsg.Info, "RawURL using: ", sUrl)

   
    sStructname = "RawStruct"
    sFilename = "RawFile"
  
  }else{
    logmsg.Print(logmsg.Info, "Using Map with Index: ", sEndpointName)

    if(!rs.CheckEndpoint(sEndpointName)){
      logmsg.Print(logmsg.Error, "Endpoint not found: ", sEndpointName)
      return nil
    }
 
    sStructname = rs.mEndpoints[sEndpointName].Structname
    sFilename = rs.mEndpoints[sEndpointName].Filename
    sUrl = rs.GetUrl(rs.mEndpoints[sEndpointName].Endpoint)


  }

  logmsg.Print(logmsg.Info, "sStructName: ", sStructname)
  logmsg.Print(logmsg.Info, "sFilename: ", sFilename)
  logmsg.Print(logmsg.Info, "sUrl: ", sUrl)

  r := restapi.NewGet(sEndpointName, sUrl)

  //r.DebugOn()

  return r
}

func (rs *RestStruct) SaveResponseBody(r *restapi.Restapi) bool{

  // fmt.Println("SaveResponseBody - " + "makedir: " + rs.sDirName)

  err := os.MkdirAll(rs.sDirName, 0777)

  if(err != nil){
    msg := fmt.Sprintf("Error creating directory[%s]", err.Error())
    logmsg.Print(logmsg.Error, msg)
    // fmt.Println(msg)
    return false
  }

  var fullname string
  var sStructname string
  var sEndpointname string

  if(rs.bRawURL){
    fullname = rs.sDirName + "/" + "RawFile"
    sStructname = "RawStruct"
    sEndpointname = "RawEndpoint"

  }else{

    fullname = rs.sDirName + "/" + rs.mEndpoints[r.GetName()].Filename
    sStructname = rs.mEndpoints[r.GetName()].Structname
    sEndpointname = rs.mEndpoints[r.GetName()].Endpoint
  }

  // r.Dump()

  // fmt.Println("SaveResponseBody - " + "fullname: " + fullname)
  // fmt.Println("SaveResponseBody - " + "sStructname: " + sStructname)
  // fmt.Println("SaveResponseBody - " + "sEndpointname: " + sEndpointname)

  r.SaveResponseBody(fullname, sStructname, true)

  if(rs.PackageInfo.GetSave()){ // they took the time to setup the package info

    fullname = fullname+".go"

    input, err := ioutil.ReadFile(fullname)
    if err != nil {
      logmsg.Print(logmsg.Error, err)
      return false
    }

    lines := strings.Split(string(input), "\n")

    //newstuff := fmt.Sprintf("package powerwall\n\n"+
    newstuff := fmt.Sprintf("\n\nimport \"github.com/seldonsmule/logmsg\"\n\n"+
                  "func (pP *%s) Get%s() (bool, %s){\n\n"+
                  "  var s %s\n\n"+
                  "  pP.SetObject(&s)\n\n"+
                  "  if(!pP.GetStruct(\"%s\", false)){\n"+
                  "    logmsg.Print(logmsg.Error, \"GetStruct(%s) failed\")\n"+
                  "    return false, s\n"+
                  "  }\n\n"+
                  "  return true, s\n\n"+
                  "}\n\n", rs.PackageInfo.GetClassName(), sStructname, sStructname, sStructname, sEndpointname, sEndpointname)

    structstart := fmt.Sprintf("type %s ", sStructname)


    for i, line := range lines {
      if strings.Contains(line, "package") {
        lines[i] = "package " + rs.PackageInfo.GetPackageName()
      }

      if strings.Contains(line, structstart){
        lines[i-1] = newstuff
      }
    }

    output := strings.Join(lines, "\n")
    err = ioutil.WriteFile(fullname, []byte(output), 0644)
    if err != nil {
      logmsg.Print(logmsg.Error, err)
      return false
    }

  } // end if !bRawGetStruct

  

  return true
}

// ListSendSave - send all the endpoints and save the response body

func (rs *RestStruct) GetIndex() []string{

  keys := make([]string, 0, len(rs.mEndpoints))

  for k := range rs.mEndpoints {

    keys = append(keys, k)

  } 

  sort.Strings(keys)

  return keys

}

func (rs *RestStruct) Send(r *restapi.Restapi, bSave bool) bool{

  r.JsonOnly()

  if(!r.Send()){
    logmsg.Print(logmsg.Error, "Error sending: ", r.GetUrl())
    return false
  }

  // fmt.Println("Done send going to save")

  // call SaveResponseBody stuff here
  if(bSave){
    if(!rs.SaveResponseBody(r)){
      logmsg.Print(logmsg.Error, "Error saving response body: ", r.GetUrl())
      return false
    }
  }

  if(rs.oObject != nil){	
    json.Unmarshal(r.BodyBytes, rs.oObject)
  }

  return true

} 


func (rs *RestStruct) CheckEndpoint(sEndpointName string) bool{

	_, ok := rs.mEndpoints[sEndpointName]

	return ok

}

