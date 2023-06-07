package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	fmt.Println("Running...")
	http.HandleFunc("/insert", insert)
	http.HandleFunc("/request", request)
	http.HandleFunc("/update", update)
	http.HandleFunc("/sendtoapproval", sendtoapproval)
	http.HandleFunc("/delete", delete)
	http.ListenAndServe(":8001", nil)

}

func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	defer cancel()
	defer func() {

		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// Mongodb Connection----------------------------------------------------------------------------------------------------------------

func connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

func ping(client *mongo.Client, ctx context.Context) error {

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil
}

func mongoconnect() {
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer close(client, ctx, cancel)
	ping(client, ctx)
}

// Inserting Data---------------------------------------------------------------------------------------------------------------------------

func query(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {
	collection := client.Database("Employee_Info").Collection("Leave_Details")
	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

func insertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {
	collection := client.Database("Employee_Info").Collection("Leave_Details")
	result, err := collection.InsertOne(ctx, doc)
	return result, err

}

func insert(w http.ResponseWriter, r *http.Request) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	if r.Body == nil {
		errors.New("HttpRequest object is not properly setup")
	}
	r.ParseForm()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "OPTIONS" {
		return
	}
	bs, e := ioutil.ReadAll(r.Body)
	if e != nil {
		fmt.Errorf("Unable to read body: " + e.Error())
	}
	defer r.Body.Close()
	var br *bytes.Reader
	br = bytes.NewReader(bs)
	decoder := json.NewDecoder(br)

	payload := struct {
		EmployeeId   string
		Gender       string
		LeaveType    string
		Leavebalance int
		FromDate     string
		ToDate       string
		NoofDays     int
		Status       string
	}{}

	edecode := decoder.Decode(&payload)
	if edecode != nil {
		fmt.Errorf("Payload Decode Error: " + edecode.Error() + " .Bytes Data: " + string(bs))
	}
    dateString := payload.FromDate
    fdate, err := time.Parse(time.RFC3339, dateString)
    if err != nil {
        panic(err)
    }
    fdate = fdate.Add(5*time.Hour + 30*time.Minute)
    //----
    tdateString := payload.ToDate
    tdate, err := time.Parse(time.RFC3339, tdateString)
    if err != nil {
        panic(err)
    }
    tdate = tdate.Add(5*time.Hour + 30*time.Minute)

	fmt.Println(payload)
	userscollection := client.Database("Employee_Info").Collection("Leave_Details")
	userscollectionsecond := client.Database("Employee_Info").Collection("Leavedbs_Details")

	var documents interface{} = bson.D{
		{"EmployeeId", payload.EmployeeId},
		{"Gender", payload.Gender},
		{"LeaveType", payload.LeaveType},
		{"Leavebalance", payload.Leavebalance},
		{"FromDate", fdate},
		{"ToDate", tdate},
		{"NoofDays", payload.NoofDays},
		{"Status", payload.Status},
	}
	result, err := userscollection.InsertOne(context.TODO(), documents)
	if err != nil {
		panic(err)
	}
	resultsecond, err := userscollectionsecond.InsertOne(context.TODO(), documents)
	if err != nil {
		panic(err)
	}
	fmt.Println("Result of InsertOne")
	fmt.Println(result.InsertedID)
	findval, _ := json.Marshal(payload)
	fmt.Fprintf(w, "%v", string(findval))

	fmt.Println("Result of InsertOne")
	fmt.Println(resultsecond.InsertedID)
	findvals, _ := json.Marshal(payload)
	fmt.Fprintf(w, "%v", string(findvals))

}

// Get--------------------------------------------------------------------------------------------------------------------------------------

func request(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		return
	}
	type payload struct {
		EmployeeId   string
		Gender       string
		LeaveType    string
		Leavebalance int
		FromDate     time.Time
		ToDate       time.Time
		NoofDays     int
		Status       string
	}
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer close(client, ctx, cancel)

	var filter, option interface{}
	filter = bson.D{{}}
	option = bson.D{{"_id", 0}}

	var results []payload

	cursor, err := query(client, ctx, "Employee_Info", "Leave_Details", filter, option)

	if err != nil {
		panic(err)
	}
	if err := cursor.All(ctx, &results); err != nil {
		panic(err)
	}
	senddata, _ := json.Marshal(results)
	fmt.Fprintf(w, string(senddata))

}

// update --------------------------------------------------------------------------------------------------------------------------------------

func Updateone(client *mongo.Client, ctx context.Context, dataBase, col string, filter, update interface{}) (result *mongo.UpdateResult, err error) {

	collection := client.Database("Employee_Info").Collection("Leave_Details")

	result, err = collection.UpdateOne(ctx, filter, update)
	return
}
func Updateoness(client *mongo.Client, ctx context.Context, dataBase, col string, filter, update interface{}) (result *mongo.UpdateResult, err error) {

	collection := client.Database("Employee_Info").Collection("Leavedbs_Details")

	result, err = collection.UpdateOne(ctx, filter, update)
	return
}

func Updateones(client *mongo.Client, ctx context.Context, dataBase, col string, filter, update interface{}) (result *mongo.UpdateResult, err error) {

	collection := client.Database("Employee_Info").Collection("emp_leave_info")

	result, err = collection.UpdateOne(ctx, filter, update)
	return
}

func update(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		return
	}
	payload := struct {
		EmployeeId   string
		Gender       string
		LeaveType    string
		Leavebalance int
		FromDate     time.Time
		ToDate       time.Time
		NoofDays     int
		Status       string
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)
	filter := bson.D{
		{"FromDate", bson.D{{"$eq", payload.FromDate}}},
	}

	var leavebal int
	if payload.Status == "Approved!" {
		leavebal = payload.Leavebalance - payload.NoofDays
	} else {
		leavebal = payload.Leavebalance
	}

	update := bson.D{
		{"$set", bson.D{
			{"Leavebalance", leavebal},
			{"FromDate", payload.FromDate},
			{"ToDate", payload.ToDate},
			{"Status", payload.Status},
		}},
	}

	result, err := Updateone(client, ctx, "Employee_Info", "Leave_Details", filter, update)
	if err != nil {
		panic(err)
	}

//--------------------------------------------------------------------------------------------------------
	
filterss := bson.D{
		{"FromDate", bson.D{{"$eq", payload.FromDate}}},
	}

	updatess := bson.D{
		{"$set", bson.D{
			{"EmployeeId", payload.EmployeeId},
			{"Gender", payload.Gender},
			{"LeaveType", payload.LeaveType},
			{"Leavebalance", leavebal},
			{"FromDate", payload.FromDate},
			{"ToDate", payload.ToDate},
			{"NoofDays", payload.NoofDays},
			{"Status", payload.Status},
		}},
	}

	resultss, err := Updateoness(client, ctx, "Employee_Info", "Leavedbs_Details", filterss, updatess)
	if err != nil {
		panic(err)
	}

//-------------------------------------------------------------------------------------------------------------
	
	filters := bson.D{
		{"Empid", bson.D{{"$eq", payload.EmployeeId}}},
	}

	updateLeave := bson.D{}

	if payload.LeaveType == "Casual Leave" {
		updateLeave = bson.D{
			{"$set", bson.D{
				{"CL", leavebal},
			}},
		}
	} else if payload.LeaveType == "Sick Leave" {
		updateLeave = bson.D{
			{"$set", bson.D{
				{"SL", leavebal},
			}},
		}
	} else {
		updateLeave = bson.D{
			{"$set", bson.D{
				{"PL", leavebal},
			}},
		}
	}

	result1, err := Updateones(client, ctx, "Employee_Info", "emp_leave_info", filters, updateLeave)
	if err != nil {
		panic(err)
	}

	fmt.Println("updated single document")
	findval, _ := json.Marshal(result)
	fmt.Println(string(findval))
	fmt.Fprintf(w, "Status Submitted... ")

	fmt.Println("updated single document")
	findvalss, _ := json.Marshal(resultss)
	fmt.Println(string(findvalss))
	fmt.Fprintf(w, "Status Submitted... ")

	fmt.Println("updated single document")
	findvals, _ := json.Marshal(result1)
	fmt.Println(string(findvals))
	fmt.Fprintf(w, "Status Submitted... ")

}

// ------------------------------------------------------------------------------------------------------------------------

func querys(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {
	collection := client.Database("Employee_Info").Collection("Leavedbs_Details")
	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

func sendtoapproval(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		return
	}
	type payload struct {
		EmployeeId   string
		Gender       string
		LeaveType    string
		Leavebalance int
		FromDate     time.Time
		ToDate       time.Time
		NoofDays     int
		Status       string
	}
	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer close(client, ctx, cancel)

	var filter, option interface{}
	filter = bson.D{{}}
	option = bson.D{{"_id", 0}}

	var results []payload

	cursor, err := querys(client, ctx, "Employee_Info", "Leavedbs_Details", filter, option)

	if err != nil {
		panic(err)
	}
	if err := cursor.All(ctx, &results); err != nil {
		panic(err)
	}
	senddata, _ := json.Marshal(results)
	fmt.Fprintf(w, string(senddata))

}

//------------------------------------------------------------------------------------------------

func deleteOne(client *mongo.Client, ctx context.Context, dataBase, col string, query interface{}) (result *mongo.DeleteResult, err error) {

	collection := client.Database("Employee_Info").Collection("Leavedbs_Details")

	result, err = collection.DeleteOne(ctx, query)
	return
}

func delete(w http.ResponseWriter, r *http.Request) {

	if r.Body == nil {
		errors.New("HttpRequest object is not properly setup")
	}
	r.ParseForm()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	if r.Method == "OPTIONS" {
		return
	}

	bs, e := ioutil.ReadAll(r.Body)
	if e != nil {
		fmt.Errorf("Unable to read body: " + e.Error())
	}
	defer r.Body.Close()
	var br *bytes.Reader
	br = bytes.NewReader(bs)
	decoder := json.NewDecoder(br)

	payload := struct {
		EmployeeId string
		Status     string
	}{}
	edecode := decoder.Decode(&payload)
	if edecode != nil {
		fmt.Errorf("Payload Decode Error: " + edecode.Error() + " .Bytes Data: " + string(bs))
	}

	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	defer close(client, ctx, cancel)
	if payload.Status == "Approved!" || payload.Status == "Rejected!" {
		query := bson.D{
			{"EmployeeId", bson.D{{"$eq", payload.EmployeeId}}},
		}
		fmt.Println(query)

		result, err := deleteOne(client, ctx, "Employee_Info", "Leavedbs_Details", query)

		if err != nil {
			panic(err)
		}

		fmt.Println("No.of rows affected by DeleteOne()", result.DeletedCount)
	} else {
		fmt.Println("No()")
	}

}

// -----------------------------------------------------------------------------------------------------------------------------------------