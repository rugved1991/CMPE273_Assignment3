package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "encoding/json"
    "strings"
    "io/ioutil"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "strconv"
    "sort"
    "bytes"
    "log"
)


type Products struct {
    Lat  float64
    Lng float64
    Products  []Product `json:"products"`
}


type Product struct {
    Prod_id   string `json:"product_id"`
    Description string `json:"description"`
    Disp_name string `json:"display_name"`
    Capacity    int    `json:"capacity"`
    Image       string `json:"image"`
}



func (pe *PriceEsts) get(c *Client) error {
    priceEstimateParams := map[string]string{
        "start_latitude":  strconv.FormatFloat(pe.StartLat, 'f', 2, 32),
        "start_longitude": strconv.FormatFloat(pe.StartLong, 'f', 2, 32),
        "end_latitude":    strconv.FormatFloat(pe.EndLat, 'f', 2, 32),
        "end_longitude":   strconv.FormatFloat(pe.EndLong, 'f', 2, 32),
    }

    data := c.getRequest("estimates/price", priceEstimateParams)
    if e := json.Unmarshal(data, &pe); e != nil {
        return e
    }
    return nil
}

const (
    
    APIUrl string = "https://sandbox-api.uber.com/v1/%s%s"
)


type Getter interface {
    get(c *Client) error
}


type RequestOptions struct {
    ServerToken    string
    ClientId       string
    ClientSecret   string
    AppName        string
    AuthorizeUrl   string
    AccessTokenUrl string
    AccessToken string
    BaseUrl        string
}

type PriceEsts struct {
    StartLat  float64
    StartLong float64
    EndLat    float64
    EndLong   float64
    Prices         []PriceEst `json:"prices"`
}


type PriceEst struct {
    Prod_id       string  `json:"product_id"`
    Currency_code    string  `json:"currency_code"`
    Disp_name     string  `json:"display_name"`
    Estimate        string  `json:"estimate"`
    LowEstim     int     `json:"low_estimate"`
    HighEstim    int     `json:"high_estimate"`
    SurgeMultiplier float64 `json:"surge_multiplier"`
    Duration        int     `json:"duration"`
    Distance        float64 `json:"distance"`
}

type Client struct {
    Options *RequestOptions
}


func Create(options *RequestOptions) *Client {
    return &Client{options}
}


func (c *Client) Get(getter Getter) error {
    if e := getter.get(c); e != nil {
        return e
    }

    return nil
}


func (c *Client) getRequest(endpoint string, params map[string]string) []byte {
    urlParams := "?"
    params["server_token"] = c.Options.ServerToken
    for k, v := range params {
        if len(urlParams) > 1 {
            urlParams += "&"
        }
        urlParams += fmt.Sprintf("%s=%s", k, v)
    }

    url := fmt.Sprintf(APIUrl, endpoint, urlParams)

    res, err := http.Get(url)
    if err != nil {
        
    }

    data, err := ioutil.ReadAll(res.Body)
    res.Body.Close()

    return data
}





func (pl *Products) get(c *Client) error {
    productParams := map[string]string{
        "latitude":  strconv.FormatFloat(pl.Lat, 'f', 2, 32),
        "longitude": strconv.FormatFloat(pl.Lng, 'f', 2, 32),
    }

    data := c.getRequest("products", productParams)
    if e := json.Unmarshal(data, &pl); e != nil {
        return e
    }
    return nil
}



type reqObj struct{
Id int
Name string `json:"Name"`
Address string `json:"Address"`
City string `json:"City"`
State string `json:"State"`
Zip string `json:"Zip"`
Coordinates struct{
    Lat float64
    Lng float64
}
}

var id int;
var tripId int;

type TripResponse struct {
    ID                     string   `json:"id"`
    Status                 string   `json:"status"`
    StartingFromLocationID string   `json:"starting_from_location_id"`
    BestRouteLocationIds   []string `json:"best_route_location_ids"`
    TotalUberCosts         int      `json:"total_uber_costs"`
    TotalUberDuration      int      `json:"total_uber_duration"`
    TotalDistance          float64  `json:"total_distance"`
}

type OnGoingTrip struct {
    ID                        string   `json:"id"`
    Status                    string   `json:"status"`
    StartingFromLocationID    string   `json:"starting_from_location_id"`
    NextDestinationLocationID string   `json:"next_destination_location_id"`
    BestRouteLocationIds      []string `json:"best_route_location_ids"`
    TotalUberCosts            int      `json:"total_uber_costs"`
    TotalUberDuration         int      `json:"total_uber_duration"`
    TotalDistance             float64  `json:"total_distance"`
    UberWaitTimeEta           int      `json:"uber_wait_time_eta"`
}


type Responz struct {
    Results []struct {
        AddressComponents []struct {
            LongName  string   `json:"long_name"`
            ShortName string   `json:"short_name"`
            Types     []string `json:"types"`
        } `json:"address_components"`
        FormattedAddress string `json:"formatted_address"`
        Geometry         struct {
            Location struct {
                Lat float64 `json:"lat"`
                Lng float64 `json:"lng"`
            } `json:"location"`
            LocationType string `json:"location_type"`
            Viewport     struct {
                Northeast struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"northeast"`
                Southwest struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"southwest"`
            } `json:"viewport"`
        } `json:"geometry"`
        PartialMatch bool     `json:"partial_match"`
        PlaceID      string   `json:"place_id"`
        Types        []string `json:"types"`
    } `json:"results"`
    Status string `json:"status"`
}



type RideRequest struct {
    EndLat    string `json:"end_latitude"`
    EndLong   string `json:"end_longitude"`
    Prod_id      string `json:"product_id"`
    StartLat  string `json:"start_latitude"`
    StartLong string `json:"start_longitude"`
}


type ReqResponse struct {
    Driver          interface{} `json:"driver"`
    Eta             int         `json:"eta"`
    Location        interface{} `json:"location"`
    RequestID       string      `json:"request_id"`
    Status          string      `json:"status"`
    SurgeMultiplier int         `json:"surge_multiplier"`
    Vehicle         interface{} `json:"vehicle"`
}




type resObj struct{
Greeting string
}

func main() {
    mux := httprouter.New()
    

    id=0;
    tripId=0;
    currentPos=0;
    ogtID=0;
    mux.POST("/locations",createlocation)
    mux.POST("/trips",plantrip)
    mux.GET("/locations/:locid",getloc)
    mux.GET("/trips/:tripid",gettrip)
    mux.PUT("/locations/:locid",updateloc)
    mux.PUT("/trips/:tripid/request",requesttrip)
    mux.DELETE("/locations/:locid",deleteloc)
    server := http.Server{
            Addr:        "0.0.0.0:8083",
            Handler: mux,
    }

    server.ListenAndServe()
}



func createlocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    id=id+1;


    decoder := json.NewDecoder(req.Body)
    var t reqObj 
    t.Id = id; 
    err := decoder.Decode(&t)
    if err != nil {
        fmt.Println("Error")
    }


    
    st:=strings.Join(strings.Split(t.Address," "),"+");
    fmt.Println(st);
    constr := []string {strings.Join(strings.Split(t.Address," "),"+"),strings.Join(strings.Split(t.City," "),"+"),t.State}
    lstringplus := strings.Join(constr,"+")
    locstr := []string{"http://maps.google.com/maps/api/geocode/json?address=",lstringplus}
    
    resp, err := http.Get(strings.Join(locstr,""))
    
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
       fmt.Println("Error: Wrong address");
     }
     var data Responz
    err = json.Unmarshal(body, &data)
    fmt.Println(data.Status)
    
    t.Coordinates.Lat=data.Results[0].Geometry.Location.Lat;
    t.Coordinates.Lng=data.Results[0].Geometry.Location.Lng;




 conn, err := mgo.Dial("mongodb://rugved:rugved@ds045454.mongolab.com:45454/rugved")

    
    if err != nil {
        panic(err)
    }
    defer conn.Close();

conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("rugved").C("assignment2");
err = c.Insert(t);


    js,err := json.Marshal(t)
    if err != nil{
	   fmt.Println("Error")
	   return
	}
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}



type modReqObj struct{
    Address string `json:"address"`
    City string `json:"city"`
    State string `json:"state"`
    Zip string `json:"zip"`
}

func updateloc(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    
 id ,err1:= strconv.Atoi(p.ByName("locid"))
 
 if err1 != nil {
         panic(err1)
     }
  conn, err := mgo.Dial("mongodb://rugved:rugved@ds045454.mongolab.com:45454/rugved")


     if err != nil {
         panic(err)
     }
     defer conn.Close();

conn.SetMode(mgo.Monotonic,true);
 c:=conn.DB("rugved").C("assignment2");


     decoder := json.NewDecoder(req.Body)
     var t modReqObj  
     err = decoder.Decode(&t)
     if err != nil {
         fmt.Println("Error")
     }


     colQuerier := bson.M{"id": id}
     change := bson.M{"$set": bson.M{"address": t.Address, "city":t.City,"state":t.State,"zip":t.Zip}}
     err = c.Update(colQuerier, change)
     if err != nil {
         panic(err)
     }

}


func deleteloc(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
     id ,err1:= strconv.Atoi(p.ByName("locid"))
 
 if err1 != nil {
         panic(err1)
     }
  conn, err := mgo.Dial("mongodb://rugved:rugved@ds045454.mongolab.com:45454/rugved")
  conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("rugved").C("assignment2");


     if err != nil {
         panic(err)
     }
     defer conn.Close();
     err=c.Remove(bson.M{"id":id})
     if err != nil { fmt.Printf("Could not find kitten %s to delete", id)}
    rw.WriteHeader(http.StatusNoContent)
}




type userUber struct {
    LocationIds            []string `json:"location_ids"`
    StartingFromLocationID string   `json:"starting_from_location_id"`
}

func plantrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

    decoder := json.NewDecoder(req.Body)
    var uUD userUber 
    err := decoder.Decode(&uUD)
    if err != nil {
        log.Println("Error")
    }

        log.Println(uUD.StartingFromLocationID);



    var options RequestOptions;
    options.ServerToken= "";
    options.ClientId= "";
    options.ClientSecret= "";
    options.AppName= "CMPE273TransitApp";
    options.BaseUrl= "https://sandbox-api.uber.com/v1/";
    

    client :=Create(&options); 


        sid ,err1:= strconv.Atoi(uUD.StartingFromLocationID)
 fmt.Println(uUD.StartingFromLocationID);
 fmt.Println(sid);
 if err1 != nil {
         panic(err1)
     }

    conn, err := mgo.Dial("mongodb://rugved:rugved@ds045454.mongolab.com:45454/rugved");

    
    if err != nil {
        panic(err)
    }
    defer conn.Close();

    conn.SetMode(mgo.Monotonic,true);
    c:=conn.DB("rugved").C("assignment2");
    result:=reqObj{}
    err = c.Find(bson.M{"id":sid}).One(&result)
    if err != nil {
                fmt.Println(err)
        }

    
    index:=0;
    totalPrice := 0;
    totalDistance :=0.0;
    totalDuration :=0;
    bestroute:=make([]float64,len(uUD.LocationIds));
    m := make(map[float64]string)

    for _,ids := range uUD.LocationIds{
    
        lid,err1:= strconv.Atoi(ids)
            
        if err1 != nil {
            panic(err1)
        }
        

        resultLID:=reqObj{}
        err = c.Find(bson.M{"id":lid}).One(&resultLID)
        if err != nil {
             fmt.Println(err)
        }
        pe := &PriceEsts{}
        pe.StartLat = result.Coordinates.Lat;
        pe.StartLong = result.Coordinates.Lng;
        pe.EndLat = resultLID.Coordinates.Lat;
        pe.EndLong = resultLID.Coordinates.Lng;

        if e := client.Get(pe); e != nil {
            fmt.Println(e);
        }
        totalDistance=totalDistance+pe.Prices[0].Distance;
        totalDuration=totalDuration+pe.Prices[0].Duration;
        totalPrice=totalPrice+pe.Prices[0].LowEstim;
        bestroute[index]=pe.Prices[0].Distance;
        m[pe.Prices[0].Distance]=ids;
        index=index+1;
    }
    
    sort.Float64s(bestroute);
    


    var tripres TripResponse;

    tripId=tripId+1;

     tripres.ID=strconv.Itoa(tripId);
     tripres.TotalDistance=totalDistance;
     tripres.TotalUberCosts=totalPrice;
     tripres.TotalUberDuration=totalDuration;
     tripres.Status="Planning";
     tripres.StartingFromLocationID=strconv.Itoa(sid);
     tripres.BestRouteLocationIds=make([]string,len(uUD.LocationIds));
     index=0;
     for _, ind := range bestroute{
        tripres.BestRouteLocationIds[index]=m[ind];
        index=index+1;
     }
     fmt.Println(tripres.BestRouteLocationIds[1]);

     

    c1:=conn.DB("rugved").C("trips");
    err = c1.Insert(tripres);

    
        js,err := json.Marshal(tripres)
    if err != nil{
       fmt.Println("Error")
       return
    }
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)

    }


func requesttrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    

    
    kid ,err1:= strconv.Atoi(p.ByName("tripid"))
    var siD int;
    
    if err1 != nil {
         panic(err1)
     }
    var ogt OnGoingTrip;
    result1:=reqObj{}
    result2:=reqObj{}
    conn, err := mgo.Dial("mongodb://rugved:rugved@ds045454.mongolab.com:45454/rugved")

    
    if err != nil {
        panic(err)
    }
    defer conn.Close();

    conn.SetMode(mgo.Monotonic,true);
    c:=conn.DB("rugved").C("trips");
    result:=TripResponse{}

    err = c.Find(bson.M{"id":strconv.Itoa(kid)}).One(&result)
    if err != nil {
        fmt.Println(err)
    }else{

    var iD int;

    c1:=conn.DB("rugved").C("assignment2");
    if currentPos==0{
        iD, err = strconv.Atoi(result.StartingFromLocationID)
        siD=iD;
        if err != nil {
        
            fmt.Println(err)
        }
    }else
    {
        iD, err = strconv.Atoi(result.BestRouteLocationIds[currentPos-1])
        siD=iD;
        if err != nil {
        
            fmt.Println(err)
        }
    }

    err = c1.Find(bson.M{"id":iD}).One(&result1)
    if err != nil {
                fmt.Println(err)
        }
    iD, err = strconv.Atoi(result.BestRouteLocationIds[currentPos])
    if err != nil {
        
        fmt.Println(err)
    }
    err = c1.Find(bson.M{"id":iD}).One(&result2)
    if err != nil {
                fmt.Println(err)
        }


        fmt.Println(result2.Coordinates.Lat);
    }

    ogt.ID=strconv.Itoa(ogtID);
    ogt.BestRouteLocationIds=result.BestRouteLocationIds;
    ogt.StartingFromLocationID=strconv.Itoa(siD);
    ogt.NextDestinationLocationID=result.BestRouteLocationIds[currentPos];
    ogt.TotalDistance=result.TotalDistance;
    ogt.TotalUberCosts=result.TotalUberCosts;
    ogt.TotalUberDuration=result.TotalUberDuration;
    ogt.Status="requesting";
    
    

    var options RequestOptions;
    options.ServerToken= "";
    options.ClientId= "";
    options.ClientSecret= "";
    options.AppName= "CMPE273TransitApp";
    options.BaseUrl= "https://sandbox-api.uber.com/v1/";
    client :=Create(&options);

    pl:=Products{};
    pl.Lat=result1.Coordinates.Lat;
    pl.Lng=result1.Coordinates.Lng;
    if e := pl.get(client); e != nil {
         fmt.Println(e)
    }
    var prodid string;
    i:=0
    for _, product := range pl.Products {
         if(i == 0){
             prodid = product.Prod_id
        }
    }



    var rr RideRequest;

    rr.StartLat=strconv.FormatFloat(result1.Coordinates.Lat, 'f', 6, 64);
    rr.StartLong=strconv.FormatFloat(result1.Coordinates.Lng, 'f', 6, 64);
    rr.EndLat=strconv.FormatFloat(result2.Coordinates.Lat, 'f', 6, 64);
    rr.EndLong=strconv.FormatFloat(result2.Coordinates.Lng, 'f', 6, 64);
    rr.Prod_id=prodid;
    buf, _ := json.Marshal(rr)
    body := bytes.NewBuffer(buf)
    url := fmt.Sprintf(APIUrl, "requests?","access_token=")
    res, err := http.Post(url,"application/json",body)
    if err != nil {
        fmt.Println(err)
    }
    data, err := ioutil.ReadAll(res.Body)
    var rRes ReqResponse;
    err = json.Unmarshal(data, &rRes)
    ogt.UberWaitTimeEta=rRes.Eta;

    js,err := json.Marshal(ogt)
    if err != nil{
       fmt.Println("Error")
       return
    }
    ogtID=ogtID+1;
    currentPos=currentPos+1;
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)

}



func gettrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

    conn, err := mgo.Dial("mongodb://rugved:rugved@ds045454.mongolab.com:45454/rugved")

    
    if err != nil {
        panic(err)
    }
    defer conn.Close();

    conn.SetMode(mgo.Monotonic,true);
    c:=conn.DB("rugved").C("trips");
    result:=TripResponse{}
    err = c.Find(bson.M{"id":p.ByName("tripid")}).One(&result)
    if err != nil {
        fmt.Println(err)
    }

    
    js,err := json.Marshal(result)
    if err != nil{
       fmt.Println("Error")
       return
    }
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}


var currentPos int;
var ogtID int;






func getloc(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
fmt.Println(p.ByName("locid"));
id ,err1:= strconv.Atoi(p.ByName("locid"))
if err1 != nil {
        panic(err1)
    }
 conn, err := mgo.Dial("mongodb://rugved:rugved@ds045454.mongolab.com:45454/rugved")

    
    if err != nil {
        panic(err)
    }
    defer conn.Close();

conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("rugved").C("assignment2");
result:=reqObj{}
err = c.Find(bson.M{"id":id}).One(&result)
if err != nil {
                fmt.Println(err)
        }

        
        js,err := json.Marshal(result)
    if err != nil{
       fmt.Println("Error")
       return
    }
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}