package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"myapp/models"
	"myapp/repositories"
	"myapp/utils/logic"
	"myapp/utils/validation"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

/*
 Todoリストを取得
*/
func GetAllTodos(w http.ResponseWriter, userId int) ([]models.BaseTodoResponse, error) {
	var todos []models.Todo
	// todoリストデータ取得
	if err := repositories.GetAllTodos(&todos, userId); err != nil {
		logic.SendResponse(w, logic.CreateErrorStringResponse("データ取得に失敗"), http.StatusInternalServerError)
		return nil, err
	}
	// レスポンス用の構造体に変換
	responseTodos := logic.CreateAllTodoResponse(&todos)

	return responseTodos, nil
}

/*
 IDに紐づくTodoを取得
*/
func GetTodoById(w http.ResponseWriter, r *http.Request, userId int) (models.BaseTodoResponse, error) {
	// getパラメータからIDを取得
	vars := mux.Vars(r)
    id := vars["id"]
	var todo models.Todo
	// todoデータ取得処理
	if err := repositories.GetTodoById(&todo, id, userId); err != nil {
		var errMessage string
		var statusCode int
		// https://gorm.io/ja_JP/docs/error_handling.html
		if errors.Is(err, gorm.ErrRecordNotFound) {
			statusCode = http.StatusBadRequest
			errMessage = "該当データは存在しません。"
		} else {
			statusCode = http.StatusInternalServerError
			errMessage = "データ取得に失敗しました。"
		}
		// エラーレスポンス送信
		logic.SendResponse(w, logic.CreateErrorStringResponse(errMessage), statusCode)
		return models.BaseTodoResponse{}, err
	}

	// レスポンス用の構造体に変換
	responseTodos := logic.CreateTodoResponse(&todo)

	return responseTodos, nil
}

/*
 Todo新規登録処理
*/
func CreateTodo(w http.ResponseWriter, r *http.Request, userId int) (models.BaseTodoResponse, error) {
	// ioutil: ioに特化したパッケージ
    reqBody,_ := ioutil.ReadAll(r.Body)
	var mutationTodoRequest models.MutationTodoRequest
	if err := json.Unmarshal(reqBody, &mutationTodoRequest); err != nil {
        log.Fatal(err)
        errMessage := "リクエストパラメータを構造体へ変換処理でエラー発生"
		logic.SendResponse(w, logic.CreateErrorStringResponse(errMessage), http.StatusInternalServerError)
		return models.BaseTodoResponse{}, err
	}
	// バリデーション
	if err := validation.MutationTodoValidate(mutationTodoRequest); err != nil {
		// バリデーションエラーのレスポンスを送信
		logic.SendResponse(w, logic.CreateErrorResponse(err), http.StatusBadRequest)
		return models.BaseTodoResponse{}, err
	}

	var todo models.Todo
    todo.Title = mutationTodoRequest.Title
    todo.Comment = mutationTodoRequest.Comment
    todo.UserId = userId

	// todoデータ新規登録処理
	if err := repositories.CreateTodo(&todo); err != nil {
		logic.SendResponse(w, logic.CreateErrorStringResponse("データ取得に失敗しました。"), http.StatusInternalServerError)
		return models.BaseTodoResponse{}, err
	}

	// 登録したtodoデータ取得処理
	if err := repositories.GetTodoLastByUserId(&todo, userId); err != nil {
		var errMessage string
		var statusCode int
		// https://gorm.io/ja_JP/docs/error_handling.html
		if errors.Is(err, gorm.ErrRecordNotFound) {
			statusCode = http.StatusBadRequest
			errMessage = "該当データは存在しません。"
		} else {
			statusCode = http.StatusInternalServerError
			errMessage = "データ取得に失敗しました。"
		}
		// エラーレスポンス送信
		logic.SendResponse(w, logic.CreateErrorStringResponse(errMessage), statusCode)
		return models.BaseTodoResponse{}, err
	}

	// レスポンス用の構造体に変換
	responseTodos := logic.CreateTodoResponse(&todo)

	return responseTodos, nil
}

/*
 Todo削除処理
*/
func DeleteTodo(w http.ResponseWriter, r *http.Request, userId int) error {
	// getパラメータからIDを取得
	vars := mux.Vars(r)
    id := vars["id"]
	// データ削除処理
	if err := repositories.DeleteTodo(id, userId); err != nil {
		logic.SendResponse(w, logic.CreateErrorStringResponse("データ削除に失敗"), http.StatusInternalServerError)
		return err
	}
	return nil
}

/*
 Todo更新処理
*/
func UpdateTodo(w http.ResponseWriter, r *http.Request, userId int) (models.BaseTodoResponse, error) {
	// GetパラメータからIDを取得
	vars := mux.Vars(r)
    id := vars["id"]
	// request bodyから値を取得
    reqBody, _ := ioutil.ReadAll(r.Body)

    var mutationTodoRequest models.MutationTodoRequest
	if err := json.Unmarshal(reqBody, &mutationTodoRequest); err != nil {
        fmt.Print("======")
		log.Fatal(err)
        errMessage := "リクエストパラメータを構造体へ変換処理でエラー発生"
		logic.SendResponse(w, logic.CreateErrorStringResponse(errMessage), http.StatusInternalServerError)
		return models.BaseTodoResponse{}, err
	}
	// バリデーション
	if err := validation.MutationTodoValidate(mutationTodoRequest); err != nil {
		// バリデーションエラーのレスポンスを送信
		logic.SendResponse(w, logic.CreateErrorResponse(err), http.StatusBadRequest)
		return models.BaseTodoResponse{}, err
	}

	// 更新用データ用意
	var updateTodo models.Todo
	updateTodo.Title = mutationTodoRequest.Title
	updateTodo.Comment = mutationTodoRequest.Comment

	// todoデータ新規登録処理
	if err := repositories.UpdateTodo(&updateTodo, id, userId); err != nil {
		logic.SendResponse(w, logic.CreateErrorStringResponse("データ更新に失敗しました。"), http.StatusInternalServerError)
		return models.BaseTodoResponse{}, err
	}

	// 更新データを取得
	var todo models.Todo
	if err := repositories.GetTodoById(&todo, id, userId); err != nil {
		var errMessage string
		var statusCode int
		if errors.Is(err, gorm.ErrRecordNotFound) {
			statusCode = http.StatusBadRequest
			errMessage = "該当データは存在しません。"
		} else {
			statusCode = http.StatusInternalServerError
			errMessage = "データ取得に失敗しました。"
		}
		// エラーレスポンス送信
		logic.SendResponse(w, logic.CreateErrorStringResponse(errMessage), statusCode)
		return models.BaseTodoResponse{}, err
	}

	// レスポンス用の構造体に変換
	responseTodos := logic.CreateTodoResponse(&todo)

	return responseTodos, nil
}


/*
 Todoリスト取得APIのレスポンス送信処理
*/
func SendAllTodoResponse(w http.ResponseWriter, todos *[]models.BaseTodoResponse) {
	var response models.AllTodoResponse
	response.Todos = *todos
	// レスポンスデータ作成
	responseBody, _ := json.Marshal(response)

	// レスポンス送信
	logic.SendResponse(w, responseBody, http.StatusOK)
}

/*
 Todoデータのレスポンス送信処理
*/
func SendTodoResponse(w http.ResponseWriter, todo *models.BaseTodoResponse) {
	var response models.TodoResponse
	response.Todo = *todo
	// レスポンスデータ作成
	responseBody, _ := json.Marshal(response)
	// レスポンス送信
	logic.SendResponse(w, responseBody, http.StatusOK)
}

/*
 CreateTodoAPIのレスポンス送信処理
*/
func SendCreateTodoResponse(w http.ResponseWriter, todo *models.BaseTodoResponse) {
	var response models.TodoResponse
	response.Todo = *todo
	// レスポンスデータ作成
	responseBody, _ := json.Marshal(response)
	// レスポンス送信
	logic.SendResponse(w, responseBody, http.StatusCreated)
}

/*
 DeleteTodoAPIのレスポンス送信処理
*/
func SendDeleteTodoResponse(w http.ResponseWriter) {
	// レスポンス送信
	logic.SendNotBodyResponse(w)
}