package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	gomail "gopkg.in/gomail.v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	fmt.Println("Running...")
	http.HandleFunc("/request", request)
	http.HandleFunc("/getempid", getempid)
	http.HandleFunc("/requestempdetails", requestempdetails)
	http.HandleFunc("/getemail", getemail)
	http.HandleFunc("/getgridview", getgridview)

	http.ListenAndServe(":8007", nil)
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

// Mongodb Connection ----------------------------------------------------------------------------------------------------------------

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

// Get --------------------------------------------------------------------------------------------------------------------------------

func query(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {
	collection := client.Database("Employee_Info").Collection("emp_leave_info")
	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

func request(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		return
	}
	type payload struct {
		Empid  string
		Name   string
		Gender string
		CL     int
		SL     int
		PL     int
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

	cursor, err := query(client, ctx, "Employee_Info", "emp_leave_info", filter, option)

	if err != nil {
		panic(err)
	}
	if err := cursor.All(ctx, &results); err != nil {
		panic(err)
	}
	senddata, _ := json.Marshal(results)
	fmt.Fprintf(w, string(senddata))

}

//---------------------------------------------------------------------------------------------------------------------------------------------

func querys(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {
	collection := client.Database("Employee_Info").Collection("emp_leave_info")
	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

func getempid(w http.ResponseWriter, r *http.Request) {

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

	// type receivedata struct {
	// 	Empid  string
	// 	Name   string
	// 	Gender string
	// 	CL     string
	// 	SL     string
	// 	PL     string
	// }

	client, ctx, cancel, err := connect("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	defer close(client, ctx, cancel)
	ping(client, ctx)

	var filter, option interface{}

	payload := struct {
		Empid string
	}{}

	fmt.Println("@#", payload.Empid)
	edecode := decoder.Decode(&payload)
	if edecode != nil {
		fmt.Errorf("Payload Decode Error: " + edecode.Error() + " .Bytes Data: " + string(bs))
	}

	filter = bson.D{{"Empid", bson.D{{"$eq", payload.Empid}}}}
	// filter = bson.D{{}}

	// filter = bson.D{{"Empid",payload.EmployeeId}}

	// filterCursor, err := emp_leave_info.Find(ctx, bson.M{"Empid": payload.EmployeeId})

	option = bson.D{{"_id", 0}}

	cursor, err := querys(client, ctx, "Employee_Info", "emp_leave_info", filter, option)

	if err != nil {
		panic(err)
	}

	var results []bson.D
	fmt.Println("Result", results)

	if err := cursor.All(ctx, &results); err != nil {
		panic(err)
	}
	fmt.Println("Query Result")
	for _, doc := range results {
		fmt.Println(doc)
	}
	findval, _ := json.Marshal(results)
	fmt.Fprintf(w, string(findval))

}

// --------------------------------------------------------------------------------------------------------------------------------------------------
func queryss(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {
	collection := client.Database("Employee_Info").Collection("Emp_details")
	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

func requestempdetails(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		return
	}
	type payload struct {
		Name    string
		Gender  string
		Skills  []string
		Emailid string
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

	cursor, err := queryss(client, ctx, "Employee_Info", "Emp_details", filter, option)

	if err != nil {
		panic(err)
	}
	if err := cursor.All(ctx, &results); err != nil {
		panic(err)
	}
	senddata, _ := json.Marshal(results)
	fmt.Fprintf(w, string(senddata))

}

// --------------------------------------------------------------------------------------------------------------------------------------

// func querye(client *mongo.Client, ctx context.Context,
// 	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {
// 	collection := client.Database("Employee_Info").Collection("Emp_details")
// 	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
// 	return
// }

// func queryy(client *mongo.Client, ctx context.Context,
// 	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {
// 	collection := client.Database("Employee_Info").Collection("Grid_View")
// 	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
// 	return
// }

// func insertOnee(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {
// 	collection := client.Database("Employee_Info").Collection("Grid_View")
// 	result, err := collection.InsertOne(ctx, doc)
// 	return result, err

// }

// func getemail(w http.ResponseWriter, r *http.Request) {

// 	r.ParseForm()
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

// 	if r.Method == "OPTIONS" {
// 		return
// 	}

// 	bs, e := ioutil.ReadAll(r.Body)
// 	if e != nil {
// 		fmt.Errorf("Unable to read body: " + e.Error())
// 	}

// 	defer r.Body.Close()
// 	var br *bytes.Reader
// 	br = bytes.NewReader(bs)
// 	decoder := json.NewDecoder(br)

// 	client, ctx, cancel, err := connect("mongodb://localhost:27017")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer close(client, ctx, cancel)
// 	ping(client, ctx)

// 	var filter, option interface{}

// 	payload := struct {
// 		Audience    string
// 		Selecttype  string
// 		Title       string
// 		Postedon    string
// 		Attachments string
// 	}{}

// 	edecode := decoder.Decode(&payload)
// 	if edecode != nil {
// 		fmt.Errorf("Payload Decode Error: " + edecode.Error() + " .Bytes Data: " + string(bs))
// 	}

// 	if payload.Audience == "Gender" {
// 		filter = bson.D{{"Gender", bson.D{{"$eq", payload.Selecttype}}}}
// 	} else if payload.Audience == "Skills" {
// 		filter = bson.D{{"Skills", bson.D{{"$eq", payload.Selecttype}}}}
// 	}

// 	option = bson.D{{"_id", 0}, {"Gender", 0}, {"Name", 0}, {"Skills", 0}}

// 	cursor, err := querye(client, ctx, "Employee_Info", "Emp_details", filter, option)

// 	if err != nil {
// 		panic(err)
// 	}

// 	var results []bson.D
// 	fmt.Println("Result", results)

// 	if err := cursor.All(ctx, &results); err != nil {
// 		panic(err)
// 	}

// 	var mailarr []string
// 	fmt.Println("Query Result")
// 	for _, doc := range results {
// 		fmt.Println(doc)

// 		var i interface{}
// 		i = doc

// 		mailid := fmt.Sprintf("%v", i)
// 		split := strings.Split(mailid, " ")
// 		replace := strings.Replace(split[1], "}]", "", 1)
// 		mailarr = append(mailarr, replace)
// 		fmt.Println("mai##", mailid)
// 	}
// 	// fmt.Println("###" , mailarr[1])
// 	fmt.Println("###", mailarr)

// 	findval, _ := json.Marshal(results)
// 	fmt.Fprintf(w, string(findval))

// 	// Email sent function ---------------------------------------------------------------------------------------------------------------------

// 	for i := 0; i < len(mailarr); i++ {
// 		msg := gomail.NewMessage()
// 		msg.SetHeader("From", "<amizhthan@gmail.com>")
// 		msg.SetHeader("To", mailarr[i])
// 		msg.SetHeader("Subject", payload.Title)
// 		msg.SetBody("text/html", "<h1>mail send task</h1><p>finish the task this evening</p>")
// 		msg.Attach(payload.Attachments)
// 		n := gomail.NewDialer("smtp.freesmtpservers.com", 25, "<amizhthan@gmail.com>", "<password>")
// 		// Send the email
// 		if err := n.DialAndSend(msg); err != nil {
// 			panic(err)
// 		} else {
// 			fmt.Println("Mail Sent Successfully to", mailarr[i])
// 		}
// 	}

// Grid View Get ----------------------------------------------------------------------------------------------------------------------------------

func queryys(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {
	collection := client.Database("Employee_Info").Collection("Grid_View")
	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

func getgridview(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "OPTIONS" {
		return
	}
	type payload struct {
		Audience   string
		Types      string
		Title      string
		Postedon   string
		Postedby   string
		Attachment string
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

	cursor, err := queryys(client, ctx, "Employee_Info", "Grid_View", filter, option)

	if err != nil {
		panic(err)
	}
	if err := cursor.All(ctx, &results); err != nil {
		panic(err)
	}
	senddata, _ := json.Marshal(results)
	fmt.Fprintf(w, string(senddata))

}

// -----------------------------------------------------------------------------------------------------------------------------------

func insertOnee(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {
	collection := client.Database("Employee_Info").Collection("Grid_View")
	result, err := collection.InsertOne(ctx, doc)
	return result, err
}

func querye(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {
	collection := client.Database("Employee_Info").Collection("Emp_details")
	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

var filter, option interface{}

func getemail(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			fmt.Println("errroorr", err)
			return
		}
		defer file.Close()
		f, err := os.OpenFile("C:/Users/Amizhthan.MURUGAN/OneDrive - Talentpro INDIA HR Pvt Ltd/Vuejs_Tasks/my_app/public/images/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("errooooooooo", err)
			return
		}

		defer f.Close()
		io.Copy(f, file)
		Filename := handler.Filename
		// fmt.Println("Filename = ", Filename)

		Title := r.PostForm.Get("title")
		Audience := r.PostForm.Get("audience")
		Types := r.PostForm.Get("selecttype")
		CC := r.PostForm.Get("cc")
		Body := r.PostForm.Get("body")

		client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
		if err != nil {
			log.Fatal(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)

		if Audience == "Gender" {
			filter = bson.D{{"Gender", bson.D{{"$eq", Types}}}}
		} else if Audience == "Skills" {
			filter = bson.D{{"Skills", bson.D{{"$eq", Types}}}}
		}

		option = bson.D{{"_id", 0}, {"Gender", 0}, {"Name", 0}, {"Skills", 0}}

		cursor, err := querye(client, ctx, "Employee_Info", "Emp_details", filter, option)

		if err != nil {
			panic(err)
		}

		var results []bson.D
		fmt.Println("Result", results)

		if err := cursor.All(ctx, &results); err != nil {
			panic(err)
		}

		var mailarr []string
		fmt.Println("Query Result")
		for _, doc := range results {
			fmt.Println(doc)

			var i interface{}
			i = doc

			mailid := fmt.Sprintf("%v", i)
			split := strings.Split(mailid, " ")
			replace := strings.Replace(split[1], "}]", "", 1)
			mailarr = append(mailarr, replace)
			fmt.Println("##mail##", mailid)
		}

// Mail Send--------------------------------------------------------------------------------------------------------------------------

		for i := 0; i < len(mailarr); i++ {

			m := gomail.NewMessage()
			m.SetHeader("From", "amizhthan@gmail.com")
			m.SetHeader("To", mailarr[i], CC)
			m.SetHeader("Subject", Title)
			m.SetBody("text/html", Body)
			m.Attach("C:/Users/Amizhthan.MURUGAN/OneDrive - Talentpro INDIA HR Pvt Ltd/Vuejs_Tasks/my_app/public/images/" + Filename)

			n := gomail.NewDialer("smtp.freesmtpservers.com", 25, "<amizhthan@gmail.com>", "<pass1234>")

			if err := n.DialAndSend(m); err != nil {
				panic(err)
			} else {
				fmt.Println("Mail Sent Successfully to", mailarr[i])
			}
		}

// Insert Into Grid View----------------------------------------------------------------------------------------------------------------

		userscollection := client.Database("Employee_Info").Collection("Grid_View")

		time := time.Now()
		PostDate := time.Format("2006-01-02")
		fmt.Println("########", PostDate)
		Postedby := "Amizhthan"

		var documents interface{} = bson.D{
			{"Title", Title},
			{"Audience", Audience},
			{"Types", Types},
			{"Postedon", PostDate},
			{"Postedby", Postedby},
			{"Attachment", Filename},
		}
		result, err := userscollection.InsertOne(context.TODO(), documents)
		if err != nil {
			panic(err)
		}
		fmt.Println("Result of InsertOne")
		fmt.Println(result.InsertedID)
	}
}

