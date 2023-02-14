package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"personal-web/connec"
	"personal-web/middleware"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

type MetaData struct {
	Title     string
	IsLogin   bool
	UserName  string
	FlashData string
}

var Data = map[string]interface{}{
	"Title":   "Personal Web",
	"IsLogin": false,
}

type User struct {
	Id       int
	name     string
	email    string
	password string
}

type Blog struct {
	Id           int
	StartDate    time.Time
	EndDate      time.Time
	StartDateStr string
	EndDateStr   string
	Title        string
	Author       string
	Content      string
	Deference    string
	Images       string
	//
	Thecnologies []string
	TechOne      bool
	TechTwo      bool
	TechTre      bool
	TechFor      bool
}

var BlogData = []Blog{}

func main() {

	route := mux.NewRouter()

	connec.DatabaseConnect()

	// crated static files
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	// route.HandleFunc("/helloworld", helloWorld).Methods("GET")
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("GET")
	route.HandleFunc("/add-project/blog", middleware.UploadFile(projectBlog)).Methods("POST")
	route.HandleFunc("/contact-me", contact).Methods("GET")
	route.HandleFunc("/blog-content/{id}", blogContent).Methods("GET")
	route.HandleFunc("/delete-content/{id}", deleteBlog).Methods("GET")
	route.HandleFunc("/edit-content/{id}", editBlog).Methods("GET")
	route.HandleFunc("/update-content/blog{id}", updateButton).Methods("POST")
	route.HandleFunc("/register", register).Methods("GET")
	route.HandleFunc("/register-button", registerButton).Methods("POST")
	route.HandleFunc("/login", login).Methods("GET")
	route.HandleFunc("/login", loginButton).Methods("POST")
	route.HandleFunc("/logout/clear", logoutButton).Methods("GET")

	fmt.Println("server running on port : 5000")
	http.ListenAndServe("localhost:5000", route)
}

// logout button
func logoutButton(w http.ResponseWriter, r *http.Request) {

	store := sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func loginButton(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		log.Fatal(err)
		return
	}

	User := User{}
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	sqlUser := "SELECT id, name, email, password FROM public.tb_users WHERE email=$1"

	rows := connec.Conn.QueryRow(context.Background(), sqlUser, email).Scan(&User.Id, &User.name, &User.email, &User.password)
	if rows != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		log.Fatal(err)
		return
	}
	comparePw := bcrypt.CompareHashAndPassword([]byte(User.password), []byte(password))
	if comparePw != nil {
		w.Write([]byte("message :" + err.Error()))
		log.Fatal(comparePw)
		fmt.Println("Password not Macth")
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.Values["IsLogin"] = true
	session.Values["Id"] = User.Id
	session.Values["Name"] = User.name
	session.Options.MaxAge = 1000

	session.AddFlash("login succes")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

// register button handle
func registerButton(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	fmt.Println(name, email, password)
	// hasing password
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	sqlInsertUser := "INSERT INTO public.tb_users( name, email, password) VALUES ( $1, $2, $3);"
	_, eror := connec.Conn.Exec(context.Background(), sqlInsertUser, name, email, passwordHash)

	if eror != nil {
		log.Fatal(eror)
		fmt.Println("can't input data")
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

// login handle func
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text-html;charset=utf-8")
	tmpl, err := template.ParseFiles("./views/login.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}
	resp := map[string]interface{}{
		"Title": Data,
	}
	tmpl.Execute(w, resp)
}

// register handle func
func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text-html;charset=utf-8")
	tmpl, err := template.ParseFiles("./views/register.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return

	}
	resp := map[string]interface{}{
		"Title": Data,
	}
	tmpl.Execute(w, resp)
}

// update Button func handler
func updateButton(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)
	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// getting value form
	projectName := r.PostForm.Get("projectName")
	projectBlog := r.PostForm.Get("projectBlog")
	// getting Date input
	startDateInput := r.PostForm.Get("startDate")
	endDateInput := r.PostForm.Get("endDate")
	// getting checkBox input
	tech := r.Form["thecnologies"]

	sqlUpdate := "UPDATE public.tb_project SET project_name=$1, start_date=$2, end_date=$3, content=$4, technologies=$5, images='img-2' WHERE id=$6;"

	_, Error := connec.Conn.Exec(context.Background(), sqlUpdate, projectName, startDateInput, endDateInput, projectBlog, tech, id)

	if Error != nil {
		fmt.Println("data failed to update")
		log.Fatal(Error)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// editblog Handler
func editBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text-html;charset=utf-8")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	fmt.Println(id)
	// parsing template html
	tmpl, err := template.ParseFiles("./views/edit-project.html")

	BlogEdit := Blog{}
	rows := connec.Conn.QueryRow(context.Background(), "SELECT id, project_name, start_date, end_date, content, technologies FROM public.tb_project WHERE id=$1", id).Scan(&BlogEdit.Id, &BlogEdit.Title, &BlogEdit.StartDate, &BlogEdit.EndDate, &BlogEdit.Content, &BlogEdit.Thecnologies)

	BlogEdit.StartDateStr = BlogEdit.StartDate.Format("2006-01-02")
	BlogEdit.EndDateStr = BlogEdit.EndDate.Format("2006-01-02")

	for _, data := range BlogEdit.Thecnologies {
		if data == "reactLogo" {
			BlogEdit.TechOne = true
		} else if data == "nodeLogo" {

			BlogEdit.TechTwo = true
		} else if data == "goLogo" {
			BlogEdit.TechTre = true

		} else if data == "laravelLogo" {
			BlogEdit.TechFor = true

		}
	}

	if rows != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return

	}

	resp := map[string]interface{}{
		"Title": Data,
		"data":  BlogEdit,
	}
	tmpl.Execute(w, resp)

}

// delete func
func deleteBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text-html;charset=utf-8")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	_, err := connec.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Delete data succes")

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func projectBlog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		return
	}
	projectName := r.PostForm.Get("projectName")
	projectBlog := r.PostForm.Get("projectBlog")
	startDateInput := r.PostForm.Get("startDate")
	endDateInput := r.PostForm.Get("endDate")
	tech := r.Form["thecnologies"]

	store := sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	author := session.Values["Id"].(int)

	dataContext := r.Context().Value("dataFile")
	image := dataContext.(string)

	sqlInsert := "INSERT INTO public.tb_project( project_name, start_date, end_date, content, technologies, images, author_id ) VALUES ( $1, $2, $3, $4,$5,$6 ,$7);"

	_, Error := connec.Conn.Exec(context.Background(), sqlInsert, projectName, startDateInput, endDateInput, projectBlog, tech, image, author)
	if Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}
	fmt.Println("Data berhasil dimasukan")
	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

// blog content handler
func blogContent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text-html;charset=utf-8")

	// parsing template html
	var tmpl, err = template.ParseFiles("./views/blog-content.html")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}

	BlogDetail := Blog{}
	rows := connec.Conn.QueryRow(context.Background(), "SELECT id, project_name, start_date, end_date, content,technologies,images FROM public.tb_project WHERE id=$1", id).Scan(&BlogDetail.Id, &BlogDetail.Title, &BlogDetail.StartDate, &BlogDetail.EndDate, &BlogDetail.Content, &BlogDetail.Thecnologies, &BlogDetail.Images)

	BlogDetail.StartDateStr = BlogDetail.StartDate.Format("02-01-2006")
	BlogDetail.EndDateStr = BlogDetail.EndDate.Format("02-01-2006")

	if rows != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return

	}

	for _, thc := range BlogDetail.Thecnologies {
		if thc == "reactLogo" {
			BlogDetail.TechOne = true
		} else if thc == "nodeLogo" {
			BlogDetail.TechTwo = true
		} else if thc == "goLogo" {
			BlogDetail.TechTre = true
		} else if thc == "laravelLogo" {
			BlogDetail.TechFor = true
		}
	}

	resp := map[string]interface{}{
		"Title": Data,
		"data":  BlogDetail,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

// add project handler
func addProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text-html;charset=utf-8")

	var tmpl, err = template.ParseFiles("./views/add-project.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}
	resp := map[string]interface{}{
		"Title": Data,
		"Data":  Data,
	}
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

// contact me handler
func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html;charset=utf-8")

	var tmpl, err = template.ParseFiles("./views/contact-me.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}
	fmt.Println(Data["IsLogin"])

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)

}

// home handler
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	store := sessions.NewCookieStore([]byte("SESSION_ID"))
	Session, _ := store.Get(r, "SESSION_ID")

	if Session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
	} else {
		Data["IsLogin"] = Session.Values["IsLogin"].(bool)
		Data["UserName"] = Session.Values["Name"].(string)
	}

	var tmpl, err = template.ParseFiles("./views/home.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message :" + err.Error()))
		return
	}
	var result []Blog
	var rows pgx.Rows
	var sqlLeftJoin string
	if Data["IsLogin"] == true {
		sqlLeftJoin = "SELECT tb_project.id,project_name, start_date,end_date,content,technologies,images,tb_users.name as author_id FROM public.tb_project LEFT JOIN tb_users ON tb_project.author_id = tb_users.id WHERE tb_users.name = $1"
		rows, _ = connec.Conn.Query(context.Background(), sqlLeftJoin, Data["UserName"])
	} else {
		sqlLeftJoin = "SELECT tb_project.id,project_name, start_date,end_date,content,technologies,images,tb_users.name as author_id FROM public.tb_project LEFT JOIN tb_users ON tb_project.author_id = tb_users.id"
		rows, _ = connec.Conn.Query(context.Background(), sqlLeftJoin)
	}

	for rows.Next() {
		each := Blog{}
		var err = rows.Scan(&each.Id, &each.Title, &each.StartDate, &each.EndDate, &each.Content, &each.Thecnologies, &each.Images, &each.Author)
		each.Deference = deferentTime(each.StartDate, each.EndDate)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		result = append(result, each)

	}

	resp := map[string]interface{}{
		"Title": Data,
		"Data":  result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)

}

// func defenrent time
func deferentTime(start, end time.Time) string {

	// final value
	var finalValue string
	//
	deferent := end.Sub(start)
	valueDay := deferent.Hours() / 24

	if valueDay > 0 && valueDay < 30 {
		fd := strconv.FormatFloat(deferent.Hours()/(24), 'f', 0, 64)
		finalValue = "Duration :" + fd + " Day"
	} else if valueDay > 30 && valueDay < 365 {
		fd := strconv.FormatFloat(deferent.Hours()/(24*30), 'f', 0, 64)
		finalValue = "Duration :" + fd + " Bulan"
	} else if valueDay > 365 {
		fd := strconv.FormatFloat(deferent.Hours()/(24*30*12), 'f', 0, 64)
		finalValue = "Duration :" + fd + " tahun"
	} else {
		finalValue = "Your time is up !!!"
	}
	return finalValue
}
