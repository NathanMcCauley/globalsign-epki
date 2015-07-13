package main

import (
        "fmt"
        "os"
        "os/user"
        "io/ioutil"

        "github.com/Sirupsen/logrus"

        gowsdl "github.com/hooklift/gowsdl/generator"
        epki "./epki_service"
        "gopkg.in/inconshreveable/log15.v2"

)

func init() {
        logrus.SetLevel(logrus.DebugLevel)
        logrus.SetOutput(os.Stderr)
        // Retrieve current user to get home directory
        usr, err := user.Current()
        if err != nil {
                fatalf("cannot get current user: %v", err)
        }

        // Get home directory for current user
        homeDir := usr.HomeDir
        if homeDir == "" {
                fatalf("cannot get current user home directory")
        }
}

func getOrderRequestHeader() (*epki.OrderRequestHeader ) {
        return &epki.OrderRequestHeader{
          AuthToken: &epki.AuthToken{
            UserName: "username",
            Password: "password",
          },
        };
}

func getCsr() (string) {
        bytes, _ := ioutil.ReadFile("fixtures/test.csr")
        return string(bytes);
}

func getOrderAndIssueCertificateRequest() (*epki.OrderAndIssueCertificate) {
      oaiRequest := &epki.OrderAndIssueRequest{
         OrderRequestHeader: getOrderRequestHeader(),
         ProfileID: "MP201507110518", //test ProfileID assigned by globalsign
         ProductCode: "EPKIPSPersonal",
         Year: 1,
         CSR: getCsr(),
         EFSOption: false,
         UPN: "blahblahblah", //Microsoft smartcard thing, probably not needed?
         DnAttributes: &epki.DnAttributes{
           CommonName: "default_common_name",
           OrganizationUnit: []string {"default_ou"},
           //Email: "default@example.com",
         },
         PickupPassword: "ponies11",
         CertificateTemplate: &epki.CertificateTemplate{
           Template: "blah",
           MajorVersion: "1",
           MinorVersion: "2",
         },
       }
       oaicRequest := &epki.OrderAndIssueCertificate{Request: oaiRequest}
       return oaicRequest;
}

func getOrderPkcs12Request() (*epki.OrderPkcs12) {
      p12oRequest := &epki.Pkcs12OrderRequest{
         OrderRequestHeader: getOrderRequestHeader(),
         ProfileID: "MP201507110518", //test ProfileID assigned by globalsign
         ProductCode: "EPKIPSPersonal",
         PKCS12PIN: "ponies11",
         Year: 1,
         EFSOption: false,
         Renew: false,
         UPN: "", //Microsoft smartcard thing, probably not needed?
         DnAttributes: &epki.DnAttributes{
           CommonName: "default_common_name",
           OrganizationUnit: []string {"default_ou", "secondary_ou"},
           //Email: "default@example.com",
         },
       }
       op12Request := &epki.OrderPkcs12{Request: p12oRequest}
       return op12Request;
}

func main() {
        gowsdl.Log.SetHandler(log15.StreamHandler(os.Stdout, log15.TerminalFormat()))

        service := epki.NewGasOrderService("https://testsystem.globalsign.com/cr/ws/GasOrderService", false)

        //fmt.Println("service.OrderAndIssueCertificate------------------------");
        //oaicr, err := service.OrderAndIssueCertificate(getOrderAndIssueCertificateRequest());

        //if err != nil {
        //  //fatalf("\n->%s\n", err.(*gowsdl.SoapFault).Error);
        //  fatalf("\n->%s\n", err);
        //}

        //fmt.Printf("oaicr %v\n"   , *oaicr);
        //fmt.Printf("oair %x\n"    , oaicr.Response);
        //fmt.Printf("oair %+v\n"    , oaicr.Response);
        //fmt.Printf("Code %s\n"    , oaicr.Response.OrderResponseHeader.SuccessCode);
        //fmt.Println("Errors");
        //for _,err := range oaicr.Response.OrderResponseHeader.Errors.Error {
        //    // element is the element from someSlice for where we are
        //    fmt.Println("code "  + err.ErrorCode);
        //    fmt.Println("field " + err.ErrorField);
        //    fmt.Println("msg "   + err.ErrorMessage);
        //}
        //fmt.Printf("Time %s\n"    , oaicr.Response.OrderResponseHeader.Timestamp);
        //fmt.Printf("OrderID %s\n" , oaicr.Response.OrderID);
        //fmt.Printf("CERT %s\n"    , oaicr.Response.CERT);

        fmt.Println("service.OrderPkcs12 ------------------------");

        op12r, err := service.OrderPkcs12(getOrderPkcs12Request());

        if err != nil {
          //fatalf("\n->%s\n", err.(*gowsdl.SoapFault).Error);
          fatalf("\n->%s\n", err);
        }

        fmt.Printf("op12r %v\n"    , *op12r);
        fmt.Printf("p12or %x\n"    ,  op12r.Response);
        fmt.Printf("p12or %+v\n"   ,  op12r.Response);
        fmt.Printf("Code  %s\n"    ,  op12r.Response.OrderResponseHeader.SuccessCode);
        fmt.Println("Errors");
        for _,err := range op12r.Response.OrderResponseHeader.Errors.Error {
            // element is the element from someSlice for where we are
            fmt.Println("code "  + err.ErrorCode);
            fmt.Println("field " + err.ErrorField);
            fmt.Println("msg "   + err.ErrorMessage);
        }
        fmt.Printf("Time %s\n"    , op12r.Response.OrderResponseHeader.Timestamp);
        fmt.Printf("OrderID %s\n" , op12r.Response.OrderID);
        fmt.Printf("CERT %s\n"    , op12r.Response.BASE64PKCS12);
}

func fatalf(format string, args ...interface{}) {
        fmt.Printf("* fatal: "+format+"\n", args...)
        os.Exit(1)
}
